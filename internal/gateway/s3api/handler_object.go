package s3api

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/go-chi/chi/v5"
	"github.com/minio/minio/cmd"

	"github.com/ocean-rw/ocean/internal/lib/httputil"
	"github.com/ocean-rw/ocean/pkg/proto"
)

func (m *Mgr) PutObject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bucketID := chi.URLParam(r, "bucket_id")
	objectID := chi.URLParam(r, "object_id")

	h := md5.New()
	tr := io.TeeReader(r.Body, h)
	defer r.Body.Close()
	fd, size, err := m.storage.Put(r.Context(), tr)
	if err != nil {
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}

	mt, err := mimetype.DetectReader(r.Body)
	if err != nil {
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}

	fileInfo := &proto.FileInfo{
		BID:      bucketID,
		FID:      objectID,
		Size:     size,
		Status:   proto.FileStatusNormal,
		Hash:     hex.EncodeToString(h.Sum(nil)),
		MimeType: mt.String(),
		Meta:     nil,
		FD:       fd,
	}

	err = m.db.FileTable.Upsert(ctx, fileInfo)
	if err != nil {
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("ETag", fileInfo.Hash)
	httputil.ReplyXML(w, http.StatusOK, nil)
}

func (m *Mgr) GetObject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bucketID := chi.URLParam(r, "bucket_id")
	objectID := chi.URLParam(r, "object_id")

	fileInfo, err := m.db.FileTable.Get(ctx, bucketID, objectID)
	if err != nil {
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}

	data, err := m.storage.Get(ctx, fileInfo.FD)
	if err != nil {
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}
	defer data.Close()

	w.Header().Set("ETag", fileInfo.Hash)
	w.Header().Set(httputil.ContentLength, strconv.FormatInt(fileInfo.Size, 10))
	httputil.ReplyBinary(w, http.StatusOK, data)
}

func (m Mgr) ListObjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bucketID := chi.URLParam(r, "bucket_id")

	files, err := m.db.FileTable.List(ctx, bucketID)
	if err != nil {
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}

	data := generateListObjectsV2Response(bucketID, "", "", "", "", "", "", false, 0, files, nil)
	httputil.ReplyXML(w, http.StatusOK, data)
}

// generates an ListObjectsV2 response for the said bucket with other enumerated options.
func generateListObjectsV2Response(
	bucket, prefix, token, nextToken, startAfter,
	delimiter, encodingType string, isTruncated bool, maxKeys int,
	objects []*proto.FileInfo, prefixes []string,
) *cmd.ListObjectsV2Response {
	contents := make([]cmd.Object, 0, len(objects))
	owner := cmd.Owner{ID: "ocean", DisplayName: "ocean"}
	data := new(cmd.ListObjectsV2Response)

	for _, object := range objects {
		content := cmd.Object{}
		if object.FID == "" {
			continue
		}
		content.Key = s3EncodeName(object.FID, encodingType)
		content.LastModified = object.PutTime.UTC().Format(iso8601TimeFormat)
		if object.ETag != "" {
			content.ETag = "\"" + object.ETag + "\""
		}
		content.Size = object.Size
		content.StorageClass = "STANDARD"
		content.Owner = owner
		contents = append(contents, content)
	}
	data.Name = bucket
	data.Contents = contents

	data.EncodingType = encodingType
	data.StartAfter = s3EncodeName(startAfter, encodingType)
	data.Delimiter = s3EncodeName(delimiter, encodingType)
	data.Prefix = s3EncodeName(prefix, encodingType)
	data.MaxKeys = maxKeys
	data.ContinuationToken = base64.StdEncoding.EncodeToString([]byte(token))
	data.NextContinuationToken = base64.StdEncoding.EncodeToString([]byte(nextToken))
	data.IsTruncated = isTruncated

	commonPrefixes := make([]cmd.CommonPrefix, 0, len(prefixes))
	for _, prefix := range prefixes {
		prefixItem := cmd.CommonPrefix{}
		prefixItem.Prefix = s3EncodeName(prefix, encodingType)
		commonPrefixes = append(commonPrefixes, prefixItem)
	}
	data.CommonPrefixes = commonPrefixes
	data.KeyCount = len(data.Contents) + len(data.CommonPrefixes)
	return data
}

func s3EncodeName(name string, encodingType string) (result string) {
	// Quick path to exit
	if encodingType == "" {
		return name
	}
	encodingType = strings.ToLower(encodingType)
	switch encodingType {
	case "url":
		return s3URLEncode(name)
	}
	return name
}

// s3URLEncode is based on Golang's url.QueryEscape() code,
// while considering some S3 exceptions:
//	- Avoid encoding '/' and '*'
//	- Force encoding of '~'
func s3URLEncode(s string) string {
	spaceCount, hexCount := 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c) {
			if c == ' ' {
				spaceCount++
			} else {
				hexCount++
			}
		}
	}

	if spaceCount == 0 && hexCount == 0 {
		return s
	}

	var buf [64]byte
	var t []byte

	required := len(s) + 2*hexCount
	if required <= len(buf) {
		t = buf[:required]
	} else {
		t = make([]byte, required)
	}

	if hexCount == 0 {
		copy(t, s)
		for i := 0; i < len(s); i++ {
			if s[i] == ' ' {
				t[i] = '+'
			}
		}
		return string(t)
	}

	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == ' ':
			t[j] = '+'
			j++
		case shouldEscape(c):
			t[j] = '%'
			t[j+1] = "0123456789ABCDEF"[c>>4]
			t[j+2] = "0123456789ABCDEF"[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}

func shouldEscape(c byte) bool {
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return false
	}

	switch c {
	case '-', '_', '.', '/', '*':
		return false
	}
	return true
}

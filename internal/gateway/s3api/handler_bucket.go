package s3api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/minio/minio/cmd"

	"github.com/ocean-rw/ocean/internal/lib/httputil"
	"github.com/ocean-rw/ocean/pkg/proto"
)

func (m *Mgr) ListBuckets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	buckets, err := m.db.BucketTable.List(ctx, 0)
	if err != nil {
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}
	httputil.ReplyXML(w, http.StatusOK, generateListBucketsResponse(buckets))
}

// generates ListBucketsResponse from array of BucketInfo which can be
// serialized to match XML and JSON API spec output.
func generateListBucketsResponse(buckets []*proto.BucketInfo) *cmd.ListBucketsResponse {
	listbuckets := make([]cmd.Bucket, 0, len(buckets))

	owner := cmd.Owner{ID: "ocean", DisplayName: "ocean"}

	for _, bucket := range buckets {
		listbuckets = append(listbuckets, cmd.Bucket{
			Name:         bucket.BID,
			CreationDate: bucket.PutTime.UTC().Format(iso8601TimeFormat),
		})
	}

	data := new(cmd.ListBucketsResponse)
	data.Owner = owner
	data.Buckets.Buckets = listbuckets

	return data
}

func (m *Mgr) CreateBucket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bucketID := chi.URLParam(r, "bucket_id")

	bucketInfo := &proto.BucketInfo{
		BID:    bucketID,
		UID:    0,
		Status: proto.BucketStatusNormal,
	}

	err := m.db.BucketTable.Insert(ctx, bucketInfo)
	if err != nil {
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (m *Mgr) DeleteBucket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bucketID := chi.URLParam(r, "bucket_id")
	err := m.db.BucketTable.Delete(ctx, bucketID)
	if err != nil {
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}
	httputil.ReplyXML(w, http.StatusOK, nil)
}

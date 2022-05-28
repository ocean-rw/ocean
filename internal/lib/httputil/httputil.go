package httputil

import (
	"encoding/xml"
	"io"
	"net/http"
	"strconv"

	json "github.com/json-iterator/go"
)

func ReplyErr(w http.ResponseWriter, code int, err error) {
	if err != nil {
		http.Error(w, err.Error(), code)
	}
	http.Error(w, http.StatusText(code), code)
}

func ReplyBinary(w http.ResponseWriter, code int, data io.Reader) error {
	if code != 0 {
		code = http.StatusOK
	}

	w.Header().Set(ContentType, "application/octet-stream")
	w.WriteHeader(code)
	_, err := io.Copy(w, data)
	return err
}

func ReplyJSON(w http.ResponseWriter, code int, data interface{}) error {
	if code != 0 {
		code = http.StatusOK
	}

	ret, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set(ContentType, "application/json")
	w.Header().Set(ContentLength, strconv.Itoa(len(ret)))
	w.WriteHeader(code)
	_, err = w.Write(ret)
	return err
}

func ReplyXML(w http.ResponseWriter, code int, data interface{}) error {
	if code == 0 {
		code = http.StatusOK
	}

	ret, err := xml.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set(ContentType, "application/xml")
	w.Header().Set(ContentLength, strconv.Itoa(len(ret)))
	w.WriteHeader(code)
	_, err = w.Write(ret)
	return err
}

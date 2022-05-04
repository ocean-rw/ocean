package httputil

import (
	"encoding/json"
	"net/http"
)

func ReplyErr(w http.ResponseWriter, code int, err error) {
	if err != nil {
		http.Error(w, err.Error(), code)
	}
	http.Error(w, http.StatusText(code), code)
}

func ReplyData(w http.ResponseWriter, code int, data interface{}) error {
	ret, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if code != 0 {
		w.WriteHeader(code)
	}
	_, err = w.Write(ret)
	return err
}

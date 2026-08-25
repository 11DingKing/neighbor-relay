package transport

import (
	"encoding/json"
	"io"
	"net/http"
)

func Decode(r *http.Request, dst any) error {
	if r.Body == nil {
		return io.EOF
	}
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}
func Encode(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func NoContent(w http.ResponseWriter) { w.WriteHeader(http.StatusNoContent) }

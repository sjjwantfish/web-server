package web

import (
	"testing"
)

func TestServer(t *testing.T) {
	var h Server = &HTTPServer{Addr: "8080"}
	// http.ListenAndServe(":8080", h)

	h.Start()
}

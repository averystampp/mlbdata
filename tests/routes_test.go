package mlbtests

import (
	"net/http/httptest"
	"testing"

	"github.com/averystampp/mlb"
	"github.com/averystampp/sesame"
)

func TestIndex(t *testing.T) {
	server := httptest.NewServer(sesame.Handler(mlb.Index))
	defer server.Close()

	resp, err := server.Client().Get(server.URL + "/")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Server responded with status" + resp.Status)
}

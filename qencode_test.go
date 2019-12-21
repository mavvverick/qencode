package qencode

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
)

var (
	mux *http.ServeMux

	ctx = context.TODO()

	client *Client

	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient(nil, os.Getenv("QUENCODE_KEY"))
	url, _ := url.Parse("https://api.qencode.com")
	client.BaseURL = url
	client.StorageKey = os.Getenv("GOOGLE_STORAGE_KEY")
	client.StorageSecret = os.Getenv("GOOGLE_STORAGE_SECRET")
	client.CallbackURL = os.Getenv("QUENCODE_CALLBACK_URL")
	client.Bucket = os.Getenv("QUENCODE_BUCKET")
	client.WATERMARK = os.Getenv("QUENCODE_WATERMARK")
	fmt.Println(os.Getenv("QUENCODE_BUCKET"), os.Getenv("QUENCODE_CALLBACK_URL"))
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}

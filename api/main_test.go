package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	key "github.com/fairhive-labs/ethkeygen/pkg"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// countRows returns the number of generated key rows in the rendered index
// page. The template emits one "(#" marker per row.
func countRows(body string) int {
	return strings.Count(body, "(#")
}

func TestHealth(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status, got %d, want %d", w.Code, http.StatusOK)
	}
	if w.Body.String() != "ok" {
		t.Fatalf("unexpected body, got %q, want %q", w.Body.String(), "ok")
	}
}

func TestGenerateCount(t *testing.T) {
	tt := []struct {
		name string
		url  string
		want int
	}{
		{"default", "/", 10},
		{"explicit", "/?l=5", 5},
		{"floor", "/?l=0", 1},
		{"negative floor", "/?l=-3", 1},
		{"clamped to max", "/?l=99999999", maxKeys},
	}
	r := setupRouter()
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("unexpected status, got %d, want %d", w.Code, http.StatusOK)
			}
			if got := countRows(w.Body.String()); got != tc.want {
				t.Fatalf("unexpected row count, got %d, want %d", got, tc.want)
			}
		})
	}
}

func TestGenerateInvalid(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?l=not-a-number", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status, got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestBulkSignature(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/bulk-signature?n=2",
		strings.NewReader(`{"message":"hello"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("unexpected status, got %d, want %d", w.Code, http.StatusAccepted)
	}

	var resp struct {
		Total      int         `json:"total"`
		Message    string      `json:"message"`
		Signatures []signature `json:"signatures"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("cannot decode response: %v", err)
	}
	if resp.Total != 2 || len(resp.Signatures) != 2 {
		t.Fatalf("unexpected count, total=%d signatures=%d, want 2/2", resp.Total, len(resp.Signatures))
	}
	if resp.Message != "hello" {
		t.Fatalf("unexpected message, got %q, want %q", resp.Message, "hello")
	}
	for i, s := range resp.Signatures {
		if !key.ValidPrivateKey(s.PrivateKey) {
			t.Errorf("signature %d: %q is not a valid private key", i, s.PrivateKey)
		}
		if !key.ValidPublicAddress(s.Address) {
			t.Errorf("signature %d: %q is not a valid public address", i, s.Address)
		}
		if !strings.HasPrefix(s.Signature, "0x") {
			t.Errorf("signature %d: %q is not 0x-prefixed", i, s.Signature)
		}
	}
}

func TestBulkSignatureBadJSON(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/bulk-signature", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status, got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCORS(t *testing.T) {
	r := setupRouter()

	// Preflight-style OPTIONS request is short-circuited with 204.
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("unexpected OPTIONS status, got %d, want %d", w.Code, http.StatusNoContent)
	}
	if origin := w.Header().Get("Access-Control-Allow-Origin"); origin != "*" {
		t.Errorf("unexpected Allow-Origin, got %q, want %q", origin, "*")
	}
	if methods := w.Header().Get("Access-Control-Allow-Methods"); methods != "GET, PUT, OPTIONS" {
		t.Errorf("unexpected Allow-Methods, got %q, want %q", methods, "GET, PUT, OPTIONS")
	}
	// Credentials must never be advertised alongside a wildcard origin.
	if cred := w.Header().Get("Access-Control-Allow-Credentials"); cred != "" {
		t.Errorf("Allow-Credentials must be unset, got %q", cred)
	}
}

func TestFaviconETag(t *testing.T) {
	r := setupRouter()

	// First request returns the icon and an ETag.
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status, got %d, want %d", w.Code, http.StatusOK)
	}
	etag := w.Header().Get("ETag")
	if etag == "" {
		t.Fatal("missing ETag header")
	}

	// A conditional request with a matching ETag returns 304.
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)
	req.Header.Set("If-None-Match", etag)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotModified {
		t.Fatalf("unexpected conditional status, got %d, want %d", w.Code, http.StatusNotModified)
	}
}

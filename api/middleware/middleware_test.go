package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestJSONRecovery(t *testing.T) {
	r := gin.New()
	r.Use(JSONRecovery())

	r.GET("/panic", func(c *gin.Context) {
		panic("simulated critical runtime panic!")
	})

	req, _ := http.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 Internal Server Error, got %d", w.Code)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal recovery response: %v", err)
	}

	errObj, ok := res["error"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected error object, got %v", res)
	}
	if errObj["code"] != "INTERNAL_SERVER_ERROR" {
		t.Errorf("expected code INTERNAL_SERVER_ERROR, got %v", errObj["code"])
	}
}

func TestCORS(t *testing.T) {
	allowed := []string{"http://localhost:3000", "https://app.example.com"}

	r := gin.New()
	r.Use(CORS(allowed))

	r.GET("/data", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	t.Run("Allowed origin receives CORS headers", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/data", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
			t.Errorf("expected Allow-Origin http://localhost:3000, got %s", w.Header().Get("Access-Control-Allow-Origin"))
		}
	})

	t.Run("Disallowed origin does not get Allow-Origin header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/data", nil)
		req.Header.Set("Origin", "http://malicious-site.com")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Header().Get("Access-Control-Allow-Origin") != "" {
			t.Errorf("expected empty Allow-Origin for unauthorized origin, got %s", w.Header().Get("Access-Control-Allow-Origin"))
		}
	})

	t.Run("OPTIONS preflight returns 204", func(t *testing.T) {
		req, _ := http.NewRequest("OPTIONS", "/data", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204 No Content for OPTIONS, got %d", w.Code)
		}
	})
}

func TestMaxBodySize(t *testing.T) {
	r := gin.New()
	r.Use(MaxBodySize(100)) // 100 bytes limit

	r.POST("/upload", func(c *gin.Context) {
		var body map[string]interface{}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	t.Run("Request within limit succeeds", func(t *testing.T) {
		smallBody := []byte(`{"msg":"hi"}`)
		req, _ := http.NewRequest("POST", "/upload", bytes.NewBuffer(smallBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Request exceeding limit fails", func(t *testing.T) {
		largeBody := []byte(`{"msg":"` + strings.Repeat("A", 200) + `"}`)
		req, _ := http.NewRequest("POST", "/upload", bytes.NewBuffer(largeBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusRequestEntityTooLarge {
			t.Errorf("expected 413 Payload Too Large, got %d. Body: %s", w.Code, w.Body.String())
		}
	})
}

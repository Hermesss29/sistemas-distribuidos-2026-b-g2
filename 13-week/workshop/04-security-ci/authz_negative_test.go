// Negative authorization test for lms-membership-api (W13 CP4, story library-docs#41).
// Goes in internal/adapter/in/httpapi/middleware/ next to auth.go.
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testInternalSecret = "test-internal-secret"

func serviceToken(t *testing.T) string {
	t.Helper()
	now := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "circulation-service",
		"iat": now.Unix(),
		"exp": now.Add(time.Minute).Unix(),
	})
	signed, err := token.SignedString([]byte(testInternalSecret))
	if err != nil {
		t.Fatal(err)
	}
	return signed
}

// A token minted by circulation-api may read and suspend a student (HU-06, HU-08),
// but it must never deactivate, edit or list students: those are Administrator actions (HU-03).
func TestServiceTokenCannotReachAdminRoutes(t *testing.T) {
	adminRoutes := []struct{ method, path string }{
		{http.MethodPost, "/api/v1/students/5f2c8e1a-7b3d-4c9e-8a2f-1d6b0e4c3a97/deactivate"},
		{http.MethodPatch, "/api/v1/students/5f2c8e1a-7b3d-4c9e-8a2f-1d6b0e4c3a97"},
		{http.MethodGet, "/api/v1/students"},
	}
	for _, rt := range adminRoutes {
		t.Run(rt.method+" "+rt.path, func(t *testing.T) {
			reached := false
			h := RequireAuth(nil, testInternalSecret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				reached = true
				w.WriteHeader(http.StatusOK)
			}))
			req := httptest.NewRequest(rt.method, rt.path, nil)
			req.Header.Set("Authorization", "Bearer "+serviceToken(t))
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			if rec.Code != http.StatusForbidden || reached {
				t.Fatalf("service token reached %s %s: status %d, want 403", rt.method, rt.path, rec.Code)
			}
		})
	}
}

package web

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAuthHeader(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		authToken  string
	}{
		{
			"Empty Input",
			"",
			"",
		},
		{
			"Space Only Input",
			"   ",
			"",
		},
		{
			"Bearer Without Token",
			"Bearer",
			"",
		},
		{
			"Bearer Prefix With Empty Token",
			"Bearer   ",
			"",
		},
		{
			"Wrong Prefix With Token",
			"Basic token",
			"",
		},
		{
			"Wrong Prefix Case With Token",
			"BEARER token",
			"",
		},
		{
			"Bearer Prefix With New Line Token",
			"Bearer \n",
			"",
		},
		{
			"Valid Request, Odd spacing",
			"   Bearer    token   ",
			"token",
		},
		{
			"ValidRequest",
			"Bearer token",
			"token",
		},
	}
	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			token := getAuthHeader(c.authHeader)
			if c.authToken != token {
				t.Fatal("expected", c.authToken, "but got", token)
			}
		})
	}
}

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name        string
		wantedToken string
		headerToken string
		want        int
	}{
		{
			name:        "Fail",
			wantedToken: "token",
			headerToken: "BLAH",
			want:        http.StatusUnauthorized,
		},
		{
			name:        "Pass",
			wantedToken: "BLAH",
			headerToken: "BLAH",
			want:        http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.headerToken))
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			AuthMiddleware(tt.wantedToken)(nextHandler).ServeHTTP(rec, req)
			if tt.want != rec.Code {
				t.Errorf("Wanted %d, got %d", tt.want, rec.Code)
			}
		})
	}
}

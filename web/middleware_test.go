package web

import (
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

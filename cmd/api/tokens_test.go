package main

import (
	"encoding/json"
	"greenlight.bcc/internal/assert"
	"net/http"
	"testing"
)

func TestCreateToken(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	tests := []struct {
		name     string
		email    string
		password string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid",
			email:    "mus123@gmail.com",
			password: "paSSword123",
			wantCode: http.StatusCreated,
		},
		{
			name:     "test for wrong input",
			email:    "mus123@gmail.com",
			password: "paSSword123",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "failed Validation",
			email:    "mus123@gmail.com",
			password: "paSSword123",
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "email not found",
			email:    "notfound@gmail.com",
			password: "paSSword123",
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "password didn't match",
			email:    "mus123@gmail.com",
			password: "paSSword123",
			wantCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			inputData := struct {
				email    string `json:"email"`
				password string `json:"password"`
			}{
				email:    tt.email,
				password: tt.password,
			}

			b, err := json.Marshal(&inputData)
			if err != nil {
				t.Fatal("wrong input data")
			}
			if tt.name == "test for wrong input" {
				b = append(b, 'a')
			}

			code, _, _ := ts.postForm(t, "/v1/tokens/authentication", b)

			assert.Equal(t, code, tt.wantCode)

		})
	}

}

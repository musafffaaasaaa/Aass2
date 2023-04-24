package main

import (
	"encoding/json"
	"greenlight.bcc/internal/assert"
	"net/http"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	tests := []struct {
		name     string
		username string
		email    string
		password string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			username: "/v1/movies/1",
			email:    "useremail@gmail.com",
			password: "qwer123456",
			wantCode: http.StatusOK,
		},
		{
			name:     "Non-existent ID",
			username: "/v1/movies/2",
			email:    "useremail@gmail.com",
			password: "qwer123456",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			username: "/v1/movies/-1",
			email:    "useremail@gmail.com",
			password: "qwer123456",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			username: "/v1/movies/1.23",
			email:    "useremail@gmail.com",
			password: "qwer123456",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			username: "/v1/movies/foo",
			email:    "useremail@gmail.com",
			password: "qwer123456",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputData := struct {
				Name     string `json:"name"`
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Name:     tt.username,
				Email:    tt.email,
				Password: tt.password,
			}
			b, err := json.Marshal(&inputData)
			if err != nil {
				t.Fatal("wrong input data")
			}
			if tt.name == "test for wrong input" {
				b = append(b, 'a')
			}

			code, _, body := ts.postForm(t, "/v1/users", b)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}

		})
	}
}
func TestActivateUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()
	tests := []struct {
		name     string
		token    string
		wantCode int
		wantBody string
	}{
		{
			name:     "inValid",
			token:    "qweqweqwe",
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "test for wrong input",
			token:    "qweqweqwe",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Valid",
			token:    "1231231231223123123123123",
			wantCode: http.StatusOK,
		},
		{
			name:     "ErrRecordNotFound",
			token:    "1231231231223123123123123",
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "unable to update",
			token:    "1231231231223123123123123",
			wantCode: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := struct {
				token string `json:"token"`
			}{
				token: tt.token,
			}

			b, err := json.Marshal(&input)
			if err != nil {
				t.Fatal("wrong input data")
			}
			if tt.name == "test for wrong input" {
				b = append(b, 'a')
			}

			code, _, body := ts.putReq(t, "/v1/users/activated", b)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}

		})
	}
}

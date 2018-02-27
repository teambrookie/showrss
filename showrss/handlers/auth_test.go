package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/teambrookie/hermes/showrss/dao"
)

type FakeUser struct {
	Username string
	Password string
	Token    string
}

type FakeEpisodeProvider struct {
	users map[string]FakeUser
}

func (fake FakeEpisodeProvider) Auth(login, password string) (string, error) {
	user, ok := fake.users[login]
	if !ok {
		return "", errors.New("No user found")
	}
	if user.Password == password {
		return user.Token, nil
	}
	return "", errors.New("Invalid password")
}

func (fake FakeEpisodeProvider) Episodes(token string) ([]dao.Episode, error) {
	return nil, nil
}

func TestAuthHandler(t *testing.T) {
	episodeProvider := FakeEpisodeProvider{
		users: map[string]FakeUser{
			"binou": FakeUser{
				Username: "binou",
				Password: "binette",
				Token:    "tokenbinou",
			},
		},
	}
	authRequest, err := http.NewRequest("POST", "/auth", nil)
	cases := []struct {
		w                    *httptest.ResponseRecorder
		r                    *http.Request
		login                string
		password             string
		expectedResponseCode int
		expectedResponseBody []byte
	}{
		{
			w:                    httptest.NewRecorder(),
			r:                    authRequest,
			login:                "",
			password:             "",
			expectedResponseCode: http.StatusUnauthorized,
			expectedResponseBody: []byte("empty login/password\n"),
		},
		{
			w:                    httptest.NewRecorder(),
			r:                    authRequest,
			login:                "thierry",
			password:             "",
			expectedResponseCode: http.StatusUnauthorized,
			expectedResponseBody: []byte("empty login/password\n"),
		},
		{
			w:                    httptest.NewRecorder(),
			r:                    authRequest,
			login:                "brigitte",
			password:             "bardot",
			expectedResponseCode: http.StatusUnauthorized,
			expectedResponseBody: []byte("No user found\n"),
		},
		{
			w:                    httptest.NewRecorder(),
			r:                    authRequest,
			login:                "binou",
			password:             "bigornot",
			expectedResponseCode: http.StatusUnauthorized,
			expectedResponseBody: []byte("Invalid password\n"),
		},
		{
			w:                    httptest.NewRecorder(),
			r:                    authRequest,
			login:                "binou",
			password:             "binette",
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: []byte("{\"username\":\"binou\",\"token\":\"tokenbinou\"}\n"),
		},
	}

	for _, c := range cases {
		data := url.Values{}
		data.Set("username", c.login)
		data.Add("password", c.password)

		c.r.Form = data
		if err != nil {
			t.Fatalf("Auth request error: %v", err)
		}
		AuthHandler(episodeProvider).ServeHTTP(c.w, c.r)

		if c.expectedResponseCode != c.w.Code {
			t.Errorf("Expected status : code %d got : %d", c.expectedResponseCode, c.w.Code)
		}
		if !bytes.Equal(c.expectedResponseBody, c.w.Body.Bytes()) {
			t.Errorf("Expected body : %s got : %s", string(c.expectedResponseBody), string(c.w.Body.Bytes()))
		}

	}

}

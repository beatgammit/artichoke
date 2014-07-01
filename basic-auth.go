package artichoke

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"strings"
)

type Auth struct {
	User          string
	Pass          string
	Authenticated bool
}

func NewAuth(u string, p string, a bool) *Auth {
	auth := new(Auth)
	auth.User = u
	auth.Pass = p
	auth.Authenticated = a

	return auth
}

type AuthError struct {
	err string
}

func NewError(msg string) *AuthError {
	e := new(AuthError)
	e.err = msg
	return e
}

func (e *AuthError) String() string {
	return e.err
}

func GetAuth(r *http.Request) *Auth {
	if a, ok := Get(r, "auth"); ok {
		return a.(*Auth)
	}

	return nil
}

func Authenticated(r *http.Request) bool {
	if auth := GetAuth(r); auth != nil {
		return auth.Authenticated
	}

	return false
}

func BasicAuth(auth map[string]string, required bool) Middleware {
	return func(w http.ResponseWriter, r *http.Request) {
		buf := bytes.Buffer{}
		str := r.Header.Get("authorization")

		if len(str) == 0 {
			if required {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Authorization required"))
				w.Write([]byte(""))
			} else {
				Continue(r)
			}
			return
		}

		// just get the auth part
		str = strings.Split(str, " ")[1]

		i := len(str)/4*3 - strings.Count(str, "=")
		outBuf := make([]byte, len(str)/4*3)

		dec := base64.NewDecoder(base64.StdEncoding, &buf)
		buf.WriteString(str)
		dec.Read(outBuf)

		cAuth := strings.Split(string(outBuf[:i]), ":")

		user := cAuth[0]
		pass := cAuth[1]

		success := auth[user] == pass
		Set(r, "auth", NewAuth(user, pass, success))

		if success {
			if required {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Not authorized"))
				w.Write([]byte(""))
			} else {
				Continue(r)
			}
		}
	}
}

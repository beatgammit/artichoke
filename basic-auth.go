package artichoke

import (
	"net/http"
	"encoding/base64"
	"bytes"
	"strings"
)

func BasicAuth(auth map[string]string, required bool) Middleware {
	return func (w http.ResponseWriter, r *http.Request, m Data) bool {
		buf := bytes.Buffer{}
		str := r.Header.Get("authorization")

		if len(str) == 0 {
			if required {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Authorization required"))
				w.Write([]byte(""))
			}

			return required
		}

		// just get the auth part
		str = strings.Split(str, " ")[1]

		i := len(str) / 4 * 3 - strings.Count(str, "=")
		outBuf := make([]byte, len(str) / 4 * 3)

		dec := base64.NewDecoder(base64.StdEncoding, &buf)
		buf.WriteString(str)
		dec.Read(outBuf)

		cAuth := strings.Split(string(outBuf[:i]), ":")

		user := cAuth[0]
		pass := cAuth[1]

		success := auth[user] == pass
		m["auth"] = map[string]interface{} {
			"user": user,
			"pass": pass,
			"authenticated": success,
		}

		if success {
			if required {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Not authorized"))
				w.Write([]byte(""))
			}

			return required
		}

		return false
	}
}

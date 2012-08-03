package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/beatgammit/artichoke"
	"io/ioutil"
	"net/http"
	"strings"
)

type Body struct {
	Data interface{}
	Raw  []byte
	Err  error
}

func NewBody(d interface{}, r []byte, e error) *Body {
	body := new(Body)
	body.Data = d
	body.Raw = r
	body.Err = e

	return body
}

func GetBody(d artichoke.Data) *Body {
	if b, ok := d.Get("body"); ok {
		return b.(*Body)
	}

	return nil
}

func BodyParser(maxMemory int64) artichoke.Middleware {
	return func(w http.ResponseWriter, r *http.Request, d artichoke.Data) bool {
		// ignore GET and HEAD requests, they don't have useful data
		if r.Method != "PUT" && r.Method != "POST" {
			return false
		}

		s := r.Header.Get("Content-Type")
		switch {
		// parse as JSON
		case strings.Contains(s, "application/json"):
			s, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Println("Error reading body: " + err.Error())
				return false
			}

			var body interface{}
			err = json.Unmarshal(s, &body)
			if err != nil {
				fmt.Println("Error parsing JSON: " + err.Error())
			}

			d.Set("body", NewBody(body, s, err))

		case strings.Contains(s, "multipart/form-data"):
			err := r.ParseMultipartForm(maxMemory)
			if err != nil {
				fmt.Println("Error parsing body as multi-part form: " + err.Error())
			}
			d.Set("body", NewBody("", nil, err))

		case strings.Contains(s, "application/x-www-form-encoded"):
			err := r.ParseForm()
			if err != nil {
				fmt.Println("Error parsing body as form")
			}
			d.Set("body", NewBody("", nil, err))

		// treat as default handler
		case strings.Contains(s, "text/plain"):
			fallthrough
		// last resort, just read the body and pass it
		default:
			s, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Println("Error reading body: " + err.Error())
				return false
			}

			d.Set("body", NewBody(string(s), s, err))
		}
		return false
	}
}

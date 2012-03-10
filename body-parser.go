package artichoke

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

func BodyParser(maxMemory int64) Middleware {
	return func(w http.ResponseWriter, r *http.Request, d Data) bool {
		// ignore GET and HEAD requests, they don't have useful data
		if r.Method != "PUT" && r.Method != "POST" {
			return false
		}

		switch r.Header.Get("Content-Type") {
			// parse as JSON
			case "application/json":
				s, err := ioutil.ReadAll(r.Body)
				if err != nil {
					fmt.Println("Error reading body: " + err.Error())
					return false
				}

				var body interface{}
				err = json.Unmarshal(s, &body)
				if err != nil {
					fmt.Println("Error parsing JSON: " + err.Error())
					d["bodyParseError"] = err
				} else {
					d["bodyJson"] = body
				}

				d["body"] = string(s)

			case "multipart/form-data":
				err := r.ParseMultipartForm(maxMemory)
				if err != nil {
					fmt.Println("Error parsing body as multi-part form: " + err.Error())
					d["bodyParseError"] = err
				}

			case "application/x-www-form-encoded":
				err := r.ParseForm()
				if err != nil {
					fmt.Println("Error parsing body as form")
					d["bodyParseError"] = err
				}

			// treat as default handler
			case "text/plain":
				fallthrough
			// last resort, just read the body and pass it
			default:
				s, err := ioutil.ReadAll(r.Body)
				if err != nil {
					fmt.Println("Error reading body: " + err.Error())
					return false
				}

				d["body"] = string(s)
		}
		return false
	}
}

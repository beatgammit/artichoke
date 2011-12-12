package artichoke

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

/*
func printJson(j interface{}) string {
	switch j.(type) {
		case bool:
			return strconv.FormatBool(j.(bool))

		case float64:
			return strconv.FormatFloat(j.(float64), 'g', -1, 64)

		case string:
			return "\"" + j.(string) + "\""

		case []interface{}:
			t := j.([]interface{})
			s := make([]string, len(t))

			for i, val := range t {
				s[i + 1] = printJson(val)
			}

			return fmt.Sprintf("[%s]", strings.Join(s, ","))

		case map[string]interface{}:
			m := j.(map[string]interface{})
			s := make([]string, len(m))

			i := 0
			for k, v := range m {
				s[i] = fmt.Sprintf("\"%s\":%s", k, printJson(v))
				i += 1
			}

			return fmt.Sprintf("{%s}", strings.Join(s, ","))

		case nil:
			return "null"
	}

	return ""
}
*/

func BodyParser(maxMemory int64) Middleware {
	return func(w http.ResponseWriter, r *http.Request, d Data) bool {
		// ignore GET and HEAD requests
		if r.Method == "GET" || r.Method == "HEAD" {
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
				}

				d["bodyJson"] = body
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

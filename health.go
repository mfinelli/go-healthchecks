/*!
 * Copyright 2024 Mario Finelli
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package healthcheck

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

// The Health endpoint returns a more detailed state of the application health
// in a JSON response. It includes the name of the node where the application
// is running (hostname) as well as a current timestamp. If any of the
// configured healthchecks fails it will return detailed error information.
// It will always return a 200 OK response (except in the case of panic), so
// it may not be suitable for machine consumption (the Ping endpoint can be
// used for machine consumption which will return a proper status code).
func (c *Config) Health() handler {
	node := hostnameOrPanic()

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodHead, http.MethodGet:
			w.Header().Set("Content-Type", "application/json")

			err := c.Check(r.Context())
			var info string

			if err == nil {
				data, e := json.Marshal(struct {
					Node      string `json:"node"`
					Status    string `json:"status"`
					Timestamp string `json:"timestamp"`
				}{
					Node:   node,
					Status: "success",
					Timestamp: time.Now().Format(
						time.RFC3339),
				})

				if e != nil {
					panic(e)
				}

				info = string(data)
			} else {
				data, e := json.Marshal(struct {
					Debug     string `json:"debug"`
					Message   string `json:"message"`
					Node      string `json:"node"`
					Status    string `json:"status"`
					Timestamp string `json:"timestamp"`
				}{
					Debug: fmt.Sprintf("%v",
						reflect.TypeOf(err)),
					Message: err.Error(),
					Node:    node,
					Status:  "failure",
					Timestamp: time.Now().Format(
						time.RFC3339),
				})

				if e != nil {
					panic(e)
				}

				info = string(data)
			}

			if r.Method == http.MethodHead {
				w.Header().Set("Content-Length",
					strconv.Itoa(len(info)))
			}

			w.WriteHeader(http.StatusOK)

			if r.Method == http.MethodGet {
				fmt.Fprint(w, info)
			}
		case http.MethodOptions:
			w.Header().Set("Allow", "HEAD, GET, OPTIONS")
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "405 method not allowed",
				http.StatusMethodNotAllowed)
		}
	}
}

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

import "net/http"

// The Ping endping returns a status code based on the health of the service.
// It's intended to be used by machines operating on the status code and so it
// doesn't return any response body. It returns a 200 OK if the service is
// healthy and a 500 Internal Server Error otherwise. If you need
// human-readable output you should use the Health endpoint.
func (c *Config) Ping() handler {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodHead, http.MethodGet:
			w.Header().Set("Content-Type",
				"text/plain; charset=utf-8")

			err := c.Check(r.Context())

			if err == nil {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
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

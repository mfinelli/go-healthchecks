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
	"strconv"
)

// The Info route provides simple information about the service.
func (c *Config) Info() handler {
	data, err := json.Marshal(struct {
		Commit  string `json:"commit,omitempty"`
		Date    string `json:"date,omitempty"`
		Version string `json:"version,omitempty"`
	}{
		Commit:  c.GitSha,
		Date:    c.BuildDate,
		Version: c.Version,
	})

	if err != nil {
		panic(err)
	}

	info := string(data)

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodHead:
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Content-Length",
				strconv.Itoa(len(info)))
			w.WriteHeader(http.StatusOK)
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, info)
		case http.MethodOptions:
			w.Header().Set("Allow", "HEAD, GET, OPTIONS")
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "405 method not allowed",
				http.StatusMethodNotAllowed)
		}
	}
}

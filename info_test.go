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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	t.Run("non-GET requests", func(t *testing.T) {
		c := &Config{}

		for _, m := range nonGetMethods {
			req, err := http.NewRequest(m, "/", nil)
			assert.Nil(t, err)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(c.Info())
			h.ServeHTTP(w, req)

			assert.Equal(t, w.Code, http.StatusMethodNotAllowed)
		}
	})

	t.Run("the basics", func(t *testing.T) {
		c := &Config{
			Version:   "theversion",
			GitSha:    "commitsha",
			BuildDate: "thedate",
		}

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		h := http.HandlerFunc(c.Info())
		h.ServeHTTP(w, req)

		assert.Equal(t, w.Code, http.StatusOK)
		assert.Equal(t, w.Header().Get("Content-Type"),
			"application/json")
		assert.Equal(t, w.Body.String(), "{\"commit\":\"commitsha\","+
			"\"date\":\"thedate\",\"version\":\"theversion\"}")
	})

	t.Run("omits empty", func(t *testing.T) {
		tests := []struct {
			config *Config
			body   string
		}{
			{
				config: &Config{
					GitSha:    "commitsha",
					BuildDate: "thedate",
				},
				body: "{\"commit\":\"commitsha\"," +
					"\"date\":\"thedate\"}",
			},
			{
				config: &Config{
					Version:   "theversion",
					BuildDate: "thedate",
				},
				body: "{\"date\":\"thedate\"," +
					"\"version\":\"theversion\"}",
			},
			{
				config: &Config{
					Version: "theversion",
					GitSha:  "commitsha",
				},
				body: "{\"commit\":\"commitsha\"," +
					"\"version\":\"theversion\"}",
			},
		}

		for _, test := range tests {
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			assert.Nil(t, err)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(test.config.Info())
			h.ServeHTTP(w, req)

			assert.Equal(t, w.Body.String(), test.body)
		}
	})
}

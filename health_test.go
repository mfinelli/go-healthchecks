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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	t.Run("non-GET requests", func(t *testing.T) {
		c := &Config{}

		for _, m := range nonGetMethods {
			req, err := http.NewRequest(m, "/", nil)
			assert.Nil(t, err)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(c.Health())
			h.ServeHTTP(w, req)

			assert.Equal(t, w.Code, http.StatusMethodNotAllowed)
		}
	})

	t.Run("successful healthcheck", func(t *testing.T) {
		c := &Config{
			Check: func(ctx context.Context) error {
				return nil
			},
		}

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		h := http.HandlerFunc(c.Health())
		h.ServeHTTP(w, req)

		type resp struct {
			Node      string `json:"node"`
			Status    string `json:"status"`
			Timestamp string `json:"timestamp"`
		}

		var j resp
		err = json.NewDecoder(w.Body).Decode(&j)
		assert.Nil(t, err)

		assert.Equal(t, w.Code, http.StatusOK)
		assert.Equal(t, j.Status, "success")
		assert.NotEmpty(t, j.Node)
		assert.NotEmpty(t, j.Timestamp)
	})

	t.Run("unsuccessful healthcheck", func(t *testing.T) {
		c := &Config{
			Check: func(ctx context.Context) error {
				return fmt.Errorf("error!")
			},
		}

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		h := http.HandlerFunc(c.Health())
		h.ServeHTTP(w, req)

		type resp struct {
			Debug     string `json:"debug"`
			Message   string `json:"message"`
			Node      string `json:"node"`
			Status    string `json:"status"`
			Timestamp string `json:"timestamp"`
		}

		var j resp
		err = json.NewDecoder(w.Body).Decode(&j)
		assert.Nil(t, err)

		assert.Equal(t, w.Code, http.StatusOK)
		assert.Equal(t, j.Status, "failure")
		assert.Equal(t, j.Debug, "*errors.errorString")
		assert.Equal(t, j.Message, "error!")
		assert.NotEmpty(t, j.Node)
		assert.NotEmpty(t, j.Timestamp)
	})
}

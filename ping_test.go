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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	t.Run("OPTIONS request", func(t *testing.T) {
		c := &Config{}

		req, err := http.NewRequest(http.MethodOptions, "/", nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		h := http.HandlerFunc(c.Ping())
		h.ServeHTTP(w, req)

		assert.Equal(t, w.Code, http.StatusNoContent)
		assert.Equal(t, w.Header().Get("Allow"), "HEAD, GET, OPTIONS")
		assert.Empty(t, w.Body.String())
	})

	t.Run("non-GET requests", func(t *testing.T) {
		c := &Config{}

		for _, m := range nonGetMethods {
			req, err := http.NewRequest(m, "/", nil)
			assert.Nil(t, err)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(c.Ping())
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

		for _, m := range []string{http.MethodHead, http.MethodGet} {
			req, err := http.NewRequest(m, "/", nil)
			assert.Nil(t, err)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(c.Ping())
			h.ServeHTTP(w, req)

			assert.Equal(t, w.Code, http.StatusOK)
			assert.Empty(t, w.Body.String())
		}
	})

	t.Run("unsuccessful healthcheck", func(t *testing.T) {
		c := &Config{
			Check: func(ctx context.Context) error {
				return fmt.Errorf("error!")
			},
		}

		for _, m := range []string{http.MethodHead, http.MethodGet} {
			req, err := http.NewRequest(m, "/", nil)
			assert.Nil(t, err)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(c.Ping())
			h.ServeHTTP(w, req)

			assert.Equal(t, w.Code, http.StatusInternalServerError)
			assert.Empty(t, w.Body.String())
		}
	})
}

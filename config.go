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

// The healthchecks package gives functions and http handlers for implementing
// service healthchecks.
package healthcheck

import (
	"context"
	"net/http"
	"os"
)

// The Healthcheck type represents the function that is called when determining
// if the service is healthy or not.
type Healthcheck func(context.Context) error

type handler func(http.ResponseWriter, *http.Request)

// The Config object allows configuring the healthchecks and informational
// routes.
type Config struct {
	BuildDate string
	GitSha    string
	Version   string

	Check Healthcheck
}

// A simple wrapper around the stdlib hostname function that panics if it
// receives an error.
func hostnameOrPanic() string {
	hn, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return hn
}

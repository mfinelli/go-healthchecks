# go-healthchecks

[![Default](https://github.com/mfinelli/go-healthchecks/actions/workflows/default.yml/badge.svg)](https://github.com/mfinelli/go-healthchecks/actions/workflows/default.yml)
[![Go Reference](https://pkg.go.dev/badge/go.finelli.dev/healthchecks.svg)](https://pkg.go.dev/go.finelli.dev/healthchecks)

Healthechecks and info routes for your golang-based webservice.

## usage

```go
package main

import (
        "context"
        "log"
        "net/http"

        "github.com/jackc/pgx/v5/pgxpool"
        "go.finelli.dev/healthchecks"
)

var commit string
var date string
var version string

func main() {
        ctx := context.Background()

        // set up your database(s) here...
        db, err := pgxpool.New(ctx, "postgres://...")
        if err != nil {
                panic(err)
        }

        hc := healthchecks.Config{
                // you can omit any of these that you don't wish to expose
                Version: version,
                BuildDate: date,
                GitSha: commit,

                // You can also define a function somewhere else that returns
                // the healthchecks.Healthcheck type:
                // func MyHealthcheckSetup(...) healthchecks.Healthcheck {}
                // Check: MyHealthcheckSetup(...),
                Check: func (ctx context.Context) error {
                        if err := db.Ping(ctx); err != nil {
                                return err
                        }

                        // more healthchecks...

                        return nil
                },
        }

        http.HandleFunc("/health", hc.Health())
        http.HandleFunc("/info", hc.Info())
        http.HandleFunc("/ping", hc.Ping())

        log.Fatal(http.ListenAndServe(":8000", nil))
}
```

Populate the info variables at compile-time:

```shell
go build -o main \
    -ldflags "-X main.commit=$(git rev-parse --short HEAD) \
        -X main.date=$(date --utc --iso-8601=seconds) \
        -X main.version=1.0.0" \
    main.go
```

## license

```
Copyright 2024 Mario Finelli

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

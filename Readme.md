[![GoDoc](https://godoc.org/github.com/carlmjohnson/gatweay?status.svg)](https://godoc.org/github.com/carlmjohnson/gateway)
![](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/status-stable-green.svg)
 [![Calver v0.YY.Minor](https://img.shields.io/badge/calver-v0.YY.Minor-22bfda.svg)](https://calver.org)

Package gateway provides a drop-in replacement for net/http's `ListenAndServe` for use in AWS Lambda & API Gateway, simply swap it out for `gateway.ListenAndServe`. Extracted from [Up](https://github.com/apex/up) which provides additional middleware features and operational functionality.

This version is forked from [Apex/gateway](https://github.com/apex/gateway), which tended to merge pull requests very infrequently. Another fork, by [piotrkubisa](https://github.com/piotrkubisa/apigo), was good, but he has discontinued work.

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/carlmjohnson/gateway"
)

func main() {
	http.HandleFunc("/", hello)
	log.Fatal(gateway.ListenAndServe("n/a", nil))
}

func hello(w http.ResponseWriter, r *http.Request) {
	// example retrieving values from the api gateway proxy request context.
	requestContext, ok := gateway.RequestContext(r.Context())
	if !ok || requestContext.Authorizer["sub"] == nil {
		fmt.Fprint(w, "Hello World from Go")
		return
	}

	userID := requestContext.Authorizer["sub"].(string)
	fmt.Fprintf(w, "Hello %s from Go", userID)
}
```

# envloader

Parse environment variables into typed Go values, with required and conditional settings. Uses `flag.Value` parsers.

Status: the API is a work in progress.

## Install

Requires Go 1.19 or newer.

```sh
go get github.com/andreyvit/envloader@latest
```

## Example

```go
package main

import (
	"fmt"
	"os"

	"github.com/andreyvit/envloader"
)

func main() {
	var port int
	var vars envloader.VarSet
	vars.Var("APP_PORT", envloader.Required, envloader.IntVar(&port), "HTTP port")
	// Use vars.TryParse() to read the process environment.
	err := vars.TryParseFrom(func(name string) string { return map[string]string{"APP_PORT": "8080"}[name] })
	if err != nil {
		envloader.PrintError(err, os.Stderr)
		os.Exit(2)
	}
	fmt.Println(port) // 8080
}
```

Use `Optional` for optional values and `WhenTrue` for conditional requirements. `Parse` prints errors and exits; `TryParse` lets the caller decide. Avoid printing variable values when they contain secrets.

## License

See [LICENSE](LICENSE).

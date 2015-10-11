## Why

This project is a fork of https://github.com/caarlos0/env, which provides a
neat way of managing config from the environment, to support 12factor apps
and the like. However, on the basis that explicit is better than implicit, I
wanted to receive an error when a variable wasn't set, rather than yield the
default type. This means that if a config variable is renamed in a big
deployment the desired value isn't silently having no effect, but the code
otherwise runs.

Sometimes, a bit of config is truly optional, so this should be allowed, but again,
explicitly. Having defaults implicitly hidden in the code also breaks this philosophy,
so that upstream feature is removed.

## Example

A very basic example (check the `examples` folder):

```go
package main

import (
	"fmt"
	"os"

	"github.com/richleigh/env"
)

type config struct {
	Home         string `env:"HOME"`
	Port         int    `env:"PORT"`
	IsProduction bool   `env:"PRODUCTION,optional"`
}

func main() {
	os.Setenv("HOME", "/tmp/fakehome")
	cfg := config{}
	err := env.Parse(&cfg)
	if err == nil {
		fmt.Println(cfg)
	} else {
		fmt.Println(err)
	}

}
```

You can run it like this:

```sh
$ PORT=3000 go run examples/first.go
{/tmp/fakehome 3000 false}
```

And to see what happens when we don't set a PORT explicitly:
```sh
$ go run examples/first.go
Missing config environment variable 'PORT'
```

## Supported types and defaults

Currently we only support `string`, `bool` and `int`.

For optional fields, if the environment doesn't provide an explicit
value, the zero-value of the type will be used: empty for `string`s, `false`
for `bool`s and `0` for `int`s.


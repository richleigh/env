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

Additionally, there may be sensitive information, such as passwords or keys in
the environment. It may be possible for an attacker to gain sufficient access
to a running application that they can retrieve environment variables, so this
fork also adds the ability to erase them once we have parsed them, making them in
effect read once.

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
	Password     string `env:"PASSWORD,sensitive,optional"`
}

func main() {
	os.Setenv("HOME", "/tmp/fakehome")
	fmt.Printf("PASSWORD before: '%s'\n", os.Getenv("PASSWORD"))
	cfg := config{}
	err := env.Parse(&cfg)
	if err == nil {
		fmt.Println(cfg)
	} else {
		fmt.Println(err)
	}
	fmt.Printf("PASSWORD after: '%s'\n", os.Getenv("PASSWORD"))
}
```

You can run it like this:

```sh
$ PORT=8000 go run examples/first.go
PASSWORD before: ''
{/tmp/fakehome 8000 false }
PASSWORD after: ''
```

And to see what happens when we don't set a PORT explicitly:
```sh
$ go run examples/first.go
Missing config environment variable 'PORT'
```

And if we set PASSWORD:
```sh
$ PORT=8000 PASSWORD=secr3t go run examples/first.go
PASSWORD before: 'secr3t'
{/tmp/fakehome 8000 false secr3t}
PASSWORD after: ''
```

## Supported types and defaults

Currently we only support `string`, `bool` and `int`.

For optional fields, if the environment doesn't provide an explicit
value, the zero-value of the type will be used: empty for `string`s, `false`
for `bool`s and `0` for `int`s.


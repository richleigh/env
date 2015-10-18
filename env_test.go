package env_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/richleigh/env"
	"github.com/stretchr/testify/assert"
)

type Config struct {
	Some        string `env:"somevar"`
	Other       bool   `env:"othervar"`
	Port        int    `env:"PORT"`
	NotAnEnv    string
	DatabaseURL string `env:"DATABASE_URL,optional"`
	Password    string `env:"PASSWORD,optional,sensitive"`
}

func TestParsesEnv(t *testing.T) {
	os.Setenv("somevar", "somevalue")
	os.Setenv("othervar", "true")
	os.Setenv("PORT", "8080")
	defer os.Setenv("somevar", "")
	defer os.Setenv("othervar", "")
	defer os.Setenv("PORT", "")

	cfg := Config{}
	assert.NoError(t, env.Parse(&cfg))
	assert.Equal(t, "somevalue", cfg.Some)
	assert.Equal(t, true, cfg.Other)
	assert.Equal(t, 8080, cfg.Port)
}

func TestEmptyVars(t *testing.T) {
	cfg := Config{}
	assert.Error(t, env.Parse(&cfg))
}

func TestPassAnInvalidPtr(t *testing.T) {
	var thisShouldBreak int
	assert.Error(t, env.Parse(&thisShouldBreak))
}

func TestPassReference(t *testing.T) {
	cfg := Config{}
	assert.Error(t, env.Parse(cfg))
}

func TestInvalidBool(t *testing.T) {
	os.Setenv("somevar", "somevalue")
	os.Setenv("PORT", "8080")
	defer os.Setenv("somevar", "")
	defer os.Setenv("PORT", "")

	os.Setenv("othervar", "should-be-a-bool")
	defer os.Setenv("othervar", "")

	cfg := Config{}
	assert.Error(t, env.Parse(&cfg))
}

func TestInvalidInt(t *testing.T) {
	os.Setenv("somevar", "somevalue")
	os.Setenv("othervar", "true")
	defer os.Setenv("somevar", "")
	defer os.Setenv("othervar", "")

	os.Setenv("PORT", "should-be-an-int")
	defer os.Setenv("PORT", "")

	cfg := Config{}
	assert.Error(t, env.Parse(&cfg))
}

func TestParsesOptionalMissing(t *testing.T) {
	os.Setenv("somevar", "somevalue")
	os.Setenv("othervar", "true")
	os.Setenv("PORT", "8080")
	defer os.Setenv("somevar", "")
	defer os.Setenv("othervar", "")
	defer os.Setenv("PORT", "")

	cfg := Config{}
	assert.NoError(t, env.Parse(&cfg))
	assert.Equal(t, "", cfg.DatabaseURL)
}

func TestParsesOptionalPresent(t *testing.T) {
	os.Setenv("somevar", "somevalue")
	os.Setenv("othervar", "true")
	os.Setenv("PORT", "8080")
	defer os.Setenv("somevar", "")
	defer os.Setenv("othervar", "")
	defer os.Setenv("PORT", "")

	cfg := Config{}
	db := "postgres://localhost:5432/db"
	os.Setenv("DATABASE_URL", db)
	assert.NoError(t, env.Parse(&cfg))
	assert.Equal(t, db, cfg.DatabaseURL)
}

func TestParsesSensitive(t *testing.T) {
	os.Setenv("somevar", "somevalue")
	os.Setenv("othervar", "true")
	os.Setenv("PORT", "8080")
	defer os.Setenv("somevar", "")
	defer os.Setenv("othervar", "")
	defer os.Setenv("PORT", "")

	cfg := Config{}
	pw := "MrGoodbytes"
	os.Setenv("PASSWORD", pw)
	assert.NoError(t, env.Parse(&cfg))
	assert.Equal(t, pw, cfg.Password)
	assert.Equal(t, "", os.Getenv("PASSWORD"))
}

func TestClearSensitiveOnError(t *testing.T) {
	os.Setenv("somevar", "somevalue")
	os.Setenv("othervar", "true")
	defer os.Setenv("somevar", "")
	defer os.Setenv("othervar", "")

	cfg := Config{}
	pw := "MrGoodbytes"
	os.Setenv("PASSWORD", pw)
	// Should get an error, since we didn't pass in PORT
	assert.Error(t, env.Parse(&cfg))
	assert.Equal(t, "", os.Getenv("PASSWORD"))
}

func TestParseStructWithoutEnvTag(t *testing.T) {
	os.Setenv("somevar", "somevalue")
	os.Setenv("othervar", "true")
	os.Setenv("PORT", "8080")
	defer os.Setenv("somevar", "")
	defer os.Setenv("othervar", "")
	defer os.Setenv("PORT", "")

	cfg := Config{}
	assert.NoError(t, env.Parse(&cfg))
	assert.Empty(t, cfg.NotAnEnv)
}

func TestParseStructWithInvalidFieldKind(t *testing.T) {
	type config struct {
		WontWork int64 `env:"BLAH"`
	}
	os.Setenv("BLAH", "10")
	cfg := config{}
	assert.Error(t, env.Parse(&cfg))
}

func ExampleParse() {
	type config struct {
		Home         string `env:"HOME"`
		Port         int    `env:"PORT"`
		IsProduction bool   `env:"PRODUCTION,optional"`
	}
	os.Setenv("HOME", "/tmp/fakehome")
	os.Setenv("PORT", "3000")
	defer os.Setenv("HOME", "")
	defer os.Setenv("PORT", "")

	cfg := config{}
	err := env.Parse(&cfg)
	if err == nil {
		fmt.Println(cfg)
	} else {
		fmt.Println(err)
	}
	// Output: {/tmp/fakehome 3000 false}
}

func ExampleFail() {
	type config struct {
		Home         string `env:"HOME"`
		Port         int    `env:"PORT"`
		IsProduction bool   `env:"PRODUCTION,optional"`
	}
	os.Setenv("HOME", "/tmp/fakehome")
	defer os.Setenv("HOME", "")

	cfg := config{}
	err := env.Parse(&cfg)
	if err == nil {
		fmt.Println(cfg)
	} else {
		fmt.Println(err)
	}
	// Output: Missing config environment variable 'PORT'
}

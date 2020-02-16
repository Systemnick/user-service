package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type Config struct {
	DatabaseUrl string `env:"DATABASE_URL"`
	LogLevel    string `env:"LOG_LEVEL"`
}

func GetConfigFromEnv() *Config {
	c := &Config{}

	v := reflect.ValueOf(c)

	// Dereference if needed.
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		if v.IsNil() {
			v.Set(reflect.New(t))
		}
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		// Don't even bother for unexported fields.
		if !v.Field(i).CanSet() {
			continue
		}

		field := v.Field(i)
		envVar := v.Type().Field(i).Tag.Get("env")

		switch field.Type().String() {
		case "string":
			value := readEnvStr(envVar, "")
			field.SetString(value)
		case "int":
			value := readEnvInt64(envVar, 0)
			field.SetInt(value)
		}
	}

	return c
}

func readEnvInt64(key string, defaultValue int64) int64 {
	s := os.Getenv(key)
	if s == "" {
		return defaultValue
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		fmt.Printf("Environment variable %s: bad int64 value: %s\n", key, err.Error())
		i = defaultValue
	}

	return i
}

func readEnvStr(key string, defaultValue string) string {
	s := os.Getenv(key)
	if s == "" {
		s = defaultValue
	}

	return s
}

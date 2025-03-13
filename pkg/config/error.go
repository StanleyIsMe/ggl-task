package config

import "fmt"

type MissingEnvConfigError struct {
	Env string
	Err error
}

func (e *MissingEnvConfigError) Error() string {
	return fmt.Sprintf("missing config %s: %v", e.Env, e.Err)
}

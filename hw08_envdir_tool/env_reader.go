package main

import (
	"errors"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var errUnsupported = errors.New("для описания переменных окружения поддерживаются только обычные файлы")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirEntryes, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	env := make(Environment)
	for _, entry := range dirEntryes {
		if !entry.Type().IsRegular() {
			return nil, errUnsupported
		}
		content, err := os.ReadFile(dir + string(os.PathSeparator) + entry.Name())
		if err != nil {
			return nil, err
		}
		if len(content) == 0 {
			env[entry.Name()] = EnvValue{NeedRemove: true}
			continue
		}
		lines := strings.Split(string(content), "\n")
		if len(lines) > 0 {
			env[entry.Name()] = EnvValue{
				Value: strings.ReplaceAll(
					strings.TrimRight(lines[0], " \t"),
					"\x00", "\n"),
			}
		}
	}
	return env, nil
}

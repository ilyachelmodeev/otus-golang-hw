package main

import (
	"fmt"
	"os"
	"regexp"
)

func main() {
	if len(os.Args) < 3 {
		os.Exit(125)
	}
	dirEnv, err := ReadDir(os.Args[1])
	if err != nil {
		_, err := fmt.Fprint(os.Stderr, err.Error())
		if err != nil {
			return
		}
		os.Exit(1)
	}
	nameRegex := regexp.MustCompile("^[^=]+")
	valueRegexp := regexp.MustCompile("^[^=]+=(.*)")
	osEnv := make(Environment)
	for _, s := range os.Environ() {
		n := nameRegex.FindString(s)
		v := valueRegexp.FindStringSubmatch(s)
		osEnv[n] = EnvValue{Value: v[1]}
	}
	for n, v := range dirEnv {
		if v.NeedRemove {
			delete(osEnv, n)
			continue
		}
		osEnv[n] = v
	}
	os.Exit(RunCmd(os.Args[2:], osEnv))
}

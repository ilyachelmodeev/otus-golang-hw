package hw09structvalidator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var vFuncRegexp = regexp.MustCompile(`(^[^:]+):(.+)`)

func (v ValidationErrors) Error() string {
	b := strings.Builder{}
	b.WriteRune('[')
	for i, ve := range v {
		b.WriteString(ve.Field)
		b.WriteRune(':')
		b.WriteString(ve.Err.Error())
		if i < len(v)-1 {
			b.WriteRune('\n')
		}
	}
	b.WriteRune(']')
	return b.String()
}

func Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return errors.New("ожидается структура")
	}
	var vErrors ValidationErrors
	for i := 0; i < rv.NumField(); i++ {
		tags := reflect.TypeOf(v).Field(i).Tag
		vtag := tags.Get("validate")
		if vtag == "" {
			continue
		}

		name := reflect.TypeOf(v).Field(i).Name
		var err error
		f := rv.Field(i)
		switch f.Kind() { //nolint:exhaustive
		case reflect.String:
			err = validateString(f, vtag)
		case reflect.Slice:
			for j := 0; j < f.Len(); j++ {
				item := f.Index(j)
				switch item.Kind() { //nolint:exhaustive
				case reflect.String:
					err = validateString(item, vtag)
				case reflect.Int, reflect.Int64, reflect.Int32, reflect.Uint, reflect.Uint64,
					reflect.Uint32, reflect.Int8, reflect.Uint8:
					err = validateInt(f, vtag)
				default:
				}
			}
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Uint, reflect.Uint64,
			reflect.Uint32, reflect.Int8, reflect.Uint8:
			err = validateInt(f, vtag)
		default:
		}
		if vErr, ok := err.(ValidationErrors); ok { //nolint:errorlint
			for j := range vErr {
				vErr[j].Field = name
			}
			vErrors = append(vErrors, vErr...)
			continue
		}
		if err != nil {
			return err
		}
	}
	if len(vErrors) > 0 {
		return vErrors
	}
	return nil
}

func validateString(rS reflect.Value, tag string) error {
	vErrs := ValidationErrors{}
	validationFunctions := strings.Split(tag, "|")
	for _, f := range validationFunctions {
		fields := vFuncRegexp.FindStringSubmatch(f)
		if len(fields) != 3 {
			return errors.New("wrong validate syntax: " + f)
		}
		switch fields[1] {
		case "len":
			arg, ok := strconv.Atoi(fields[2])
			if ok != nil {
				return errors.New("wrong len validate syntax" + f)
			}
			if len(rS.String()) != arg {
				vErrs = append(vErrs, ValidationError{Err: errors.New("wrong length")})
			}
		case "regexp":
			regexp, err := regexp.Compile(fields[2])
			if err != nil {
				return errors.New("wrong regexp validate syntax " + f + ": " + err.Error())
			}
			regexpResult := regexp.FindString(rS.String())
			if regexpResult != rS.String() {
				vErrs = append(vErrs, ValidationError{Err: errors.New("regexp mismatch")})
			}
		case "in":
			args := strings.Split(fields[2], ",")
			if !sContains(rS.String(), args) {
				vErrs = append(vErrs, ValidationError{Err: errors.New("value not found in list")})
			}
		default:
			return errors.New("wrong validate tag:" + f)
		}
	}
	if len(vErrs) > 0 {
		return vErrs
	}
	return nil
}

func sContains(s string, list []string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}
	return false
}

func validateInt(rS reflect.Value, tag string) error {
	vErrs := ValidationErrors{}
	validationFunctions := strings.Split(tag, "|")
	for _, f := range validationFunctions {
		fields := vFuncRegexp.FindStringSubmatch(f)
		if len(fields) != 3 {
			return errors.New("wrong validate syntax: " + f)
		}
		switch fields[1] {
		case "min":
			arg, ok := strconv.Atoi(fields[2])
			if ok != nil {
				return errors.New("wrong min validate syntax" + f)
			}
			if int(rS.Int()) < arg {
				vErrs = append(vErrs, ValidationError{Field: "", Err: errors.New("value is less than the lower limit")})
			}
		case "max":
			arg, ok := strconv.Atoi(fields[2])
			if ok != nil {
				return errors.New("wrong max validate syntax" + f)
			}
			if int(rS.Int()) > arg {
				vErrs = append(vErrs, ValidationError{Field: "", Err: errors.New("value exceeds upper limit")})
			}
		case "in":
			args := strings.Split(fields[2], ",")
			found, err := intContains(int(rS.Int()), args)
			if err != nil {
				return err
			}
			if !found {
				vErrs = append(vErrs, ValidationError{Err: errors.New("value not found in list")})
			}
		default:
			return errors.New("wrong validate tag:" + f)
		}
	}
	if len(vErrs) > 0 {
		return vErrs
	}
	return nil
}

func intContains(value int, sList []string) (bool, error) {
	for _, s := range sList {
		listValue, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return false, errors.New("wrong int list:" + s)
		}
		if value == int(listValue) {
			return true, nil
		}
	}
	return false, nil
}

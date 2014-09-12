package cfg

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"text/template"

	"launchpad.net/goyaml"
)

func LoadConfig(configPath string, configStruct interface{}) error {
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	bytes, err = substitute(bytes)
	if err != nil {
		return err
	}

	if err = goyaml.Unmarshal(bytes, configStruct); err != nil {
		return err
	}

	if err = validate(configStruct); err != nil {
		return err
	}

	return nil
}

type templateData struct {
	Env map[string]string
}

func substitute(in []byte) ([]byte, error) {
	t, err := template.New("config").Parse(string(in))
	if err != nil {
		return nil, err
	}

	data := &templateData{
		Env: make(map[string]string),
	}

	values := os.Environ()
	for _, val := range values {
		keyval := strings.SplitN(val, "=", 2)
		if len(keyval) != 2 {
			continue
		}
		data.Env[keyval[0]] = keyval[1]
	}

	buffer := &bytes.Buffer{}
	if err = t.Execute(buffer, data); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func validate(configStruct interface{}) error {
	return validateStruct(
		reflect.TypeOf(configStruct).Elem(),
		reflect.ValueOf(configStruct).Elem())
}

func validateStruct(typ reflect.Type, val reflect.Value) error {
	for idx := 0; idx < val.NumField(); idx++ {
		field := typ.Field(idx)
		if field.Type.Kind() == reflect.Struct {
			if err := validateStruct(val.Field(idx).Type(), val.Field(idx)); err != nil {
				return err
			}
		} else if field.Type.Kind() == reflect.Bool || field.Type.Kind() == reflect.Int { // no way to tell if boolean field was provided or not
			continue
		} else {
			if field.Tag.Get("config") != "optional" {
				if val.Field(idx).Len() == 0 {
					return errors.New(
						fmt.Sprintf("Missing required config field: %v", field.Name))
				}
			}
		}
	}
	return nil
}

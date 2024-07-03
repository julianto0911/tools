package tools

import (
	"reflect"
	"strings"
)

func Validate(item interface{}) error {
	val := reflect.ValueOf(item)
	for i := 0; i < val.NumField(); i++ {
		var err error
		fieldType := val.Type().Field(i)
		rule := fieldType.Tag.Get("validate")
		rules := strings.Split(rule, ";")

		field := val.Field(i)
		value := field.Interface()

		//get name of variable
		name := fieldType.Tag.Get("json")
		//get value of variable
		//check type of validation
		switch field.Type().Kind() {
		case reflect.String:
			err = validateString(rules, name, value)
		case reflect.Int:
			err = validateInt(rules, name, value)
		case reflect.Int64:
			err = validateInt64(rules, name, value)
		case reflect.Float64:
			err = validateFloat64(rules, name, value)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

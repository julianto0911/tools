package lib_validator

import (
	"fmt"
	"strconv"
	"strings"
)

func validateInt64(rules []string, name string, v interface{}) error {
	value, _ := v.(int64)
	for _, rule := range rules {
		if err := int64MinValue(rule, name, value); err != nil {
			return err
		}
		if err := int64MaxValue(rule, name, value); err != nil {
			return err
		}
		if err := int64Range(rule, name, value); err != nil {
			return err
		}
		if err := int64Required(rule, name, value); err != nil {
			return err
		}
	}

	return nil
}

func int64Required(rule, name string, value int64) error {
	if !strings.Contains(rule, "required") {
		return nil
	}
	if value == 0 {
		return fmt.Errorf("field %v must not zero", name)
	}

	return nil
}
func int64Range(rule, name string, value int64) error {
	if !strings.Contains(rule, "range") {
		return nil
	}

	temp := strings.ReplaceAll(rule, "range", "")
	temp = strings.ReplaceAll(temp, "[", "")
	temp = strings.ReplaceAll(temp, "]", "")

	array := strings.Split(temp, ",")
	found := false
	for _, val := range array {
		t, _ := strconv.ParseInt(val, 10, 64)
		if t == value {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("field %v value must in [%v]", name, temp)
	}

	return nil
}

func int64MinValue(rule, name string, value int64) error {
	if !strings.Contains(rule, "min") {
		return nil
	}
	r := strings.Split(rule, "=")
	if len(r) < 2 {
		return fmt.Errorf("min-value invalid rule definition, name : " + name)
	}

	limit, err := strconv.ParseInt(strings.TrimSpace(r[1]), 10, 64)
	if err != nil {
		return fmt.Errorf("min-value invalid rule:(%v) %w", name, err)
	}

	if value < limit {
		return fmt.Errorf("field %v must not less than %v", name, limit)
	}

	return nil
}

func int64MaxValue(rule, name string, value int64) error {
	if !strings.Contains(rule, "max") {
		return nil
	}
	r := strings.Split(rule, "=")
	if len(r) < 2 {
		return fmt.Errorf("max-value invalid rule definition, name : " + name)
	}

	limit, err := strconv.ParseInt(strings.TrimSpace(r[1]), 10, 64)
	if err != nil {
		return fmt.Errorf("max-value invalid rule:(%v) %w", name, err)
	}

	if value > limit {
		return fmt.Errorf("field %v must not greater than %v", name, limit)
	}

	return nil
}

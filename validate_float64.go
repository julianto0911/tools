package tools

import (
	"fmt"
	"strconv"
	"strings"
)

func validateFloat64(rules []string, name string, v interface{}) error {
	value, _ := v.(float64)
	for _, rule := range rules {
		if err := float64MinValue(rule, name, value); err != nil {
			return err
		}
		if err := float64MaxValue(rule, name, value); err != nil {
			return err
		}
		if err := float64Range(rule, name, value); err != nil {
			return err
		}
		if err := float64Required(rule, name, value); err != nil {
			return err
		}
	}

	return nil
}

func float64Required(rule, name string, value float64) error {
	if !strings.Contains(rule, "required") {
		return nil
	}
	if value == 0 {
		return fmt.Errorf("field %v must not zero", name)
	}

	return nil
}

func float64Range(rule, name string, value float64) error {
	if !strings.Contains(rule, "range") {
		return nil
	}

	temp := strings.ReplaceAll(rule, "range", "")
	temp = strings.ReplaceAll(temp, "[", "")
	temp = strings.ReplaceAll(temp, "]", "")

	array := strings.Split(temp, ",")
	found := false
	for _, val := range array {
		t, _ := strconv.ParseFloat(val, 64)
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

func float64MinValue(rule, name string, value float64) error {
	if !strings.Contains(rule, "min") {
		return nil
	}
	r := strings.Split(rule, "=")
	if len(r) < 2 {
		return fmt.Errorf("min-value invalid rule definition, name : " + name)
	}

	limit, err := strconv.ParseFloat(r[1], 64)
	if err != nil {
		return fmt.Errorf("min-value invalid rule:(%v) %w", name, err)
	}

	if value < limit {
		return fmt.Errorf("field %v must not less than %v", name, limit)
	}

	return nil
}

func float64MaxValue(rule, name string, value float64) error {
	if !strings.Contains(rule, "max") {
		return nil
	}
	r := strings.Split(rule, "=")
	if len(r) < 2 {
		return fmt.Errorf("max-value invalid rule definition, name : " + name)
	}

	limit, err := strconv.ParseFloat(r[1], 64)
	if err != nil {
		return fmt.Errorf("max-value invalid rule:(%v) %w", name, err)
	}

	if value > limit {
		return fmt.Errorf("field %v must not greater than %v", name, limit)
	}

	return nil
}

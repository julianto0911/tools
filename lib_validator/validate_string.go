package lib_validator

import (
	"fmt"
	"strconv"
	"strings"
)

func validateString(rules []string, name string, v interface{}) error {
	value := v.(string)
	for _, rule := range rules {
		if err := strRequired(rule, name, value); err != nil {
			return err
		}
		if err := strMinLength(rule, name, value); err != nil {
			return err
		}
		if err := strMaxLength(rule, name, value); err != nil {
			return err
		}
		if err := strLength(rule, name, value); err != nil {
			return err
		}
		if err := strRange(rule, name, value); err != nil {
			return err
		}
	}

	return nil
}

func strRange(rule, name string, value string) error {
	if !strings.Contains(rule, "range") {
		return nil
	}

	temp := strings.ReplaceAll(rule, "range", "")
	temp = strings.ReplaceAll(temp, "[", "")
	temp = strings.ReplaceAll(temp, "]", "")

	array := strings.Split(temp, ",")
	found := false
	for _, val := range array {
		if val == value {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("field %v value must in [%v]", name, temp)
	}

	return nil
}

func strLength(rule, name, value string) error {
	if !strings.Contains(rule, "length") {
		return nil
	}

	r := strings.Split(rule, "=")
	if len(r) < 2 {
		return fmt.Errorf("length invalid rule definition, name : " + name)
	}

	limit, err := strconv.Atoi(strings.TrimSpace(r[1]))
	if err != nil {
		return fmt.Errorf("invalid rule:(%v) %w", name, err)
	}

	if !(len(value) == limit) {
		return fmt.Errorf("field %v must have %v character(s)", name, limit)
	}

	return nil
}

func strRequired(rule, name, value string) error {
	if !strings.Contains(rule, "required") {
		return nil
	}
	if value == "" {
		return fmt.Errorf("field %v must be filled", name)
	}

	return nil
}

func strMinLength(rule, name, value string) error {
	if !strings.Contains(rule, "min") {
		return nil
	}

	r := strings.Split(rule, "=")
	if len(r) < 2 {
		return fmt.Errorf("min-length invalid rule definition, name : " + name)
	}

	limit, err := strconv.Atoi(strings.TrimSpace(r[1]))
	if err != nil {
		return fmt.Errorf("invalid rule:(%v) %w", name, err)
	}

	if len(value) < limit {
		return fmt.Errorf("field %v must have at least %v character(s)", name, limit)
	}

	return nil
}

func strMaxLength(rule, name, value string) error {
	if !strings.Contains(rule, "max") {
		return nil
	}

	r := strings.Split(rule, "=")
	if len(r) < 2 {
		return fmt.Errorf("max-length invalid rule definition, name : " + name)
	}

	limit, err := strconv.Atoi(strings.TrimSpace(r[1]))
	if err != nil {
		return fmt.Errorf("invalid rule:(%v) %w", name, err)
	}

	if len(value) > limit {
		return fmt.Errorf("total characters for field %v must be less or same than %v character(s)", name, limit)
	}

	return nil
}

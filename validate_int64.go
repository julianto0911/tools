package tools

import (
	"errors"
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
		if err := int64OptX(rule, name, value); err != nil {
			return err
		}
		if err := int64OptY(rule, name, value); err != nil {
			return err
		}
	}

	return nil
}

func int64OptX(rule, name string, value int64) error {
	if value == 0 {
		return nil
	}

	if !strings.Contains(rule, "optx") {
		return nil
	}

	r := strings.Split(rule, "=")
	if len(r) < 2 {
		return fmt.Errorf("length invalid rule definition, name : " + name)
	}

	limit := StringToInt64(strings.TrimSpace(r[1]))
	if limit == 0 {
		return fmt.Errorf("invalid rule:(%v) %w", name, errors.New("cannot define zero in rule"))
	}

	if value < limit {
		return fmt.Errorf("when have value, field %v must have at least %v character(s)", name, limit)
	}

	return nil
}

func int64OptY(rule, name string, value int64) error {
	if value == 0 {
		return nil
	}

	if !strings.Contains(rule, "opty") {
		return nil
	}

	r := strings.Split(rule, "=")
	if len(r) < 2 {
		return fmt.Errorf("length invalid rule definition, name : " + name)
	}

	limit := StringToInt64(strings.TrimSpace(r[1]))
	if limit == 0 {
		return fmt.Errorf("invalid rule:(%v) %w", name, errors.New("cannot define zero in rule"))
	}

	if value > limit {
		return fmt.Errorf("when have value, total characters for field %v must be less or same than %v character(s)", name, limit)
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

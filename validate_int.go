package tools

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func validateInt(rules []string, name string, v interface{}) error {
	value, _ := v.(int)
	for _, rule := range rules {
		if err := intMinValue(rule, name, value); err != nil {
			return err
		}
		if err := intMaxValue(rule, name, value); err != nil {
			return err
		}
		if err := intRange(rule, name, value); err != nil {
			return err
		}
		if err := intRequired(rule, name, value); err != nil {
			return err
		}
		if err := intOptX(rule, name, value); err != nil {
			return err
		}
		if err := intOptY(rule, name, value); err != nil {
			return err
		}
	}

	return nil
}

func intOptX(rule, name string, value int) error {
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

	limit := StringToInt(strings.TrimSpace(r[1]))
	if limit == 0 {
		return fmt.Errorf("invalid rule:(%v) %w", name, errors.New("cannot define zero in rule"))
	}

	if value < limit {
		return fmt.Errorf("when have value, field %v must have at least %v character(s)", name, limit)
	}

	return nil
}

func intOptY(rule, name string, value int) error {
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

	limit := StringToInt(strings.TrimSpace(r[1]))
	if limit == 0 {
		return fmt.Errorf("invalid rule:(%v) %w", name, errors.New("cannot define zero in rule"))
	}

	if value > limit {
		return fmt.Errorf("when have value, total characters for field %v must be less or same than %v character(s)", name, limit)
	}

	return nil
}

func intRequired(rule, name string, value int) error {
	if !strings.Contains(rule, "required") {
		return nil
	}
	if value == 0 {
		return fmt.Errorf("field %v must not zero", name)
	}

	return nil
}
func intRange(rule, name string, value int) error {
	if !strings.Contains(rule, "range") {
		return nil
	}

	temp := strings.ReplaceAll(rule, "range", "")
	temp = strings.ReplaceAll(temp, "[", "")
	temp = strings.ReplaceAll(temp, "]", "")

	array := strings.Split(temp, ",")
	found := false
	for _, val := range array {
		t, _ := strconv.Atoi(val)
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

func intMinValue(rule, name string, value int) error {
	if !strings.Contains(rule, "min") {
		return nil
	}
	r := strings.Split(rule, "=")
	if len(r) < 2 {
		return fmt.Errorf("min-value invalid rule definition, name : " + name)
	}

	limit, err := strconv.Atoi(strings.TrimSpace(r[1]))
	if err != nil {
		return fmt.Errorf("min-value invalid rule:(%v) %w", name, err)
	}

	if value < limit {
		return fmt.Errorf("field %v must not less than %v", name, limit)
	}

	return nil
}

func intMaxValue(rule, name string, value int) error {
	if !strings.Contains(rule, "max") {
		return nil
	}
	r := strings.Split(rule, "=")
	if len(r) < 2 {
		return fmt.Errorf("max-value invalid rule definition, name : " + name)
	}

	limit, err := strconv.Atoi(strings.TrimSpace(r[1]))
	if err != nil {
		return fmt.Errorf("max-value invalid rule:(%v) %w", name, err)
	}

	if value > limit {
		return fmt.Errorf("field %v must not greater than %v", name, limit)
	}

	return nil
}

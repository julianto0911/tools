package tools

import (
	"fmt"
	"strconv"
	"strings"
)

/*
validate string rules
required ("validate:required") , must have value
min ("validate:min=[length]") , value typed length must be at least of value
max ("validate:max=[length]") , value typed length must same or less length than value
length ("validate:length=[length]") , value typed length must be in exact length value
range ("validate:range[val1,val2]"), value typed must be in range of value declared
optx ("validate:optx=[length]") , if have value, length must be at least of value
opty ("validate:opty=[length]") , if have value, length must be same or less than value
*/
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
		if err := strOptX(rule, name, value); err != nil {
			return err
		}
		if err := strOptY(rule, name, value); err != nil {
			return err
		}
	}

	return nil
}

func strOptX(rule, name, value string) error {
	if value == "" {
		return nil
	}

	if !strings.Contains(rule, "optx") {
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

	if len(value) < limit {
		return fmt.Errorf("when have value, field %v must have at least %v character(s)", name, limit)
	}

	return nil
}

func strOptY(rule, name, value string) error {
	if value == "" {
		return nil
	}

	if !strings.Contains(rule, "opty") {
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

	if len(value) > limit {
		return fmt.Errorf("when have value, total characters for field %v must be less or same than %v character(s)", name, limit)
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

package lib_validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type strMockValidator struct {
	Name  string `json:"name" validate:"required;min=2;max=15"`
	Value string `json:"value" validate:"length=5;range[abcde,bcdef,cdefg]"`
}

func TestValidateString(t *testing.T) {
	data := strMockValidator{}
	//test fail for required
	data.Name = ""
	err := Validate(data)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "field name must be filled")
	}

	//test fail for min length
	data.Name = "1"
	err = Validate(data)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "field name must have at least 2 character(s)")
	}

	//test fail max length
	data.Name = "1234567891234567"
	err = Validate(data)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "total characters for field name must be less or same than 15 character(s)")
	}

	//test fail for exact length
	data.Name = "1234567890"
	data.Value = "123"
	err = Validate(data)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "field value must have 5 character(s)")
	}

	//test fail for range
	data.Value = "12345"
	err = Validate(data)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "field value value must in [abcde,bcdef,cdefg]")
	}

	//test success validate
	data.Value = "abcde"
	err = Validate(data)
	assert.Nil(t, err)
}

type intMockValidator struct {
	Value1 int `json:"value1" validate:"required;min=2;max=5"`
	Value2 int `json:"value2" validate:"range[1,2,3]"`
}

func TestValidateInt(t *testing.T) {
	data := intMockValidator{}

	//test fail required
	err := Validate(data)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "field value1 must not zero")
	}

	//test fail min value
	data.Value1 = 1
	err = Validate(data)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "field value1 must not less than 2")
	}

	//test fail max value
	data.Value1 = 6
	err = Validate(data)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "field value1 must not greater than 5")
	}

	//test fail range
	data.Value1 = 5
	data.Value2 = 5
	err = Validate(data)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "field value2 value must in [1,2,3]")
	}

	//test success validate
	data.Value2 = 2
	err = Validate(data)
	assert.Nil(t, err)

}

type int64MockValidator struct {
	Value1 int64 `json:"value1" validate:"required;min=2;max=5"`
	Value2 int64 `json:"value2" validate:"range[1,2,3]"`
}

func TestValidateInt64(t *testing.T) {
	data := int64MockValidator{}

	//test fail required
	err := Validate(data)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "field value1 must not zero")
	}

	//test fail min value
	data.Value1 = 1
	err = Validate(data)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "field value1 must not less than 2")
	}

	//test fail max value
	data.Value1 = 6
	err = Validate(data)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "field value1 must not greater than 5")
	}

	//test fail range
	data.Value1 = 5
	data.Value2 = 5
	err = Validate(data)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "field value2 value must in [1,2,3]")
	}

	//test success validate
	data.Value2 = 2
	err = Validate(data)
	assert.Nil(t, err)

}

type float64MockValidator struct {
	Value1 float64 `json:"value1" validate:"required;min=2;max=5"`
	Value2 float64 `json:"value2" validate:"range[1,2,3]"`
}

func TestValidateFloat64(t *testing.T) {
	data := float64MockValidator{}

	//test fail required
	err := Validate(data)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "field value1 must not zero")
	}

	//test fail min value
	data.Value1 = 1.21
	err = Validate(data)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "field value1 must not less than 2")
	}

	//test fail max value
	data.Value1 = 6.21
	err = Validate(data)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "field value1 must not greater than 5")
	}

	//test fail range
	data.Value1 = 5
	data.Value2 = 5
	err = Validate(data)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "field value2 value must in [1,2,3]")
	}

	//test success validate
	data.Value2 = 2
	err = Validate(data)
	assert.Nil(t, err)

}

func TestValidatorForStrLengthTag(t *testing.T) {
	//test for no validation
	novalidation := struct {
		Value string `json:"value"`
	}{Value: "1"}

	err := Validate(novalidation)
	assert.Nil(t, err)

	//test for invalid tag, definition should use "=", not ":"
	invalidTag := struct {
		Value string `json:"value" validate:"length:1"`
	}{Value: "1"}
	err = Validate(invalidTag)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "length invalid rule definition, name : value")
	}

	//test for invalid tag, definition should use "=", not ":" , use "-"
	invalidRule := struct {
		Value string `json:"value" validate:"length=-"`
	}{Value: "1"}
	err = Validate(invalidRule)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "invalid rule:(value)")
	}
}

func TestValidatorForStrMinTag(t *testing.T) {
	//test for invalid tag, definition should use "=", not ":"
	invalidTag := struct {
		Value string `json:"value" validate:"min:1"`
	}{Value: "1"}

	err := Validate(invalidTag)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "length invalid rule definition, name : value")
	}

	//test for invalid tag, definition should use "=", not ":" , use "-"
	invalidRule := struct {
		Value string `json:"value" validate:"min=-"`
	}{Value: "1"}
	err = Validate(invalidRule)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "invalid rule:(value)")
	}
}

func TestValidatorForStrMaxTag(t *testing.T) {
	//test for invalid tag, definition should use "=", not ":"
	invalidTag := struct {
		Value string `json:"value" validate:"max:1"`
	}{Value: "1"}

	err := Validate(invalidTag)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "length invalid rule definition, name : value")
	}

	//test for invalid tag, definition should use "=", not ":" , use "-"
	invalidRule := struct {
		Value string `json:"value" validate:"max=-"`
	}{Value: "1"}
	err = Validate(invalidRule)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "invalid rule:(value)")
	}
}

func TestValidatorForFloat64MaxTag(t *testing.T) {
	//test for invalid tag, definition should use "=", not ":"
	invalidTag := struct {
		Value float64 `json:"value" validate:"max:1"`
	}{Value: 1}

	err := Validate(invalidTag)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "max-value invalid rule definition, name : value")
	}

	//test for invalid tag, definition should use "=", not ":" , use "-"
	invalidRule := struct {
		Value float64 `json:"value" validate:"max=-"`
	}{Value: 1}
	err = Validate(invalidRule)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "max-value invalid rule:(value)")
	}
}

func TestValidatorForFloat64MinTag(t *testing.T) {
	//test for invalid tag, definition should use "=", not ":"
	invalidTag := struct {
		Value float64 `json:"value" validate:"min:1"`
	}{Value: 1}

	err := Validate(invalidTag)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "min-value invalid rule definition, name : value")
	}

	//test for invalid tag, definition should use "=", not ":" , use "-"
	invalidRule := struct {
		Value float64 `json:"value" validate:"min=-"`
	}{Value: 1}
	err = Validate(invalidRule)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "min-value invalid rule:(value)")
	}
}

func TestValidatorForIntMinTag(t *testing.T) {
	//test for invalid tag, definition should use "=", not ":"
	invalidTag := struct {
		Value int `json:"value" validate:"min:1"`
	}{Value: 1}

	err := Validate(invalidTag)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "min-value invalid rule definition, name : value")
	}

	//test for invalid tag, definition should use "=", not ":" , use "-"
	invalidRule := struct {
		Value int `json:"value" validate:"min=-"`
	}{Value: 1}
	err = Validate(invalidRule)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "min-value invalid rule:(value)")
	}
}

func TestValidatorForIntMaxTag(t *testing.T) {
	//test for invalid tag, definition should use "=", not ":"
	invalidTag := struct {
		Value int `json:"value" validate:"max:1"`
	}{Value: 1}

	err := Validate(invalidTag)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "max-value invalid rule definition, name : value")
	}

	//test for invalid tag, definition should use "=", not ":" , use "-"
	invalidRule := struct {
		Value int `json:"value" validate:"max=-"`
	}{Value: 1}
	err = Validate(invalidRule)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "max-value invalid rule:(value)")
	}
}

func TestValidatorForInt64MinTag(t *testing.T) {
	//test for invalid tag, definition should use "=", not ":"
	invalidTag := struct {
		Value int64 `json:"value" validate:"min:1"`
	}{Value: 1}

	err := Validate(invalidTag)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "min-value invalid rule definition, name : value")
	}

	//test for invalid tag, definition should use "=", not ":" , use "-"
	invalidRule := struct {
		Value int64 `json:"value" validate:"min=-"`
	}{Value: 1}
	err = Validate(invalidRule)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "min-value invalid rule:(value)")
	}
}

func TestValidatorForInt64MaxTag(t *testing.T) {
	//test for invalid tag, definition should use "=", not ":"
	invalidTag := struct {
		Value int64 `json:"value" validate:"max:1"`
	}{Value: 1}

	err := Validate(invalidTag)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "max-value invalid rule definition, name : value")
	}

	//test for invalid tag, definition should use "=", not ":" , use "-"
	invalidRule := struct {
		Value int64 `json:"value" validate:"max=-"`
	}{Value: 1}
	err = Validate(invalidRule)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "max-value invalid rule:(value)")
	}
}

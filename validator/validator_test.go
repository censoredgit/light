package validator

import (
	"github.com/censoredgit/light/validator/input"
	"github.com/censoredgit/light/validator/rules"
	"mime/multipart"
	"net/url"
	"testing"
)

func TestEmptyValidator(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})

	validator := New()
	validator.Validate(inputData)
	if validator.IsFailed() {
		t.Error("The empty validator failed")
	}
}

func TestRequiredRuleFail(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})

	validator := New()
	validator.AddRule("req_field", rules.Required())
	validator.Validate(inputData)
	if !validator.IsFailed() {
		t.Error("The fail required rule has pass")
	}
}

func TestRequiredRuleSuccessWithEmptyInputData(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})
	inputData.SetValue("req_field", "")

	validator := New()
	validator.AddRule("req_field", rules.Required().AllowEmpty())
	validator.Validate(inputData)
	if validator.IsFailed() {
		t.Error("The success required rule has pass")
	}
}

func TestRequiredRuleFailEmptyInputData(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})
	inputData.SetValue("req_field", "")

	validator := New()
	validator.AddRule("req_field", rules.Required())
	validator.Validate(inputData)
	if !validator.IsFailed() {
		t.Error("The fail required rule has pass")
	}
}

func TestMinRuleFail(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})
	inputData.SetValue("min_field_int", "2")
	inputData.SetValue("min_field_int_negative", "-2")
	inputData.SetValue("min_field_float", "2.0")

	validator := New()
	validator.AddRule("min_field_int", rules.Float(), rules.Min(5))
	validator.AddRule("min_field_int_negative", rules.Float(), rules.Min(5))
	validator.AddRule("min_field_float", rules.Float(), rules.Min(5))
	validator.Validate(inputData)
	if !validator.IsFailed() {
		t.Error("The fail min rule has pass")
	}
}

func TestMinRuleSuccess(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})
	inputData.SetValue("min_field_int", "2")
	inputData.SetValue("min_field_int_negative", "-2")
	inputData.SetValue("min_field_float", "2.0")

	validator := New()
	validator.AddRule("min_field_int", rules.Float(), rules.Min(2))
	validator.AddRule("min_field_int_negative", rules.Float(), rules.Min(-2))
	validator.AddRule("min_field_float", rules.Float(), rules.Min(2))
	validator.Validate(inputData)
	if validator.IsFailed() {
		t.Error("The success min rule failed")
	}
}

func TestMinRuleSkip(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})

	validator := New()
	validator.AddRule("min_field_str", rules.Min(2))
	validator.AddRule("min_field_int", rules.Min(2))
	validator.AddRule("min_field_int_negative", rules.Min(-2))
	validator.AddRule("min_field_float", rules.Min(2))
	validator.Validate(inputData)
	if validator.IsFailed() {
		t.Error("The min rule has pass")
	}
}

// ------------
func TestMaxRuleFail(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})
	inputData.SetValue("max_field_int", "2")
	inputData.SetValue("max_field_int_negative", "-2")
	inputData.SetValue("max_field_float", "2.0")

	validator := New()
	validator.AddRule("max_field_int", rules.Float(), rules.Max(1))
	validator.AddRule("max_field_int_negative", rules.Float(), rules.Max(1))
	validator.AddRule("max_field_float", rules.Float(), rules.Max(1))
	validator.Validate(inputData)
	if !validator.IsFailed() {
		t.Error("The fail max rule has pass")
	}
}

func TestMaxRuleSuccess(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})
	inputData.SetValue("max_field_int", "2")
	inputData.SetValue("max_field_int_negative", "-2")
	inputData.SetValue("max_field_float", "2.0")

	validator := New()
	validator.AddRule("max_field_int", rules.Float(), rules.Max(5))
	validator.AddRule("max_field_int_negative", rules.Float(), rules.Max(-1))
	validator.AddRule("max_field_float", rules.Float(), rules.Max(5))
	validator.Validate(inputData)
	if validator.IsFailed() {
		t.Error("The success max rule failed")
	}
}

func TestMaxRuleSkip(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})

	validator := New()
	validator.AddRule("max_field_str", rules.Max(5))
	validator.AddRule("max_field_int", rules.Float(), rules.Max(5))
	validator.AddRule("max_field_int_negative", rules.Float(), rules.Max(5))
	validator.AddRule("max_field_float", rules.Float(), rules.Max(5))
	validator.Validate(inputData)
	if validator.IsFailed() {
		t.Error("The max rule has pass")
	}
}

func TestRegexpRuleSuccess(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})
	inputData.SetValue("field", "1222345")

	validator := New()
	validator.AddRule("field", rules.Regexp(`^\d+$`))
	validator.Validate(inputData)
	if validator.IsFailed() {
		t.Error("The success regexp rule failed")
	}
}

func TestRegexpRuleFail(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})
	inputData.SetValue("field", "12a22345")

	validator := New()
	validator.AddRule("field", rules.Regexp(`^\d+$`))
	validator.Validate(inputData)
	if !validator.IsFailed() {
		t.Error("The fail regexp rule has pass")
	}
}

func TestEnumRuleSuccess(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})
	inputData.SetValue("field", "enum1")

	validator := New()
	validator.AddRule("field", rules.Enum("enum2", "enum1"))
	validator.Validate(inputData)
	if validator.IsFailed() {
		t.Error("The success enum rule failed")
	}
}

func TestEnumRuleFail(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})
	inputData.SetValue("field", "enum1")

	validator := New()
	validator.AddRule("field", rules.Enum("enum2", "enum2"))
	validator.Validate(inputData)
	if !validator.IsFailed() {
		t.Error("The fail enum rule has pass")
	}
}

func TestConfirmedRuleSuccess(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})
	inputData.SetValue("field", "test")
	inputData.SetValue("field2", "test")

	validator := New()
	validator.AddRule("field", rules.Confirmed("field2"))
	validator.Validate(inputData)
	if validator.IsFailed() {
		t.Error("The success confirmed rule failed")
	}
}

func TestConfirmedRuleFail(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})
	inputData.SetValue("field", "test1")
	inputData.SetValue("field2", "test2")

	validator := New()
	validator.AddRule("field", rules.Confirmed("field2"))
	validator.Validate(inputData)
	if !validator.IsFailed() {
		t.Error("The fail confirmed rule has pass")
	}
}

func TestIntegerRuleSuccess(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})
	inputData.SetValue("field", "1")

	validator := New()
	validator.AddRule("field", rules.Integer())
	validator.Validate(inputData)
	if validator.IsFailed() {
		t.Error("The success integer rule failed")
	}
}

func TestIntegerRuleFail(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})
	inputData.SetValue("field", "33a11", "33")

	validator := New()
	validator.AddRule("field", rules.Integer())
	validator.Validate(inputData)
	if !validator.IsFailed() {
		t.Error("The fail integer rule has pass")
	}
}

func TestAcceptedRuleSuccess(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})
	inputData.SetValue("field", "on")

	validator := New()
	validator.AddRule("field", rules.Accepted())
	validator.Validate(inputData)
	if validator.IsFailed() {
		t.Error("The success accepted rule failed")
	}
}

func TestAcceptedRuleFail(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})
	inputData.SetValue("field", "no")

	validator := New()
	validator.AddRule("field", rules.Accepted())
	validator.Validate(inputData)
	if !validator.IsFailed() {
		t.Error("The fail accepted rule has pass")
	}
}

func TestLengthRuleSuccess(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})
	inputData.SetValue("field", "light")

	validator := New()
	validator.AddRule("field", rules.Length(5))
	validator.Validate(inputData)
	if validator.IsFailed() {
		t.Error("The success length rule failed")
	}
}

func TestLengthMinRuleSuccess(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})
	inputData.SetValue("field", "light")

	validator := New()
	validator.AddRule("field", rules.Length(5).SetMin(3))
	validator.Validate(inputData)
	if validator.IsFailed() {
		t.Error("The success length rule failed")
	}
}

func TestLengthRuleFail(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})
	inputData.SetValue("field", "light")

	validator := New()
	validator.AddRule("field", rules.Length(4))
	validator.Validate(inputData)
	if !validator.IsFailed() {
		t.Error("The fail length rule has pass")
	}
}

func TestLengthCyrillicRuleFail(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})
	inputData.SetValue("field", "проверка")

	validator := New()
	validator.AddRule("field", rules.Length(8))
	validator.Validate(inputData)
	if !validator.IsFailed() {
		t.Error("The fail length rule has pass")
	}
}

func TestLengthCyrillicRuleSuccess(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})
	inputData.SetValue("field", "проверка")

	validator := New()
	validator.AddRule("field", rules.Length(8).AsRunes())
	validator.Validate(inputData)
	if validator.IsFailed() {
		t.Error("The success length rule failed")
	}
}

func TestRuleCollectionSuccess(t *testing.T) {
	inputData := input.FormsToInputData(&url.Values{}, &multipart.Form{})
	inputData.SetValue("field", "test")

	ruleCollection := NewRuleCollection()
	ruleCollection.AddRule("field", rules.Required())
	ruleCollection.AddRule("field", rules.Length(4).AsRunes())
	ruleCollection.AddRule("field", rules.Regexp(`.{4}`))

	validator := New()
	validator.SetRuleCollection(ruleCollection)
	validator.Validate(inputData)
	if validator.IsFailed() {
		t.Error("The success rule collection failed")
	}
}

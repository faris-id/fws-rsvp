package validator

import (
	"reflect"
	"strings"

	"github.com/faris-arifiansyah/fws-rsvp/response"
)

const (
	JsonKey     = "json"
	OmitTag     = "-"
	RequiredTag = "required"
)

func Validate(s interface{}) []error {
	var errors []error

	st := reflect.TypeOf(s)
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		if err := validateField(field, reflect.ValueOf(s).Field(i)); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func validateField(f reflect.StructField, v reflect.Value) error {
	json := f.Tag.Get(JsonKey)
	required := strings.Contains(json, RequiredTag)
	if json == "" || json == OmitTag || !required {
		return nil
	}

	field := strings.Split(json, ",")[0]

	if v.Interface() == reflect.Zero(f.Type).Interface() {
		err := response.BadRequestError
		err.Field = field
		return err
	}

	return nil
}

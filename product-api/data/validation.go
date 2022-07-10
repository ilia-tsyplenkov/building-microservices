package data

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	validator.FieldError
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf(
		"Key: %q Error: key validation for %s failed on %s tag",
		v.Namespace(),
		v.Field(),
		v.Tag(),
	)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Errors() []string {
	var errs []string
	for _, vErr := range v {
		errs = append(errs, vErr.Error())
	}
	return errs
}

type Validation struct {
	validate *validator.Validate
}

func NewValidation() *Validation {
	validate := validator.New()
	validate.RegisterValidation("sku", skuValidation)
	return &Validation{validate}
}
func (v Validation) Validate(i interface{}) ValidationErrors {
	err := v.validate.Struct(i)
	if err == nil {
		return nil
	}
	errs := err.(validator.ValidationErrors)
	if len(errs) == 0 {
		return nil
	}
	var returnErrs ValidationErrors
	for _, err := range errs {
		ve := ValidationError{err.(validator.FieldError)}
		returnErrs = append(returnErrs, ve)
	}
	return returnErrs

}

func skuValidation(fl validator.FieldLevel) bool {
	reg := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := reg.FindAllStringSubmatch(fl.Field().String(), -1)
	if len(matches) != 1 {
		return false
	}
	return true
}

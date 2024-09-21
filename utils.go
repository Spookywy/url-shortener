package main

import "github.com/go-playground/validator/v10"

func validateStruct(structToValidate any) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(structToValidate)
}

package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator представляет валидатор для структур
type Validator struct {
	validate *validator.Validate
}

// New создает новый экземпляр валидатора
func New() *Validator {
	v := validator.New()

	// Регистрируем кастомные валидации если нужно
	// v.RegisterValidation("custom_validation", customValidationFunc)

	return &Validator{
		validate: v,
	}
}

// Validate проверяет структуру на соответствие тегам валидации
func (v *Validator) Validate(s interface{}) error {
	if err := v.validate.Struct(s); err != nil {
		return v.formatValidationError(err)
	}
	return nil
}

// formatValidationError форматирует ошибки валидации в читаемый вид
func (v *Validator) formatValidationError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string

		for _, fieldError := range validationErrors {
			field := fieldError.Field()
			tag := fieldError.Tag()
			param := fieldError.Param()

			var message string
			switch tag {
			case "required":
				message = fmt.Sprintf("%s is required", field)
			case "min":
				message = fmt.Sprintf("%s must be at least %s characters long", field, param)
			case "max":
				message = fmt.Sprintf("%s must not exceed %s characters", field, param)
			case "gt":
				message = fmt.Sprintf("%s must be greater than %s", field, param)
			case "gte":
				message = fmt.Sprintf("%s must be greater than or equal to %s", field, param)
			case "lt":
				message = fmt.Sprintf("%s must be less than %s", field, param)
			case "lte":
				message = fmt.Sprintf("%s must be less than or equal to %s", field, param)
			case "email":
				message = fmt.Sprintf("%s must be a valid email address", field)
			case "url":
				message = fmt.Sprintf("%s must be a valid URL", field)
			default:
				message = fmt.Sprintf("%s failed validation: %s", field, tag)
			}

			errorMessages = append(errorMessages, message)
		}

		return fmt.Errorf("validation failed: %s", strings.Join(errorMessages, "; "))
	}

	return err
}

// ValidateStruct проверяет структуру и возвращает ошибки валидации
func (v *Validator) ValidateStruct(s interface{}) []string {
	var errors []string

	if err := v.validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				field := fieldError.Field()
				tag := fieldError.Tag()
				param := fieldError.Param()

				var message string
				switch tag {
				case "required":
					message = fmt.Sprintf("%s is required", field)
				case "min":
					message = fmt.Sprintf("%s must be at least %s characters long", field, param)
				case "max":
					message = fmt.Sprintf("%s must not exceed %s characters", field, param)
				case "gt":
					message = fmt.Sprintf("%s must be greater than %s", field, param)
				case "gte":
					message = fmt.Sprintf("%s must be greater than or equal to %s", field, param)
				case "lt":
					message = fmt.Sprintf("%s must be less than %s", field, param)
				case "lte":
					message = fmt.Sprintf("%s must be less than or equal to %s", field, param)
				case "email":
					message = fmt.Sprintf("%s must be a valid email address", field)
				case "url":
					message = fmt.Sprintf("%s must be a valid URL", field)
				default:
					message = fmt.Sprintf("%s failed validation: %s", field, tag)
				}

				errors = append(errors, message)
			}
		}
	}

	return errors
}

// GetFieldValue получает значение поля структуры по имени
func (v *Validator) GetFieldValue(s interface{}, fieldName string) interface{} {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	field := val.FieldByName(fieldName)
	if field.IsValid() {
		return field.Interface()
	}

	return nil
}

package dto

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type UpdatePriceDto struct {
	DateRange *DateRangeDTO `json:"date_range" validate:"omitempty"`
	Price     *float64      `json:"price" validate:"omitempty,min=0"`
	Type      *string       `json:"type" validate:"omitempty,oneof=PerApartmentUnit PerGuest"`
}

type DateRangeDTO struct {
	Start time.Time `json:"start" validate:"required"`
	End   time.Time `json:"end" validate:"required,gtefield=Start"`
}

func ValidateUpdatePriceDto(dto *UpdatePriceDto) error {
	validate := validator.New()

	// Register custom validators
	validate.RegisterValidation("validateDateRange", validateDateRange)

	// Validate UpdatePriceDto
	if err := validate.Struct(dto); err != nil {
		return err
	}

	// Validate DateRangeDTO if present
	if dto.DateRange != nil {
		if err := validate.Struct(dto.DateRange); err != nil {
			return err
		}
	}

	return nil
}

// Define a custom validator for DateRangeDTO
func validateDateRange(fl validator.FieldLevel) bool {
	start, ok := fl.Parent().FieldByName("Start").Interface().(time.Time)
	if !ok {
		return false
	}
	end := fl.Field().Interface().(time.Time)
	return end.After(start)
}

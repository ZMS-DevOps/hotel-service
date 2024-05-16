package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/mmmajder/zms-devops-hotel-service/domain"
)

type AccommodationDto struct {
	Name         string          `json:"name" validate:"required"`
	Location     string          `json:"location" validate:"required"`
	Benefits     []string        `json:"benefits"`
	Photos       []string        `json:"photos"`
	GuestNumber  GuestNumberDto  `json:"guest_number" validate:"required"`
	DefaultPrice DefaultPriceDto `json:"default_price"  validate:"required"`
	//Email     string `json:"email" validate:"required,email"`
	//Age       int    `json:"age" validate:"gte=18,lte=120"`
}

type GuestNumberDto struct {
	Min int `json:"min" validate:"required,min=1"`
	Max int `json:"max" validate:"required,min=1,gtefield=Min"`
}

type DefaultPriceDto struct {
	Price float64 `json:"price" validate:"min=0"`
	Type  string  `json:"type"`
}

func ValidateAccommodationDto(dto *AccommodationDto) error {
	validate := validator.New()
	validate.RegisterStructValidation(validateGuestNumberDto, GuestNumberDto{})
	validate.RegisterStructValidation(validateDefaultPriceDto, DefaultPriceDto{})
	validate.RegisterValidation("pricetype", validatePricingType)
	return validate.Struct(dto)
}

func validateGuestNumberDto(sl validator.StructLevel) {
	dto := sl.Current().Interface().(GuestNumberDto)
	if dto.Min > dto.Max {
		sl.ReportError(dto.Min, "Min", "Min", "ltefield", "")
		sl.ReportError(dto.Max, "Max", "Max", "gtefield", "")
	}
}

func validateDefaultPriceDto(sl validator.StructLevel) {
	dto := sl.Current().Interface().(DefaultPriceDto)
	if dto.Type == "" {
		sl.ReportError(dto.Type, "Type", "Type", "required", "")
	}
}

func validatePricingType(fl validator.FieldLevel) bool {
	val := fl.Field().Interface().(domain.PricingType)
	switch val {
	case domain.PerApartmentUnit, domain.PerGuest:
		return true
	default:
		return false
	}
}
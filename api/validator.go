package api

import (
	"github.com/Ritik1101-ux/simplebank/utils"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func=func(fieldLevel validator.FieldLevel) bool {
	if currency,ok := fieldLevel.Field().Interface().(string); ok{
        return utils.IsSupportedCurrency(currency)
	}
	return false;
}

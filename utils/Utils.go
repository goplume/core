package utils

import (
	"fmt"
	"github.com/goplume/core/fault"
	"strconv"
)

// method check valid format msisdn
// if msisdn is valid then return true and err nil
// if msisdn is invalid  return false and error struct
func IsValidMsisdn(msisdn string) (
	bool,             // valid
	fault.TypedError, //err
) {
	if msisdn == "" {
		return false, fault.NewValidatorError("Msisdn is empty")
	}

	if len(msisdn) != 11 {
		return false, fault.NewValidatorError(
			fmt.Sprintf("Msisdn '%v is not valid format ", msisdn))
	}

	msisdnInt, err := strconv.ParseUint(msisdn, 10, 64)
	if err != nil {
		return false, fault.NewValidatorError(
			fmt.Sprintf("Msisdn '%v' is not valid format ", msisdn))
	}

	if msisdnInt == 0 {
		return false, fault.NewValidatorError(
			fmt.Sprintf("Msisdn '%v' is not valid format ", msisdn))
	}

	return true, nil
}

func IsValidMpan(mpan string) (bool, fault.TypedError) {
	if mpan == "" {
		return false, fault.ExceptionIllegalArgument("mpan", " Mpan not must be empty")
	}

	if len(mpan) != 16 {
		return false, fault.NewValidatorError(
			fmt.Sprintf("MPAN %v is not valid format ", mpan))
	}

	msisdnInt, err := strconv.ParseUint(mpan, 10, 64)
	if err != nil {
		return false, fault.NewValidatorError(
			fmt.Sprintf("MPAN '%v' is not valid format ", mpan))
	}

	if msisdnInt == 0 {
		return false, fault.NewValidatorError(
			fmt.Sprintf("MPAN '%v' is not valid format ", mpan))
	}

	return true, nil
}

func ToString(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%s", value)
}

func IsEmptyString(string string) bool {
	if len(string) == 0 {
		return true
	}
	return false
}

func IsBlankString(string string) bool {
	if len(string) == 0 {
		return true
	}

	for _, v := range string {
		if v != ' ' && v != '\n' && v != '\t' && v != '\r' {
			return false
		}
	}
	return true
}
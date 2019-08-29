package fault

// See https://jsonapi.org/format/#errors
type ApiErrors struct {
	Errors []*TypedErrorStr `json_off:"errors"`
}

//func (errors *ApiErrors) State() int {
//	return errors.Errors[0].State
//}

type TypedError interface {
	Error() string
	Type() ClassError
	Message() string
}

//type ApiError interface {
//	Error() string
//	Type() ClassError
//	Message() string
//}

type TypedErrorStr struct {
	Class ClassError `json:"Class,omitempty"`
	//TraceId string
	Msg string `json:"Message,omitempty" example:"status bad request" || "Not found"`
	// the HTTP status code applicable to this problem, expressed as a string value.
	Status int `json:"State,omitempty"  example:"5xx" or "4xxx"`
	//  an application-specific error code, expressed as a string value
	//Code string `json:"Code"  example:"database_error" or "invalid_group"`
	Err error `json:"-"`
	// a short, human-readable summary of the problem that SHOULD NOT change from occurrence to occurrence of the problem, except for purposes of localization.
	Title string `json:"Title,omitempty"`
	// a human-readable explanation specific to this occurrence of the problem. Like title, this field’s value can be localized.
	Details interface{} `json:"Details,omitempty"`
	Caller  Caller      `json:"-"`
}

//func IsError(err error, class string) bool {
//	if err == nil {
//		return false
//	}
//
//	switch err.(type) {
//	case TypedError:
//		{
//			return strings.EqualFold(err.(TypedError).Class, class)
//		}
//	default:
//		return false
//
//	}
//	return false
//}

type Caller struct {
	Pc   uintptr
	File string
	Line int
	Ok   bool
}

func (this TypedErrorStr) Error() string {
	return string(this.Class) + ": " + this.Msg
}

func (this TypedErrorStr) Type() ClassError {
	return this.Class
}

func (this TypedErrorStr) Message() string {
	return this.Msg
}

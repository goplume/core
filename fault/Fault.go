package fault

import (
    "github.com/go-resty/resty/v2"
    "runtime"
)

type ClassError string

const (
    EXCEPTION_UNAUTHORIZED             ClassError = "EXCEPTION_UNAUTHORIZED"
    ERROR_PARSE_REQUEST_BODY           ClassError = "ERROR_PARSE_REQUEST_BODY"
    INTEGRATION_EXCEPTION_SERVER_ERROR ClassError = "INTEGRATION_EXCEPTION_SERVER_ERROR"
    INTEGRATION_EXCEPTION_CLIENT_ERROR ClassError = "INTEGRATION_EXCEPTION_CLIENT_ERROR"
    EXCEPTION_ILLEGAL_ARGUMENT         ClassError = "EXCEPTION_ILLEGAL_ARGUMENT"
    EXCEPTION_ILLEGAL_STATE            ClassError = "EXCEPTION_ILLEGAL_STATE"
    EXCEPTION_INVALID_VALUE            ClassError = "EXCEPTION_INVALID_VALUE"
    EXCEPTION_FRAUD_DETECT             ClassError = "EXCEPTION_FRAUD_DETECT"
    EXCEPTION_INTERNAL_ERROR           ClassError = "EXCEPTION_INTERNAL_ERROR"
    EXCEPTION_FORMAT_ERROR             ClassError = "EXCEPTION_FORMAT_ERROR"
    EXCEPTION_ENTITY_NOT_FOUND         ClassError = "EXCEPTION_ENTITY_NOT_FOUND"
    PERSISTENCE_ERROR                  ClassError = "PERSISTENCE_ERROR"
)

var ErrorPrefix = map[ClassError]string{
    EXCEPTION_UNAUTHORIZED:             "Unauthorized: ",
    ERROR_PARSE_REQUEST_BODY:           "Error parse request body: ",
    INTEGRATION_EXCEPTION_SERVER_ERROR: "Integration Server Error: ",
    INTEGRATION_EXCEPTION_CLIENT_ERROR: "Integration Client Error: ",
    EXCEPTION_ILLEGAL_ARGUMENT:         "Illegal argument: ",
    EXCEPTION_ILLEGAL_STATE:            "Illegal state: ",
    EXCEPTION_INTERNAL_ERROR:           "Internal Error: ",
    EXCEPTION_ENTITY_NOT_FOUND:         "Entity Not Found: ",
    PERSISTENCE_ERROR:                  "Persistence error: ",
}

func buildCaller(pc uintptr, file string, line int, ok bool) Caller {
    return Caller{
        Pc:   pc,
        File: file,
        Line: line,
        Ok:   ok,
    }
}

func CreateErrorParseRequestBody(err error) error {
    if err == nil {
        return nil
    }
    return &TypedErrorStr{
        Class:  ERROR_PARSE_REQUEST_BODY,
        Msg:    "Parse body: " + err.Error(),
        Err:    err,
        Caller: buildCaller(runtime.Caller(1)),
    }
}

func ExceptionIllegalArgument(arg, msg string) TypedError {
    return TypedErrorStr{
        Class:  EXCEPTION_ILLEGAL_ARGUMENT,
        Msg:    msg,
        Caller: buildCaller(runtime.Caller(1)),
    }
}

func ExceptionEntityNotFound(msg string) TypedError {
    return TypedErrorStr{
        Class:  EXCEPTION_ENTITY_NOT_FOUND,
        Msg:    msg,
        Caller: buildCaller(runtime.Caller(1)),
    }
}

// todo wrap all Err in all repositories
func NewPersisnteceError(err error) TypedError {
    return TypedErrorStr{
        Class:  PERSISTENCE_ERROR,
        Msg:    err.Error(),
        Err:    err,
        Caller: buildCaller(runtime.Caller(1)),
    }
}

func NewPersisnteceErrorM(message string, err error) TypedError {
    return TypedErrorStr{
        Class:  PERSISTENCE_ERROR,
        Msg:    message + err.Error(),
        Err:    err,
        Caller: buildCaller(runtime.Caller(1)),
    }
}

func ExceptionIllegalState(msg string) TypedError {
    return TypedErrorStr{
        Class:  EXCEPTION_ILLEGAL_STATE,
        Msg:    msg,
        Caller: buildCaller(runtime.Caller(1)),
    }
}

func ExceptionInvalidValue(msg string) TypedError {
    return TypedErrorStr{
        Class:  EXCEPTION_INVALID_VALUE,
        Msg:    msg,
        Caller: buildCaller(runtime.Caller(1)),
    }
}

func ExceptionFraudDetect(msg string) TypedError {
    return TypedErrorStr{
        Class:  EXCEPTION_FRAUD_DETECT,
        Msg:    msg,
        Caller: buildCaller(runtime.Caller(1)),
    }
}

func ExceptionUnauthorized(msg string) TypedError {
    return TypedErrorStr{
        Class:  EXCEPTION_UNAUTHORIZED,
        Msg:    msg,
        Caller: buildCaller(runtime.Caller(1)),
    }
}

func ExceptionInternalError(msg string) TypedError {
    return TypedErrorStr{
        Class:  EXCEPTION_INTERNAL_ERROR,
        Msg:    msg,
        Caller: buildCaller(runtime.Caller(1)),
    }
}

func ExceptionFormatError(msg string) TypedError {
    return TypedErrorStr{
        Class:  EXCEPTION_FORMAT_ERROR,
        Msg:    msg,
        Caller: buildCaller(runtime.Caller(1)),
    }
}

func ExceptionInternalErrorE(err error) TypedError {
    return TypedErrorStr{
        Class:  EXCEPTION_INTERNAL_ERROR,
        Msg:    err.Error(),
        Err:    err,
        Caller: buildCaller(runtime.Caller(1)),
    }
}

func IntegrationExceptionServerErrorR(msg string, response *resty.Response) TypedError {
    return TypedErrorStr{
        Class:  EXCEPTION_INTERNAL_ERROR,
        Msg:    msg,
        Caller: buildCaller(runtime.Caller(1)),
    }
}

func ExceptionClientError(msg string, Status int) TypedError {
    return TypedErrorStr{
        Class:  INTEGRATION_EXCEPTION_CLIENT_ERROR,
        Status: Status,
        Msg:    msg,
        Caller: buildCaller(runtime.Caller(1)),
    }
}

func IntegratrionExceptionClientErrorR(response *resty.Response) TypedError {
    message := ""
    if len(response.Body()) > 0 {
        message = string(response.Body())
    }
    return TypedErrorStr{
        Class:  INTEGRATION_EXCEPTION_CLIENT_ERROR,
        Status: response.StatusCode(),
        Msg:    message,
        Caller: buildCaller(runtime.Caller(1)),
    }
}

// To handle the error returned by c.Bind in gin framework
// https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go
func NewValidatorError(message string) TypedError {

    return TypedErrorStr{
        Class:  EXCEPTION_ILLEGAL_ARGUMENT,
        Msg:    message,
        Caller: buildCaller(runtime.Caller(1)),
    }
}

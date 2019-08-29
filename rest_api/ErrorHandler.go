package rest_api

import (
	"github.com/goplume/core/fault"
	"github.com/goplume/core/utils/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)


func HandleError(err error, ctx *gin.Context, log *logger.Logger) bool {
	if err == nil {
		return false
	}

	log.RLog.Error(err)
	//r := err
	//errorMessage := r.Error()

	switch err.(type) {
	case fault.TypedError:
		switch err.(fault.TypedError).Type() {
		case fault.EXCEPTION_UNAUTHORIZED:
			ThrowError(ctx, http.StatusUnauthorized, err)
		case fault.ERROR_PARSE_REQUEST_BODY:
			ThrowError(ctx, http.StatusBadRequest, err)
		case fault.EXCEPTION_ILLEGAL_ARGUMENT:
			ThrowError(ctx, http.StatusBadRequest, err)
		case fault.EXCEPTION_ENTITY_NOT_FOUND:
			ThrowError(ctx, http.StatusNotFound, err)
		case fault.EXCEPTION_ILLEGAL_STATE:
			ThrowError(ctx, http.StatusConflict, err)
		case fault.INTEGRATION_EXCEPTION_SERVER_ERROR:
			ThrowError(ctx, http.StatusInternalServerError, err)
		case fault.INTEGRATION_EXCEPTION_CLIENT_ERROR:
			ThrowError(ctx, http.StatusInternalServerError, err)
		case fault.EXCEPTION_INTERNAL_ERROR:
			ThrowError(ctx, http.StatusInternalServerError, err)
		case fault.PERSISTENCE_ERROR:
			ThrowError(ctx, http.StatusInternalServerError, err)
		default:
			ThrowError(ctx, http.StatusInternalServerError, err)
		}
	default:
		//if errorMessage == http.StatusText(http.StatusNotFound) {
		//	ThrowError(ctx, http.StatusNotFound, errorMessage, nil)
		//
		//} else if errorMessage == http.StatusText(http.StatusUnauthorized) {
		//	ThrowError(ctx, http.StatusUnauthorized, errorMessage, nil)
		//} else if strings.HasPrefix(errorMessage, string(fault.EXCEPTION_UNAUTHORIZED)) {
		//	ThrowError(ctx, http.StatusUnauthorized, errorMessage, nil)
		//} else if strings.HasPrefix(errorMessage, string(fault.EXCEPTION_ILLEGAL_ARGUMENT)) {
		//	ThrowError(ctx, http.StatusBadRequest, errorMessage, nil)
		//
		//} else if strings.HasPrefix(errorMessage, string(fault.EXCEPTION_ILLEGAL_STATE)) {
		//	ThrowError(ctx, http.StatusConflict, errorMessage, nil)
		//
		//} else {
		//	// todo Don`t show internal text of error for security
		//	ThrowError(ctx, http.StatusInternalServerError, errorMessage, nil)
		//}
		ThrowError(ctx, http.StatusInternalServerError, err)
	}

	return ctx.IsAborted()
}

// Please, use сtx.AbortWithError
// Deprecated
func ThrowError(ctx *gin.Context, status int, err error) {

	//switch err.(type) {
	//
	//}
	//apiError := fault.TypedError{
	//	State:  status,
	//	Code:    http.StatusText(status),
	//	Title:   title,
	//	Details: detail,
	//}
	ErrorRestResponse(ctx, status, err)
}

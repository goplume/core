package utils

import (
	"github.com/goplume/core/fault"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

type ParamType string

const (
	PATH_PARAM   ParamType = "PATH_PARAM"
	QUERY_PARAM  ParamType = "QUERY_PARAM"
	HEADER_PARAM ParamType = "HEADER_PARAM"
)

func PathParamAsUint64(
	ctx *gin.Context, paramName string, requered bool,
) (
	value uint64, err error,
) {
	paramValue, defined := GetParam(ctx, paramName)
	return parseUint64Param(paramValue, defined, requered, paramName)
}

func CtxParamAsUint64(
	paramType ParamType, ctx *gin.Context, paramName string, requered bool,
) (
	value uint64, err error,
) {
	var paramValue string
	var defined bool
	switch paramType {
	case PATH_PARAM:
		paramValue, defined = ctx.Param(paramName), StringIsNotEmpty(ctx.GetHeader(paramName))
	case QUERY_PARAM:
		paramValue, defined = ctx.GetQuery(paramName)
	case HEADER_PARAM:
		paramValue, defined = ctx.GetHeader(paramName), StringIsNotEmpty(ctx.GetHeader(paramName))
	default:
		err = fault.ExceptionIllegalArgument("",
			"Unknow request context param type '"+string(paramType)+"' ")
	}
	return parseUint64Param(paramValue, defined, requered, paramName)
}

func QueryParamAsUnint64(
	ctx *gin.Context, paramName string, requered bool,
) (
	value uint64, err error,
) {
	paramValue, defined := ctx.GetQuery(paramName)
	return parseUint64Param(paramValue, defined, requered, paramName)
}

func PathParam(
	ctx *gin.Context, paramName string, requered bool,
) (
	value string, err error,
) {
	paramValue, defined := GetParam(ctx, paramName)
	return parseStringParam(paramValue, defined, requered, paramName)
}

func QueryParam(
	ctx *gin.Context, paramName string, requered bool,
) (
	value string, err error,
) {
	paramValue, defined := ctx.GetQuery(paramName)
	return parseStringParam(paramValue, defined, requered, paramName)
}

func PathParamAsTime(
	ctx *gin.Context, paramName string, requered bool,
) (
	time.Time, bool, error,
) {
	paramValue, defined := GetParam(ctx, paramName)
	return parseTimeParam(paramValue, defined, requered, paramName)
}

func QueryParamAsTime(
	ctx *gin.Context, paramName string, requered bool,
) (
	time.Time, bool, error,
) {
	paramValue, defined := ctx.GetQuery(paramName)
	return parseTimeParam(paramValue, defined, requered, paramName)
}

func parseUint64Param(
	paramValue string, defined bool, requered bool, paramName string,
) (
	value uint64, err error,
) {
	if defined == false && requered == true {
		err = fault.ExceptionIllegalArgument(paramName,
			"Mandatory params "+paramName+" missing")
	} else {
		if StringIsBlank(paramValue) {
			value = 0
		} else {
			value, err = strconv.ParseUint(paramValue, 10, 32)
			if err != nil {
				err = fault.ExceptionIllegalArgument(paramName,
					"Param "+paramName+" invalid "+err.Error())
			}
		}
	}
	return
}

func parseStringParam(
	paramValue string, defined bool, requered bool, paramName string,
) (
	value string, err error,
) {
	if defined == false && requered == true {
		err = fault.ExceptionIllegalArgument(paramName,
			"Mandatory params "+paramName+" missing")
	} else {
		if StringIsBlank(paramValue) && requered == true {
			err = fault.ExceptionIllegalArgument(paramName,
				"Mandatory params "+paramName+" missing")
		}
	}
	value = paramValue
	return
}

func parseTimeParam(
	paramValue string, defined bool, requered bool, paramName string,
) (
	value time.Time, def bool, err error,
) {
	def = defined
	if defined == false && requered == true {
		err = fault.ExceptionIllegalArgument(paramName,
			"Mandatory params "+paramName+" missing")
	} else if defined == false && requered == false {
		return time.Time{}, false, nil
	} else {
		value, err = ParseTime(paramValue)
		if err != nil {
			err = fault.ExceptionIllegalArgument(paramName,
				"Invalid format '"+paramName+"': "+err.Error())
		}
	}
	return value, def, err
}

func ParseTime(paramValue string) (time.Time, error) {
	paramValue = strings.ReplaceAll(paramValue, " ","+")
	value, err := time.Parse(DATETIME_PARSE_LAYOUT, paramValue)
	return value, err
}

func PathParamAsString(
	ctx *gin.Context, pathParamName string, requered bool,
) (
	string, error,
) {
	value := ctx.Param(pathParamName)
	if value == "" && requered {
		return "", fault.ExceptionIllegalArgument(pathParamName,
			"Mandatory params "+pathParamName+" missing")
	}
	return value, nil
}

func HeaderParam(
	ctx *gin.Context, pathParamName string, requered bool,
) (
	string, error,
) {
	value := ctx.GetHeader(pathParamName)
	if value == "" && requered {
		return "", fault.ExceptionIllegalArgument(pathParamName,
			"Mandatory params "+pathParamName+" missing")
	}
	return value, nil
}

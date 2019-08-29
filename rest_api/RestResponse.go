package rest_api

import (
    "github.com/goplume/core/rest_api/context"
    "github.com/gin-gonic/gin"
    "time"
)

type ResponseStatus string
type EntityName string

const (
    // for all http code from [100 : 400)
    ResponseStatusSuccess ResponseStatus = "success"
    // for all http code from [400 : 500)
    ResponseStatusFail ResponseStatus = "fail"
    // for all http code from [500 : 600)
    ResponseStatusError ResponseStatus = "error"
)

// Rest Response
// See link https://github.com/omniti-labs/jsend
// See link https://technick.net/guides/software/software_json_api_format/
type RestResponse struct {
	Version   string         `json:"version,omitempty"`   // Api Version
    DateTime  string         `json:"datetime,omitempty"`  // Human-readable date and time when the event occurred
    Timestamp int64          `json:"timestamp,omitempty"` // Machine-readable UTC timestamp in nanoseconds since EPOCH
    Status    ResponseStatus `json:"status"`              // State code (error|fail|success)
    Code      int            `json:"code"`                // HTTP status code
    Message   string         `json:"message"`             // Error or status message
    Data      interface{}    `json:"data,omitempty"`      // Data payload
    Meta      interface{}    `json:"meta,omitempty"`      // Meta information about reguest
    Errors    interface{}    `json:",omitempty"`          // List errors
}


//Program   string         `json:"-"`                   // Program name
//Release   string         `json:"-"`                   // Program release number
//URL       string         `json:"-"`                   // Public URL of this service


func FailRestResponse(ctx *gin.Context, httpCode int, message string) RestResponse {
    response := buildResponse(ctx, httpCode)
    response.Message = message
    ctx.AbortWithStatusJSON(httpCode, response)
    return response
}

func ErrorRestResponse(ctx *gin.Context, httpCode int, err error) RestResponse {
    response := buildResponse(ctx, httpCode)

    //response.Message = err.Error()
    response.Errors = []interface{}{err}
    ctx.AbortWithStatusJSON(httpCode, response)
    return response
}

func SuccessRestResponse(
    ctx *gin.Context,
    httpCode int,
    message string,
    data interface{},
) RestResponse {
    response := buildResponse(ctx, httpCode)
    response.Message = message
    response.Data = data
    ctx.JSON(httpCode, response)
    return response
}

func SuccessRestResponseEntity(
    ctx *gin.Context,
    httpCode int,
    message string,
    entityName EntityName,
    entity interface{},
) RestResponse {
    response := buildResponse(ctx, httpCode)
    response.Message = message
    response.Data = map[EntityName]interface{}{entityName: entity}
    ctx.JSON(httpCode, response)
    return response
}

func buildMetaInfo(ctx *gin.Context) map[string]interface{} {
    meta := map[string]interface{}{}
    //requestId, exists := ctx.Get(context.CTX_REQUEST_ID)
    //if exists {
    //	meta[context.CTX_REQUEST_ID] = requestId
    //}
    state, exists := ctx.Get(context.CTX_CURRENT_STATE)
    if exists {
        meta["request_state"] = state
    }
    meta["request_query_params"] = ctx.Request.URL.Query()

    return meta
}

func buildResponse(ctx *gin.Context, httpCode int) RestResponse {
    var requestTime time.Time
    value, exist := ctx.Get("request_time")
    if exist {
        requestTime = value.(time.Time)
    } else {
        requestTime = time.Now()
    }

    metaInfo := buildMetaInfo(ctx)
    response := RestResponse{

        Status:    getStatus(httpCode),
        Code:      httpCode,
        Meta:      metaInfo,
        DateTime:  requestTime.Format(time.RFC3339),
        Timestamp: requestTime.Unix(),
    }
    response.Version = ctx.GetString("api_version")

    return response
}

// convert the HTTP status code into JSend status
func getStatus(code int) ResponseStatus {
    if code >= 500 {
        return ResponseStatusError
    }
    if code >= 400 {
        return ResponseStatusFail
    }
    return ResponseStatusSuccess
}

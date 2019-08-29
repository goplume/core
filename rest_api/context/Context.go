package context

import (
    "github.com/goplume/core/types"
    "github.com/goplume/core/fault"
    utils2 "github.com/goplume/core/utils"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "net/http"
    "time"
)

const (
    CTX_MPAN            = "CTX_MPAN"
    CTX_MSISDN          = "CTX_MSISDN"
    CTX_TOKEN_INFO      = "CTX_TOKEN_INFO"
    CTX_CURRENT_STATE   = "CTX_CURRENT_STATE"
    CTX_MERCHANT_ID     = "CTX_MERCHANT_ID"
    CTX_REQUEST_CHANNEL = "channel"
    CTX_REQUEST_ID      = "request_id"
    CTX_REQUEST_TIME    = "current_time"
    CTX_ACTION          = "action"
)

func GetReqeustId(ctx *gin.Context) (types.RequestId) {
    requestId :=  ctx.GetString(CTX_REQUEST_ID)
    if len(requestId) != 0 {
        return types.RequestId(requestId)
    }

    requestId = ctx.GetHeader(CTX_REQUEST_ID)
    if len(requestId) == 0 {
        requestId, _ = ctx.GetQuery(CTX_REQUEST_ID)
    }

    if len(requestId) == 0 {
        // default behavior
        requestId = string(GenerateRequestId())
    }
    ctx.Set(CTX_REQUEST_ID, requestId)
    return types.RequestId(requestId)

}

func GetRequestChannel(ctx *gin.Context) types.Channel {
    requestChanel := ctx.GetString(CTX_REQUEST_CHANNEL)
    if len(requestChanel) != 0 {
        return types.Channel(requestChanel)
    }

    requestChanel = ctx.GetHeader(CTX_REQUEST_CHANNEL)
    if len(requestChanel) == 0 {
        requestChanel, _ = ctx.GetQuery(CTX_REQUEST_CHANNEL)
    }

    if len(requestChanel) == 0 {
        // default value
        requestChanel = ""
    }
    ctx.Set(CTX_REQUEST_CHANNEL, requestChanel)
    return types.Channel(requestChanel)

}

func GetCurrentTime(ctx *gin.Context) time.Time {

    if currentTimeTime, ok := ctx.Get(CTX_REQUEST_TIME); ok && currentTimeTime != nil {
        return currentTimeTime.(time.Time)
    }

    // default value
    currentTime := time.Now()

    currentTimeStr := ctx.GetHeader(CTX_REQUEST_TIME)
    if len(currentTimeStr) == 0 {
        currentTimeStr, _ = ctx.GetQuery(CTX_REQUEST_TIME)
    }

    if !utils2.StringIsBlank(currentTimeStr) {
        parseTime, parseError := utils2.ParseTime(currentTimeStr)
        if parseError != nil {
            parseTimeErr := fault.ExceptionInternalError("Error parse " + CTX_REQUEST_TIME + " " + parseError.Error())
            ctx.AbortWithError(http.StatusInternalServerError, parseTimeErr)
        }
        currentTime = parseTime
        ctx.Set(CTX_REQUEST_TIME, currentTime)
    }

    return currentTime
}

func GenerateRequestId() types.RequestId {
    //u1 := uuid.Must(uuid.NewV5())
    // todo add log generate uuid for payment
    uuids, _ := uuid.NewRandom()
    requestId := uuids.String()
    return types.RequestId(requestId)
}

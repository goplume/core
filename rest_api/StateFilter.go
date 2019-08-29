package rest_api

import (
    "github.com/goplume/core/rest_api/context"
    "github.com/goplume/core/rules"
    "github.com/gin-gonic/gin"
)

type StateFilter struct {

}

func (this StateFilter) DoCurrentSatet(
    ctx *gin.Context,
) {
    reqeustId := context.GetReqeustId(ctx)
    requestChannel := context.GetRequestChannel(ctx)
    currentTime := context.GetCurrentTime(ctx)

    state := rules.CurrentState{
        RequestId:      reqeustId,
        RequestChannel: requestChannel,
        CurrentTime:    currentTime,
        InputParams:    make(map[string]interface{}, 10),
    }

    ctx.Set(context.CTX_CURRENT_STATE, &state)
}

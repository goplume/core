package rules

import (
    "github.com/goplume/core/types"
    ctypes "github.com/goplume/core/types"
    "time"
)

// Структура определяет текущее состояние
type CurrentState struct {
    // todo to ref
    Action         types.Action
    CurrentTime    time.Time
    RequestChannel types.Channel
    RequestId      types.RequestId
    PaymentsDayId  ctypes.PaymentsDayId
    InputParams    map[string]interface{}
}

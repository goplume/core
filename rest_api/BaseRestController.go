package rest_api

import (
    "github.com/goplume/core/oaut_client"
    "github.com/goplume/core/rest_api/context"
    "github.com/goplume/core/rest_client"
    "github.com/goplume/core/rules"
    "github.com/goplume/core/types"
    utils2 "github.com/goplume/core/utils"
    "github.com/goplume/core/utils/logger"
    "github.com/gin-gonic/gin"
    "github.com/go-resty/resty/v2"
    "github.com/pkg/errors"
    "github.com/sirupsen/logrus"
    "net/http"
    "strings"
)

type BaseRestController struct {
    Log           *logger.Logger
    ServicePrefix string
}

// Deprecated
func (this *BaseRestController) GetBearerAuthToken(ctx *gin.Context) string {
    authValue := ctx.GetHeader("Authorization")
    bearerEndPrefix := len("Bearer ")
    if len(authValue) <= bearerEndPrefix {
        return ""
    }

    return authValue[bearerEndPrefix:]
}

func (this *BaseRestController) RedirectToService(
    ctx *gin.Context,
    restClient *rest_client.RestClient,
) {
    // todo get from ctx
    rlog := this.Log.RLog
    buf := make([]byte, ctx.Request.ContentLength)
    ctx.Request.Body.Read(buf)
    reqBody := string(buf)

    //prefix := "/merchant_api"

    bearerAuthToken := this.GetBearerAuthToken(ctx)
    rlog.Info("Request url: " + ctx.Request.RequestURI)
    redirectUrl := ctx.Request.RequestURI[len(this.ServicePrefix):]
    rlog.Info("Redirect url: " + redirectUrl)
    restClient.HttpClient.SetLogger(rlog)
    response, err := restClient.HttpClient.R().
        SetAuthToken(bearerAuthToken).
        SetBody(reqBody).
        Execute(ctx.Request.Method, redirectUrl)

    if this.ValidateError(err, ctx) {
        this.PushResponse(ctx, response)
    }
}

// Deprecated
func (this *BaseRestController) CallService(
    ctx *gin.Context,
    restClient *resty.Client,
    servicePath string,
) {

    pathParams := make(map[string]string)
    for _, p := range ctx.Params {
        pathParams[p.Key] = p.Value
    }

    queryParams := make(map[string]string)
    for key, _ := range ctx.Request.URL.Query() {
        queryParams[key] = ctx.Request.URL.Query().Get(key)
    }

    buf := make([]byte, ctx.Request.ContentLength)
    ctx.Request.Body.Read(buf)
    reqBody := string(buf)

    response, err := restClient.R().
        SetPathParams(pathParams).
        SetQueryParams(queryParams).
        SetBody(reqBody).
        Execute(ctx.Request.Method, servicePath)

    if this.ValidateError(err, ctx) {
        this.PushResponse(ctx, response)
    }
}

// Deprecated
func (this *BaseRestController) checkToken(token string) (err error) {
    if len(token) < 8 {
        err = errors.New(http.StatusText(http.StatusUnauthorized))
    }
    return err
}

// Deprecated
func (this *BaseRestController) ValidateError(err error, ctx *gin.Context) bool {
    if err != nil {
        HandleError(err, ctx, this.Log)
    }
    return err == nil
}

func (this *BaseRestController) PushResponse(ctx *gin.Context, response *resty.Response) {
    ctx.Data(
        response.RawResponse.StatusCode,
        response.RawResponse.Header.Get("Content-Type"),
        response.Body(),
    )
}

func (this *BaseRestController) ListEntityesResponse(
    ctx *gin.Context,
    entityName EntityName,
    listIsEmptyMessage string,
    entityId func(interface{}) interface{},
    convertor func(interface{}) interface{},
//producerEntityList func() ([]interface{}, fault.TypedError),
    entitiesList ...interface{},
) {
    //list, err := producerEntityList()
    //if HandleError(err, ctx, this.Log) {
    //    return
    //}

    //entitiesList := list.([]interface{})
    if entitiesList == nil || len(entitiesList) == 0 {
        FailRestResponse(ctx, http.StatusNotFound, listIsEmptyMessage)
        return
    }

    if this.IsVerboseDetails(ctx) {
        // return verbose details
        var response []interface{}

        for _, entity := range entitiesList {
            toResponse := convertor(entity)
            response = append(response, toResponse)
        }

        SuccessRestResponseEntity(
            ctx, http.StatusOK, "", entityName,
            response,
        )
    } else {
        // return only ids
        response := make([]interface{}, len(entitiesList))

        for i, entity := range entitiesList {
            response[i] = entityId(entity)
        }

        SuccessRestResponseEntity(ctx, http.StatusOK, "", entityName, response, )
    }
}

func (this *BaseRestController) IsVerboseDetails(context *gin.Context) bool {
    details, _ := utils2.QueryParam(context, Param_details, false)

    return strings.EqualFold(details, "verbose")
}

func (this *BaseRestController) GetMpan(ctx *gin.Context) string {
    return ctx.GetString(context.CTX_MPAN)
}

func (this *BaseRestController) GetTokenInfo(ctx *gin.Context) (token interface{}, exists bool) {
    return ctx.Get(context.CTX_TOKEN_INFO)
}

func (this *BaseRestController) GetCurrentStatePtr(ctx *gin.Context) (*rules.CurrentState) {
    state, _ := ctx.Get(context.CTX_CURRENT_STATE)
    currentState := state.(*rules.CurrentState)
    //action, _ := ctx.Get(context.CTX_ACTION)
    //currentState.Action = action.(types.Action)
    // ctx.Set(context.CTX_CURRENT_STATE, currentState)
    return currentState
}

func (this *BaseRestController) RLogAction(ctx *gin.Context, action types.Action) logrus.FieldLogger {
    mpan := this.GetMpan(ctx)
    clientID := ""
    publicID := ""
    tokenInfo, exists := this.GetTokenInfo(ctx)
    ctx.Set(context.CTX_ACTION, action)
    if exists && tokenInfo != nil {
        //clientID = tokenInfo.(*oaut_client.TokenInfo).ClientID
        //publicID = tokenInfo.(*oaut_client.TokenInfo).PublicID
        clientID = tokenInfo.(*oaut_client.ParamsTokenInfo).ClientID
        publicID = tokenInfo.(*oaut_client.ParamsTokenInfo).PublicID
    }

    state := this.GetCurrentStatePtr(ctx)
    state.Action = action
    //requestChannel := this.GetRequestChannel(ctx)
    var rlog logrus.FieldLogger
    if this.Log != nil && this.Log.RLog != nil {
        rlog = this.Log.RLog
    }

    ctxRlog, rlogFounded := ctx.Get("rlog")
    if rlogFounded {
        rlog = ctxRlog.(logrus.FieldLogger)
    }

    rlog = rlog.WithFields(map[string]interface{}{
        "mpan":       mpan,
        "action":     action,
        "channel":    state.RequestChannel,
        "client_id":  clientID,
        "public_id":  publicID,
        "request_id": state.RequestId,
    })
    ctx.Set("rlog", rlog)

    return rlog
}

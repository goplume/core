package rest_api

import (
    "fmt"
    "github.com/goplume/core/rest_api/context"
    "github.com/goplume/core/fault"
    "github.com/goplume/core/oaut_client"
	"github.com/goplume/core/utils/logger"
    "github.com/gin-gonic/gin"
    "net/http"
    "strings"
)

type AuthFilter struct {
	OAuthClient           *oaut_client.OAuthClient
	Log                   *logger.Logger
	Scope                 string
	Mock                  *AuthFilterMock
	ExcludePath           []string
	SkipCheckMpanPath     []string
	PermanentAccessTokens map[string]interface{}
}

type AuthFilterMock struct {
	Mpan string
}

func (this AuthFilter) DoAuth(
	ctx *gin.Context,
) {

	if this.SkipFilterRules(ctx) {
		return
	}
	var tokenInfo *oaut_client.ParamsTokenInfo
	var err error

	token := this.extractBearerAuthToken(ctx)
	if len(token) == 0 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Required Bearer authorization header")
		return
	}

	for k, v := range this.PermanentAccessTokens {
		if token == k {
			ctx.Set(context.CTX_MPAN, v)
			return
		}
	}

	this.Log.RLog.Info("Check token on scope " + this.Scope)

	tokenInfo, err = this.checkToken(token)
	if HandleError(err, ctx, this.Log) {
		return
	}
	if tokenInfo == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Token info missing")
		return
	}

	ctx.Set(context.CTX_TOKEN_INFO, tokenInfo)

	// validate Mpan
	this.doValidateMpanAuth(tokenInfo, ctx)
	return

}

func (this AuthFilter) doValidateMpanAuth(tokenInfo *oaut_client.ParamsTokenInfo, ctx *gin.Context) {
	if this.SkipMpanCheckRules(ctx) {
		return
	}

	if len(tokenInfo.Params.Mpan) != 16 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Invalid MPAN")
		return
	}

	ctx.Set(context.CTX_MPAN, tokenInfo.Params.Mpan)
	ctx.Set(context.CTX_MSISDN, tokenInfo.Params.Msisdn)

}

func (this *AuthFilter) checkToken(
	token string,
) (tokenInfo *oaut_client.ParamsTokenInfo, err fault.TypedError, ) {
	//if len(token) < 8 {
	//	err = errors.New(http.StatusText(http.StatusUnauthorized))
	//}

	if this.Mock != nil {
		tokenInfo := oaut_client.ParamsTokenInfo{}
		tokenInfo.Params.Mpan = string(this.Mock.Mpan)
		//if tokenInfo.Scope == "qrapi" {
		//	tokenInfo.ClientID = "4598959951692349"
		//}
		return &tokenInfo, nil

	}

	tokenInfo, err = this.OAuthClient.CheckToken(token)
	if err != nil {
		return nil, err
	}
	// test token
	if tokenInfo.Scope != this.Scope {
		return nil, fault.ExceptionUnauthorized("Violation Access By Scope")
	}

	//if tokenInfo.Scope == "qrapi" {
	//	tokenInfo.ClientID = "4598959951692349"
	//}
	// // test time expired
	return
}

func (this *AuthFilter) extractBearerAuthToken(ctx *gin.Context) string {
	authValue := ctx.GetHeader("Authorization")
	bearerEndPrefix := len("Bearer ")
	if len(authValue) <= bearerEndPrefix {
		return ""
	}

	return authValue[bearerEndPrefix:]
}

func (this *AuthFilter) HasError(err error, ctx *gin.Context) bool {
	if err != nil {
		HandleError(err, ctx, this.Log)
	}
	return ctx.IsAborted()
}

func (this AuthFilter) SkipFilterRules(context *gin.Context) bool {
	if len(this.ExcludePath) > 0 {
		for _, v := range this.ExcludePath {
			if strings.Compare(v, context.FullPath()) == 0 {
				this.Log.RLog.Debug("skip auth filter for: %s ", context.FullPath())
				return true
			}
		}
	}
	this.Log.RLog.Info(fmt.Sprintf("Enable auth filter for: %s ", context.FullPath()))
	return false
}

func (this AuthFilter) SkipMpanCheckRules(context *gin.Context) bool {
	if len(this.SkipCheckMpanPath) > 0 {
		for _, v := range this.SkipCheckMpanPath {
			if strings.Compare(v, context.FullPath()) == 0 {
				this.Log.RLog.Debug("skip check mpan for: %s ", context.FullPath())
				return true
			}
		}
	}
	this.Log.RLog.Debug("Enable check mpan for: %s ", context.FullPath())
	return false
}

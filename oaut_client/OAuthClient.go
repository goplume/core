package oaut_client

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/goplume/core/cache"
	"github.com/goplume/core/fault"
	utils2 "github.com/goplume/core/utils"
	"github.com/goplume/core/utils/logger"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type OAuthClient struct {
	Scopes         map[string]ClientSecret
	GrantType      string
	CheckTokenURL  string
	GetTokenURL    string
	RemoveTokenURL string
	AuthEnabled    bool
	Log            *logger.Logger
	Cache          *cache.Cache
	HttpClient     *resty.Client
}

type ClientSecret struct {
	ClientID     string
	ClientSecret string
	Scope        string
	GrantType    string
	Enable       string
}

//  {
//    "access_token": "MZZCRLXYORO38HWPV7NUKW",
//    "expires_in": "604800",
//    "refresh_token": "",
//    "scope": "payment",
//    "token_type": "Bearer"
//  }
type OAuthToken struct {
	Access_token  string `example:"5COKJFUJPEQN4X0RY1R7AQ"`
	Expires_in    string `example:"604800"`
	Scope         string `example:"qrapi"`
	Token_type    string `example:"Bearer"`
	Refresh_token string `example:"MZZCRLXYORO38HWPV7NUKW"`
}

// {
//    "ClientID": "test",
//    "PublicID": "",
//    "RedirectURI": "",
//    "Scope": "webapi",
//    "Access": "1V1FUUGFO9GJZQOZUAMCFA",
//    "AccessCreateAt": "2019-06-04T16:35:10.688420232+06:00",
//    "AccessExpiresIn": 7200000000000,
//    "Refresh": "",
//    "RefreshCreateAt": "0001-01-01T00:00:00Z",
//    "RefreshExpiresIn": 0
// }
//
type TokenInfo struct {
	ClientID         string
	PublicID         string
	Scope            string
	Access           string
	AccessExpiresIn  time.Duration
	Role             string
	RedirectURI      string
	AccessCreateAt   string
	Refresh          string
	RefreshCreateAt  string
	RefreshExpiresIn uint64
}

type ParamsTokenInfo struct {
	ClientID         string
	PublicID         string
	Scope            string
	Access           string
	AccessExpiresIn  time.Duration
	Role             string
	RedirectURI      string
	AccessCreateAt   string
	Refresh          string
	RefreshCreateAt  string
	RefreshExpiresIn uint64
	Params           struct {
		Amount    string
		Curency   string
		InvoiceID string
		Mpan      string
		Msisdn    string
		Terminal  string
	}
}

// {
// "ClientID": "testqrapiclient",
// "PublicID": "",
// "RedirectURI": "",
// "Scope": "qrapi",
// "Access": "QVYQADQPPVQMZGKDG1ETAW",
// "AccessCreateAt": "2019-07-23T13:15:18.32982433+06:00",
// "AccessExpiresIn": 3600000000000,
// "Refresh": "K6PWVXV3WXEBGGPSOW14HQ",
// "RefreshCreateAt": "2019-07-23T13:15:18.32982433+06:00",
// "RefreshExpiresIn": 14400000000000,
// "Params": {
//     "amount": "",
//     "curency": "",
//     "invoiceID": "",
//     "mpan": "4598959951692349",
//     "msisdn": "77775477557",
//     "terminal": ""
// }
// }

/*
func NewOAuthClient(
    clientID string,
    clientSecret string,
    checkTokenURL string,
    getTokenURL string,
    removeTokenURL string,
) *OAuthClient {
    return &OAuthClient{
        ClientID:       clientID,
        ClientSecret:   clientSecret,
        CheckTokenURL:  checkTokenURL,
        GetTokenURL:    getTokenURL,
        RemoveTokenURL: removeTokenURL,
    }
}
*/

func (this *OAuthClient) InvalidateToken(
	rlog logrus.FieldLogger,
	scope string,
	mdc *utils2.MDC,
) {
	if this.Cache != nil {
		this.Cache.Evict(rlog, scope)
	}
}

func (this *OAuthClient) GetAuthorizationHeader(
	rlog logrus.FieldLogger,
	scope string,
	mdc *utils2.MDC,
) (authorizationHeader string, err error) {
	token, cachedToken, err := this.GetToken(rlog, scope)
	if err != nil {
		return "", err
	}

	if mdc != nil {
		mdc.Ψ("token_"+scope+"_is_cached", strconv.FormatBool(cachedToken))
	}
	//authorizationHeader = token.Token_type + " " + token.Access_token
	authorizationHeader = token.Access_token
	return authorizationHeader, nil
}

/*
func (this *OAuthClient) GetTokenScope() (token OAuthToken, err error) {
    return this.GetToken(this.Scope)
}
*/

// Retrieve  OAuthToken
// Lookup token in cache by scope name
// If not found then do auth
func (this *OAuthClient) GetToken(
	rlog logrus.FieldLogger,
	scope string,
) (
	tokenInfo *OAuthToken, fromCache bool, err error,
) {
	if this.Cache == nil {
		return nil, false, fault.ExceptionIllegalState("Cache not setting")
	}

	tokenInfo = &OAuthToken{}
	foundToken, err := this.Cache.Get(rlog, scope, tokenInfo)
	if err != nil {
		return nil, false, err
	}

	if foundToken {
		return tokenInfo, true, nil
	}

	tokenInfo, err = this.CreateToken(rlog, scope)
	if err != nil {
		return nil, false, err
	}

	if this.Cache != nil && this.Cache.Enabled {
		this.Cache.PutEntity(rlog, scope, tokenInfo)
	}

	return tokenInfo, false, nil

}

func (this *OAuthClient) CreateToken(
	rlog logrus.FieldLogger,
	scope string,
) (
	token *OAuthToken, err error,
) {

	authToken, err := this.CreateMToken(rlog, OAuthParams{
		Grant_type:    this.Scopes[scope].GrantType,
		Scope:         this.Scopes[scope].Scope,
		Client_id:     this.Scopes[scope].ClientID,
		Client_secret: this.Scopes[scope].ClientSecret,
	})

	return authToken, err
}

type OAuthParams struct {
	Grant_type    string
	Scope         string
	Client_id     string
	Client_secret string
	Extention     map[string]string
}

func validateOAuthParams(request OAuthParams) fault.TypedError {
	if len(request.Grant_type) == 0 {
		validatorError := fault.NewValidatorError("OAuth request Grant_type parameter empty")
		return validatorError
	}
	if len(request.Scope) == 0 {
		validatorError := fault.NewValidatorError("OAuth request Scope parameter empty")
		return validatorError
	}
	if len(request.Client_id) == 0 {
		validatorError := fault.NewValidatorError("OAuth request Client_id parameter empty")
		return validatorError
	}
	if len(request.Client_secret) == 0 {
		validatorError := fault.NewValidatorError("OAuth request Client_secret parameter empty")
		return validatorError
	}
	return nil
}

type ErrorCodeResponse struct {
	Code    int64
	Message string
}

func (this *OAuthClient) CreateMToken(
	rlog logrus.FieldLogger,
	oauthParams OAuthParams,
) (
	token *OAuthToken, err fault.TypedError,
) {
	if this.AuthEnabled == false {
		token = &OAuthToken{}
		token.Access_token = "OAuth disabled"
		rlog.Info("OAuth disabled")
		return token, nil
	}

	err = validateOAuthParams(oauthParams)
	if err != nil {
		return nil, err
	}

	oauthParamsForm := map[string]string{
		"grant_type":    oauthParams.Grant_type,
		"scope":         oauthParams.Scope,
		"client_id":     oauthParams.Client_id,
		"client_secret": oauthParams.Client_secret,
	}

	// append additional params
	if oauthParams.Extention != nil {
		for k, v := range oauthParams.Extention {
			oauthParamsForm[k] = v
		}
	}

	body := this.formdataToBody(oauthParamsForm)

	this.HttpClient.SetLogger(rlog)
	response, erri := this.HttpClient.R().
		EnableTrace().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		//SetHeader("Content-Type", "multipart/form-data").
		SetBody(body).
        Post(this.GetTokenURL)

    //response, erri := this.HttpClient.R().
    //    EnableTrace().
    //    SetHeader("Content-Type", "application/x-www-form-urlencoded").
		//SetFormData(oauthParamsForm).
		//Post(this.GetTokenURL)

	//rlog.Debug(response.Request.TraceInfo())

	if erri != nil {
		err := fault.ExceptionInternalErrorE(erri)
		return nil, err
	}

	if response.StatusCode() != http.StatusOK {
		errResponse := ErrorCodeResponse{}
		errResponseErrParse := json.Unmarshal(response.Body(), &errResponse)
		if errResponseErrParse != nil {
			rlog.Error(errResponseErrParse)
		}
		if errResponseErrParse == nil {
			if errResponse.Code == 399 {
				return nil, fault.ExceptionUnauthorized("")
			}
		}
		err := fault.IntegratrionExceptionClientErrorR(response)
		//err = errors.New(http.StatusText(response.StatusCode()))
		return nil, err
	}

	erri = json.Unmarshal(response.Body(), &token)
	if erri != nil {
		err := fault.ExceptionInternalErrorE(erri)
		return nil, err
	}

	return token, nil
}

func (this *OAuthClient) formdataToBody(data map[string]string) string {
	if data == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := data[k]
		keyEscaped := url.QueryEscape(k)

		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(keyEscaped)
		buf.WriteByte('=')

		if strings.Compare(k, "client_secret") != 0 {
			vs = url.QueryEscape(vs)
		}

		buf.WriteString(vs)
	}
	return buf.String()

}

func (this *OAuthClient) CheckToken(
	token string,
) (
	*ParamsTokenInfo, fault.TypedError,
) {
	this.Log.RLog.Info("Check token")
	if this.AuthEnabled == false {
		this.Log.RLog.Info("OAuth disabled")
		return nil, nil
	}

	response, rerr := this.HttpClient.R().
		SetAuthToken(token).
		Get(this.CheckTokenURL)

	if rerr != nil {
		return nil, fault.ExceptionInternalErrorE(rerr)
	}

	if response.StatusCode() == http.StatusUnauthorized {
		return nil, fault.ExceptionUnauthorized("")
	}

	if response.StatusCode() != http.StatusOK {
		return nil, fault.IntegrationExceptionServerErrorR(
			"OAuth fail", response)
	}

	//  status if OK
	tokenInfo := ParamsTokenInfo{}
	rerr = json.Unmarshal(response.Body(), &tokenInfo)
	if rerr != nil {
		return nil, fault.ExceptionInternalErrorE(rerr)
	}

	return &tokenInfo, nil
}
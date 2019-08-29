package rest_client

import (
	"github.com/goplume/core/oaut_client"
	utils2 "github.com/goplume/core/utils"
	"github.com/goplume/core/utils/logger"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type RestClient struct {
	OAuthClient *oaut_client.OAuthClient
	HttpClient  *resty.Client
	Log         *logger.Logger
	AuthScope   string
}

//this.OAuthClient.GetToken()

func (this *RestClient) Get(
	rlog logrus.FieldLogger,
	url string,
	pathParams map[string]string,
	body interface{},
) (
	*resty.Response, error,
) {
	return this.Get_(rlog, url, pathParams, body, nil)
}

func (this *RestClient) Get_(
	rlog logrus.FieldLogger,
	url string,
	pathParams map[string]string,
	body interface{},
	mdc *utils2.MDC,
) (
	*resty.Response, error,
) {
	return this.Request(rlog, "GET", url, pathParams, body, mdc)
}

func (this *RestClient) Post(
	rlog logrus.FieldLogger,
	url string,
	pathParams map[string]string,
	body interface{},
) (
	*resty.Response, error,
) {
	return this.Post_(rlog, url, pathParams, body, nil)
}

func (this *RestClient) Post_(
	rlog logrus.FieldLogger,
	url string,
	pathParams map[string]string,
	body interface{},
	mdc *utils2.MDC,
) (
	*resty.Response, error,
) {
	return this.Request(rlog, "POST", url, pathParams, body, mdc)
}

func (this *RestClient) Request(
	rlog logrus.FieldLogger,
	method string,
	url string,
	pathParams map[string]string,
	body interface{},
	mdc *utils2.MDC,
) (
	*resty.Response, error,
) {
	response, err := this.doRequest(rlog, method, url, pathParams, body, mdc)
	if err != nil {
		return nil, err
	}
	if response.StatusCode() == 401 {
		// repeat request with clear old token from cache
		this.OAuthClient.InvalidateToken(rlog, this.AuthScope, mdc)
		response, err = this.doRequest(rlog, method, url, pathParams, body, mdc)
		if err != nil {
			return nil, err
		}
	}
	return response, err
}

func (this *RestClient) doRequest(
	rlog logrus.FieldLogger,
	method string,
	url string,
	pathParams map[string]string,
	body interface{},
	mdc *utils2.MDC,
) (
	*resty.Response, error,
) {
	this.HttpClient.SetLogger(rlog)
	request := this.HttpClient.R()
	if this.RequiredAuthorization() {
		token, err := this.OAuthClient.GetAuthorizationHeader(rlog, this.AuthScope, mdc)
		if err != nil {
			return nil, err
		}
		request.SetAuthToken(token)
	}
	request.SetPathParams(pathParams)
	request.SetBody(body)
	resp, err := request.Execute(method, url)
	return resp, err
}

func (this *RestClient) RequiredAuthorization() bool {
	return this.OAuthClient != nil
}

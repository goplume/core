package health_checker

import (
    "encoding/json"
    "github.com/goplume/core/rest_client"
    "github.com/goplume/core/health"
    "github.com/go-resty/resty/v2"
    "net/http"
)

type Checker struct {
	RestClient   *rest_client.RestClient
	GoodHttpCode int
}

func NewRestClientChecker(RestClient *rest_client.RestClient) Checker {
	return Checker{
		RestClient:   RestClient,
		GoodHttpCode: http.StatusOK,
	}
}

func NewRestClientCheckerCode(RestClient *rest_client.RestClient, httpCode int) Checker {
	return Checker{
		RestClient:   RestClient,
		GoodHttpCode: httpCode,
	}
}

type helath struct {
	data string `json_off:"" `
}

func (this Checker) Check() (healthData health.Health) {

	healthData = health.NewHealth()

	if this.RestClient == nil {
		healthData.Up()
		healthData.AddInfo("err", "Not configured MerchantServiceRestClient or disabled")
		return
	}

	if this.RestClient.HttpClient == nil {
		healthData.Down()
		healthData.AddInfo("err", "Not configured MerchantServiceRestClient.HttpClient")
		return
	}

	if this.RestClient != nil {
		healthData.AddInfo("url", this.RestClient.HttpClient.HostURL)
	}

	health.TelnetCheck("", this.RestClient.HttpClient.HostURL, &healthData)

	data, resp, err := this.Health()

	if err != nil {
		healthData.Down()
		healthData.AddInfo("err", err.Error())
	}

	if data != "" {
		var f interface{}
		err = json.Unmarshal([]byte(data), &f)
		if err != nil {
			healthData.AddInfo("data", data)
		} else {
			healthData.AddInfo("data", f)
		}
	}

	if resp != nil {
		if resp.RawResponse != nil {
			healthData.AddInfo("http-status-code", resp.RawResponse.StatusCode)
			if resp.RawResponse.StatusCode == this.GoodHttpCode {
				healthData.Up()
			}
		}
		if resp.Request != nil {
			healthData.AddInfo("URL", resp.Request.URL)
		}
	}

	return healthData
}

func (this *Checker) Health() (data string, resp *resty.Response, err error) {

	resp, err = this.RestClient.HttpClient.R().Get("/health")

	json.Unmarshal(resp.Body(), &data)
	//return data, resp, err;
	return string(resp.Body()), resp, err
}

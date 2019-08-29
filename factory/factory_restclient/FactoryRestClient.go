package factory_restclient

import (
    "crypto/tls"
    "fmt"
    "github.com/goplume/core/rest_client"
    "github.com/goplume/core/configuration"
    "github.com/goplume/core/oaut_client"
    utils2 "github.com/goplume/core/utils"
    "github.com/goplume/core/utils/logger"
    "github.com/go-resty/resty/v2"
    "time"
)

type FactoryRestClient struct {
	Log *logger.Logger
}

func (this *FactoryRestClient) InitFactory() {
}

func (this *FactoryRestClient) CreateRestClient(
	serviceName string,
	clientName string,
) *rest_client.RestClient {
	return this.CreateOAuthRestClient(serviceName, clientName, nil)
}

func (this *FactoryRestClient) CreateOAuthRestClient(
	serviceName string,
	clientName string,
	oauthClient *oaut_client.OAuthClient,
) *rest_client.RestClient {
	config := configuration.NewServiceConfiguration(serviceName, "restclient."+clientName+".%s", this.Log)
	if this.Log != nil && this.Log.RLog != nil {
		this.Log.RLog.Info("Read configuration from context " + config.ConfigurationContext)
	}

	serviceEnable := config.GetBoolD("enable", false)
	if serviceEnable == false {
		if this.Log != nil && this.Log.RLog != nil {
			this.Log.RLog.Warn(fmt.Sprintf("TransactionsService %s disabled", serviceName))
		}
		return nil
	}

	serviceUrl := config.GetString("url")
	if len(serviceUrl) == 0 {
		return nil
	}
	timeout := config.GetInt("timeout-sec")
	debugEnable := config.GetBool("debug-enable")
	logEnable := config.GetBool("log-enable")
	oauthScope := config.GetString("oauth.scope")
	print(config.GetString("tls.insecure-skip-verify"))
	var InsecureSkipVerify utils2.Bool
	if config.IsSet("tls.insecure-skip-verify") {
		InsecureSkipVerify.Set(config.GetBool("tls.insecure-skip-verify"))
		if this.Log != nil && this.Log.RLog != nil {
			this.Log.RLog.Info(fmt.Sprintf("Read tls.insecure-skip-verify as %s", InsecureSkipVerify.Get()))
		}
	}
	// todo restore
	//logPath := config.GetString("log-file")
	// "/home/x/tmp/merchantapi/adapter/logs/go-resty.log"

	client := resty.New()
	//// Unique settings at Client level
	////--------------------------------
	//// Enable debug mode
	client.SetDebug(debugEnable)

	//// Using you custom log writer
	if logEnable {
		//client.SetLogger(this.Log.RLog)
		// todo restore
		//logFile, _ := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		// todo restore
		//client.SetLogger(*this.Log.RLog.Out())
		//logFile, err := os.OpenFile("service.log", os.O_APPEND|os.O_WRONLY, 0666)
		//if err == nil {
		//    client.SetLogger(logFile)
		//} else {
		//    log.Info("Failed to log to file, using default stderr")
		//}
	}

	//// Assign Client TLSClientConfig
	//// One can set custom root-certificate. Refer: http://golang.org/pkg/crypto/tls/#example_Dial
	// client.SetTLSClientConfig(&tls.Config{ RootCAs: roots })

	//// or One can disable security check (https)
	if InsecureSkipVerify.IsSet() {
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: InsecureSkipVerify.Get()})
	}

	//// Set client timeout as per your need
	if timeout != 0 {
		client.SetTimeout(time.Duration(timeout) * time.Second)
	}

	//// You can override all below settings and options at request level if you want to
	////--------------------------------------------------------------------------------
	//// Host URL for all request. So you can use relative URL in the request
	client.SetHostURL(serviceUrl)

	//// Headers for all request
	client.SetHeader("Accept", "application/json")
	client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "mvisa-merchant", //
	})

	client.RemoveProxy()

	//// Cookies for all request
	/*    client.SetCookie(&http.Cookie{
	          Name:     "go-resty",
	          Value:    "This is cookie value",
	          Path:     "/",
	          Domain:   "sample.com",
	          MaxAge:   36000,
	          HttpOnly: true,
	          Secure:   false,
	      })
	*///client.SetCookies(cookies) ????

	//// URL query parameters for all request
	//client.SetQueryParam("user_id", "00001")
	//client.SetQueryParams(map[string]string{ // sample of those who use this manner
	//    "api_key": "api-key-here",
	//    "api_secert": "api-secert",
	//})
	//client.R().SetQueryString("productId=232&template=fresh-sample&cat=resty&source=google&kw=buy a lot more")

	//// Form data for all request. Typically used with POST and PUT
	//client.SetFormData(map[string]string{
	//    "access_token": "BC594900-518B-4F7E-AC75-BD37F019E08F",
	//})

	//// Basic Auth for all request
	//client.SetBasicAuth("myuser", "mypass")

	//// Bearer Auth Token for all request
	//client.SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F")

	//// Enabling Content length value for all request
	client.SetContentLength(true)

	//// Registering global Error object structure for JSON/XML request
	//client.SetError(&Error{})    // or resty.SetError(Error{})

	restClient := &rest_client.RestClient{
		HttpClient: client,
		Log:        this.Log,
	}

	if len(oauthScope) > 0 {
		restClient.OAuthClient = oauthClient
		restClient.AuthScope = oauthScope
	}
	return restClient
}

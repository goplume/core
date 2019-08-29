package oauth_checker

import (
	"github.com/goplume/core/health"
	"github.com/goplume/core/oaut_client"
	"github.com/goplume/core/utils/logger"
	"github.com/sirupsen/logrus"
)

// Checker is a checker that check a given URL
type Checker struct {
	OAuthClient *oaut_client.OAuthClient
	Log         *logger.Logger
}

// NewChecker returns a new url.Checker with the given URL
func NewChecker(
	log *logger.Logger,
	OAuthClient *oaut_client.OAuthClient,
) *Checker {
	return &Checker{
		OAuthClient: OAuthClient,
		Log:         log,
	}
}

func (this *Checker) Check() health.Health {

	healthInfo := health.NewHealth()
	healthInfo.Up()
	healthInfo.AddInfo("get-token.url", this.OAuthClient.GetTokenURL)
	healthInfo.AddInfo("check-token.url", this.OAuthClient.CheckTokenURL)
	healthInfo.AddInfo("remove-token.url", this.OAuthClient.RemoveTokenURL)

	health.TelnetCheck("get-token.url.", this.OAuthClient.GetTokenURL, &healthInfo)
	health.TelnetCheck("check-token.url.", this.OAuthClient.CheckTokenURL, &healthInfo)
	health.TelnetCheck("remove-token.url.", this.OAuthClient.RemoveTokenURL, &healthInfo)

	for k, v := range this.OAuthClient.Scopes {
		// validate create token
		healthInfo.AddInfo("scope["+k+"].grand-type", v.GrantType)
		healthInfo.AddInfo("scope["+k+"].cliendId", v.ClientID)
		createdNewToken := this.validateCreateToken(this.Log.RLog, k, healthInfo)

		if createdNewToken == nil {
			continue
		}

		// validate check token
		this.validateCheckToken(createdNewToken, healthInfo, k)

		// validate get token from cache
		if this.OAuthClient.Cache != nil {
			this.validateCachedTokent(this.Log.RLog, v, healthInfo, k)
		}
	}
	return healthInfo
}

func (this *Checker) validateCachedTokent(
	rlog logrus.FieldLogger,
	v oaut_client.ClientSecret,
	health health.Health,
	k string,
) {
	createdNewToken, tokenCached, err := this.OAuthClient.GetToken(rlog, v.Scope)
	if err != nil {
		health.Down()
		health.AddInfo("scope["+k+"].cache-token.err", err.Error())
	} else {
		if createdNewToken != nil {
			health.AddInfo("scope["+k+"].cache-token.status", "successfuly")
			health.AddInfo("scope["+k+"].cache-token.from-cache", tokenCached)
		}
		if createdNewToken == nil {
			health.AddInfo("scope["+k+"].create-token.err", "token empty")
			health.AddInfo("scope["+k+"].cache-token.status", "failed")
		}
	}
}

func (this *Checker) validateCheckToken(
	createdNewToken *oaut_client.OAuthToken, health health.Health, k string) {
	tokenInfo, err := this.OAuthClient.CheckToken(createdNewToken.Access_token)
	if err != nil {
		health.Down()
		health.AddInfo("scope["+k+"].check-token.err", err.Error())
		health.AddInfo("scope["+k+"].check-token.status", "failed")
	} else {
		if tokenInfo != nil && len(tokenInfo.Access) > 0 {
			health.AddInfo("scope["+k+"].check-token.status", "authorized")
		} else {
			health.Down()
			health.AddInfo("scope["+k+"].create-token.err", "infotoken empty")
			health.AddInfo("scope["+k+"].check-token.status", "failed")
		}
	}
}

func (this *Checker) validateCreateToken(
	rlog logrus.FieldLogger,
	k string, health health.Health,
) (*oaut_client.OAuthToken) {
	createdNewToken, err := this.OAuthClient.CreateToken(rlog, k)
	if err != nil {
		health.Down()
		health.AddInfo("scope["+k+"].create-token.err", err.Error())
		health.AddInfo("scope["+k+"].create-token.status", "failed")
	} else {
		if createdNewToken != nil {
			health.AddInfo("scope["+k+"].create-token.status", "authorized")
		} else {
			health.Down()
			health.AddInfo("scope["+k+"].create-token.err", "token empty")
			health.AddInfo("scope["+k+"].create-token.status", "failed")
		}
	}
	return createdNewToken
}

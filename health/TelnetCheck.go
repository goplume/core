package health

import (
	"fmt"
	"github.com/reiver/go-telnet"
	"net"
	"net/url"
)

func TelnetCheck(prefix string, checkedUrl string, health *Health) {
	telnetHost, err := calcTelnetHost(prefix, checkedUrl, health)
	if err != nil {
		return
	}
	_, err = telnet.DialTo(telnetHost)
	health.AddInfo(prefix+"telnet-host", telnetHost)
	if err != nil {
		health.AddInfo(prefix+"telnet-ping", err.Error())
		health.Down()
	} else {
		health.AddInfo(prefix+"telnet-ping", "success")
	}
}

func calcTelnetHost(
	prefix string,
	checkedUrl string,
	health *Health,
) (
	telnetHost string,
	err error,
) {
	parse, err := url.Parse(checkedUrl)
	if err != nil {
		health.Down()
		health.AddInfo(prefix+"parse.url.error", err.Error())
		return "", err
	}
	host, port, err := net.SplitHostPort(parse.Host)
	if err != nil {
		port = ":80"
		if parse.Scheme == "https" {
			port = ":443"
		}
		host, port, err = net.SplitHostPort(parse.Host + port)
		if err != nil {
			health.AddInfo(prefix+"telnet-ping", err.Error())
			health.Down()
			return "", err
		}
	}

	telnetHost = fmt.Sprintf("%s:%s", host, port)
	return telnetHost, nil
}

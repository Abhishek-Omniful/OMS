package httpclient

import (
	nethttp "net/http"

	"github.com/Abhishek-Omniful/OMS/mycontext"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
)

var client *http.Client
var logger = log.DefaultLogger()

func InitHttpClient() {
	ctx := mycontext.GetContext()
	serviceName := config.GetString(ctx, "client.serviceName")
	baseURL := config.GetString(ctx, "client.baseURL")
	timeout := config.GetDuration(ctx, "http.timeout")
	maxIdleConns := config.GetInt(ctx, "client.maxIdleConns")
	maxIdleConnsPerHost := config.GetInt(ctx, "client.maxIdleConnsPerHost")

	transport := &nethttp.Transport{
		MaxIdleConns:        maxIdleConns,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
	}

	var err error
	client, err = http.NewHTTPClient(
		serviceName,
		baseURL,
		transport,
		http.WithTimeout(timeout),
	)
	if err != nil {
		ctx := mycontext.GetContext()
		logger.Errorf(i18n.Translate(ctx, "Failed to initialize HTTP client: %v"), err)
	}
	logger.Infof(i18n.Translate(ctx, "HTTP client initialized successfully"))
}

func GetHttpClient() *http.Client {
	return client
}

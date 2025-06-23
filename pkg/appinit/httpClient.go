package appinit

// import (
// 	nethttp "net/http"

// 	"github.com/Abhishek-Omniful/OMS/mycontext"
// 	"github.com/omniful/go_commons/config"
// 	"github.com/omniful/go_commons/http"
// )

// var client *http.Client

// func InitHttpClient() {
// 	// Initialize client with base URL
// 	ctx := mycontext.GetContext()
// 	serviceName := config.GetString(ctx, "client.serviceName")
// 	baseURL := config.GetString(ctx, "client.baseURL")
// 	timeout := config.GetDuration(ctx, "http.timeout")
// 	maxIdleConns := config.GetInt(ctx, "client.maxIdleConns")
// 	maxIdleConnsPerHost := config.GetInt(ctx, "client.maxIdleConnsPerHost")

// 	transport := &nethttp.Transport{
// 		MaxIdleConns:        maxIdleConns,
// 		MaxIdleConnsPerHost: maxIdleConnsPerHost,
// 	}

// 	client, err = http.NewHTTPClient(
// 		serviceName, // client service name
// 		baseURL,     // base URL
// 		transport,
// 		http.WithTimeout(timeout), // optional timeout
// 	)
// 	if err != nil {
// 		logger.Errorf("Failed to initialize HTTP client: %v", err)
// 	}
// }

// func GetHttpClient() *http.Client {
// 	return client
// }

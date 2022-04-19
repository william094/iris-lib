package iris_lib

import (
	"bytes"
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var client *http.Client
var proxyClient *http.Client

const (
	CONTENTTYPE_JSON     = "application/json;charset=utf8"
	CONTENTTYPE_FORM_URL = "application/x-www-form-urlencoded"
	TrackId              = "TrackId"
)

func init() {
	client = &http.Client{
		Timeout: 60 * time.Second,
	}
	dialSocksProxy, err := proxy.SOCKS5("tcp", "127.0.0.1:6005", nil, proxy.Direct)
	if err != nil {
		SystemLogger.Error("Error connecting to proxy", zap.Error(err))
	}
	proxyClient = &http.Client{Transport: &http.Transport{Dial: dialSocksProxy.Dial}}
}

func Get(url string, ctx context.Context) ([]byte, error) {
	logger := ctx.Value("log").(*zap.Logger)
	if logger != nil {
		logger = WithContext(ctx)
	} else {
		logger = SystemLogger
	}
	logger.Info("GET Request", zap.String("url", url))
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	if ctx != nil && ctx.Value(TrackId) != nil {
		req.Header.Set(TrackId, ctx.Value(TrackId).(string))
	}
	resp, err := client.Do(req.Clone(ctx))
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	logger.Info("接口响应:", zap.String("url", url), zap.String("response", string(response)))
	if err != nil {
		logger.Error("请求异常", zap.String("url", url), zap.Any("response", resp), zap.Error(err))
		return nil, err
	}
	return response, err
}

func PostByJson(url string, params interface{}, ctx context.Context) ([]byte, error) {
	logger := ctx.Value("log").(*zap.Logger)
	if logger != nil {
		logger = WithContext(ctx)
	} else {
		logger = SystemLogger
	}
	logger.Info("POST Request", zap.String("url", url), zap.Any("params", params))
	body, err := json.Marshal(params)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if ctx != nil && ctx.Value(TrackId) != nil {
		req.Header.Set(TrackId, ctx.Value(TrackId).(string))
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		logger.Error("请求异常", zap.String("url", url), zap.Error(err))
		return nil, err
	}
	response, err := ioutil.ReadAll(resp.Body)
	logger.Info("接口响应:", zap.String("url", url), zap.String("response", string(response)))
	if err != nil {
		logger.Error("请求异常", zap.String("url", url), zap.Any("response", resp), zap.Error(err))
		return nil, err
	}
	return response, err
}

func PostByFormUrl(URL string, params map[string]string, ctx context.Context) ([]byte, error) {
	logger := ctx.Value("log").(*zap.Logger)
	if logger != nil {
		logger = WithContext(ctx)
	} else {
		logger = SystemLogger
	}
	logger.Info("POST Request", zap.String("url", URL), zap.Any("params", params))
	data := url.Values{}
	for key, value := range params {
		data.Set(key, value)
	}
	resp, err := client.PostForm(URL, data)
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	logger.Info("接口响应:", zap.String("url", URL), zap.String("response", string(response)))
	if err != nil {
		logger.Error("请求异常", zap.String("url", URL), zap.Any("response", resp), zap.Error(err))
		return nil, err
	}
	return response, err
}

func GetByProxy(url string, ctx context.Context) ([]byte, error) {
	logger := ctx.Value("log").(*zap.Logger)
	if logger != nil {
		logger = WithContext(ctx)
	} else {
		logger = SystemLogger
	}
	logger.Info("GET Request", zap.String("url", url))
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	if ctx != nil && ctx.Value(TrackId) != nil {
		req.Header.Set(TrackId, ctx.Value(TrackId).(string))
	}
	resp, err := proxyClient.Do(req.Clone(ctx))
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	logger.Info("接口响应:", zap.String("url", url), zap.String("response", string(response)))
	if err != nil {
		logger.Error("请求异常", zap.String("url", url), zap.Error(err))
		return nil, err
	}
	return response, err
}

func PostProxyByJson(url string, params interface{}, ctx context.Context) ([]byte, error) {
	logger := ctx.Value("log").(*zap.Logger)
	if logger != nil {
		logger = WithContext(ctx)
	} else {
		logger = SystemLogger
	}
	logger.Info("POST Request", zap.String("url", url), zap.Any("params", params))
	body, err := json.Marshal(params)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if ctx != nil && ctx.Value(TrackId) != nil {
		req.Header.Set(TrackId, ctx.Value(TrackId).(string))
	}
	resp, err := proxyClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		logger.Error("请求异常", zap.String("url", url), zap.Error(err))
		return nil, err
	}
	response, err := ioutil.ReadAll(resp.Body)
	logger.Info("接口响应:", zap.String("url", url), zap.String("response", string(response)))
	if err != nil {
		logger.Error("请求异常", zap.String("url", url), zap.Error(err))
		return nil, err
	}
	return response, err
}

func PostProxyByFormUrl(URL string, params map[string]string, ctx context.Context) ([]byte, error) {
	logger := ctx.Value("log").(*zap.Logger)
	if logger != nil {
		logger = WithContext(ctx)
	} else {
		logger = SystemLogger
	}
	logger.Info("POST Request", zap.String("url", URL), zap.Any("params", params))
	data := url.Values{}
	for key, value := range params {
		data.Set(key, value)
	}
	resp, err := proxyClient.PostForm(URL, data)
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	logger.Info("接口响应:", zap.String("url", URL), zap.String("response", string(response)))
	if err != nil {
		logger.Error("请求异常", zap.String("url", URL), zap.Error(err))
		return nil, err
	}
	return response, err
}

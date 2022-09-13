package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/william094/iris-lib/logx"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var client *http.Client

const (
	CONTENTTYPE_JSON     = "application/json;charset=utf8"
	CONTENTTYPE_FORM_URL = "application/x-www-form-urlencoded"
	TrackId              = "TrackId"
)

func init() {
	client = &http.Client{
		Timeout: 10 * time.Second,
	}
}

func Get(ctx context.Context, url string, params url.Values) ([]byte, error) {
	if params != nil {
		url = url + "?" + params.Encode()
	}
	logx.SystemLogger.Info("GET Request", zap.String("url", url))
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	if ctx != nil && ctx.Value(TrackId) != nil {
		req.Header.Set(TrackId, ctx.Value(TrackId).(string))
		req = req.Clone(ctx)
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		logx.SystemLogger.Error("Get Request Exception", zap.String("url", url), zap.Any("response", resp), zap.Error(err))
		return nil, err
	}
	response, err := ioutil.ReadAll(resp.Body)
	logx.SystemLogger.Info("Get Response:", zap.String("url", url), zap.String("response", string(response)))
	if err != nil {
		logx.SystemLogger.Error("Get Response Read Exception", zap.String("url", url), zap.Any("response", resp), zap.Error(err))
		return nil, err
	}
	return response, err
}

func PostByJson(ctx context.Context, url string, params interface{}) ([]byte, error) {
	logx.SystemLogger.Info("POST Request", zap.String("url", url), zap.Any("params", params))
	body, err := json.Marshal(params)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", CONTENTTYPE_JSON)
	if ctx != nil && ctx.Value(TrackId) != nil {
		req.Header.Set(TrackId, ctx.Value(TrackId).(string))
		req = req.Clone(ctx)
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		logx.SystemLogger.Error("POST Request Exception", zap.String("url", url), zap.Error(err))
		return nil, err
	}
	response, err := ioutil.ReadAll(resp.Body)
	logx.SystemLogger.Info("POST Response:", zap.String("url", url), zap.String("response", string(response)))
	if err != nil {
		logx.SystemLogger.Error("POST Response Read Exception", zap.String("url", url), zap.Any("response", resp), zap.Error(err))
		return nil, err
	}
	return response, err
}

func PostByFormUrl(ctx context.Context, URL string, params map[string]string) ([]byte, error) {
	logx.SystemLogger.Info("POST Request", zap.String("url", URL), zap.Any("params", params))
	data := url.Values{}
	for key, value := range params {
		data.Set(key, value)
	}
	resp, err := client.PostForm(URL, data)
	defer resp.Body.Close()
	if err != nil {
		logx.SystemLogger.Error("POST Request Exception", zap.String("url", URL), zap.Any("response", resp), zap.Error(err))
		return nil, err
	}
	response, err := ioutil.ReadAll(resp.Body)
	logx.SystemLogger.Info("POST Response:", zap.String("url", URL), zap.String("response", string(response)))
	if err != nil {
		logx.SystemLogger.Error("POST Response Read Exception", zap.String("url", URL), zap.Any("response", resp), zap.Error(err))
		return nil, err
	}
	return response, err
}

package http

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hercules_compiler/rce-executor/log"
	"io/ioutil"
	goHttp "net/http"
	"strings"

	"github.com/astaxie/beego/context"
)

//定义http方法
const (
	HttpMethodGet    = "get"
	HttpMethodPost   = "post"
	HttpMethodDelete = "delete"
)

//定义http头信息
type HttpHead map[string]string

//NewHttpHead 新建http头信息
func NewHttpHead() HttpHead {
	httpHead := make(map[string]string)
	return httpHead
}

//AddHttpHead 添加httphead
func AddHttpHead(head HttpHead, key, value string) HttpHead {
	if head == nil {
		head = NewHttpHead()
	}
	head[key] = value
	return head
}

//定义http参数 针对http get请求的封装
type HttpParam HttpHead

//NewHttpHead 新建http参数
func NewHttpParam() HttpParam {
	httpParam := make(map[string]string)
	return httpParam
}

//AddHttpParam 添加HttpParam
func AddHttpParam(param HttpParam, key, value string) HttpParam {
	if param == nil {
		param = NewHttpParam()
	}
	param[key] = value
	return param
}

//httpDo http请求数据
func httpDo(reqURL, method, param string, head HttpHead) (string, error) {
	var (
		err error           = nil
		req *goHttp.Request = nil
	)
	client := &goHttp.Client{}
	method = strings.ToLower(method)
	switch method {
	case HttpMethodGet:
		req, err = goHttp.NewRequest("GET", reqURL, nil)
		if err != nil {
			return "", err
		}
	case HttpMethodPost:
		req, err = goHttp.NewRequest("POST", reqURL, strings.NewReader(param))
		if err != nil {
			return "", err
		}
	case HttpMethodDelete:
		log.Debug("url----->", reqURL)
		req, err = goHttp.NewRequest("DELETE", reqURL, strings.NewReader(param))
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unsupport http method %s", method)
	}

	//设置http 头
	for key, value := range head {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	response := string(body)
	//只有返回200 才能正确返回
	if resp.StatusCode != goHttp.StatusOK {
		return "", errors.New(response)
	}

	log.Debug("response = ", response)
	return response, nil
}

//httpDo http请求获取流文件
func httpDoStream(reqURL, method, param string, head HttpHead, myreq *context.Response) (string, error) {
	var (
		err error           = nil
		req *goHttp.Request = nil
	)
	client := &goHttp.Client{}
	method = strings.ToLower(method)
	switch method {
	case HttpMethodGet:
		req, err = goHttp.NewRequest("GET", reqURL, nil)
		if err != nil {
			return "", err
		}
	case HttpMethodPost:
		req, err = goHttp.NewRequest("POST", reqURL, strings.NewReader(param))
		if err != nil {
			return "", err
		}
	case HttpMethodDelete:
		log.Debug("url----->", reqURL)
		req, err = goHttp.NewRequest("DELETE", reqURL, strings.NewReader(param))
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unsupport http method %s", method)
	}

	//设置http 头
	for key, value := range head {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	for k, v := range resp.Header {
		myreq.Header().Set(k, v[0])
	}
	myreq.Write(body)
	return "", nil
}

//HttpGet 提交http get 方法
func HttpGetBasic(reqURL string, params HttpParam, head HttpHead) (string, error) {
	urlAppend := ""
	index := 0
	for key, value := range params {
		if index == 0 {
			urlAppend += fmt.Sprintf("?%s=%s", key, value)
		} else {
			urlAppend += fmt.Sprintf("&%s=%s", key, value)
		}
		index++
	}
	reqURL += urlAppend
	log.Debug("reqURL=", reqURL)
	return httpDo(reqURL, HttpMethodGet, "", head)
}

//HttpGet 提交http get 方法
func HttpPostStream(reqURL string, param interface{}, head HttpHead, myres *context.Response) (string, error) {
	reqURL = strings.TrimSpace(reqURL)
	bResult, err := json.Marshal(param)
	if err != nil {
		return "", err
	}
	log.Debug("request url = ", reqURL)
	return httpDoStream(reqURL, HttpMethodPost, string(bResult), head, myres)
}

//HttpGet 提交http get 方法
func HttpGetOneResource(reqURL string, params HttpParam, head HttpHead) (string, error) {
	urlAppend := ""
	for _, value := range params {
		urlAppend += fmt.Sprintf("/%s", value)
	}
	reqURL += urlAppend
	return httpDo(reqURL, HttpMethodGet, "", head)
}

//HttpPost 提交http post 方法
func HttpPostOneResource(reqURL string, params HttpParam, head HttpHead) (string, error) {
	urlAppend := ""
	for _, value := range params {
		urlAppend += fmt.Sprintf("/%s", value)
	}
	reqURL += urlAppend
	return httpDo(reqURL, HttpMethodPost, "", head)
}

//HttpPost 提交http post方法
func HttpPost(reqURL, param string, head HttpHead) (string, error) {
	reqURL = strings.TrimSpace(reqURL)
	return httpDo(reqURL, HttpMethodPost, param, head)
}

//HttpPostDoJson 提交http post 方法 body 为json
func HttpPostDoJson(reqURL string, param interface{}, head HttpHead) (string, error) {
	reqURL = strings.TrimSpace(reqURL)
	bResult, err := json.Marshal(param)
	if err != nil {
		return "", err
	}
	log.Debug("request url = ", reqURL)
	return httpDo(reqURL, HttpMethodPost, string(bResult), head)
}

//HttpDelete 提交http delete 方法
func HttpDelete(reqURL, param string, head HttpHead) (string, error) {
	reqURL = strings.TrimSpace(reqURL)
	return httpDo(reqURL, HttpMethodDelete, param, head)
}

func GetBasicAuthHeadInfo(username, password string) (head string) {
	head = fmt.Sprintf("Authorization:Basic %s", base64.StdEncoding.EncodeToString([]byte(username+":"+password)))
	return
}

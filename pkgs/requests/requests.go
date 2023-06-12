package requests

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"tebot/pkgs/initfunc"
	"tebot/pkgs/logtool"

	"time"

	//"fmt"
	//"os"
	"github.com/go-resty/resty/v2"
)

const pioTimeout = 300
const pTimeOut = 100
const retry = 3

var Client *ClientStruct

func init() {
	initfunc.RegisterInitFunc(
		func() {
			Client = getclinet()
		},
	)
}

type ClientStruct struct {
	*resty.Client
}

func (c *ClientStruct) SetNotls() *ClientStruct {
	newclint := getclinet()
	newclint.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	return newclint
}

type RequestStruct struct {
	*resty.Request
}

func (r *RequestStruct) Notparse() *RequestStruct {
	r.SetDoNotParseResponse(true)
	return r
}

func (c *ClientStruct) R() *RequestStruct {
	return &RequestStruct{
		c.Client.R(),
	}
}

func getclinet() *ClientStruct {
	clientConfiguration := &ClientStruct{
		resty.New(),
	}
	clientConfiguration.
		SetTransport(&http.Transport{
			Proxy:             http.ProxyFromEnvironment,
			ForceAttemptHTTP2: true,
			Dial: TimeoutDialer(
				time.Duration(pTimeOut)*time.Second,    //建立连接的超时
				time.Duration(pioTimeout)*time.Second), //单次网络读写的超时
		}).
		//SetTimeout(time.Duration(pTimeOut) * time.Second)// 总的请求时间,必须大于Transport.Dial: TimeoutDialer(30 * time.Second, 1 * time.Minute)
		SetRetryCount(retry).
		SetRetryWaitTime(100 * time.Nanosecond).
		AddRetryCondition(
			func(response *resty.Response, err error) bool {
				return !response.IsSuccess() || err != nil
			},
		).
		OnAfterResponse(
			func(c *resty.Client, resp *resty.Response) error {
				// Now you have access to Client and current Response object
				// manipulate it as per your need
				if !resp.IsSuccess() {
					return errors.New("request failed,http code is " + resp.Status())

				}
				return nil // if its success otherwise return error
			}).
		SetLogger(logtool.SugLog)
	return clientConfiguration
}

func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		//conn, err := net.DialTimeout(netw, addr, cTimeout)
		d := net.Dialer{
			Timeout:   cTimeout,
			DualStack: true}
		conn, err := d.Dial(netw, addr)

		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}

func Parsebody_to_json(resp *resty.Response) map[string]interface{} {
	//if err != nil{
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	var v interface{}
	json.Unmarshal(resp.Body(), &v)
	return v.(map[string]interface{})
}

func Ecocde_json(v any) ([]byte, error) {
	e, err := json.Marshal(v)
	return e, err
}

//func tojson(resp *resty.Response,err error) interface{} {
//	if err != nil{
//		fmt.Println(err)
//		os.Exit(1)
//	}
//	fmt.Println(resp.RawBody())
//	var v interface{}
//	json.NewDecoder(resp.RawBody()).Decode(&v)
//	return v
//}

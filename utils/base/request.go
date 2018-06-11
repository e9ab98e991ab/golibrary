/**
 * @author : godfeer@aliyun.com
 * @date : 2018/6/11/011
 **/
//

package base

import (
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Request struct {
	Url string
}

func (reqs *Request) Fetch(userAgent string, headRequest bool) (*http.Response, error) {
	u, _ := url.Parse(reqs.Url)
	var reqType string
	// Prepare the request with the right user agent
	if headRequest {
		reqType = "HEAD"
	} else {
		reqType = "GET"
	}
	req, e := http.NewRequest(reqType, u.String(), nil)
	if e != nil {
		return nil, e
	}
	req.Header.Set("User-Agent", userAgent)
	DefaultClient := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(10 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*10)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}
	return DefaultClient.Do(req)
}

func (reqs *Request) FetchPost(userAgent string, params map[string]string) (*http.Response, error) {
	u, _ := url.Parse(reqs.Url)
	vv := url.Values{}
	for k, v := range params {
		vv.Add(k, v)
	}
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(vv.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	DefaultClient := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(10 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*10)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}
	return DefaultClient.Do(req)
}

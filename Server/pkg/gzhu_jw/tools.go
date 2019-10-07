package gzhu_jw

import (
	"github.com/astaxie/beego/logs"
	"io"
	"net/http"
	"time"
)

func (c *JWClient) doRequest(method, url string, header http.Header, body io.Reader) (*http.Response, error) {
	t1 := time.Now().UnixNano() / 1000000

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header[k] = v
	}
	resp, err := c.Client.Do(req)

	logs.Info(time.Now().UnixNano()/1000000-t1, "ms", url)
	return resp, err
}

/**
 * @File: cosfs
 * @Author: Shaw
 * @Date: 2020/5/9 12:49 AM
 * @Desc: 腾讯云cos存储

 */

package cosfs

import (
	"bytes"
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/spf13/viper"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
)

var CosCli *cos.Client

var (
	cosUrl     = "https://shaw-1256261760.cos.ap-guangzhou.myqcloud.com"
	piBasePath = "/gzhu-pi/upload/"
	secretID   = "COS_SECRETID"
	secretKey  = "COS_SECRETKEY"
)

func initCos() error {

	if viper.IsSet("cos.secret_id") {
		secretID = viper.GetString("cos.secret")
	} else {
		return fmt.Errorf("cos.secret_id not set")
	}
	if viper.IsSet("cos.secret_key") {
		secretKey = viper.GetString("cos.secret_key")
	} else {
		return fmt.Errorf("cos.secret_key not set")
	}

	u, _ := url.Parse(cosUrl)
	b := &cos.BaseURL{BucketURL: u}
	CosCli = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})
	return nil
}

func SaveToCos(data []byte, filename string) (url string, err error) {

	if CosCli == nil {
		err = initCos()
		if err != nil {
			logs.Error(err)
			return
		}
	}

	// 对象键（Key）是对象在存储桶中的唯一标识。
	// 例如，在对象的访问域名 `examplebucket-1250000000.cos.COS_REGION.myqcloud.com/test/objectPut.go` 中，对象键为 test/objectPut.go
	target := piBasePath + filename

	f := bytes.NewBuffer(data)
	_, err = CosCli.Object.Put(context.Background(), target, f, nil)
	if err != nil {
		logs.Error(err)
		return
	}
	url = cosUrl  + target
	return
}

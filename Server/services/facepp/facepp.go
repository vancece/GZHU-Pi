/**
 * @File: facepp
 * @Author: Shaw
 * @Date: 2020/5/9 1:43 AM
 * @Desc: face++ api

 */

package facepp

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/spf13/viper"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

type SegmentResp struct {
	ImageID      string `json:"image_id"`
	Result       string `json:"result"`
	BodyImage    string `json:"body_image"`
	RequestID    string `json:"request_id"`
	TimeUsed     int    `json:"time_used"`
	ErrorMessage string `json:"error_message"`
}

//face++人像抠图api HumanBody Segment
//https://console.faceplusplus.com.cn/documents/40608240
//传入原图数据，返回抠图后的数据
func HumanBodySegment(data []byte) (segmentData []byte, err error) {

	var apiKey, apiSecret string
	if viper.IsSet("facepp.api_key") {
		apiKey = viper.GetString("facepp.api_key")
	}
	if viper.IsSet("facepp.api_secret") {
		apiSecret = viper.GetString("facepp.api_secret")
	}

	if apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("facepp.api_key or facepp.api_secret not set")
	}

	url := "https://api-cn.faceplusplus.com/humanbodypp/v2/segment"

	var postData = map[string]string{
		"api_key":          apiKey,
		"api_secret":       apiSecret,
		"return_grayscale": "0",
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	//写入 form-data
	for k, v := range postData {
		err = w.WriteField(k, v)
		if err != nil {
			logs.Error(err)
			return
		}
	}

	//添加文件
	fw, err := w.CreateFormFile("image_file", fmt.Sprint(time.Now().Unix()))
	if err != nil {
		logs.Error(err)
		return
	}
	_, err = fw.Write(data)
	if err != nil {
		logs.Error(err)
		return
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	//发送请求
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		logs.Error(err)
		return
	}
	// Set the content type = multipart/form-data , this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		logs.Error(err)
		return
	}

	//解析响应
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logs.Error(err)
		return
	}
	var respStruct SegmentResp
	err = json.Unmarshal(resBody, &respStruct)
	if err != nil {
		logs.Error(err)
		return
	}
	if respStruct.ErrorMessage != "" || res.StatusCode != http.StatusOK {
		err = fmt.Errorf("statu_code %d %s", res.StatusCode, respStruct.ErrorMessage)
		logs.Error(err)
		return
	}
	//解析base64
	segmentData, err = base64.StdEncoding.DecodeString(respStruct.BodyImage)
	if err != nil {
		logs.Error(err)
		return
	}
	return
}

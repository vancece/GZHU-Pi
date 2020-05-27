/**
 * @File: cet
 * @Author: Shaw
 * @Date: 2020/5/27 11:19 AM
 * @Desc

 */

package routers

import (
	"GZHU-Pi/pkg/cet"
	"fmt"
	"net/http"
)

func GetCet(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	name := r.URL.Query().Get("name")
	captcha := r.URL.Query().Get("captcha")

	if id == "" || name == "" {
		err := fmt.Errorf("准考证/姓名不能为空")
		Response(w, r, nil, http.StatusUnauthorized, err.Error())
		return
	}

	client := cet.NewCetClient(id, name, captcha)

	//获取验证码
	if captcha == "" {
		imgUrl, err := client.GetCaptcha()
		if err != nil {
			Response(w, r, nil, http.StatusBadRequest, err.Error())
			return
		}
		basImg, err := imgUrlToBase64(imgUrl)
		if err != nil {
			Response(w, r, nil, http.StatusBadRequest, err.Error())
			return
		}
		Response(w, r, basImg, http.StatusOK, "request ok")
		return
	}

	//查询结果
	err := client.GetCetInfo()
	if err != nil {
		Response(w, r, nil, http.StatusBadRequest, err.Error())
		return
	}
	client.DelCache()

	Response(w, r, client, http.StatusOK, "request ok")
}

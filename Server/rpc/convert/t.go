/**
 * @File: t
 * @Author: Shaw
 * @Date: 2020/4/20 11:46 PM
 * @Desc

 */

package main

import (
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/rpc"
)

type Convert struct {
	Token        string `json:"token"`
	Body         []byte `json:"body"`
	ConvertType  string `json:"convert_type"`
	ConvertType1 string `json:"convert_type1"`
}

func main0() {

	rpcClient, err := rpc.DialHTTP("tcp", "localhost:7201")
	if err != nil {
		logs.Error(err)
		return
	}
	var reply []byte

	body, err := ioutil.ReadFile("/Users/Shaw/Desktop/保单模板/1.校责险-投保单_模板.xlsx")

	var req = Convert{
		Token:       "123456",
		Body:        body,
		ConvertType: "pdf",
	}

	err = rpcClient.Call("ConvertService.Convert", req, &reply)
	if err != nil {
		logs.Error(err)
		return
	}

	ioutil.WriteFile("a.pdf", reply, 0666)
}

/**
 * @File: t
 * @Author: Shaw
 * @Date: 2020/4/20 11:46 PM
 * @Desc

 */

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
)

// Convert 示例结构
type Convert struct {
	Token       string `json:"token"`
	Body        []byte `json:"body"`
	ConvertType string `json:"convert_type"`
}

func example() {

	addr := "localhost:6618"
	rpcClient, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	var reply []byte
	body, err := ioutil.ReadFile("/Users/Shaw/Desktop/保单模板/1.校责险-投保单_模板.xlsx")

	var req = Convert{
		Token:       "123456",
		Body:        body,
		ConvertType: "pdf",
	}

	fmt.Println("Ping rpc server ", addr)
	var pong string
	err = rpcClient.Call("ConvertService.Ping", req, &pong)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pong)

	err = rpcClient.Call("ConvertService.Convert", req, &reply)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("convert.pdf", reply, 0666)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("convert succeed")
}

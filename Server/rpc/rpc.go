package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/otiai10/gosseract"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {

	err := rpc.RegisterName("OcrService", new(OcrService))
	if err != nil {
		log.Fatal("Register error:", err)
	}
	//将Rpc绑定到HTTP协议上
	rpc.HandleHTTP()

	var port = ":7201"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	logs.Info("RPC Server listening on port", port)
	//启动http服务，处理连接请求
	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal("Error serving: ", err)
	}
}

type OcrService struct{}

func (p *OcrService) Capture(request []byte, reply *string) (err error) {

	client := gosseract.NewClient()
	defer client.Close()

	err = client.SetImageFromBytes(request)
	if err != nil {
		logs.Error(err)
		return
	}
	*reply, err = client.Text()
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Debug("识别结果：", *reply)

	return
}

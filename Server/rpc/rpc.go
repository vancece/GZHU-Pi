package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/otiai10/gosseract"
	"log"
	"net"
	"net/rpc"
)

func main() {

	err := rpc.RegisterName("OcrService", new(OcrService))
	if err != nil {
		log.Fatal("RegisterName error:", err)
	}
	var port = ":1234"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	logs.Info("RPC Server listening on port", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}
		rpc.ServeConn(conn)
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

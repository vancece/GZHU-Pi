package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

//全局实例
var ConvertCli = ConvertService{
	Token:      "123456",
	supported:  []string{"pdf"},
	delInfile:  true,
	delOutfile: true,
	outDir:     "files",
}

const port = "6618"

func main() {
	err := InitLogger("./log")
	if err != nil {
		log.Fatal(err)
	}

	err = rpc.RegisterName("ConvertService", &ConvertCli)
	if err != nil {
		log.Fatal("Register error:", err)
	}
	//将Rpc绑定到HTTP协议上
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":"+port)
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

type ConvertService struct {
	Token       string `json:"token"`        //秘钥，防止非法调用
	Body        []byte `json:"body"`         //文件数据
	ConvertType string `json:"convert_type"` //转为的目标类型，如pdf
	Converted   []byte `json:"converted"`    //转换后的数据

	Lock      sync.Mutex //libreoffice的soffice不支持多进程
	supported []string   //支持的ConvertType

	delInfile  bool //是否删除源文件
	delOutfile bool //是否删除转换后的文件

	outDir string //保存文件的路径
}

//读取body保存成文件 filename=md5
//调用系统命令 soffice --headless --convert-to {{ConvertType}} filename
func (p *ConvertService) Convert(request *ConvertService, reply *[]byte) (err error) {
	defer logs.Info("====== convert done ======")
	if request == nil {
		err = fmt.Errorf("call convert with  nil request argument")
		return
	}
	if request.Token != p.Token {
		err = fmt.Errorf("service auth failed with wrong token")
		return
	}
	if len(request.Body) == 0 {
		err = fmt.Errorf("empty request body")
		return
	}
	support := false
	request.ConvertType = strings.ToLower(request.ConvertType)
	for _, v := range p.supported {
		if v == request.ConvertType {
			support = true
		}
	}
	if !support {
		err = fmt.Errorf("convert type: %s not in supported list: %v",
			request.ConvertType, p.supported)
		logs.Error(err)
		return
	}

	if p.outDir == "" {
		p.outDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logs.Error(err)
			return
		}
	} else {
		err = os.MkdirAll(p.outDir, 0666)
		if err != nil {
			logs.Error(err)
			return
		}
	}
	p.outDir = strings.TrimRight(p.outDir, "/") + "/"
	infile, _, err := saveFileByMd5(request.Body, p.outDir)
	if err != nil {
		logs.Error(err)
		return
	}
	if p.delInfile {
		defer os.Remove(infile)
	}

	command := fmt.Sprintf("soffice --headless --convert-to %s %s --outdir %s",
		request.ConvertType, infile, strings.TrimRight(p.outDir, "/"))
	logs.Info(command)

	p.Lock.Lock()
	stdout, _, err := shellOut(command)
	p.Lock.Unlock()
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Info("stdout below: \n", stdout)

	//匹配输出文件路径 示例 "convert /home/example.xlsx -> /home/outdir/example.pdf using filter : calc_pdf_Export"
	reg := regexp.MustCompile(`convert [\s\S]* -> (.*) using filter`)
	out := reg.FindStringSubmatch(stdout)
	if len(out) < 2 || out[1] == "" {
		err = fmt.Errorf("cannot match outfile path in stdout: \n%s", stdout)
		logs.Error(err)
		return
	}
	outfile := out[1]
	logs.Info("outfile:%s", out[1])

	p.Converted, err = ioutil.ReadFile(outfile)
	if err != nil {
		logs.Error(err)
		return
	}
	if p.delOutfile {
		defer os.Remove(outfile)
	}

	*reply = p.Converted

	return nil
}

//文件数据保存到指定目录，md5作为文件名
func saveFileByMd5(data []byte, dir string) (filepath, MD5 string, err error) {

	if len(data) == 0 || dir == "" {
		err = fmt.Errorf("zero value or arguments")
		logs.Error(err)
		return
	}

	dir = strings.TrimRight(dir, "/") + "/"
	err = os.MkdirAll(dir, 0666)
	if err != nil {
		logs.Error(err)
		return
	}

	hash := md5.New()
	_, err = hash.Write(data)
	if err != nil {
		logs.Error(err)
		return
	}
	MD5 = hex.EncodeToString(hash.Sum(nil))
	filepath = dir + MD5

	err = ioutil.WriteFile(filepath, data, 0666)
	if err != nil {
		logs.Error(err)
		return
	}
	return
}

func shellOut(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func InitLogger(path string) error {
	if path == "" {
		path = "/tmp/log/"
	}
	//创建日志目录
	path = strings.TrimRight(path, "/") + "/"
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatal("创建日志目录失败 ", err)
		return err
	}
	fileName := path + "log.log"
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("创建日志文件失败 ", err)
		return err
	}
	defer f.Close()

	logs.SetLogFuncCallDepth(3)    //调用层级
	logs.EnableFuncCallDepth(true) //输出文件名和行号
	//logs.Async()                   //提升性能, 可以设置异步输出

	config := make(map[string]interface{})
	config["filename"] = fileName

	logs.SetLevel(logs.LevelDebug)

	configStr, err := json.Marshal(config)
	if err != nil {
		log.Fatal("initLogger failed, marshal err:", err)
		return err
	}
	err = logs.SetLogger(logs.AdapterConsole, "") //控制台输出
	if err != nil {
		log.Fatal("SetLogger failed, err:", err)
		return err
	}
	err = logs.SetLogger(logs.AdapterFile, string(configStr)) //文件输出
	if err != nil {
		log.Fatal("SetLogger failed, err:", err)
		return err
	}
	//err = logs.SetLogger(logs.AdapterEs, `{"dsn":"http://localhost:9200/","level":1}`)
	//if err != nil {
	//	log.Fatal("SetLogger failed, err:", err)
	//	return err
	//}
	return nil
}

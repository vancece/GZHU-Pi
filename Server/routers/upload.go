/**
 * @File: upload
 * @Author: Shaw
 * @Date: 2020/5/8 9:44 PM
 * @Desc

 */

package routers

import (
	"GZHU-Pi/services/cosfs"
	"GZHU-Pi/services/facepp"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var maxUploadFileSize = int64(10 * 1024 * 1024)
var storagePath string = "/tmp/"

func Upload(w http.ResponseWriter, r *http.Request) {

	_ = r.ParseMultipartForm(128)
	//fmt.Println("r.Form: ", r.Form)

	file, fHeader, err := r.FormFile("file")
	if err != nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}
	defer file.Close()
	if fHeader.Size == 0 || fHeader.Filename == "" {
		err = fmt.Errorf("zero size or empty filename")
		Response(w, r, nil, http.StatusBadRequest, err.Error())
		return
	}

	d, err := ioutil.ReadAll(file)
	if err != nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}
	segmentData, err := facepp.HumanBodySegment(d)
	if err != nil {
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}

	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}
	MD5 := hex.EncodeToString(hash.Sum(nil))

	url, err := cosfs.SaveToCos(segmentData, MD5+".png")
	if err != nil {
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}
	logs.Info(url)

	data := map[string]interface{}{"url": url}

	Response(w, r, data, http.StatusOK, "upload success")
}

func saveFile(r *http.Request) (err error) {
	//解析form-data
	_ = r.ParseMultipartForm(128)

	file, fHeader, err := r.FormFile("file")
	if err != nil {
		logs.Error(err)
		return
	}
	defer file.Close()

	if fHeader.Filename == "" {
		err = fmt.Errorf("文件名不能为空")
		logs.Error(err)
		return
	}
	if fHeader.Size == 0 {
		err = fmt.Errorf("文件大小不能为零")
		logs.Error(err)
		return
	}
	if fHeader.Size > maxUploadFileSize {
		err = fmt.Errorf("文件 %s 的大小为 %4.2f兆，超过了系统的 %4.2f兆限制", fHeader.Filename,
			float64(fHeader.Size)/(1024*1024), float64(maxUploadFileSize)/(1024*1024))
		logs.Error(err)
		return
	}

	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		logs.Error(err)
		return
	}
	MD5 := hex.EncodeToString(hash.Sum(nil))

	filename := fmt.Sprintf("%s%s_%s", storagePath, MD5, fHeader.Filename)
	logs.Info(filename)

	var dst *os.File
	dst, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logs.Error(err)
		return
	}
	defer dst.Close()

	var written int64
	written, err = file.Seek(0, io.SeekStart)
	if err != nil {
		logs.Error(err)
		return
	}
	if written != 0 {
		err = fmt.Errorf("rewind to start failed")
		logs.Error(err)
		return
	}
	written, err = io.Copy(dst, file)
	if err != nil {
		logs.Error(err)
		return
	}
	if written != fHeader.Size {
		err = fmt.Errorf("short written")
		logs.Error(err)
		return
	}
	return
}

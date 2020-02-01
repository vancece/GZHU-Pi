package routers

import (
	"GZHU-Pi/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func TableAccessHandle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	ctx := InitCtx(w, r)

	switch strings.ToUpper(r.Method) {

	case "GET":
	case "POST":
		if strings.Contains(r.URL.Path, "t_topic") {

			if err := TopicCheck(ctx); err != nil {
				Response(w, r, nil, http.StatusBadRequest, err.Error())
				return
			}

		}

	case "PUT", "PATCH":
	case "DELETE":
		p := getCtxValue(ctx)
		qry := strings.ReplaceAll(p.r.URL.Query().Get("id"), "$eq.", "")
		id, err := strconv.ParseInt(qry, 10, 64)
		if err != nil {
			logs.Error(err)
			Response(w, r, nil, http.StatusBadRequest, err.Error())
			return
		}
		if strings.Contains(r.URL.Path, "t_topic") {
			var t models.TTopic
			p.gormDB.First(&t, id)
			if t.CreatedBy.Int64 != p.user.ID {
				err := fmt.Errorf("permission denied")
				Response(w, r, nil, http.StatusBadRequest, err.Error())
				return
			}
		}

	default:
		_, _ = w.Write([]byte("unsupported method: " + r.Method))
		return
	}

	next(w, r)
}

func TopicCheck(ctx context.Context) (err error) {
	p := getCtxValue(ctx)

	body, err := ioutil.ReadAll(p.r.Body)
	if err != nil {
		logs.Error(err)
		return
	}
	defer p.r.Body.Close()
	if len(body) == 0 {
		err = fmt.Errorf("Call api by post with empty body ")
		logs.Error(err)
		return
	}
	var t models.TTopic
	err = json.Unmarshal(body, &t)
	if err != nil {
		logs.Error(err)
		return
	}
	if t.Type.String == "" || t.Title.String == "" || t.Content.String == "" {
		err = fmt.Errorf("必要字段咋能为空")
		logs.Error(err)
		return
	}
	if t.Anonymous.Bool == true && t.Anonymity.String == "" {
		err = fmt.Errorf("请指定 Anonymity 的值")
		logs.Error(err)
		return
	}
	if t.CreatedBy.Valid {
		err = fmt.Errorf("不能手动指定created_by")
		logs.Error(err)
		return
	}

	newBodyStr := fmt.Sprintf(`%s,"created_by":%d}`, strings.TrimSuffix(string(body), "}"), p.user.ID)

	body = []byte(newBodyStr)
	p.r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return
}

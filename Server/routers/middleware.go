package routers

import (
	"GZHU-Pi/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func TableAccessHandle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	//====== 无需校验token的接口 =======

	if strings.Contains(r.URL.Path, "/auth") ||
		strings.Contains(r.URL.Path, "/jwxt") ||
		strings.Contains(r.URL.Path, "/library") ||
		strings.Contains(r.URL.Path, "/second") {
		next(w, r)
		return
	}
	if strings.ToUpper(r.Method) == "GET" && strings.Contains(r.URL.Path, "_topic") {
		topicViewCounter(r.URL)
	}

	//======= 数据库可以找到对应用户、需要检查token =========

	ctx, err := InitCtx(w, r)
	if err != nil {
		return
	}
	switch strings.ToUpper(r.Method) {

	case "GET":
	case "POST":
		if strings.Contains(r.URL.Path, "t_topic") {
			if err := topicCheck(ctx); err != nil {
				Response(w, r, nil, http.StatusBadRequest, err.Error())
				return
			}
		}
		if strings.Contains(r.URL.Path, "t_discuss") {
			if err := discussCheck(ctx); err != nil {
				Response(w, r, nil, http.StatusBadRequest, err.Error())
				return
			}
		}
		if strings.Contains(r.URL.Path, "t_relation") {
			if err := relationCheck(ctx); err != nil {
				Response(w, r, nil, http.StatusBadRequest, err.Error())
				return
			}
		}
	case "PUT", "PATCH":
		p := getCtxValue(ctx)
		if strings.Contains(r.URL.Path, "t_user") {
			logs.Info(r.URL.RawQuery)
			r.URL.RawQuery = fmt.Sprintf("id=%d", p.user.ID)
			if err := userCheck(ctx); err != nil {
				Response(w, r, nil, http.StatusBadRequest, err.Error())
				return
			}
		}
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
		if strings.Contains(r.URL.Path, "t_discuss") {
			var t models.TDiscuss
			p.gormDB.First(&t, id)
			if t.CreatedBy.Int64 != p.user.ID {
				err := fmt.Errorf("permission denied")
				Response(w, r, nil, http.StatusBadRequest, err.Error())
				return
			}
		}
		if strings.Contains(r.URL.Path, "t_relation") {
			var t models.TRelation
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

func topicCheck(ctx context.Context) (err error) {
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

func discussCheck(ctx context.Context) (err error) {
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
	var t models.TDiscuss
	err = json.Unmarshal(body, &t)
	if err != nil {
		logs.Error(err)
		return
	}
	if t.ObjectID <= 0 || t.Content.String == "" {
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

func relationCheck(ctx context.Context) (err error) {
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
	var t models.TRelation
	err = json.Unmarshal(body, &t)
	if err != nil {
		logs.Error(err)
		return
	}
	if t.ObjectID <= 0 {
		err = fmt.Errorf("Are you kidding me ? ")
		logs.Error(err)
		return
	}
	if t.Object.String != "t_topic" && t.Object.String != "t_discuss" {
		err = fmt.Errorf("unsupported object name: %s", t.Object.String)
		logs.Error(err)
		return
	}
	if t.Type.String != "star" && t.Type.String != "claim" && t.Type.String != "favourite" {
		err = fmt.Errorf("Are you kidding me ? ")
		logs.Error(err)
		return
	}
	if t.CreatedBy.Valid {
		err = fmt.Errorf("不能手动指定created_by")
		logs.Error(err)
		return
	}
	//根据唯一主键删除，防止写入冲突
	p.gormDB.Where("object_id=? and object=? and type=? and created_by=?",
		t.ObjectID, t.Object, t.Type, p.user.ID).Delete(models.TRelation{})

	newBodyStr := fmt.Sprintf(`%s,"created_by":%d}`, strings.TrimSuffix(string(body), "}"), p.user.ID)

	body = []byte(newBodyStr)
	p.r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return
}

func userCheck(ctx context.Context) (err error) {

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
	var u models.TUser
	err = json.Unmarshal(body, &u)
	if err != nil {
		logs.Error(err)
		return
	}
	if u.ID != 0 || u.RoleID.Int64 != 0 || u.OpenID.String != "" {
		err = fmt.Errorf("could not update id/role_id/open_id")
		logs.Error(err)
		return
	}
	if u.Phone.String != "" && !verifyPhone(u.Phone.String) {
		err = fmt.Errorf("%s not a valid phone number", u.Phone.String)
		return
	}
	return
}

//浏览人数+1
func topicViewCounter(u *url.URL) {

	if u == nil {
		return
	}
	q := u.Query().Get("id")
	idStr := strings.Trim(q, "$eq.")

	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		return
	}

	db := models.GetGorm()
	db.Model(&models.TTopic{ID: int64(id)}).
		UpdateColumn("viewed", gorm.Expr("viewed + ?", 1))

}

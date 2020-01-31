package routers

import (
	"GZHU-Pi/models"
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	"net/http"
)

type piCtx struct {
	r *http.Request
	w http.ResponseWriter

	user models.TUser

	gormDB *gorm.DB
}

type ctxKey string

func (v ctxKey) String() string {
	return string(v)
}

const piKey = ctxKey("qNearCtx")

func getCtxValue(ctx context.Context) (p *piCtx) {
	var err error
	f := ctx.Value(piKey)
	if f == nil {
		err = fmt.Errorf(`get nil from ctx.Value["%s"]`, piKey.String())
		logs.Error(err.Error())
		panic(err.Error())
	}
	var ok bool
	p, ok = f.(*piCtx)
	if !ok {
		err := fmt.Errorf("failed to type assertion for *piCtx")
		logs.Error(err.Error())
		panic(err.Error())
	}
	if p == nil {
		err := fmt.Errorf(`ctx.Value["%s"] should be non nil *piCtx`, piKey.String())
		logs.Error(err.Error())
		panic(err.Error())
	}
	return
}

func InitCtx(w http.ResponseWriter, r *http.Request) (ctx context.Context, ) {

	p := &piCtx{
		r:    r,
		w:    w,
		user: models.TUser{ID: 1},

		gormDB: models.GetGorm(),
	}
	ctx = context.Background()
	ctx = context.WithValue(ctx, piKey, p)

	//TODO 读取cookie、初始化user
	token := p.r.Header.Get("Authorization")
	if token == "" {

	}
	//刷新token
	//p.w.Header().Set("Set-Cookie", token)

	return
}

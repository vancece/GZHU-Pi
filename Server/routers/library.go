package routers

import (
	"GZHU-Pi/pkg/gzhu_library"
	"github.com/astaxie/beego/logs"
	"net/http"
)

func BookSearch(w http.ResponseWriter, r *http.Request) {
	q, _ := ReadRequestArg(r, "query")
	p, _ := ReadRequestArg(r, "page")

	query, _ := q.(string)
	page, _ := p.(string)

	if query == "" || page == "" {
		Response(w, r, nil, http.StatusBadRequest, "illegal request query")
		return
	}

	data, err := gzhu_library.BookSearch(query, page)
	if err != nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}
	Response(w, r, data, http.StatusOK, "request ok")
}

func BookHoldings(w http.ResponseWriter, r *http.Request) {
	q, _ := ReadRequestArg(r, "id")
	p, _ := ReadRequestArg(r, "source")

	query, _ := q.(string)
	page, _ := p.(string)

	if query == "" || page == "" {
		Response(w, r, nil, http.StatusBadRequest, "illegal request query")
		return
	}

	data, err := gzhu_library.GetHoldings(query, page)
	if err != nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}
	Response(w, r, data, http.StatusOK, "request ok")
}

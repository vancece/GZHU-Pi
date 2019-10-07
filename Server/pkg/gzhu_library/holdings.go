package gzhu_library

//图书馆藏信息
type Info struct {
	BarCode   string `json:"bar_code" remark:"条码"`
	Circulate string `json:"circulate" remark:"流通类型"`
	DueBack   string `json:"due_back" remark:"还书时间"`
	Explain   string `json:"explain" remark:"流通说明"`
	LoanDate  string `json:"loan_date" remark:"借出时间"`
	Location  string `json:"location" remark:"馆藏地点"`
	Status    string `json:"status" remark:"馆藏状态"`
}

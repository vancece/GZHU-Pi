package gzhu_library

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"regexp"
)

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

func GetHoldings(bookID, bookSource string) (infoList []*Info, err error) {

	resp, err := http.Get(`http://lib.gzhu.edu.cn:8080/bookle/search2/detail/` + bookID + `?source=` + bookSource)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	var info = &Info{}

	//匹配table
	re := regexp.MustCompile(`〖馆藏信息〗([\s\S]*?)</table>`)
	table := re.FindAllStringSubmatch(string(bytes), -1)
	if table == nil {
		logs.Error(`正则匹配失败，返回为空`)
		return nil, fmt.Errorf(`正则匹配失败，返回为空`)
	} else {

		//匹配tr
		re = regexp.MustCompile(`tr([\s\S]*?)</tr>`)
		tr := re.FindAllStringSubmatch(table[0][1], -1)
		if tr == nil {
			logs.Error(`正则匹配失败，返回为空`)
		} else {
			for i := 1; i < len(tr); i++ {
				//匹配td
				re = regexp.MustCompile(`td>([\s\S]*?)</td>`)
				td := re.FindAllStringSubmatch(tr[i][1], -1)
				if td == nil {
					logs.Error(`正则匹配失败，返回为空`)
				} else {
					if len(td) != 7 {
						logs.Error(`馆藏信息数不为7`)
					} else {
						info.BarCode = td[0][1]
						info.Status = td[1][1]
						info.LoanDate = td[2][1]
						info.DueBack = td[3][1]
						re = regexp.MustCompile(`[\S]+`)
						Location := re.FindAllStringSubmatch(td[4][1], -1)
						info.Location = Location[0][0]
						info.Circulate = td[5][1]
						info.Explain = td[6][1]
					}
				}
				infoList = append(infoList, info)
			}
		}

	}
	return infoList, nil
}

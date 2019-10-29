/**
广州大学图书馆网站搜索查询接口
*/

package gzhu_library

import (
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

type Book struct {
	ISBN      string `json:"ISBN" remark:"ISBN"`
	Author    string `json:"author" remark:"作者"`
	BookName  string `json:"book_name" remark:"书名"`
	CallNo    string `json:"call_No" remark:"索书号"`
	Copies    string `json:"copies" remark:"复本数量"`
	ID        string `json:"id" remark:"书籍id"`
	Image     string `json:"image" remark:"封面图片"`
	Loanable  string `json:"loanable" remark:"可借数量"`
	Publisher string `json:"publisher" remark:"出版社"`
	Source    string `json:"source" remark:"搜索源"`
}

/**
query:用户搜索的图书
searchPage:请求页数
[]*Book:返回Book类型的切片
*/
func BookSearch(query string, searchPage string) (books []*Book, err error) {

	if query == "" || searchPage == "" {
		return
	}
	//为0时会返回所有页数据
	if searchPage == "0" {
		searchPage = "1"
	}

	resp, err := http.PostForm(`http://lib.gzhu.edu.cn:8080/bookle/`, url.Values{"query": {query}, "searchPage": {searchPage}})
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

	//匹配book_info
	re := regexp.MustCompile(`<div class=book_info>([\s\S]*?)</div>`)
	bookInfo := re.FindAllStringSubmatch(string(bytes), -1)
	if bookInfo == nil {
		logs.Info(`正则匹配失败，返回为空，书名:%s`, query)
		return
	}

	var wg = sync.WaitGroup{}

	for i := 0; i < len(bookInfo); i++ {
		if len(bookInfo[i]) < 2 {
			logs.Error("非预期错误", bookInfo[i])
			continue
		}
		book := &Book{}

		//匹配book_info下的h2
		re = regexp.MustCompile(`<h2[\s\S]*?>([\s\S]*?)</h2>`)

		bookInfoH2 := re.FindAllStringSubmatch(bookInfo[i][1], -1)
		if bookInfoH2 == nil {
			logs.Error(`正则匹配失败，返回为空`)
		} else { //书名
			re = regexp.MustCompile(`<a[\s\S]*?>[\s]*?([\S][\s\S]*?)</a>`)
			bookInfoH2A := re.FindAllStringSubmatch(bookInfoH2[0][1], -1)
			if bookInfoH2A == nil {
				logs.Error(`正则匹配失败，返回为空`)
			} else {
				book.BookName = bookInfoH2A[0][1]
				book.BookName = strings.TrimSuffix(book.BookName, " /")

				//书籍id
				re = regexp.MustCompile(`[\d]{3,}`)
				bookID := re.FindAllStringSubmatch(bookInfoH2A[0][0], -1)
				if bookID == nil {
					logs.Error(`正则匹配失败，返回为空`)
				} else {
					book.ID = bookID[0][0]
				}

				//搜索源
				re = regexp.MustCompile(`source=(.+)"`)
				bookResource := re.FindAllStringSubmatch(bookInfoH2A[0][0], -1)
				if bookResource == nil {
					logs.Error(`正则匹配失败，返回为空`)
				} else {
					book.Source = bookResource[0][1]
				}
			}

			//作者
			re = regexp.MustCompile(`<span[\s\S]*?>(.*?)</span>`)
			bookAuthor := re.FindAllStringSubmatch(bookInfoH2[0][1], -1)
			if bookAuthor == nil {
				logs.Error(`正则匹配失败，返回为空`)
			} else {
				book.Author = bookAuthor[0][1]
			}
		}

		//匹配book_info下的h4
		re = regexp.MustCompile(`<h4[\s\S]*?>[\s\S]*?</h4>`)
		bookInfoH4 := re.FindAllStringSubmatch(bookInfo[i][1], -1)
		if bookInfoH4 == nil {
			logs.Error(`正则匹配失败，返回为空`)
		} else {
			//出版社
			re = regexp.MustCompile(`出版发行：([\s\S]*?)\r\n`)
			bookPublisher := re.FindAllStringSubmatch(bookInfoH4[0][0], -1)
			if bookPublisher == nil {
				logs.Error(`正则匹配失败，返回为空`)
			} else {
				book.Publisher = bookPublisher[0][1]
				book.Publisher = strings.Replace(book.Publisher, ",", "", -1)
				book.Publisher = strings.Replace(book.Publisher, "O&#39", "", -1)
				book.Publisher = strings.Replace(book.Publisher, ";", "", -1)
				book.Publisher = strings.Replace(book.Publisher, " ", "", -1)
			}

			re = regexp.MustCompile(`&nbsp;&nbsp;[\s\S]*?&nbsp;&nbsp;`)
			bookISbnRow := re.FindAllStringSubmatch(bookInfoH4[0][0], -1)
			if bookISbnRow == nil {
				logs.Error(`正则匹配失败，返回为空`)
			} else {
				//ISbn
				re = regexp.MustCompile(`[\d-]{10,}`)
				bookISbn := re.FindAllStringSubmatch(bookISbnRow[0][0], -1)
				if bookISbn == nil {
					logs.Error(`正则匹配失败，返回为空`)
				} else {
					book.ISBN = bookISbn[0][0]
				}
				book.ISBN = strings.Replace(book.ISBN, "-", "", -1)
			}

			//复本数量
			re = regexp.MustCompile(`复本数.*([\d]+).*\n`)
			bookCopyNumber := re.FindAllStringSubmatch(bookInfoH4[1][0], -1)
			if bookCopyNumber == nil {
				logs.Error(`正则匹配失败，返回为空`)
			} else {
				book.Copies = bookCopyNumber[0][1]
			}

			//可借数量
			re = regexp.MustCompile(`在馆数.*([\d]+).*\n`)
			bookCouldBeBorrow := re.FindAllStringSubmatch(bookInfoH4[1][0], -1)
			if bookCouldBeBorrow == nil {
				logs.Error(`正则匹配失败，返回为空`)
			} else {
				book.Loanable = bookCouldBeBorrow[0][1]
			}

			//索书号
			re = regexp.MustCompile(`([A-Z][\s\S]*?)\r\n`)
			bookCallNumber := re.FindAllStringSubmatch(bookInfoH4[1][0], -1)
			if bookCallNumber == nil {
				logs.Error(`正则匹配失败，返回为空`)
			} else {
				book.CallNo = bookCallNumber[0][1]
			}
		}
		//异步请求豆瓣接口获取图书封面
		wg.Add(1)
		go func() {
			book.Image = GetCover(book.ISBN)
			wg.Done()
		}()
		books = append(books, book)
	}
	wg.Wait()
	return books, nil
}

//提取豆瓣图书封面
func GetCover(ISBN string) (image string) {
	if ISBN == "" {
		return
	}
	var URL = "https://douban.uieee.com/v2/book/isbn/" + ISBN
	client := http.Client{}
	request, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		logs.Error(err)
		return
	}
	const UA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36"
	request.Header.Set("User-Agent", UA)
	resp, err := client.Do(request)
	if err != nil {
		logs.Error(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	image = jsoniter.Get(body, "image").ToString()
	return image
}

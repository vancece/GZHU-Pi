package gzhu_library

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

var (
	book     Book
	BookList []Book
)

/**
query:用户搜索的图书
searchPage:请求页数
[]Book:返回Book类型的数组
*/
func LibraryBookSearch(query string, searchPage string) ([]Book, error) {
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
		logs.Error(`正则匹配失败，返回为空`)
		return nil, fmt.Errorf(`正则匹配失败，返回为空`)
	}

	for i := 0; i < len(bookInfo); i++ {

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
				//logs.Debug(book.BookName)

				//书籍id
				re = regexp.MustCompile(`[\d]{3,}`)
				bookID := re.FindAllStringSubmatch(bookInfoH2A[0][0], -1)
				if bookID == nil {
					logs.Error(`正则匹配失败，返回为空`)
				} else {
					book.ID = bookID[0][0]
				}
				//logs.Debug(book.ID)

				//搜索源
				re = regexp.MustCompile(`source=(.+)"`)
				bookResource := re.FindAllStringSubmatch(bookInfoH2A[0][0], -1)
				if bookResource == nil {
					logs.Error(`正则匹配失败，返回为空`)
				} else {
					book.Source = bookResource[0][1]
				}
				//logs.Debug(book.Source)
			}

			//作者
			re = regexp.MustCompile(`<span[\s\S]*?>(.*?)</span>`)
			bookAuthor := re.FindAllStringSubmatch(bookInfoH2[0][1], -1)
			if bookAuthor == nil {
				logs.Error(`正则匹配失败，返回为空`)
			} else {
				book.Author = bookAuthor[0][1]
			}
			//logs.Debug(book.Author)
		}

		//匹配book_info下的h4
		re = regexp.MustCompile(`<h4[\s\S]*?>[\s\S]*?</h4>`)
		bookInfoH4 := re.FindAllStringSubmatch(bookInfo[i][1], -1)
		if bookInfoH4 == nil {
			logs.Error(`正则匹配失败，返回为空`)
		} else {
			//出版社
			re = regexp.MustCompile(`出版发行：([\s\S]*?)\n`)
			bookPublisher := re.FindAllStringSubmatch(bookInfoH4[0][0], -1)
			if bookPublisher == nil {
				logs.Error(`正则匹配失败，返回为空`)
			} else {
				book.Publisher = bookPublisher[0][1]
				//logs.Debug(book.Publisher)
			}

			re = regexp.MustCompile(`&nbsp;&nbsp;[\s\S]*?&nbsp;&nbsp;`)
			bookISbnRow := re.FindAllStringSubmatch(bookInfoH4[0][0], -1)
			if bookISbnRow == nil {
				logs.Error(`正则匹配失败，返回为空`)
			} else {
				// logs.Info(bookISbnRow)

				//ISbn
				re = regexp.MustCompile(`[\d-]{10,}`)
				bookISbn := re.FindAllStringSubmatch(bookISbnRow[0][0], -1)
				if bookISbn == nil {
					logs.Error(`正则匹配失败，返回为空`)
				} else {
					book.ISBN = bookISbn[0][0]
					//logs.Debug(book.ISBN)
				}
			}

			//复本数量
			re = regexp.MustCompile(`复本数.*([\d]+).*\n`)
			bookCopyNumber := re.FindAllStringSubmatch(bookInfoH4[1][0], -1)
			if bookCopyNumber == nil {
				logs.Error(`正则匹配失败，返回为空`)
			} else {
				book.Copies = bookCopyNumber[0][1]
				// logs.Debug(book.Copies)
			}

			//可借数量
			re = regexp.MustCompile(`在馆数.*([\d]+).*\n`)
			bookCouldBeBorrow := re.FindAllStringSubmatch(bookInfoH4[1][0], -1)
			if bookCouldBeBorrow == nil {
				logs.Error(`正则匹配失败，返回为空`)
			} else {
				book.Loanable = bookCouldBeBorrow[0][1]
				//logs.Debug(book.Loanable)
			}

			//索书号
			re = regexp.MustCompile(`([A-Z][\s\S]*?)\n`)
			bookCallNumber := re.FindAllStringSubmatch(bookInfoH4[1][0], -1)
			if bookCallNumber == nil {
				logs.Error(`正则匹配失败，返回为空`)
			} else {
				book.CallNo = bookCallNumber[0][1]
				//logs.Debug(book.CallNo)
			}
		}

		//封面图片
		resp, err = http.Get(`https://douban.uieee.com/v2/book/isbn/` + book.ISBN)
		if err != nil {
			logs.Error(err)
			return nil, err
		}
		defer resp.Body.Close()
		bytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			logs.Error(err)
			return nil, err
		}
		var s interface{}
		json.Unmarshal(bytes, &s)
		imgs := s.(map[string]interface{})
		img := imgs["images"].(map[string]interface{})
		book.Image = img["small"].(string)
		//logs.Debug(book.Image)

		logs.Info("\n%+v\n", book)
		BookList = append(BookList, book)
	}
	return BookList, nil
}

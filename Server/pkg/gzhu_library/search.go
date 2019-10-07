package gzhu_library


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
package forms

import (
	"net/http"
	"net/url"
)

//自創的form 結構
type Form struct {
	url.Values
	Errors errors
}


//初始化form
func New(data url.Values) *Form{
	return &Form{
		data,
		errors(map[string][]string{}),		
	}
}

//Has 檢查表單是否有填寫
func (f *Form) Has(field string,r *http.Request) bool {
	x := r.Form.Get(field)
	//表單必須要填寫
	if x == "" {
		return false
	}
	return true
}
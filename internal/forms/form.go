package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

//自創的form 結構
type Form struct {
	url.Values
	Errors errors
}
//Valid 如果有錯return False 如果正確return true
func(f *Form) Valid() bool{
	return len(f.Errors) == 0
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

//Required to validate form data 
func (f *Form) Required(fields ...string){
	for _, field := range fields{
		//從form object get field 
		value := f.Get(field)
		//如果表單內容未填寫
		if strings.TrimSpace(value) == ""{
			f.Errors.Add(field,"必須填寫表單")
		}

	}
}

//Minlength to validate form data length
func(f *Form) MinLength(field string,length int,r *http.Request) bool{
	x := r.Form.Get(field)
	if len(x) < length {
		f.Errors.Add(field,fmt.Sprintf("This field must be at least %d characters long",length))
		return false
	}
	return true
}

//Isemail use Govalidator to check email addres
func (f *Form) Isemail(field string){
	if !govalidator.IsEmail(f.Get(field)){
		f.Errors.Add(field,"電子郵件錯誤")
	}
	
}


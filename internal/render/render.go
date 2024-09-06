package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/donbaking/booking/internal/config"
	"github.com/donbaking/booking/internal/models"
	"github.com/justinas/nosurf"
)
var functions = template.FuncMap{
	"humanDate" : HumanDate,
}

var app *config.AppConfig
//建立templates的path,提供給測試用
var pathToTemplates = "./templates"

//NewRenderer 從template package 設置config
func  NewRenderer(a *config.AppConfig){
	app = a
}

//
func HumanDate(t time.Time)string{
	return t.Format("2006-01-02")
}

//AddDefaultData 用在如果要在每個頁面加上相同資料時可以使用
func AddDefaultData(td *models.TemplateData,r *http.Request) *models.TemplateData{
	//用seesion讓頁面獲得Flash
	td.Flash = app.Session.PopString(r.Context(),"flash")
	td.Error = app.Session.PopString(r.Context(),"error")
	td.Warning = app.Session.PopString(r.Context(),"warning")
	
	//讓頁面獲得CSTF的token,從nosurf package拿token 
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(),"user_id"){
		td.IsAuthenticated = true
	}
	return td
}


// 渲染前端用的Func 將頁面名稱傳入函式後渲染指定的html頁面
func Template(w http.ResponseWriter, r *http.Request,tmpl string,td *models.TemplateData) error{
	var tempcache map[string]*template.Template
	
	//
	if app.UseCache{
		//從appconfig獲得template cache
		tempcache = app.TemplateCache
	}else{
		tempcache,_ = CreateTemplateCache()
	}
	
	//get request template from cache
	t, ok:= tempcache[tmpl]
	if !ok{
		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td , r)
	
    _ = t.Execute(buf,td)
	//render the template
	_,err :=buf.WriteTo(w)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}


func CreateTemplateCache() (map[string]*template.Template,error) {
	//創建一個空的map 在後面加上一個{}代表為空
	myCache := map[string]*template.Template{}
	//從templates資料夾中取得所有資料
	pages , err := filepath.Glob(fmt.Sprintf("%s/*page.tmpl",pathToTemplates))
	if err!= nil {
		return myCache,err
	}
	//遍歷pages取得的所有資料
	for _,page := range pages {
		//name會取得tmpL的檔名
		name := filepath.Base(page)
        ts, err := template.New(name).Funcs(functions).ParseFiles(page)
        if err!= nil {
            return myCache,err
        }
		
		matches ,err := filepath.Glob(fmt.Sprintf("%s/*layout.tmpl",pathToTemplates))
		if err!=nil{
			return myCache ,err
		}
		
		if len(matches)>0{
			ts,err = ts.ParseGlob(fmt.Sprintf("%s/*layout.tmpl",pathToTemplates))
			if err!=nil{
				return myCache ,err
			}
		}
		//將目前的模板存到緩存
		myCache[name] = ts
	}
	return myCache,nil
}

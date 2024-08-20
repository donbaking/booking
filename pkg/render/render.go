package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/donbaking/booking/pkg/config"
	"github.com/donbaking/booking/pkg/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig
//NewTemplates 從template package 設置config
func  NewTemplates(a *config.AppConfig){
	app = a
}

//AddDefaultData 用在如果要在每個頁面加上相同資料時可以使用
func AddDefaultData(td *models.TemplateData,r *http.Request) *models.TemplateData{
	//讓頁面獲得CSTF的token,從nosurf package拿token 
	td.CSRFToken = nosurf.Token(r)
	
	return td
}


// 渲染前端用的Func 將頁面名稱傳入函式後渲染指定的html頁面
func RenderTemplate(w http.ResponseWriter, r *http.Request,tmpl string,td *models.TemplateData) {
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
		log.Fatal("Could not get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td , r)
	
    _ = t.Execute(buf,td)
	//render the template
	_,err :=buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}


func CreateTemplateCache() (map[string]*template.Template,error) {
	//創建一個空的map 在後面加上一個{}代表為空
	myCache := map[string]*template.Template{}
	//從templates資料夾中取得所有資料
	pages , err := filepath.Glob("./templates/*page.tmpl")
	if err!= nil {
		return myCache,err
	}
	//遍歷pages取得的所有資料
	for _,page := range pages {
		//name會取得tmpL的檔名
		name := filepath.Base(page)
        ts, err := template.New(name).ParseFiles(page)
        if err!= nil {
            return myCache,err
        }
		
		matches ,err := filepath.Glob("./templates/*layout.tmpl")
		if err!=nil{
			return myCache ,err
		}
		
		if len(matches)>0{
			ts,err = ts.ParseGlob("./templates/*layout.tmpl")
			if err!=nil{
				return myCache ,err
			}
		}
		//將目前的模板存到緩存
		myCache[name] = ts
	}
	return myCache,nil
}

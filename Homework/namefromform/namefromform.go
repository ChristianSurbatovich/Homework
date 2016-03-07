package main

import(
	"net/http"
	"html/template"
	"log"
)

var mytemplates *template.Template

func init() {
	var err error
	mytemplates,err = template.ParseGlob("*.gohtml")
	if(err != nil){
		log.Println(err)
	}
}



func formsubmit(res http.ResponseWriter, req *http.Request){

	err := mytemplates.ExecuteTemplate(res,"formsubmit.gohtml",nil)
	if(err != nil){
		log.Println(err)
	}
}

func formshow(res http.ResponseWriter, req *http.Request){
	mytemplates.ExecuteTemplate(res,"formshow.gohtml",req.FormValue("myname"))
}


func main(){
http.Handle("/favicon.ico",http.NotFoundHandler())
http.HandleFunc("/name",formshow)
http.HandleFunc("/",formsubmit)


http.ListenAndServe("localhost:8080",nil)
}

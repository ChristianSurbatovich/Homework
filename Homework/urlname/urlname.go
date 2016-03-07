package main

import(
	"net/http"
	"html/template"
	"log"
	"strings"
)

var mytemplate *template.Template

func init() {
	var err error
	mytemplate, err = template.ParseFiles("gotemplate.gohtml")
	if(err != nil){
		log.Println(err)
	}
}



func homepage(res http.ResponseWriter, req *http.Request){
	name := strings.Split(req.URL.Path,"/")
	err := mytemplate.Execute(res,name[1])
	if(err != nil){
		log.Println(err)
	}
}



func main(){
http.Handle("/favicon.ico",http.NotFoundHandler())
http.HandleFunc("/",homepage)

http.ListenAndServe("localhost:8080",nil)
}

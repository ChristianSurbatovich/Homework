package main

import(
	"net/http"
	"html/template"
	"log"
)

var mytemplate *template.Template

func init() {
	var err error
	mytemplate, err = template.ParseFiles("gotemplate.gohtml")
	if(err != nil){
		log.Println(err)
	}
}



func urlQuery(res http.ResponseWriter, req *http.Request){

	err := mytemplate.Execute(res,req.FormValue("n"))
	if(err != nil){
		log.Println(err)
	}
}



func main(){
http.Handle("/favicon.ico",http.NotFoundHandler())
http.HandleFunc("/",urlQuery)

http.ListenAndServe("localhost:8080",nil)
}

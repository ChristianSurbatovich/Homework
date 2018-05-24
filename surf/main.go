package main

import(
	"net/http"
	"html/template"
	"log"
)

var mytemplate *template.Template

func init() {
	var err error
	mytemplate, err = template.ParseFiles("homepage.gohtml")
	if(err != nil){
		log.Println(err)
	}
}



func homepage(res http.ResponseWriter, req *http.Request){
	err := mytemplate.Execute(res,req.URL.Host + req.URL.Path)
	if(err != nil){
		log.Println(err)
	}
}



func main(){
	http.Handle("/favicon.ico",http.NotFoundHandler())
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./css"))))
	http.Handle("/pictures/", http.StripPrefix("/pictures", http.FileServer(http.Dir("./pictures"))))
	http.HandleFunc("/",homepage)


	http.ListenAndServe("localhost:8080",nil)
}

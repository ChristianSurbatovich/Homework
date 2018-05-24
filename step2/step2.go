package main

import(
	"net/http"
	"html/template"
	"log"
	"github.com/satori/go.uuid"
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
	mycookie, err := req.Cookie("sessionfino")
	if err != nil {
		log.Println("creating a cookie")
		myuuid := uuid.NewV4()
		mycookie = &http.Cookie{
			Name: "my-id",
			Value: "0",
			HttpOnly: true,
			//Secure: true,
		}
		mycookie.Value = myuuid.String()
		http.SetCookie(res, mycookie)
	}
	err = mytemplate.Execute(res,"this is a template")
	if err != nil{
		log.Println(err)
	}
}



func main(){
	http.Handle("/favicon.ico",http.NotFoundHandler())
	http.HandleFunc("/",homepage)

	http.ListenAndServe("localhost:8080",nil)
}


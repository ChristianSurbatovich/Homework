package main

import(
	"net/http"
	"html/template"
	"log"
	"strconv"

	"github.com/satori/go.uuid"
	"time"
)

var mytemplates *template.Template
var datastuff map[string] string
var num_users int

func init() {
	var err error
	mytemplates,err = template.ParseGlob("*.gohtml")
	if(err != nil){
		log.Println(err)
	}
	num_users = 0
	datastuff = make(map[string]string,1)
}



func cookies(res http.ResponseWriter, req *http.Request){

	mycookie, err := req.Cookie("my-id")
	if(err != nil){
		log.Println("creating a cookie")
			myuuid := uuid.NewV4()
			mycookie = &http.Cookie{
			Name: "my-id",
			Value: "0",
			HttpOnly: true,
			Expires: time.Now().AddDate(0,1,1),
			//Secure: true,
		}
		mycookie.Value = myuuid.String()
		http.SetCookie(res,mycookie)
		num_users++
		datastuff[mycookie.Value] = strconv.Itoa(num_users)

	}


	err = mytemplates.ExecuteTemplate(res,"gotemplate.gohtml","you are number "+datastuff[mycookie.Value])

	if(err != nil){
		log.Println(err)
	}

}




func main(){
	http.Handle("/favicon.ico",http.NotFoundHandler())
	http.HandleFunc("/",cookies)


	http.ListenAndServe("localhost:8080",nil)
}
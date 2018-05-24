package main

import(
	"net/http"
	"html/template"
	"log"
	"strconv"
)

var mytemplates *template.Template

func init() {
	var err error
	mytemplates,err = template.ParseGlob("*.gohtml")
	if(err != nil){
		log.Println(err)
	}
}



func cookies(res http.ResponseWriter, req *http.Request){

	mycookie, err := req.Cookie("my-visits")
	if(err != nil){
		log.Println("creating a cookie")
		mycookie = &http.Cookie{
			Name: "my-visits",
			Value: "0",
		}
	}
	visits, err := strconv.Atoi(mycookie.Value)
	if(err != nil){
		log.Println(err)
	}
	visits++
	mycookie.Value =  strconv.Itoa(visits)
	http.SetCookie(res,mycookie)
	if(visits == 1){
		err = mytemplates.ExecuteTemplate(res,"gotemplate.gohtml","This is your first visit")
	}else{
		err = mytemplates.ExecuteTemplate(res,"gotemplate.gohtml","You have visited " +  mycookie.Value + " times")
	}
	if(err != nil){
		log.Println(err)
	}

}




func main(){
	http.Handle("/favicon.ico",http.NotFoundHandler())
	http.HandleFunc("/",cookies)


	http.ListenAndServe("localhost:8080",nil)
}
package main

import(
	"net/http"
	"html/template"
	"log"
	"github.com/satori/go.uuid"
	"strings"
)

var myTemplates *template.Template



func init() {
	var err error
	myTemplates,err = template.ParseGlob("*.gohtml")
	if err != nil{
		log.Println(err)
	}
}



func homepage(res http.ResponseWriter, req *http.Request){
	mycookie, err := req.Cookie("sessionfino")

	if err != nil {
		log.Println("creating a cookie")
		myuuid := uuid.NewV4()
		mycookie = &http.Cookie{
			Name: "sessionfino",
			Value: myuuid.String(),
			HttpOnly: true,
			//Secure: true,
		}
		http.SetCookie(res, mycookie)
		err = myTemplates.ExecuteTemplate(res,"step3form.gohtml",nil)
		if err != nil{
			log.Println(err)
		}
	}else {
		err = myTemplates.ExecuteTemplate(res, "show.gohtml", mycookie.Value)
		if err != nil {
			log.Println(err)
		}
	}
}

func show(res http.ResponseWriter, req *http.Request){
	mycookie, err := req.Cookie("sessionfino")

	if err != nil{
		log.Println(err)
		return
	}

	myuuid := strings.Split(mycookie.Value,"|")

	mycookie.Value = myuuid[0]+"|"+req.FormValue("myname")+"|"+req.FormValue("myage")

	http.SetCookie(res, mycookie)
	err = myTemplates.ExecuteTemplate(res,"show.gohtml",mycookie.Value)

	if err != nil{
		log.Println(err)
	}

}



func main(){
	http.Handle("/favicon.ico",http.NotFoundHandler())
	http.HandleFunc("/",homepage)
	http.HandleFunc("/show",show)

	http.ListenAndServe("localhost:8080",nil)
}


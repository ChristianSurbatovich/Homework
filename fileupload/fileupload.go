package main

import(
	"net/http"
	"html/template"
	"log"
	"io"
)

var mytemplates *template.Template

func init() {
	var err error
	mytemplates,err = template.ParseGlob("*.gohtml")
	if(err != nil){
		log.Println(err)
	}
}



func filesubmit(res http.ResponseWriter, req *http.Request){

	err := mytemplates.ExecuteTemplate(res,"formupload.gohtml",nil)
	if(err != nil){
		log.Println(err)
	}
}

func fileuploaded(res http.ResponseWriter, req *http.Request) {
	if (req.Method == "POST") {
		input, _, err := req.FormFile("myfile")
		if (err != nil) {
			log.Println(err)
		}

		io.Copy(res, input)

		input.Close()
	}

}




func main(){
	http.Handle("/favicon.ico",http.NotFoundHandler())
	http.HandleFunc("/upload",fileuploaded)
	http.HandleFunc("/",filesubmit)


	http.ListenAndServe("localhost:8080",nil)
}

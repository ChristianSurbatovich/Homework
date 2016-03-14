package main

import(
	"net/http"
	"html/template"
	"log"
	"github.com/satori/go.uuid"
	"strings"
	"encoding/json"
	"encoding/base64"
)

type user struct{
	Name string
	Age string

}



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
		err = myTemplates.ExecuteTemplate(res,"step4form.gohtml",nil)
		if err != nil{
			log.Println(err)
		}
	}else {
		http.Redirect(res,req,"/show",http.StatusFound)
	}
}

func show(res http.ResponseWriter, req *http.Request){
	var currentUser user
	mycookie, err := req.Cookie("sessionfino")
	if err != nil{
		log.Println(err)
		return
	}

	myuuid := strings.Split(mycookie.Value,"|")
	if len(myuuid) < 2{
		currentUser = user{
			Name: req.FormValue("myname"),
			Age: req.FormValue("myage"),
		}
		mycookie.Value = myuuid[0] +"|"+ toJSON64(currentUser)
	}else {
		currentUser = getUser(myuuid[1])
	}
	http.SetCookie(res, mycookie)
	err = myTemplates.ExecuteTemplate(res,"show.gohtml","cookie value = " + mycookie.Value +" decoded cookie uuid = "+myuuid[0]+" name: " + currentUser.Name + " age: "+currentUser.Age)

	if err != nil{
		log.Println(err)
	}

}

func toJSON64(us user)string{
	newStr, err := json.Marshal(us)
	if err != nil{
		log.Println(err)
		return ""
	}
	return base64.URLEncoding.EncodeToString(newStr)
}

func getUser(str string)user{
	var newUser user
	newStr,err := base64.URLEncoding.DecodeString(str)
	if err != nil{
		log.Println(err)
		return newUser
	}

	err = json.Unmarshal(newStr,&newUser)
	if err != nil{
		log.Println(err)
		return newUser
	}
	return newUser;
}

func main(){
	http.Handle("/favicon.ico",http.NotFoundHandler())
	http.HandleFunc("/",homepage)
	http.HandleFunc("/show",show)

	http.ListenAndServe("localhost:8080",nil)
}


package main

import(
	"net/http"
	"html/template"
	"log"
	"github.com/satori/go.uuid"
	"strings"
	"encoding/json"
	"encoding/base64"
	"crypto/hmac"
	"crypto/sha256"
	"io"
	"fmt"
)

type user struct{
	Name string
	Age string

}



var myTemplates *template.Template

var hmac_key = []byte("abcdedf")

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
		err = myTemplates.ExecuteTemplate(res,"step5form.gohtml",nil)
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

	userFino := strings.Split(mycookie.Value,"|")
	if len(userFino) < 2{
		currentUser = user{
			Name: req.FormValue("myname"),
			Age: req.FormValue("myage"),
		}
		str := toJSON64(currentUser)
		mycookie.Value = userFino[0] +"|"+ str +"|" + setKey(str)
		http.SetCookie(res, mycookie)
		userFino = strings.Split(mycookie.Value,"|")
	}else {
		currentUser = getUser(userFino[1])
	}

	if check(userFino[1],userFino[2]){
		err = myTemplates.ExecuteTemplate(res,"show.gohtml","cookie value = " + mycookie.Value +" decoded cookie uuid = "+userFino[0]+" name: " + currentUser.Name + " age: "+currentUser.Age +" hmac = " +userFino[2])
	}else{
		err = myTemplates.ExecuteTemplate(res,"show.gohtml","cookie was tampered with")
	}

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

func setKey(key string) string{
	hm := hmac.New(sha256.New,hmac_key)
	io.WriteString(hm,key)
	return fmt.Sprintf("%x",hm.Sum(nil))
}

func check(key, key2 string) bool{
	// check if hmac of key is key2
	hm := hmac.New(sha256.New,hmac_key)
	io.WriteString(hm,key)
	if(fmt.Sprintf("%x",hm.Sum(nil)) == key2){
		return true
	}
	return false
}

func main(){
	http.Handle("/favicon.ico",http.NotFoundHandler())
	http.HandleFunc("/",homepage)
	http.HandleFunc("/show",show)

	http.ListenAndServe("localhost:8080",nil)
}


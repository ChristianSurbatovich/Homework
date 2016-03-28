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
	"strconv"
)

type user struct{
	Name string
	Age string
	Login bool

}

type data struct{
	CookieValue string
	Good bool
	Login bool
	Logout bool
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


func form(res http.ResponseWriter, req *http.Request){
	mycookie, err := req.Cookie("sessionfino")

	if err != nil {
		http.Redirect(res,req,"/",http.StatusTemporaryRedirect)
	}else {
		userFino := strings.Split(mycookie.Value, "|")
		if len(userFino) > 3 {
			if check(userFino[2], userFino[3]) {
				currentUser := getUser(userFino[2])
				if currentUser.Login {
					err = myTemplates.ExecuteTemplate(res, "step7form.gohtml", data{
						Good: true,
						Login: true,
						Logout: false,

					})
					if err != nil {
						log.Println(err)
					}
				}else{
					http.Redirect(res,req,"/",http.StatusTemporaryRedirect)
				}
			}else{
				http.Redirect(res, req, "/tamper", http.StatusFound)
			}
		}
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
	if len(userFino) < 4{
		http.Redirect(res,req,"/",http.StatusFound)
	}
	currentUser = getUser(userFino[2])
	if currentUser.Login == false{
		http.Redirect(res,req,"/",http.StatusFound)
	}
	currentUser = user{
		Name: req.FormValue("myname"),
		Age: req.FormValue("myage"),
		Login: true,
	}
	str := toJSON64(currentUser)
	// uuid | login status | json data | hmac + 64 bit encoding of json data
	mycookie.Value = userFino[0] + "|" + strconv.FormatBool(currentUser.Login) +"|"+ str +"|" + setKey(str)
	http.SetCookie(res, mycookie)
	http.Redirect(res,req,"/",http.StatusFound)


}

func tamper(res http.ResponseWriter, req *http.Request){
	mycookie, err := req.Cookie("sessionfino")
	if err != nil{
		log.Println(err)
	}
	ss := strings.Split(mycookie.Value,"|")
	mycookie.Value = ss[0]+"|"+ss[1]+"|"+"8uyhjk9tg3hladkfunxz7k" + "|" + "alksfd"
	http.SetCookie(res,mycookie)
	http.Redirect(res,req,"/",http.StatusFound)
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

func homepage(res http.ResponseWriter, req *http.Request){
	var currentUser user
	var mydata data
	mycookie, err := req.Cookie("sessionfino")
	if err != nil{
		log.Println(err)
		mycookie = &http.Cookie{
			Name: "sessionfino",
			HttpOnly: true,
			//Secure: true,
		}
	}

	userFino := strings.Split(mycookie.Value,"|")
	if len(userFino) < 3{
		currentUser = user{
			Login: false,
		}
		str := toJSON64(currentUser)
		// uuid | login status | json data | hmac + 64 bit encoding of json data
		mycookie.Value = uuid.NewV4().String() + "|" + strconv.FormatBool(currentUser.Login) +"|"+ str +"|" + setKey(str)
		http.SetCookie(res, mycookie)
		userFino = strings.Split(mycookie.Value,"|")
	}else {
		currentUser = getUser(userFino[2])
	}

	if check(userFino[2],userFino[3]){
		if currentUser.Login == true{
			mydata = data{
				CookieValue: "you are logged in",
				Good: true,
				Login: currentUser.Login,
				Logout: !currentUser.Login,
			}
		}else{
			mydata = data{
				CookieValue: "you are logged out",
				Good: true,
				Login: currentUser.Login,
				Logout: !currentUser.Login,
			}
		}


	}else{
		mydata = data{
			CookieValue: "cookie was tampered with",
			Good: false,
			Login: false,
			Logout: true,
		}
	}
	err = myTemplates.ExecuteTemplate(res,"show.gohtml",mydata)
	if err != nil{
		log.Println(err)
	}

}

func login(res http.ResponseWriter, req *http.Request){
	mycookie, err := req.Cookie("sessionfino")
	var currentUser user
	if err != nil {
		mycookie = &http.Cookie{
			Name: "sessionfino",
			HttpOnly: true,
			//Secure: true,
		}
		currentUser = user{
			Login: true,
		}
		str := toJSON64(currentUser)
		// uuid | login status | json data | hmac + 64 bit encoding of json data
		mycookie.Value = uuid.NewV4().String() + "|" + strconv.FormatBool(currentUser.Login) +"|"+ str +"|" + setKey(str)
		http.SetCookie(res, mycookie)
	}else {
		userFino := strings.Split(mycookie.Value, "|")
		if len(userFino) > 3 {
			if check(userFino[2], userFino[3]) {
				currentUser = getUser(userFino[2])
				if currentUser.Login == true{
					http.Redirect(res,req,"/form",http.StatusTemporaryRedirect)
				}
				currentUser.Login = true;
			}
			str := toJSON64(currentUser)
			// uuid | login status | json data | hmac + 64 bit encoding of json data
			mycookie.Value = userFino[0] + "|" + strconv.FormatBool(currentUser.Login) + "|" + str + "|" + setKey(str)
			http.SetCookie(res, mycookie)
		}else{
			currentUser = user{
				Login: true,
			}
			str := toJSON64(currentUser)
			// uuid | login status | json data | hmac + 64 bit encoding of json data
			mycookie.Value = uuid.NewV4().String() + "|" + strconv.FormatBool(currentUser.Login) +"|"+ str +"|" + setKey(str)
			http.SetCookie(res, mycookie)
		}

	}
	http.Redirect(res,req,"/",http.StatusTemporaryRedirect)
}

func logout(res http.ResponseWriter, req *http.Request){
	mycookie, err := req.Cookie("sessionfino")
	var currentUser user
	if err != nil {
		mycookie = &http.Cookie{
			Name: "sessionfino",
			HttpOnly: true,
			//Secure: true,
		}
		currentUser = user{
			Login: false,
		}
		str := toJSON64(currentUser)
		// uuid | login status | json data | hmac + 64 bit encoding of json data
		mycookie.Value = uuid.NewV4().String() + "|" + strconv.FormatBool(currentUser.Login) +"|"+ str +"|" + setKey(str)
	}else {
		userFino := strings.Split(mycookie.Value, "|")
		if len(userFino) > 3 {
			if check(userFino[2], userFino[3]) {
				currentUser = getUser(userFino[2])
				if currentUser.Login == false{
					http.Redirect(res,req,"/",http.StatusTemporaryRedirect)
				}
				currentUser.Login = false;
			}
			str := toJSON64(currentUser)
			// uuid | login status | json data | hmac + 64 bit encoding of json data
			mycookie.Value = userFino[0] + "|" + strconv.FormatBool(currentUser.Login) + "|" + str + "|" + setKey(str)
		}else{
			currentUser = user{
				Login: false,
			}
			str := toJSON64(currentUser)
			// uuid | login status | json data | hmac + 64 bit encoding of json data
			mycookie.Value = uuid.NewV4().String() + "|" + strconv.FormatBool(currentUser.Login) +"|"+ str +"|" + setKey(str)
		}

	}
	http.SetCookie(res, mycookie)
	http.Redirect(res,req,"/",http.StatusTemporaryRedirect)
}

func view(res http.ResponseWriter, req *http.Request){
	var currentUser user
	var mydata data
	mycookie, err := req.Cookie("sessionfino")
	if err != nil{
		log.Println(err)
		mycookie = &http.Cookie{
			Name: "sessionfino",
			HttpOnly: true,
			//Secure: true,
		}
	}

	userFino := strings.Split(mycookie.Value,"|")
	if len(userFino) < 3{
		currentUser = user{
			Login: false,
		}
		str := toJSON64(currentUser)
		// uuid | login status | json data | hmac + 64 bit encoding of json data
		mycookie.Value = uuid.NewV4().String() + "|" + strconv.FormatBool(currentUser.Login) +"|"+ str +"|" + setKey(str)
		http.SetCookie(res, mycookie)
		userFino = strings.Split(mycookie.Value,"|")
	}else {
		currentUser = getUser(userFino[2])
	}

	if check(userFino[2],userFino[3]){
		if currentUser.Login == true{
			mydata = data{
				CookieValue: "name: "+ currentUser.Name + " age: "+currentUser.Age,
				Good: true,
				Login: currentUser.Login,
				Logout: !currentUser.Login,
			}
		}else{
			mydata = data{
				CookieValue: "you are logged out",
				Good: true,
				Login: currentUser.Login,
				Logout: !currentUser.Login,
			}
		}


	}else{
		mydata = data{
			CookieValue: "cookie was tampered with",
			Good: false,
			Login: false,
			Logout: true,
		}
	}
	err = myTemplates.ExecuteTemplate(res,"show.gohtml",mydata)
	if err != nil{
		log.Println(err)
	}

}

func main(){
	http.Handle("/favicon.ico",http.NotFoundHandler())
	http.HandleFunc("/",homepage)
	http.HandleFunc("/show",show)
	http.HandleFunc("/tampering",tamper)
	http.HandleFunc("/login",login)
	http.HandleFunc("/logout",logout)
	http.HandleFunc("/form",form)
	http.HandleFunc("/view",view)
	http.ListenAndServe("localhost:8080",nil)
}


package main

import(
	"net/http"
	"html/template"
	"log"
	"strconv"

	"github.com/satori/go.uuid"
	"time"
	"crypto/hmac"
	"crypto/sha256"
	"io"
	"strings"
	"fmt"
)

var mytemplates *template.Template
var datastuff map[string] string
var num_users int
var hmac_key []byte

func init() {
	var err error
	mytemplates,err = template.ParseGlob("*.gohtml")
	if(err != nil){
		log.Println(err)
	}
	num_users = 0
	datastuff = make(map[string]string,1)
	hmac_key = []byte("asl3ioaqwlimb")
}

func check(key string) bool{
	string_parts := strings.Split(key,"|")
	// in this case there aren't even values to check
	if(len(string_parts) < 2){
		return false
	}

	hm := hmac.New(sha256.New,hmac_key)
	io.WriteString(hm,string_parts[0])
	if(fmt.Sprintf("%x",hm.Sum(nil)) == string_parts[1]){
		return true
	}
	return false
}

func setKey(key string) string{
	hm := hmac.New(sha256.New,hmac_key)
	io.WriteString(hm,key)
	return fmt.Sprintf("%x",hm.Sum(nil))
}

func cookies(res http.ResponseWriter, req *http.Request){

	mycookie, err := req.Cookie("session-id")
	if(err != nil){
		log.Println("creating a cookie")
			myuuid := uuid.NewV4()
			mycookie = &http.Cookie{
			Name: "session-id",
			Value: "0",
			HttpOnly: true,
			Expires: time.Now().AddDate(0,1,1),
			//Secure: true,
		}
		mycookie.Value = myuuid.String() + "|" + setKey(myuuid.String())
		http.SetCookie(res,mycookie)
		num_users++
		datastuff[mycookie.Value] = strconv.Itoa(num_users)

	}

	if(check(mycookie.Value)) {
		err = mytemplates.ExecuteTemplate(res, "gotemplate.gohtml", "you are number " + datastuff[mycookie.Value])
	}else{
		err = mytemplates.ExecuteTemplate(res, "gotemplate.gohtml", "session-id was modified")
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
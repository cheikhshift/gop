package gop


import (
	"gopkg.in/mgo.v2/bson"
	"github.com/cheikhshift/db"
	"github.com/gorilla/sessions"
	"log"
	"errors"
	"crypto/sha256"
	"github.com/jtblin/go-ldap-client"
	"time"
	"strings"
)

type User struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	Username string `valid:"unique,required"`
	Pw [32]byte
	Email string `valid:"email,unique,required"`
    Created time.Time //timestamp local format
    Scopes []string
    Groups []string
    Props db.O
    LDAP bool
    Attr map[string] string
}

var (
	DB db.DB
	DbName,UPage string
	LDAPClient  *ldap.LDAPClient
	Zones []string
)

func UseLDAP(Ldap *ldap.LDAPClient){
	LDAPClient = Ldap
}

func SetUnAuthPage(path string) {
	UPage = path
}

func  AddAuthZone (path string){
	if Zones == nil {
		Zones = [] string{}
	}
	if !Contains(Zones, path){
	Zones = append(Zones, path)
	}
}

func (u User) SetPassword(pw string) {
	u.Pw = sha256.Sum256([]byte(pw))
}

func (u User) Push (ses *sessions.Session) error {
	ses.Values["passport"] = u
	if u.LDAP {
		return nil
	} else {
		return DB.Save(&u)
	}
}

func New(u User) User {
	return DB.New(&u).(User)
}

func Contains (arr []string, lookup string) bool {

	for _, v := range arr {
		if strings.Contains(lookup, v) {
			return true
		}
	}
	return false
}

func  GetUser(ses  *sessions.Session) (User, error) {
	if _, ok := ses.Values["password"]; ok  {
		return ses.Values["password"].(User), nil
	} else {
		return User{}, errors.New("No session found.")
	}
}

func  Logout(ses  *sessions.Session) (bool, error) {
	if _, ok := ses.Values["passport"]; ok  {
		delete(ses.Values, "passport")
		return true, nil
	} else {
		return false, errors.New("No session found.")
	}
}

func Join(args ...interface{}) (bool, error) {
	var email string
	var ses *sessions.Session
	if len(args) > 3 {
		email = args[2].(string)
		ses = args[3].(*sessions.Session)
	} else {
		ses = args[2].(*sessions.Session)
	}

	if _, ok := ses.Values["passport"]; !ok  {
		//delete(session.Values, "passport")
		user  := User{Username : args[0].(string)}
		user.SetPassword(args[1].(string))
		user.Email = email
		user = DB.New(&user).(User)
		err := DB.Save(&user)

		if err != nil {
			return false, err
		} else {
			ses.Values["passport"] = user
		}

		return true, nil
	} else {
		return false, errors.New("Already logged in!")
	}
}


func Login(u string, pw string, ses  *sessions.Session) (bool, error) {
	if _, ok := ses.Values["passport"]; !ok  {
		//delete(session.Values, "passport")
		user  := User{}
		err := DB.Query(user, bson.M{"username" : u , "pw" : sha256.Sum256([]byte(pw)) }).One(&user)

		if err != nil {
			return false, err
		} else if user.Username == "" {
			return false, errors.New("User not found!")
		}

		ses.Values["passport"] = user

		return true, nil
	} else {
		return false, errors.New("Already logged in!")
	}
}

func  LoginLDAP(u string, pw string, ses  *sessions.Session) (bool, error) {
	if _, ok := ses.Values["passport"]; !ok  {
		//delete(session.Values, "passport")

	ok, user, err := LDAPClient.Authenticate(u, pw)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, err
	}
	Suser := User{ Username: u , Attr : user}
	groups, err := LDAPClient.GetGroupsOfUser(u)
	if err != nil {
		return  false, err
	}
	Suser.Groups = groups
	Suser.LDAP  = true

	ses.Values["passport"] = Suser

		return true, nil
	} else {
		return false, errors.New("Already logged in!")
	}
}


func  RemoveAuthZone (path string){
	if Zones == nil {
		Zones = [] string{}
	}
	ns := []string{}
	for _,v := range Zones {
		if v != path {
			ns = append(ns, v)
		}
	}
	Zones = ns
}

func (u User) AddZone(path string){
	if u.Scopes == nil {
		u.Scopes = [] string{}
	}
	if !Contains(u.Scopes, path){
	u.Scopes = append(u.Scopes, path)
	}
}

func (u User) RemoveZone(path string){
	if u.Scopes == nil {
		u.Scopes = [] string{}
	}
	ns := []string{}
	for _,v := range u.Scopes {
		if v != path {
			ns = append(ns, v)
		}
	}
	u.Scopes = ns
}



func Connect(host string, dbname string){
	dbs,err := db.Connect(host, dbname) 
	if err != nil {
		log.Fatal(err)
	}
	DbName = dbname
	SetDb(dbs)
}


func SetDb(db  db.DB) {
	DB = db
}


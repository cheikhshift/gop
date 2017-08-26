package main 
import (
			"net/http"
			"time"
			"github.com/gorilla/sessions"
			"errors"
			"github.com/cheikhshift/db"
			"github.com/elazarl/go-bindata-assetfs"
			"bytes"
			"encoding/json"
			"fmt"
			"html"
			"html/template"
			"github.com/fatih/color"
			"strings"
			"reflect"
			"unsafe"
			"github.com/jtblin/go-ldap-client"
			"gopkg.in/mgo.v2"
			"gopkg.in/mgo.v2/bson"
		)
				var store = sessions.NewCookieStore([]byte("a very very very very secret key"))

				type NoStruct struct {
					/* emptystruct */
				}

				func net_sessionGet(key string,s *sessions.Session) string {
					return s.Values[key].(string)
				}


				func net_sessionDelete(s *sessions.Session) string {
						//keys := make([]string, len(s.Values))

						//i := 0
						for k := range s.Values {
						   // keys[i] = k.(string)
						    net_sessionRemove(k.(string), s)
						    //i++
						}

					return ""
				}

				func net_sessionRemove(key string,s *sessions.Session) string {
					delete(s.Values, key)
					return ""
				}
				func net_sessionKey(key string,s *sessions.Session) bool {					
				 if _, ok := s.Values[key]; ok {
					    //do something here
				 		return true
					}

					return false
				}

				 

			 func net_add(x,v float64) float64 {
					return v + x				   
			 }

			 func net_subs(x,v float64) float64 {
				   return v - x
			 }

			 func net_multiply(x,v float64) float64 {
				   return v * x
			 }

			 func net_divided(x,v float64) float64 {
				   return v/x
			 }

	

				func net_sessionGetInt(key string,s *sessions.Session) interface{} {
					return s.Values[key]
				}

				func net_sessionSet(key string, value string,s *sessions.Session) string {
					 s.Values[key] = value
					 return ""
				}
				func net_sessionSetInt(key string, value interface{},s *sessions.Session) string {
					 s.Values[key] = value
					 return ""
				}

				func net_importcss(s string) string {
					return "<link rel=\"stylesheet\" href=\"" + s + "\" /> "
				}

				func net_importjs(s string) string {
					return "<script type=\"text/javascript\" src=\"" + s + "\" ></script> "
				}



				func formval(s string, r*http.Request) string {
					return r.FormValue(s)
				}
			
				func renderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, p *Page)  bool {
				     filename :=  tmpl  + ".tmpl"
				    body, err := Asset(filename)
				    session, er := store.Get(r, "session-")

				 	if er != nil {
				           session,er = store.New(r,"session-")
				    }
				    p.Session = session
				    p.R = r
				    if err != nil {
				      // fmt.Print(err)
				    	return false
				    } else {
				    t := template.New("PageWrapper")
				    t = t.Funcs(template.FuncMap{"a":net_add,"s":net_subs,"m":net_multiply,"d":net_divided,"js" : net_importjs,"css" : net_importcss,"sd" : net_sessionDelete,"sr" : net_sessionRemove,"sc": net_sessionKey,"ss" : net_sessionSet,"sso": net_sessionSetInt,"sgo" : net_sessionGetInt,"sg" : net_sessionGet,"form" : formval,"eq": equalz, "neq" : nequalz, "lte" : netlt,"Add" : net_Add,"New" : net_New,"Update" : net_Update,"Delete" : net_Delete,"u" : net_u,"sf" : net_sf,"search" : net_search,"k" : net_k,"v" : net_v,"rmk" : net_rmk,"Q" : net_Q,"Limit" : net_Limit,"Skip" : net_Skip,"Sort" : net_Sort,"Count" : net_Count,"One" : net_One,"All" : net_All,"Close" : net_Close,"UserSpaceInterface" : net_UserSpaceInterface,"Sample" : net_structSample,"isSample" : net_castSample,"UserSpace" : net_structUserSpace,"isUserSpace" : net_castUserSpace})
				    t, _ = t.Parse(strings.Replace(strings.Replace(strings.Replace(BytesToString(body), "/{", "\"{",-1),"}/", "}\"",-1 ) ,"`", `\"` ,-1) )
				    outp := new(bytes.Buffer)
				    error := t.Execute(outp, p)
				    if error != nil {
				   // fmt.Print(error)
				    	 http.Redirect(w,r,"/your-500-page",301)
				    return false
				    }  else {
				    p.Session.Save(r, w)

				    fmt.Fprintf(w, html.UnescapeString(outp.String()) )
				    return true
					}
				    }
				}

				func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
				  return func(w http.ResponseWriter, r *http.Request) {
				  	if !apiAttempt(w,r) {
				      fn(w, r, "")
				  	}
				  }
				} 

				func mHandler(w http.ResponseWriter, r *http.Request) {
				  	
				  	if !apiAttempt(w,r) {
				      handler(w, r, "")
				  	}
				  
				} 
				func mResponse(v interface{}) string {
					data,_ := json.Marshal(&v)
					return string(data)
				}
				func apiAttempt(w http.ResponseWriter, r *http.Request) bool {
					session, er := store.Get(r, "session-")
					response := ""
					if er != nil {
						session,_ = store.New(r, "session-")
					}
					callmet := false

					 
				if  r.URL.Path == "/clear" && r.Method == strings.ToUpper("GET") { 
					
      	 response = "HelloWorld"
      
					callmet = true
				}
				 
				if   strings.Contains(r.URL.Path, "/")  { 
					
      		if Contains( Zones , r.URL.Path) {
      			//get user
      			user,err := GetUser(session)

      			if err != nil {
      				http.Redirect(w,r,UPage,302)
      				return true
      			}

      			if !Contains(user.Scopes, r.URL.Path) {
      				http.Redirect(w,r,UPage,302)
      				return true
      			}
      			//procceed
      		}
      
					
				}
				

					if callmet {
						session.Save(r,w)
						if response != "" {
							//Unmarshal json
							w.Header().Set("Access-Control-Allow-Origin", "*")
							w.Header().Set("Content-Type",  "application/json")
							w.Write([]byte(response))
						}
						return true
					}
					return false
				}
				func SetField(obj interface{}, name string, value interface{}) error {
					structValue := reflect.ValueOf(obj).Elem()
					structFieldValue := structValue.FieldByName(name)

					if !structFieldValue.IsValid() {
						return fmt.Errorf("No such field: %s in obj", name)
					}

					if !structFieldValue.CanSet() {
						return fmt.Errorf("Cannot set %s field value", name)
					}

					structFieldType := structFieldValue.Type()
					val := reflect.ValueOf(value)
					if structFieldType != val.Type() {
						invalidTypeError := errors.New("Provided value type didn't match obj field type")
						return invalidTypeError
					}

					structFieldValue.Set(val)
					return nil
				}
			func handler(w http.ResponseWriter, r *http.Request, context string) {
				  // fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
				  p,err := loadPage(r.URL.Path , context,r,w)
				  if err != nil {
				  	fmt.Println(err)
				        http.Redirect(w,r,"/your-404-page",307)
				        return
				  }

				   w.Header().Set("Cache-Control",  "public")
				  if !p.isResource {
				  		w.Header().Set("Content-Type",  "text/html")
				  		 defer func() {
					        if n := recover(); n != nil {
					           	color.Red("Error loading template in path : web" + r.URL.Path + ".tmpl reason :" )
					           	 fmt.Println(n)
					           	 http.Redirect(w,r,"/your-500-page",307)
					        }
					    }()
				      renderTemplate(w, r,  "web" + r.URL.Path, p)
				     
				     // fmt.Println(w)
				  } else {
				  		if strings.Contains(r.URL.Path, ".css") {
				  	  		w.Header().Add("Content-Type",  "text/css")
				  	  	} else if strings.Contains(r.URL.Path, ".js") {
				  	  		w.Header().Add("Content-Type",  "application/javascript")
				  	  	} else {
				  	  	w.Header().Add("Content-Type",  http.DetectContentType(p.Body))
				  	  	}
				  	 
				  	 
				      w.Write(p.Body)
				  }
				}

				func loadPage(title string, servlet string,r *http.Request,w http.ResponseWriter) (*Page,error) {
				    filename :=  "web" + title + ".tmpl"
				    if title == "/" {
				      http.Redirect(w,r,"/index",302)
				    }
				    body, err := Asset(filename)
				    if err != nil {
				      filename = "web" + title + ".html"
				      if title == "/" {
				    	filename = "web/index.html"
				    	}
				      body, err = Asset(filename)
				      if err != nil {
				         filename = "web" + title
				         body, err = Asset(filename)
				         if err != nil {
				            return nil, err
				         } else {
				          if strings.Contains(title, ".tmpl") || title == "/" {
				              return nil,nil
				          }
				          return &Page{Title: title, Body: body,isResource: true,request: nil}, nil
				         }
				      } else {
				         return &Page{Title: title, Body: body,isResource: true,request: nil}, nil
				      }
				    } 
				    //load custom struts
				    return &Page{Title: title, Body: body,isResource:false,request:r}, nil
				}
				func BytesToString(b []byte) string {
				    bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
				    sh := reflect.StringHeader{bh.Data, bh.Len}
				    return *(*string)(unsafe.Pointer(&sh))
				}
				func equalz(args ...interface{}) bool {
		    	    if args[0] == args[1] {
		        	return true;
				    }
				    return false;
				 }
				 func nequalz(args ...interface{}) bool {
				    if args[0] != args[1] {
				        return true;
				    }
				    return false;
				 }

				 func netlt(x,v float64) bool {
				    if x < v {
				        return true;
				    }
				    return false;
				 }
				 func netgt(x,v float64) bool {
				    if x > v {
				        return true;
				    }
				    return false;
				 }
				 func netlte(x,v float64) bool {
				    if x <= v {
				        return true;
				    }
				    return false;
				 }
				 func netgte(x,v float64) bool {
				    if x >= v {
				        return true;
				    }
				    return false;
				 }
				 type Page struct {
					    Title string
					    Body  []byte
					    request *http.Request
					    isResource bool
					    R *http.Request
					    Session *sessions.Session
					}
						var dbs db.DB
			func init(){
				 
	
			}
			type Sample struct {
	 		Id bson.ObjectId `bson:"_id,omitempty"`
			TestField string `valid:"unique"`
			FieldTwo string `valid:"email,unique,required"`
			Created time.Time //timestamp local format
		
			}

			func  net_castSample(args ...interface{}) *Sample  {
				
				s := Sample{}
				mapp := args[0].(db.O)
				if _, ok := mapp["_id"]; ok {
					mapp["Id"] = mapp["_id"]
				}
				data,_ := json.Marshal(&mapp)
				
				err := json.Unmarshal(data, &s) 
				if err != nil {
					fmt.Println(err.Error())
				}
				
				return &s
			}
			func net_structSample() *Sample{ return &Sample{} }
			type UserSpace struct {
				/* Property Type */
		
			}

			func  net_castUserSpace(args ...interface{}) *UserSpace  {
				
				s := UserSpace{}
				mapp := args[0].(db.O)
				if _, ok := mapp["_id"]; ok {
					mapp["Id"] = mapp["_id"]
				}
				data,_ := json.Marshal(&mapp)
				
				err := json.Unmarshal(data, &s) 
				if err != nil {
					fmt.Println(err.Error())
				}
				
				return &s
			}
			func net_structUserSpace() *UserSpace{ return &UserSpace{} }
			type UserSpaceInterface UserSpace
				func  net_UserSpaceInterface(args ...interface{}) (d UserSpace){
					if len(args) > 0 {
					jso := args[0].(string)
					var jsonBlob = []byte(jso)
					err := json.Unmarshal(jsonBlob, &d)
					if err != nil {
						fmt.Println("error:", err)
						return
					}
					return
					} else {
						d = UserSpace{} 
						return
					}
				}
						func net_Add(args ...interface{}) string {
							obj := args[0]
								
			//empty return assume it is string with return "" appended at end!
			err := dbs.Add(obj)
			
			if err != nil {

				return err.Error()
			}
			
		
						 return ""
						 
						}
						func net_New(args ...interface{}) interface{} {
							obj := args[0]
								
			//empty return assume it is string with return "" appended at end!
			 return dbs.New(obj)
			
			
		
						}
						func net_Update(args ...interface{}) string {
							obj := args[0]
								
			err := dbs.Save(obj)
			if err != nil {
				return err.Error()
			}
		
						 return ""
						 
						}
						func net_Delete(args ...interface{}) string {
							obj := args[0]
								
			err := dbs.Remove(obj)
			if err != nil {
				return err.Error()
			}
		
						 return ""
						 
						}
						func net_u(args ...interface{}) string {
							ke := args[0]
								val := args[1]
								obj := args[2]
								
				//args[0] = act
				reflect.ValueOf(obj).Elem().FieldByName(ke.(string)).Set(reflect.ValueOf(val))
		
						 return ""
						 
						}
						func net_sf(args ...interface{}) string {
							
				//args[0] = act
				reflect.ValueOf(&args[1]).Set(reflect.ValueOf(args[0]).Elem())
		
						 return ""
						 
						}
						func net_search(args ...interface{}) db.O {
							
			return db.O{}	
		
						}
						func net_k(args ...interface{}) string {
							ke := args[0]
								obj := args[1]
								m := args[2]
								
			
			m.(db.O)[ke.(string)] = obj
		
		
						 return ""
						 
						}
						func net_v(args ...interface{}) interface{} {
							ke := args[0]
								m := args[1]
								
			
			return m.(db.O)[ke.(string)] 
		
		
						}
						func net_rmk(args ...interface{}) db.O {
							ke := args[0]
								m := args[1]
								
			_map := m.(db.O)
			
			delete(_map, ke.(string))
			return _map
		
						}
						func net_Q(args ...interface{}) *mgo.Query {
							obj := args[0]
								req := args[1]
								
			return dbs.Q(obj).Find(req.(db.O))
		
						}
						func net_Limit(args ...interface{}) *mgo.Query {
							qry := args[0]
								val := args[1]
								
			return qry.(*mgo.Query).Limit(val.(int))
		
						}
						func net_Skip(args ...interface{}) *mgo.Query {
							qry := args[0]
								val := args[1]
								
			return  qry.(*mgo.Query).Skip(val.(int))
		
						}
						func net_Sort(args ...interface{}) *mgo.Query {
							qry := args[0]
								val := args[1]
								
			return  qry.(*mgo.Query).Sort(val.(string))
		
						}
						func net_Count(args ...interface{}) int {
							qry := args[0]
								
//with context is key to better code 
			cnt,err := qry.(*mgo.Query).Count()
			if err != nil {
				fmt.Println(err.Error())
				return -1
			}
			return cnt
		
						}
						func net_One(args ...interface{}) db.O {
							qry := args[0]
								
			one := db.O{}
			err := qry.(*mgo.Query).One(one)
			if err != nil {
				fmt.Println(err.Error())
			}
			//fmt.Println(one)
			return one
		
						}
						func net_All(args ...interface{}) []db.O {
							qry := args[0]
								
			one := []db.O{}
			err := qry.(*mgo.Query).All(&one)
			if err != nil {
				fmt.Println(err.Error())
			}
			return one
		
						}
						func net_Close(args ...interface{}) string {
							
			dbs.Close()
		
						 return ""
						 
						}
			func dummy_timer(){
				dg := time.Second *5
				fmt.Println(dg)
			}

			func main() {
				
	/* GOLANG SCOPE */
		AddAuthZone("/admin/super-user")
		AddAuthZone("/admin/guest-user")
		SetUnAuthPage("/unauth")
		UseLDAP(&ldap.LDAPClient{
		Base:         "dc=example,dc=com",
		Host:         "ldap.example.com",
		Port:         389,
		UseSSL:       false,
		BindDN:       "uid=readonlysuer,ou=People,dc=example,dc=com",
		BindPassword: "readonlypassword",
		UserFilter:   "(uid=%s)",
		GroupFilter: "(memberUid=%s)",
		Attributes:   []string{"givenName", "sn", "mail", "uid"},
		})

		defer LDAPClient.Close()

	 
	
	/* GOLANG SCOPE */
		
		
	
					 
					 fmt.Printf("Listenning on Port %v\n", "8080")
					 http.HandleFunc( "/",  makeHandler(handler))
					 http.Handle("/dist/",  http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: "web"}))
					 http.ListenAndServe(":8080", nil)
					 }
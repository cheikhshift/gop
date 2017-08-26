# Go-Passport

Go-Server authentication package. This package supports simple service security.

## Requirements

- Running Go lang workspace.
- [Go Server](https://github.com/cheikhshift/Gopher-Sauce/wiki)


## Index

1. [Mongo Authentication](#authenticating-with-mongodb)
2. [LDAP Authentication](#authenticating-with-ldap)
3. [Secure a service](#securing-a-service)

## Import 
Add this tag within the root of your `<gos>` tag in your `.gxml` file.

	<import src="github.com/cheikhshift/gop/gos.gxml" />


## Setup

Add the following `golang` statements within your `<main>` tag in your `.gxml` file. 

### Authenticating with MongoDB

#### Connect to database

	gop.Connect("host_mongo_uri", "database_name")
	defer gop.DB.Close()

#### Perform a login within `<end>` tag
  The following `<end>` tag will attempt to login a user :  If Authentication fails a text string as to why is returned.

    <end path="/login" type="POST" >
      	
      	succ, err := gop.Login(r.FormValue("username") ,
      		r.FormValue("password") ,
      		session )

      	if err != nil {
      		response = err.Error()
      	} else {
      			//redirect or return user
      			user,err := gop.GetUser(session)
      			//save in case you're redirecting
      			session.Save(r,w)

      	}

    </end>

#### Create a new user within `<end>` tag
  The following `<end>` tag will attempt to login a user :  If registration fails a text string as to why is returned.

	   <end path="/join" type="POST" >
	      	
	      succ, err := gop.Join(r.FormValue("username") ,
	      		r.FormValue("password") , 
	      		r.FormValue("email") ,
	      		session )
	
	      	if err != nil {
	      		response = err.Error()
	      	} else {
	      			//redirect or return user
	      			user,err := gop.GetUser(session)
	      			//save in case you're redirecting
	      			session.Save(r,w)
	
	      	}
	
	   </end>

	

#### Interface of Passport user : 

	type User struct {
		Id bson.ObjectId `bson:"_id,omitempty"`
		Username string `valid:"unique,required"`
		Pw [32]byte
		Email string `valid:"email,unique,required"`
	    Created time.Time //timestamp local format
	    Scopes []string
	    Attr map[string] string
	}

### Authenticating with LDAP

#### Import LDAP `pkg`

	<import src="github.com/jtblin/go-ldap-client" />

#### Connect to server
Create a connection to your LDAP server. Update the fields as needed.

		gop.UseLDAP(&ldap.LDAPClient{
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

#### Authentication with LDAP

  The following `<end>` tag will attempt to login a user :  If Authentication fails a text string as to why is returned. Once logged in, use the `pkg` function `gop.GetUser(session *sessions.Session) (User, error)` to get the current session's user interface.

    <end path="/login" type="POST" >
      	
      	succ, err := gop.LoginLDAP(r.FormValue("username") ,
      		r.FormValue("password") ,
      		session )

      	if err != nil {
      		response = err.Error()
      	} else {
      			//redirect or return user
      			user,err := gop.GetUser(session)
      			//save in case you're redirecting
      			session.Save(r,w)

      	}

    </end>


#### Interface of LDAP passport user :

	type User struct {
		Id bson.ObjectId `bson:"_id,omitempty"`
		Username string `valid:"unique,required"`
		Pw [32]byte
		Email string `valid:"email,unique,required"`
	    Created time.Time //timestamp local format
	    Scopes []string
	    Groups []string
	    Props db.O
	    Attr map[string] string
	}


### Securing a service
Use the package function `gop.AddAuthZone(path string)` to protect a request path and its subset paths. The following example will intercept any path with string `protect-resource` in it.

	gop.AddAuthZone("protect-resource")

### Give user path permission :
Use the function ` (u User)AddZone(path string)` to give user new access to a specified path. 

Use the function ` (u User)RemoveZone(path string)` to revoke user access from a specified path. 


### Update user

Use the `Push` function to update your user's session and database value.
`(u User) Push (ses *sessions.Session) error `

### Set Unauthorized page redirect

Use the `gop.SetUnAuthPage(path string)` to set the redirect path on permission error.
`gos.SetUnAuthPage(path string)`

### Issues :
Please use the Github issue tracker.

package ggsession
import (
  "encoding/gob"
  "net/http"
  "fmt"

  "github.com/gorilla/sessions"
  "github.com/cagox/gge/config"
)

//Store is the session store.
var Store *sessions.CookieStore

func init() {
  //authKeyOne := securecookie.GenerateRandomKey(64)
  //encryptionKeyOne := securecookie.GenerateRandomKey(32)

  authKeyOne := []byte(config.Config.AuthKey)
  encryptionKeyOne := []byte(config.Config.EncKey)

  Store = sessions.NewCookieStore(authKeyOne, encryptionKeyOne)

  Store.Options = &sessions.Options{  //Max Age 30 Days. This site is not exactly high risk.
    Path:   "/",
    MaxAge: 3600 * 24 * 30,
    HttpOnly: true,
  }

  //Register necessary structs.
  gob.Register(SessionData{})

}

//SessionData is a the struct to move data between the session cookie and the program.
type SessionData struct {
  UserID        uint
  Authenticated bool
  Flashes       []Flash
}

/*
Flash will be used to add flash messages to the session cookie.
Class indicates the type of message, and will be used for CSS purposes.
Message is the message itself.
*/
type Flash struct {
  Class    string
  Message  string
}

//BasePageData is the data that most pages will need. This can be used to build the data struct for templates.
type BasePageData struct {
  Page  string
  //More to Come as this thing takes form.
}


//NewSessionData returns a default SessionData struct.
func NewSessionData() SessionData {
  return SessionData{UserID: 0, Authenticated: false}
}

//GetSessionData grabs the SessionData struct from the cookie and returns it.
func GetSessionData(session *sessions.Session) SessionData {
  data := session.Values["sessiondata"]

  if data != nil {
    if page, ok := data.(SessionData); ok {
      //The cookie exists but is not ok.
      return page
    }
    //The cookie exists and is ok.
    return NewSessionData()
  }
  //The cookie doesn't exist.
  return  NewSessionData()
}

//GetSession returns the session with the open cookie file.
func GetSession(w http.ResponseWriter, r *http.Request) *sessions.Session {
  session, err := Store.Get(r, "gge-cookie")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return nil
  }
  if (session == nil) {
    session.Values["sessiondata"] = SessionData{UserID: 0, Authenticated: false}
    fmt.Println("Session was nil, now it is: ",session)
    return session
  }
  return session
}

//AddFlash wraps session.AddFlash() to more easily add flashes using the Flash struct.
func AddFlash(w http.ResponseWriter, r *http.Request, session *sessions.Session, class string, message string) {
  flash := Flash{Class: class, Message: message}
  session.AddFlash(flash)
  session.Save(r,w)
}

//AddFlash adds a flash message to the SessionData object
func (sessionData SessionData) AddFlash(class string, message string){
  flash := Flash{Class: class, Message: message}
  sessionData.Flashes = append(sessionData.Flashes, flash)
}

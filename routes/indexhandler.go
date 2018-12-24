package routes

import (
  "fmt"
  "net/http"
  "html/template"

  "github.com/cagox/gge/config"
  "github.com/cagox/gge/ggsession"
  "github.com/cagox/gge/models/user"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
  //Test for empty database.
  //If database is Empty, go to Admin User Creation Screen.
  users := user.GetUsers()
  if (len(users)==0){
    http.Redirect(w, r, "/admin/firstuser", http.StatusSeeOther)
    return
  }

  session := ggsession.GetSession(w, r)
  sessionData := ggsession.GetSessionData(session)
  fmt.Println("Index Handler Flashes: ", sessionData.Flashes)


  t := template.New("base.html")
  t, err := t.ParseFiles(config.Config.TemplateRoot+"/base/base.html")
  if err != nil {
    fmt.Println(err.Error())
  }

  session.Save(r,w)
  t.Execute(w, sessionData)
  return
}

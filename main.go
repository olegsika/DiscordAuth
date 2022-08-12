package main

import (
	"fmt"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"html/template"
	"log"
	"net/http"
)

const (
	discordKey    = "CLIENT ID"
	discordSecret = "SECRET KEY"
)

func main() {
	goth.UseProviders(
		discord.New(discordKey, discordSecret, "http://localhost:3001/auth/discord/callback", discord.ScopeIdentify, discord.ScopeEmail, discord.ScopeConnections),
	)

	http.HandleFunc("/", index)
	http.HandleFunc("/auth/discord", auth)
	http.HandleFunc("/auth/discord/callback", callback)
	http.HandleFunc("/logout/discord", logout)

	if err := http.ListenAndServe(":3001", nil); err != nil {
		log.Panic(err)
	}
}

func auth(w http.ResponseWriter, r *http.Request) {
	if user, err := gothic.CompleteUserAuth(w, r); err == nil {
		t, _ := template.New("foo").Parse(userTemplate)
		_ = t.Execute(w, user)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func callback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Println(w, err)
		return
	}
	t, _ := template.New("foo").Parse(userTemplate)
	_ = t.Execute(w, user)
}

func index(w http.ResponseWriter, _ *http.Request) {
	t, _ := template.New("foo").Parse(indexTemplate)
	_ = t.Execute(w, "text")
}

func logout(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

var indexTemplate = `
    <p><a href="/auth/discord?provider=discord">Log in with Discord</a></p>
`

var userTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
`

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"orc-server/orc"
	"time"
)

var (
	state       = "this is my random state"
	naverConfig = &oauth2.Config{
		ClientID:     "VQXXSgia6Oq1O97X6Hss",
		ClientSecret: "9YenMYeCl8",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://nid.naver.com/oauth2.0/authorize",
			TokenURL: "https://nid.naver.com/oauth2.0/token",
		},
		RedirectURL: "http://3.35.41.153/callback/naver",
	}
)

func handleLogin(w http.ResponseWriter, r *http.Request) {
	u := naverConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	s := r.FormValue("state")
	if s != state {
		fmt.Println("stats is different")
		w.WriteHeader(404)
		w.Write([]byte(http.StatusText(404)))
		return
	}

	code := r.FormValue("code")
	fmt.Println("code is :", code)
	httpClient := http.Client{Timeout: 2 * time.Second}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	token, err := naverConfig.Exchange(ctx, code)

	if err != nil {
		fmt.Println("Exchange from code to access_token err ", err)
		w.WriteHeader(404)
		w.Write([]byte(http.StatusText(404)))
		return
	}

	fmt.Println("access token is :", token)
	client := naverConfig.Client(ctx, token)
	req, err := http.NewRequest("GET", "https://openapi.naver.com/v1/nid/me", nil)
	res, err := client.Do(req)

	body, err := ioutil.ReadAll(res.Body)
	fmt.Println("body is ", string(body))

	tokenJson, err := json.Marshal(token)
	if err != nil {
		fmt.Println("Token marshal err")
		w.WriteHeader(404)
		w.Write([]byte(http.StatusText(404)))
		return
	}

	w.WriteHeader(200)
	w.Write(tokenJson)

	// fmt.Println(res)
}

func main() {
	fmt.Println("config info")
	fmt.Printf("%+v\n", orc.GlobalConfig)

	orc.StartGlobal()

	server := orc.GameServer{
		Port: orc.GlobalConfig.ServerPort,
	}
	go server.Run()
	fmt.Println("listen on port : ", orc.GlobalConfig.ServerPort)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback/naver", handleCallback)
	http.ListenAndServe(":80", nil)
}

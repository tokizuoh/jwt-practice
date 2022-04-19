package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type Model struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

var stubModel = &Model{
	Id:   1,
	Name: "Oka",
}

var secret = `hi-mi-tsu`
var globaKey *rsa.PrivateKey

func main() {
	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		t := jwt.New()
		t.Set(jwt.IssuerKey, `tokizuoh/jwt-practice`)
		t.Set(`private-claim-name`, secret)

		key, err := rsa.GenerateKey(rand.Reader, 1024)
		if err != nil {
			errMsg := fmt.Sprintf("ERROR: %s", err)
			http.Error(w, errMsg, http.StatusInternalServerError)
		}
		globaKey = key

		signed, err := jwt.Sign(t, jwt.WithKey(jwa.RS384, key))
		if err != nil {
			errMsg := fmt.Sprintf("ERROR: %s", err)
			http.Error(w, errMsg, http.StatusInternalServerError)
		}
		fmt.Fprintln(w, string(signed))
	})

	http.HandleFunc("/private", func(w http.ResponseWriter, r *http.Request) {
		if globaKey == nil {
			http.Error(w, "Unauthorized. Please access /token", http.StatusUnauthorized)
			return
		}

		h := r.Header.Get("Authorization")
		typ := "Bearer"
		arr := strings.Split(h, " ")
		if len(arr) != 2 && arr[0] != typ {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		rt := arr[1]

		t, err := jwt.Parse([]byte(rt), jwt.WithKey(jwa.RS384, globaKey))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		v, ok := t.Get(`private-claim-name`)
		if ok != true {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		if v == secret {
			json.NewEncoder(w).Encode(stubModel)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})

	http.ListenAndServe(":8080", nil)
}

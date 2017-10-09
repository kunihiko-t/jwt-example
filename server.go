package main

import (
	"flag"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"os"
	"time"
)

var mySigningKey []byte
var port *int

type IndexHandler struct {
}

type AuthHandler struct {
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tokenString := r.URL.Query().Get("token")

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	w.Header().Set("Content-Type", "text/plain")
	if err != nil {
		w.WriteHeader(401)
		fmt.Fprintf(w, err.Error())
		return
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		w.WriteHeader(200)
		w.Header().Set("JWT-Token", token.Raw)
		fmt.Fprintf(w, "%v %v %v", "Authorized", claims.Issuer, claims.ExpiresAt)
	} else {
		w.WriteHeader(401)
		fmt.Fprintf(w, "Invalid Token")
	}

}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Create the Claims
	now := time.Now().Unix()
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Unix(now, 0).Add(60 * time.Second).Unix(),
		Issuer:    "test",
		IssuedAt:  now,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	log.Printf("%v %v", ss, err)
	w.Header().Set("Content-Type", "text/plain")

	if err != nil {
		w.WriteHeader(401)
		fmt.Fprintf(w, err.Error())
		return
	}

	c := time.After(60 * time.Second)
	go func() {
		<-c
		log.Printf("Token %v has been expired", ss)
	}()

	w.Header().Set("JWT-Token", ss)
	w.WriteHeader(200)
	fmt.Fprintf(w, "Key Published\n")
	fmt.Fprintf(w, "http://localhost:%v/auth?token=%v", *port, ss)

}

func main() {

	port = flag.Int("p", 8080, "Port number")
	flag.Parse()
	mySigningKey = []byte("MyKey")

	ih := &IndexHandler{}
	http.Handle("/", ih)

	ah := &AuthHandler{}
	http.Handle("/auth", ah)

	addr := fmt.Sprintf(":%d", *port)
	err := http.ListenAndServe(addr, nil)
	checkError(err)
	log.Println("Start Listening")
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mmpg/api/engine"
	"github.com/satori/go.uuid"
)

// UserValidator tells whether a given email and password are valid
type UserValidator func(string, string) bool

// AuthToken represents the information contained in an authentication token
type AuthToken struct {
	Email    string
	Remember bool
}

var key []byte

func init() {
	key, _ := ioutil.ReadFile("key")

	if key == nil {
		key = []byte(uuid.NewV4().String())
		ioutil.WriteFile("key", key, 0660)
	}
}

// ValidateToken validates a request authorization token
func ValidateToken(r *http.Request) (authToken *AuthToken, err error) {
	tokenString := r.Header.Get("Authorization")

	// Validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return key, nil
	})

	if err != nil {
		return
	}

	if !token.Valid {
		err = errors.New("Invalid token")
		return
	}

	email, ok := token.Claims["email"].(string)

	if !ok {
		err = errors.New("Invalid email")
	}

	remember, ok := token.Claims["remember"].(bool)

	if !ok {
		remember = false
	}

	authToken = &AuthToken{
		Email:    email,
		Remember: remember,
	}

	return
}

// Auth handles authentication
func Auth(uv UserValidator) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			renewToken(w, r)
		case "POST":
			login(w, r, uv)
		default:
			w.WriteHeader(405)
		}
	}
}

type authMessage struct {
	Email    string `json:email`
	Password string `json:password`
	Remember bool   `json:remember`
}

func login(w http.ResponseWriter, r *http.Request, uv UserValidator) {
	decoder := json.NewDecoder(r.Body)
	var m authMessage

	if err := decoder.Decode(&m); err != nil {
		w.WriteHeader(403)
		return
	}

	if !uv(m.Email, m.Password) {
		w.WriteHeader(403)
		return
	}

	serveNewAuthToken(w, m.Email, m.Remember)
}

func renewToken(w http.ResponseWriter, r *http.Request) {
	token, err := ValidateToken(r)

	if err != nil {
		w.WriteHeader(400)
		return
	}

	serveNewAuthToken(w, token.Email, token.Remember)
}

func serveNewAuthToken(w http.ResponseWriter, email string, remember bool) {
	res, err := engine.PlayerExists(email)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	if res != "TRUE" {
		w.WriteHeader(400)
		return
	}

	token, err := createToken(email, remember).SignedString(key)

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Write([]byte(token))
}

func createToken(email string, remember bool) *jwt.Token {
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims["email"] = email
	token.Claims["remember"] = remember

	var duration time.Duration

	if remember {
		duration = 24 * time.Hour * 7
	} else {
		duration = time.Hour
	}

	token.Claims["exp"] = time.Now().Add(duration).Unix()

	return token
}

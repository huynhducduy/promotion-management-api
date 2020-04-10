package app

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"promotion-management-api/internal/config"
	"promotion-management-api/pkg/utils"
)

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Claims struct {
	Id        int
	ExpiresAt int64
	jwt.StandardClaims
}

type Token struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

func getToken(r *http.Request) (string, error) {
	if r.Header["Authorization"] != nil && len(strings.Split(r.Header["Authorization"][0], " ")) == 2 {
		return strings.Split(r.Header["Authorization"][0], " ")[1], nil
	} else {
		return "", errors.New("No bearer token.")
	}
}

func isAuthenticated(endpoint func(http.ResponseWriter, *http.Request, User)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		user_token, err := getToken(r)
		if err != nil {
			utils.Response(w, http.StatusUnauthorized, err.Error())
			return
		}

		token, err := jwt.ParseWithClaims(user_token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetConfig().SECRET), nil
		})

		if err == nil && token.Valid {
			var user *User
			user, err = getOneUser(token.Claims.(*Claims).Id)
			if err != nil {
				utils.ResponseInternalError(w, err)
				return
			}

			endpoint(w, r, *user)
			return
		}

		utils.ResponseMessage(w, http.StatusUnauthorized, "Invalid token!")
	}
}

func generateToken(id int) Token {
	expAt := time.Now().Unix() + 604800 // 1 week

	payload := Claims{Id: id, ExpiresAt: expAt}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), payload)
	tokenString, _ := token.SignedString([]byte(config.SECRET))

	return Token{
		Token:     tokenString,
		ExpiresAt: expAt,
	}
}

func login(w http.ResponseWriter, r *http.Request) {

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	var credential Credential
	json.Unmarshal(reqBody, &credential)

	if credential.Username == "" || credential.Password == "" {
		utils.ResponseMessage(w, http.StatusBadRequest, "Username and password must not be empty!")
		return
	}

	var id int

	results := db.QueryRow("SELECT `id` FROM `users` where `username` = ? AND `password` = ?", credential.Username, credential.Password)
	err = results.Scan(&id)
	if err == sql.ErrNoRows {
		utils.ResponseMessage(w, http.StatusNotFound, "Username and passowrd is incorrect!")
		return
	} else if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	token := generateToken(id)

	utils.Response(w, http.StatusOK, token)
}

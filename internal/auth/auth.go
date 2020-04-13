package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"promotion-management-api/internal/employee"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"promotion-management-api/internal/config"
	"promotion-management-api/internal/db"
	"promotion-management-api/pkg/utils"
)

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Claims struct {
	Id        int64
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

func generateToken(id int64) Token {
	expAt := time.Now().Unix() + 604800 // 1 week

	payload := Claims{Id: id, ExpiresAt: expAt}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), payload)
	tokenString, _ := token.SignedString([]byte(config.GetConfig().SECRET))

	return Token{
		Token:     tokenString,
		ExpiresAt: expAt,
	}
}

func hashAndSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func comparePasswords(hashed string, plain string) bool {
	byteHashed := []byte(hashed)
	bytePlain := []byte(plain)
	err := bcrypt.CompareHashAndPassword(byteHashed, bytePlain)
	if err != nil {
		return false
	}
	return true
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		empToken, err := getToken(r)
		if err != nil {
			utils.ResponseMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		token, err := jwt.ParseWithClaims(empToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetConfig().SECRET), nil
		})

		if err == nil && token.Valid {
			var emp *employee.Employee
			emp, err = employee.Read(token.Claims.(*Claims).Id)
			if err != nil {
				utils.ResponseInternalError(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), "employee", emp)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		utils.ResponseMessage(w, http.StatusUnauthorized, "Invalid token!")
	})
}

func Login(w http.ResponseWriter, r *http.Request) {

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

	var id int64
	var pass string
	db := db.GetConnection()

	results := db.QueryRow("SELECT `ID`, `Password` FROM `employee` where `username` = ?", credential.Username)
	err = results.Scan(&id, &pass)
	if err == sql.ErrNoRows {
		utils.ResponseMessage(w, http.StatusNotFound, "Username and password is incorrect!")
		return
	} else if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	if !comparePasswords(pass, credential.Password) {
		utils.ResponseMessage(w, http.StatusNotFound, "Username and password is incorrect!")
		return
	}

	token := generateToken(id)

	utils.Response(w, http.StatusOK, token)
}

func GetPwd(w http.ResponseWriter, r *http.Request) {
	utils.ResponseMessage(w, 200, hashAndSalt("password123"))
}

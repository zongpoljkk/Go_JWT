package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/zongpoljkk/go_jwt/models"
	"github.com/zongpoljkk/go_jwt/repository/userRepository"
	"github.com/zongpoljkk/go_jwt/utils"
	"golang.org/x/crypto/bcrypt"
)

func (c Controller) Signup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var error models.Error
		json.NewDecoder(r.Body).Decode(&user)

		if user.Email == "" {
			error.Message = "Email is missing"
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		if user.Password == "" {
			error.Message = "Password is missing"
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("pass text", user.Password)
		fmt.Println("hash ", hash)

		user.Password = string(hash)

		userRepo := userRepository.UserRepository{}
		user = userRepo.Signup(db, user)

		if err != nil {
			error.Message = "Server Error"
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		user.Password = ""
		w.Header().Set("Content-Type", "application/json")

		utils.ResponseJSON(w, user)
	}
}

func (c Controller) Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var jwt models.JWT
		var error models.Error

		json.NewDecoder(r.Body).Decode(&user)

		if user.Email == "" {
			error.Message = "Email is missing"
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		if user.Password == "" {
			error.Message = "Password is missing"
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		password := user.Password

		userRepo := userRepository.UserRepository{}
		user, err := userRepo.Login(db, user)

		if err != nil {
			if err == sql.ErrNoRows {
				error.Message = "The user does not exists"
				utils.RespondWithError(w, http.StatusBadRequest, error)
				return
			} else {
				log.Fatal(err)
			}
		}

		hashedPassword := user.Password

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			error.Message = "Invalid Password"
			utils.RespondWithError(w, http.StatusUnauthorized, error)
			return
		}

		token, err := utils.GenerateToken(user)

		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		jwt.Token = token

		utils.ResponseJSON(w, jwt)
	}
}

func (c Controller) TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var errorObject models.Error
		authHeader := r.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")

		fmt.Println(bearerToken)

		if len(bearerToken) == 2 {
			// extract auth token
			authToken := bearerToken[1]

			//check if the token is valid
			token, error := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}

				return []byte(os.Getenv("SECRET")), nil
			})

			if error != nil {
				errorObject.Message = error.Error()
				utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
			}

			// invoke function that getting called on
			if token.Valid {
				next.ServeHTTP(w, r)
				return
			} else {
				errorObject.Message = "Invalid Token"
				utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}
		}
	})
}

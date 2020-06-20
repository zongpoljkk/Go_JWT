package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/zongpoljkk/go_jwt/controllers"
	"github.com/zongpoljkk/go_jwt/driver"

	"github.com/gorilla/mux"
	"github.com/labstack/gommon/log"
	"github.com/subosito/gotenv"
)

var db *sql.DB

func init() {
	gotenv.Load()
}

func main() {

	db = driver.ConnectDB()
	controller := controllers.Controller{}

	router := mux.NewRouter()

	router.HandleFunc("/signup", controller.Signup(db)).Methods("POST")
	router.HandleFunc("/login", controller.Login(db)).Methods("POST")
	router.HandleFunc("/protected", controller.TokenVerifyMiddleWare(controller.ProtectedEndpoint())).Methods("GET")

	fmt.Println("Listening on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}

package main

import (
	"caringAPI/controller"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/proxy", controller.Proxy).Methods("GET")
	router.HandleFunc("/auth", controller.Authenticator).Methods("GET")
	router.HandleFunc("/user/profile", controller.GetUser).Methods("GET")

	http.Handle("/", router)
	fmt.Println("Connected to port 1234")
	log.Fatal(http.ListenAndServe(":1234", router))

}

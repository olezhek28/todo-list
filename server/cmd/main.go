package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/olezhek28/todo-list/router"
)

func main() {
	r := router.Router()

	fmt.Println("starting the server om port 9000...")

	log.Fatal(http.ListenAndServe(":9000", r))
}

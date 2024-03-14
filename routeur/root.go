package routeur

import (
	"fmt"
	"groupie/controller"
	"log"
	"net/http"
)

func InitServe() {

	FileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", FileServer))
	http.HandleFunc("/", controller.IndexHandler)

	if err := http.ListenAndServe(controller.Port, nil); err != nil {

		fmt.Printf("ERREUR LORS DE L'INITIATION DES ROUTES %v \n", err)

		log.Fatal(err)

	}
}

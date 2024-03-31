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
	http.HandleFunc("/index", controller.IndexHandler)
	http.HandleFunc("/result_search", controller.SearchAnimeHandler)
	http.HandleFunc("/anime_detail/", controller.AnimeDetailHandler)
	http.HandleFunc("/genres/", controller.GenresHandler)
	http.HandleFunc("/animes_by_genre/", controller.AnimeByGenreHandler)
	http.HandleFunc("/add_favorite", controller.AddFavoriteHandler)
	http.HandleFunc("favorites", controller.FavoritesPageHandler)

	if err := http.ListenAndServe(controller.Port, nil); err != nil {

		fmt.Printf("ERREUR LORS DE L'INITIATION DES ROUTES %v \n", err)

		log.Fatal(err)

	}
}

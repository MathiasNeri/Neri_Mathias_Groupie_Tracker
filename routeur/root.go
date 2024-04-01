package routeur

import (
	"fmt"
	"groupie/controller"
	"log"
	"net/http"
)

func NotFoundHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Créez un ResponseWriter personnalisé pour capturer le statut
		customWriter := &statusResponseWriter{ResponseWriter: w, status: http.StatusOK}

		// Appelez le gestionnaire suivant avec notre ResponseWriter personnalisé
		next.ServeHTTP(customWriter, r)

		// Vérifiez si le statut 404 a été capturé
		if customWriter.status == http.StatusNotFound {
			// Redirigez vers la page d'accueil
			http.Redirect(w, r, "/", http.StatusFound)
		}
	})
}

// statusResponseWriter est une enveloppe autour de http.ResponseWriter qui nous permet de capturer le code de statut HTTP
type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func InitServe() {
	FileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", FileServer))
	http.Handle("/index", NotFoundHandler(http.HandlerFunc(controller.IndexHandler)))
	http.Handle("/result_search", NotFoundHandler(http.HandlerFunc(controller.SearchAnimeHandler)))
	http.Handle("/anime_detail/", NotFoundHandler(http.HandlerFunc(controller.AnimeDetailHandler)))
	http.Handle("/genres/", NotFoundHandler(http.HandlerFunc(controller.GenresHandler)))
	http.Handle("/animes_by_genre/", NotFoundHandler(http.HandlerFunc(controller.AnimeByGenreHandler)))
	http.Handle("/addFavorite/", NotFoundHandler(http.HandlerFunc(controller.AddFavoriteHandler)))
	http.Handle("/removeFavorite/", NotFoundHandler(http.HandlerFunc(controller.RemoveFavoriteHandler)))
	http.Handle("/favorites", NotFoundHandler(http.HandlerFunc(controller.FavoritesPageHandler)))
	http.Handle("/about", NotFoundHandler(http.HandlerFunc(controller.AboutHandler)))

	// Définition d'une route par défaut pour les chemins non trouvés
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			controller.NotFoundPageHandler(w, r)
		} else {
			controller.IndexHandler(w, r)
		}
	})

	if err := http.ListenAndServe(controller.Port, nil); err != nil {
		fmt.Printf("ERREUR LORS DE L'INITIATION DES ROUTES %v \n", err)
		log.Fatal(err)
	}
}

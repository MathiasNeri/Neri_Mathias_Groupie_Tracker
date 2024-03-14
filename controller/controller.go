package controller

import (
	"encoding/json"
	inittemplate "groupie/templates"
	"net/http"
	"strings"
)

const Port = "localhost:8080"

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	anime, err := getRandomAnime()
	if err != nil {
		http.Error(w, "Random anime probleme", http.StatusInternalServerError)
		return
	}
	inittemplate.Temp.ExecuteTemplate(w, "index.html", anime)
}

type AnimeInfo struct {
	Title  string `json:"title"`
	Images struct {
		Jpg struct {
			ImageURL string `json:"image_url"`
		} `json:"jpg"`
	} `json:"images"`
	Synopsis string `json:"synopsis"`
	Genres   []struct {
		Name string `json:"name"`
	} `json:"genres"`
}

// Méthode pour obtenir l'URL de l'image principale
func (ai *AnimeInfo) MainImageURL() string {
	return ai.Images.Jpg.ImageURL
}

func (ai *AnimeInfo) SynopsisOrDefault() string {
	if ai.Synopsis == "" {
		return "Pas de synopsis disponible."
	}
	return ai.Synopsis
}

// Fonction pour récupérer un anime aléatoire via l'API Jikan
func isExcludedGenre(anime *AnimeInfo, excludedGenres []string) bool {
	for _, genre := range anime.Genres {
		for _, excluded := range excludedGenres {
			if strings.EqualFold(genre.Name, excluded) {
				return true
			}
		}
	}
	return false
}

// Fonction pour récupérer un anime aléatoire, excluant certains genres
func getRandomAnime() (*AnimeInfo, error) {
	excludedGenres := []string{"Hentai"}
	for {
		resp, err := http.Get("https://api.jikan.moe/v4/random/anime")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var data struct {
			Data AnimeInfo `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, err
		}

		if !isExcludedGenre(&data.Data, excludedGenres) {
			return &data.Data, nil
		}
	}
}

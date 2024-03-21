package controller

import (
	"encoding/json"
	"fmt"
	inittemplate "groupie/templates"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const Port = "localhost:8080"

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	queryParams := r.URL.Query()

	// Get the value of the 'page' parameter
	page := queryParams.Get("page")

	anime, err := getRandomAnime()
	if err != nil {
		http.Error(w, "Random anime probleme", http.StatusInternalServerError)
		return
	}
	recentAnimes, err := getRecentAnimes(page)
	if err != nil {
		http.Error(w, "Failed to fetch recent animes", http.StatusInternalServerError)
		return
	}

	data := struct {
		RandomAnime  *AnimeInfo
		RecentAnimes *AnimeResponse
		Next         int
		Before       int
	}{
		RandomAnime:  anime,
		RecentAnimes: recentAnimes,
		Next:         recentAnimes.Pagination.CurrentPage + 1,
		Before:       recentAnimes.Pagination.CurrentPage - 1,
	}

	inittemplate.Temp.ExecuteTemplate(w, "index", data)
}

type AnimeInfo struct {
	Title  string `json:"title"`
	Images struct {
		Jpg struct {
			ImageURL      string `json:"image_url"`
			LargeImageURL string `json:"large_image_url"`
		} `json:"jpg"`
	} `json:"images"`
	Synopsis string `json:"synopsis"`
	Genres   []struct {
		Name string `json:"name"`
	} `json:"genres"`
	Type string `json:"type"`
}

// Méthode pour obtenir l'URL de l'image principale
func (ai *AnimeInfo) MainImageURL() string {
	return ai.Images.Jpg.ImageURL
}
func (ai *AnimeInfo) LargeImageURL() string {
	return ai.Images.Jpg.LargeImageURL
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

		// Vérifiez si le type est "TV" pour s'assurer que c'est bien un anime série TV
		if data.Data.Type == "TV" {
			return &data.Data, nil
		}
		// Sinon, continuez à chercher
	}
}

// Fonction pour récupérer et filtrer les animes récents, excluant certains genres
func getRecentAnimes(page string) (*AnimeResponse, error) {
	// Parse 'page' parameter as an integer
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		// Handle the error (e.g., invalid 'page' parameter)
		page = ""
	}

	// Check if 'page' parameter is a valid positive integer
	if pageNum <= 1 {
		page = ""
	}

	var url string
	if page == "" {
		url = "https://api.jikan.moe/v4/anime?order_by=start_date&sort=desc&limit=20"
	} else {
		url = "https://api.jikan.moe/v4/anime?order_by=start_date&sort=desc&limit=20&page=" + page
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response AnimeResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	var filteredAnimes []AnimeInfo
	for _, anime := range response.Data {
		if !isExcludedGenre(&anime, []string{"Hentai"}) {
			filteredAnimes = append(filteredAnimes, anime)
			if len(filteredAnimes) == 25 {
				break
			}
		}
	}

	response.Data = filteredAnimes
	return &response, nil
}

// Define a struct to hold both anime data and pagination info
type AnimeResponse struct {
	Pagination struct {
		LastVisiblePage int  `json:"last_visible_page"`
		HasNextPage     bool `json:"has_next_page"`
		CurrentPage     int  `json:"current_page"`
		Items           struct {
			Count   int `json:"count"`
			Total   int `json:"total"`
			PerPage int `json:"per_page"`
		} `json:"items"`
	} `json:"pagination"`
	Data []AnimeInfo `json:"data"`
}

type Pagination struct {
	LastVisiblePage int  `json:"last_visible_page"`
	HasNextPage     bool `json:"has_next_page"`
	CurrentPage     int  `json:"current_page"`
}

type AnimeSearchResult struct {
	Query      string
	Results    []AnimeInfo
	Pagination Pagination
}

func SearchAnimeHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	page := r.URL.Query().Get("page")
	if query == "" {
		http.Error(w, "Query is required", http.StatusBadRequest)
		return
	}
	if page == "" {
		page = "1" // Default to page 1 if no page number is provided
	}

	url := fmt.Sprintf("https://api.jikan.moe/v4/anime?q=%s&page=%s", url.QueryEscape(query), url.QueryEscape(page))
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	var result struct {
		Data       []AnimeInfo `json:"data"`
		Pagination Pagination  `json:"pagination"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		http.Error(w, "Error decoding response", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	data := AnimeSearchResult{
		Query:      query,
		Results:    result.Data,
		Pagination: result.Pagination,
	}

	inittemplate.Temp.ExecuteTemplate(w, "result_search", data)
}

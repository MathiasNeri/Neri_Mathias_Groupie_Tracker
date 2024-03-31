package controller

import (
	"encoding/json"
	"fmt"
	inittemplate "groupie/templates"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

const Port = "localhost:8080"

var (
	query string
)

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
	MalId  int    `json:"mal_id"`
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
		// Check if the anime type is "TV" and it is not an excluded genre
		if anime.Type == "TV" && !isExcludedGenre(&anime, []string{"Hentai"}) {
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
	Types      []string // Nouveau champ pour stocker les types sélectionnés
}
type AnimeDetail struct {
	MalID    int         `json:"mal_id"`
	URL      string      `json:"url"`
	Title    string      `json:"title"`
	Synopsis string      `json:"synopsis"`
	Images   AnimeImages `json:"images"`
	Score    float64     `json:"score"`
	Genres   []Genre     `json:"genres"`
	Aired    Aired       `json:"aired"`
	Episodes int         `json:"episodes"`
	Status   string      `json:"status"`
	Duration string      `json:"duration"`
	Rating   string      `json:"rating"`
}

type AnimeImages struct {
	Jpg struct {
		ImageURL      string `json:"image_url"`
		SmallImageURL string `json:"small_image_url"`
		LargeImageURL string `json:"large_image_url"`
	} `json:"jpg"`
}

type Genre struct {
	MalID int    `json:"mal_id"`
	Name  string `json:"name"`
}

type Aired struct {
	From   string `json:"from"`
	To     string `json:"to"`
	String string `json:"string"`
}

func SearchAnimeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("q") != "" {
		query = r.URL.Query().Get("q")
	}
	if query == "" {
		http.Error(w, "Query is required", http.StatusBadRequest)
		return
	}

	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	types := r.URL.Query()["type"]
	var wg sync.WaitGroup
	var results []AnimeInfo
	var firstPagination Pagination

	if len(types) == 0 {
		types = append(types, "") // Handle search without any specific type
	}

	for _, t := range types {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			searchURL := "https://api.jikan.moe/v4/anime?q=" + url.QueryEscape(query) + "&page=" + url.QueryEscape(pageStr) + "&sfw=true"
			if t != "" {
				searchURL += "&type=" + url.QueryEscape(t)
			}

			resp, err := http.Get(searchURL)
			if err != nil {
				fmt.Println("Error fetching data:", err)
				return
			}
			defer resp.Body.Close()

			var result struct {
				Data       []AnimeInfo `json:"data"`
				Pagination Pagination  `json:"pagination"`
			}

			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				fmt.Println("Error decoding response:", err)
				return
			}

			if firstPagination.CurrentPage == 0 || (firstPagination.LastVisiblePage < result.Pagination.LastVisiblePage) { // Assume all types have similar pagination for simplicity
				firstPagination = result.Pagination
			}

			results = append(results, result.Data...)
		}(t)
	}

	wg.Wait() // Wait for all go routines to complete

	nextPage := page
	if firstPagination.HasNextPage {
		nextPage++
	}
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	data := AnimeSearchResult{
		Query:      query,
		Results:    results,
		Pagination: firstPagination,
		Types:      types, // Passez les types sélectionnés au template
	}

	inittemplate.Temp.ExecuteTemplate(w, "result_search", data)
}
func AnimeDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Extrait l'animeID de l'URL
	animeID := strings.TrimPrefix(r.URL.Path, "/anime_detail/")

	if animeID == "" {
		http.NotFound(w, r)
		return
	}

	// Utilisez animeID pour obtenir les détails de l'anime de l'API
	url := fmt.Sprintf("https://api.jikan.moe/v4/anime/%s", animeID)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "Failed to fetch anime details", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var detailResponse struct {
		Data AnimeDetail `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&detailResponse); err != nil {
		http.Error(w, "Failed to decode anime details", http.StatusInternalServerError)
		return
	}

	// Passez detailResponse.Data au template
	inittemplate.Temp.ExecuteTemplate(w, "anime_detail", detailResponse.Data)
}

func GenresHandler(w http.ResponseWriter, r *http.Request) {
	// Faire une requête à l'API pour obtenir les genres
	resp, err := http.Get("https://api.jikan.moe/v4/genres/anime?sfw=true")
	if err != nil {
		// Gérer l'erreur
		http.Error(w, "Erreur lors de la récupération des genres", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var genresResponse struct {
		Data []struct {
			ID   int    `json:"mal_id"`
			Name string `json:"name"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&genresResponse); err != nil {
		// Gérer l'erreur
		http.Error(w, "Erreur lors du décodage de la réponse", http.StatusInternalServerError)
		return
	}

	// Rendre une page avec la liste des genres
	inittemplate.Temp.ExecuteTemplate(w, "genres", genresResponse.Data)
}
func AnimeByGenreHandler(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")
	if len(pathSegments) < 3 || pathSegments[2] == "" {
		http.Error(w, "Genre ID is required", http.StatusBadRequest)
		return
	}
	genreID := pathSegments[2] // L'index 2 car il doit être après "/animes_by_genre/"

	// Récupérer le numéro de page depuis les paramètres de requête, avec une valeur par défaut de 1
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	// Construire l'URL de l'API pour requêter les animes du genre spécifié
	apiURL := fmt.Sprintf("https://api.jikan.moe/v4/anime/genres/%s/?page=%s&sfw=true", url.QueryEscape(genreID), url.QueryEscape(page))

	// Effectuer la requête à l'API
	resp, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		fmt.Println("Error fetching anime by genre:", err)
		return
	}
	defer resp.Body.Close()

	// Décoder la réponse
	var result struct {
		Data       []AnimeInfo `json:"data"`
		Pagination Pagination  `json:"pagination"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		http.Error(w, "Error decoding API response", http.StatusInternalServerError)
		return
	}

	// Préparation des données pour le template
	data := struct {
		GenreID    string
		Animes     []AnimeInfo
		Pagination Pagination
	}{
		GenreID:    genreID,
		Animes:     result.Data,
		Pagination: result.Pagination,
	}
	inittemplate.Temp.ExecuteTemplate(w, "animes_by_genre", data)
}

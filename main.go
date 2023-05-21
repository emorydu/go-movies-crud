package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func main() {
	r := mux.NewRouter()

	movies = append(movies,
		[]Movie{{
			ID:    "1",
			Isbn:  "428227",
			Title: "Movie One",
			Director: &Director{
				Firstname: "John",
				Lastname:  "Doe",
			},
		}, {
			ID:    "2",
			Isbn:  "454557",
			Title: "Movie Two",
			Director: &Director{
				Firstname: "Steve",
				Lastname:  "Smith",
			},
		}}...,
	)

	//  Registered Routing
	r.HandleFunc("/movies", getMovies).Methods(http.MethodGet)           // GET
	r.HandleFunc("/movies/{id}", getMovie).Methods(http.MethodGet)       // GET
	r.HandleFunc("/movies", createMovie).Methods(http.MethodPost)        // POST
	r.HandleFunc("/movies/{id}", updateMovie).Methods(http.MethodPut)    // PUT
	r.HandleFunc("/movies/{id}", deleteMovie).Methods(http.MethodDelete) // DELETE

	log.Println("Starting server at port 8000")

	log.Fatal(http.ListenAndServe(":8000", r))
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, movie := range movies {
		if params["id"] == movie.ID {
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var movie Movie

	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		fmt.Printf("decode body fail: %v\n", err)
		return
	}
	movie.ID = strconv.Itoa(rand.Intn(math.MaxInt))
	movies = append(movies, movie)

	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var newMovie Movie
	_ = json.NewDecoder(r.Body).Decode(&newMovie)

	for i, movie := range movies {
		if movie.ID == params["id"] {
			newMovie.ID = movie.ID
			movies = append(movies[:i], movies[i+1:]...)
			movies = append(movies, newMovie)
			break
		}
	}
	json.NewEncoder(w).Encode(newMovie)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for i, movie := range movies {
		if params["id"] == movie.ID {
			movies = append(movies[:i], movies[i+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(movies)
}

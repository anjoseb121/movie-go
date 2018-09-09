package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	ErrorBackend = errors.New("Somenthing went wrong")
	ErrorAPIKey  = errors.New("API Key not found")
)

// Request from client
type Request struct {
	ID int `json:"id"`
}

// MovieDBResponse response
type MovieDBResponse struct {
	Movies []Movie `json:"results"`
}

// Movie exported
type Movie struct {
	Title       string `json:"title"`
	Description string `json:"overview"`
	Cover       string `json:"poster_path"`
	ReleaseDate string `json:"release_date"`
}

// Handler function, do the Lambda exec logic
func Handler(request Request) ([]Movie, error) {
	APIKey, ok := os.LookupEnv("API_KEY")
	if !ok {
		return []Movie{}, ErrorAPIKey
	}

	movieURL := fmt.Sprintf("https://api.themoviedb.org/3/discover/movie?api_key=%s", APIKey)

	client := &http.Client{}

	req, err := http.NewRequest("GET", movieURL, nil)
	if err != nil {
		return []Movie{}, ErrorBackend
	}

	if request.ID > 0 {
		q := req.URL.Query()
		q.Add("with_genres", strconv.Itoa(request.ID))
		req.URL.RawQuery = q.Encode()
	}

	resp, err := client.Do(req)
	if err != nil {
		return []Movie{}, ErrorBackend
	}
	defer resp.Body.Close()

	var data MovieDBResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return []Movie{}, ErrorBackend
	}

	return data.Movies, nil
}

func main() {
	lambda.Start(Handler)
}

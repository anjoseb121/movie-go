package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	ErrorAPIKey        = errors.New("API Key not found")
	ErrorMakingRequest = errors.New("Error creating new request")
	ErrorBody          = errors.New("Error parsing clients request params")
	ErrorAPIReq        = errors.New("Error making movie API request")
	ErrorMovieResponse = errors.New("Error parsing movie API response")
)

// ClientBody from client
type ClientBody struct {
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

func responseError(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       err.Error(),
	}, err
}

// Handler function, do the Lambda exec logic
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	APIKey, APIKeyFound := os.LookupEnv("API_KEY")
	if !APIKeyFound {
		return responseError(ErrorAPIKey)
	}

	// Init
	movieURL := fmt.Sprintf("https://api.themoviedb.org/3/discover/movie?api_key=%s", APIKey)
	client := &http.Client{}

	req, err := http.NewRequest("GET", movieURL, nil)
	if err != nil {
		return responseError(ErrorMakingRequest)
	}

	if request.Body != "" {
		// Parse client request body to struct
		var body ClientBody
		err = json.Unmarshal([]byte(request.Body), &body)
		if err != nil {
			return responseError(ErrorBody)
		}

		// If request id param is greater than 0 add genres
		if body.ID > 0 {
			q := req.URL.Query()
			q.Add("with_genres", strconv.Itoa(body.ID))
			req.URL.RawQuery = q.Encode()
		}
	}

	// Do MovieAPI request
	resp, err := client.Do(req)
	if err != nil {
		return responseError(ErrorAPIReq)
	}
	defer resp.Body.Close()

	// Parse movieAPI response to struct
	var data MovieDBResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return responseError(ErrorMovieResponse)
	}

	// Parse movie response to json
	movies, err := json.Marshal(data.Movies)
	if err != nil {
		return responseError(ErrorMovieResponse)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(movies),
	}, nil
}

func main() {
	lambda.Start(Handler)
}

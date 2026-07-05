package themoviedb

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"nooblabs.com/video-manager/internal/model"
)

const (
	urlAuth        = "https://api.themoviedb.org/3/authentication"
	urlSearchMovie = "https://api.themoviedb.org/3/search/movie?query=%s"
)

type searchMovieRawResponse struct {
	Page    int                    `json:"page"`
	Results []searchMovieRawResult `json:"results"`
}

type searchMovieRawResult struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type TheMovieDB struct {
	key string
}

func New(key string) (*TheMovieDB, error) {
	if len(strings.TrimSpace(key)) == 0 {
		return nil, fmt.Errorf("invalid key provided - %s", key)
	}

	return &TheMovieDB{
		key: key,
	}, nil
}

func (f *TheMovieDB) SearchMovie(query string) ([]model.Movie, error) {
	slog.Info("Searching Movie via movie db", "query", query)

	url := fmt.Sprintf(urlSearchMovie, query)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", f.key))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	slog.Info(string(body))

	var rawResponse searchMovieRawResponse
	err = json.Unmarshal(body, rawResponse)

	var moviesResponse []model.Movie
	for _, result := range rawResponse.Results {
		moviesResponse = append(moviesResponse, model.Movie{
			Id:    result.Id,
			Title: result.Title,
		})
	}
	return moviesResponse, nil
}

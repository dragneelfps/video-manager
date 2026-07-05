package fetcher

import (
	"fmt"

	"nooblabs.com/video-manager/internal/model"
)

type Type int

const (
	Unknown Type = iota
	TheMovieDB
)

func (t Type) String() string {
	switch t {
	case TheMovieDB:
		return "themoviedb"
	}
	return "uknown"
}

func (t *Type) Set(s string) error {
	switch s {
	case "themoviedb":
		*t = TheMovieDB
	default:
		return fmt.Errorf("must be one of: themoviedb")
	}
	return nil
}

func (t Type) Type() string {
	return "fetcher-type"
}

type Fetcher interface {
	SearchMovie(query string) ([]model.Movie, error)
}

package model

type Config struct {
	Api ApiConfig `yaml:"api"`
}

type ApiConfig struct {
	TheMovieDB *TheMovieDBConfig `yaml:"themoviedb"`
}

type TheMovieDBConfig struct {
	Key string `yaml:"key"`
}

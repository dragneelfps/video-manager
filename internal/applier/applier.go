package applier

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"nooblabs.com/video-manager/internal/fetcher"
)

var (
	regexSquareBrackets = regexp.MustCompile(`\[[^\]]+\]`)
	regexDash           = regexp.MustCompile(`\s*-\s`)
)

type Applier struct {
	fetchers map[fetcher.Type]fetcher.Fetcher
}

type ApplyRequest struct {
	FilePaths   []string
	FetcherType fetcher.Type
}

func New(fetchers map[fetcher.Type]fetcher.Fetcher) (*Applier, error) {
	return &Applier{
		fetchers: fetchers,
	}, nil
}

func (r *ApplyRequest) Validate(availableFetchers []fetcher.Type) error {
	if r == nil {
		return fmt.Errorf("nil request provided")
	}
	if len(r.FilePaths) == 0 {
		return fmt.Errorf("empty file paths provided")
	}
	for _, path := range r.FilePaths {
		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			return fmt.Errorf("provided filepath(%s) does not exist", path)
		} else if err != nil {
			return fmt.Errorf("error checking filepath(%s) - %w", err)
		}
		if info.IsDir() {
			return fmt.Errorf("provided filepath(%s) is a directory")
		}
	}
	if !slices.Contains(availableFetchers, r.FetcherType) {
		return fmt.Errorf("fethcer(%s) not available, available fetchers - %v", r.FetcherType, availableFetchers)
	}
	return nil
}

func (a *Applier) Execute(req *ApplyRequest) error {
	if err := req.Validate(a.getAvailableFetchers()); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	slog.Info("executing applier")

	f := a.fetchers[req.FetcherType]

	for _, fp := range req.FilePaths {
		fileName := filepath.Base(fp)
		fileName = normalizeFileName(fileName)

		movies, err := f.SearchMovie(fileName)
		if err != nil {
			slog.Error("failed to fetch", "fetcher", req.FetcherType, "file_name", fileName)
			continue
		}
		slog.Info("search results", "fetcher", req.FetcherType, "movies", movies)
	}

	return nil
}

func (a *Applier) getAvailableFetchers() []fetcher.Type {
	var fetchers []fetcher.Type
	for fetcherType := range a.fetchers {
		fetchers = append(fetchers, fetcherType)
	}
	return fetchers
}

func normalizeFileName(name string) string {
	fileName := strings.TrimSuffix(name, filepath.Ext(name))
	fileName = regexSquareBrackets.ReplaceAllString(fileName, "")
	fileName = regexDash.ReplaceAllString(fileName, " ")
	return strings.TrimSpace(fileName)
}

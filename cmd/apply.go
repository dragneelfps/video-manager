package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"nooblabs.com/video-manager/internal/applier"
	"nooblabs.com/video-manager/internal/fetcher"
	"nooblabs.com/video-manager/internal/fetcher/themoviedb"
	"nooblabs.com/video-manager/internal/model"
)

var sources []string
var fetcherType fetcher.Type

var applyCmd = cobra.Command{
	Use: "apply",
	Run: func(cmd *cobra.Command, args []string) {
		fetchers := initFetchers(GlobalCfg)
		a, err := applier.New(fetchers)
		if err != nil {
			slog.Error("failed to initialize applier", "error", err)
			os.Exit(1)
		}

		err = a.Execute(&applier.ApplyRequest{
			FilePaths:   sources,
			FetcherType: fetcherType,
		})
		if err != nil {
			slog.Error("failed to execute applier", "error", err)
			os.Exit(1)
		}
	},
}

func init() {
	applyCmd.PersistentFlags().StringArrayVar(&sources, "sources", nil, "Source file paths")
	applyCmd.PersistentFlags().Var(&fetcherType, "fetcher", "metadata fetcher type")

	err := applyCmd.MarkPersistentFlagRequired("sources")
	cobra.CheckErr(err)
	err = applyCmd.MarkPersistentFlagRequired("fetcher")
	cobra.CheckErr(err)
}

func initFetchers(cfg model.Config) map[fetcher.Type]fetcher.Fetcher {
	fetchers := make(map[fetcher.Type]fetcher.Fetcher)
	if cfg.Api.TheMovieDB != nil {
		f, err := themoviedb.New(cfg.Api.TheMovieDB.Key)
		cobra.CheckErr(err)
		fetchers[fetcher.TheMovieDB] = f
	}
	return fetchers
}

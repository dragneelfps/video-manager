package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"nooblabs.com/video-manager/cmd/config"
	"nooblabs.com/video-manager/internal/model"
)

var GlobalCfg model.Config

var flagDryRun bool

var rootCmd = &cobra.Command{
	Use:   "video-manager",
	Short: "Help manage video files locally",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&flagDryRun, "dry-run", "d", false, "Run operations but do not make any changes to the video files")

	rootCmd.AddCommand(&applyCmd)
}

func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	configDir := filepath.Join(home, ".config", "video-manager")
	viper.AddConfigPath(configDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {

		cfg, err := config.LoadConfig(viper.GetViper())
		cobra.CheckErr(err)
		GlobalCfg = cfg

		// slog.Info("Loaded config", "config_file", viper.ConfigFileUsed(), "config", utils.StructToJsonString(config.Cfg))
	}

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

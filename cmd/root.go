package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "hodl",
		Short: "Hodl is a cli for fetching JSE share price data",
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .env)")
	rootCmd.PersistentFlags().StringP("apiKey", "k", "", "Alpha Vantage api key")
	rootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose logging")
	viper.BindPFlag("alpha_vantage_api_key", rootCmd.PersistentFlags().Lookup("apiKey"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.SetDefault("alpha_vantage_api_key", "demo")
	viper.SetDefault("verbose", false)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigFile(".env")
		viper.AutomaticEnv()
	}
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

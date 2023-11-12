package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(dailyCmd)
}

var dailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "Fetch daily historical OHLCV data for the provided ticker",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args[0])
		apiKey := viper.GetString("alpha_vantage_api_key")
		ticker := NewTickerProvider(resty.New(), apiKey)
		series, err := ticker.DailySeries(args[0])
		if err != nil {
			panic(err)
		}
		b, err := json.MarshalIndent(series, "", "	")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(b))
	},
}

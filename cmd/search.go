package cmd

import (
	"encoding/json"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Short: "Search for tickers based on the provided keywords",
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := viper.GetString("alpha_vantage_api_key")
		ticker := NewTickerProvider(resty.New(), apiKey)
		symbols, err := ticker.Search(args[0])
		if err != nil {
			panic(err)
		}

		b, err := json.MarshalIndent(symbols, "", "	")
		if err != nil {
			log.Println(err)
		}
		log.Println(string(b))
	},
}

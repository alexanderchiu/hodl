package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
	Use: "search",
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := viper.GetString("alpha_vantage_api_key")
		ticker := NewTickerProvider(resty.New(), apiKey)
		symbols, err := ticker.Search(args[0])
		if err != nil {
			panic(err)
		}

		b, err := json.MarshalIndent(symbols, "", "	")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(b))
	},
}

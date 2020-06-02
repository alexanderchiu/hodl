package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/mum4k/termdash/widgets/text"
	"github.com/spf13/cobra"
)

var (
	apiKey  string
	rootCmd = &cobra.Command{
		Use:   "hodl",
		Short: "Hodl is a cli for fetching JSE share price data",
		Run: func(cmd *cobra.Command, args []string) {
			// Submit query
			oracle := NewOracle(resty.New(), apiKey)
			series, err := oracle.DailySeries(args[0])
			if err != nil {
				panic(err)
			}

			// Init UI
			t, err := termbox.New()
			if err != nil {
				panic(err)
			}
			defer t.Close()
			ctx, cancel := context.WithCancel(context.Background())

			// Set up line graph and text box widgets
			data := series.Data[:20]
			latest := series.Data[0]
			x := make([]float64, len(data))
			labels := make(map[int]string)
			for i, d := range data {
				i = len(data) - i - 1
				x[i] = d.Indicators.Close
				labels[i] = d.Timestamp.Format("Jan02")
			}
			lc, err := linechart.New(
				linechart.AxesCellOpts(cell.FgColor(cell.ColorWhite)),
				linechart.YLabelCellOpts(cell.FgColor(cell.ColorWhite)),
				linechart.XLabelCellOpts(cell.FgColor(cell.ColorWhite)),
				linechart.YAxisAdaptive(),
			)
			if err != nil {
				panic(err)
			}
			err = lc.Series(args[0], x, linechart.SeriesXLabels(labels))
			if err != nil {
				panic(err)
			}
			text, err := text.New(text.DisableScrolling())
			if err != nil {
				panic(err)
			}
			err = text.Write(latest.Indicators.String())
			if err != nil {
				panic(err)
			}
			// Setup two row grid view
			g := grid.New()
			g.Add(
				grid.RowHeightPerc(70,
					grid.ColWidthPerc(90,
						grid.Widget(lc,
							container.Border(linestyle.Light),
							container.BorderTitle(fmt.Sprintf("%s Closing Prices", args[0])),
							container.FocusedColor(cell.ColorGreen),
							container.BorderTitleAlignCenter(),
							container.PlaceWidget(lc),
							container.BorderColor(cell.ColorWhite),
						),
					),
				),
				grid.RowHeightPerc(30,
					grid.ColWidthPerc(90,
						grid.Widget(text,
							container.Border(linestyle.Light),
							container.BorderTitle(fmt.Sprintf("%s Indicators %s", args[0], latest.Timestamp.Format("2006 Jan 02"))),
							container.FocusedColor(cell.ColorGreen),
							container.BorderTitleAlignCenter(),
							container.PlaceWidget(lc),
							container.BorderColor(cell.ColorWhite),
						),
					),
				),
			)

			quitter := func(k *terminalapi.Keyboard) {
				if k.Key == 'q' || k.Key == 'Q' {
					cancel()
				}
			}

			opt, err := g.Build()
			if err != nil {
				panic(err)
			}
			// Root container to hold grid
			c, err := container.New(t,
				container.ID("root"),
				container.BorderTitle("Press q to quit"))
			c.Update("root", opt...)
			if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quitter)); err != nil {
				panic(err)
			}
		},
	}
)

func init() {
	rootCmd.Flags().StringVarP(&apiKey, "apiKey", "k", "", "Alpha Vantage api key")
	rootCmd.MarkFlagRequired("apiKey")
}

// Execute executes the root command.
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

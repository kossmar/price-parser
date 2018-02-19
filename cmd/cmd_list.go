package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"sort"
)

var (
	CoinListCmd = &cobra.Command{
		Use:   "list",
		Short: "show coin list",
		RunE:  coinListCmd,
	}
)

func coinListCmd(cmd *cobra.Command, args []string) error {

	switch ApiFlag {
	case "hitbtc":
		url := "https://api.hitbtc.com/api/2/public/ticker"
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer resp.Body.Close()
		var requestInput []HitBTC
		body, err := ioutil.ReadAll(resp.Body)
		error1 := json.Unmarshal(body, &requestInput)
		if error1 != nil {
			fmt.Println(error1)
			return err
		}
		var nameListSlice []string
		for _, coin := range requestInput {
			nameListSlice = append(nameListSlice, coin.Symbol)
		}

		sort.Strings(nameListSlice)
		var i int

		for _, name := range nameListSlice {
			fmt.Print(name + "   ")
			if i%15 == 0 {
				fmt.Printf("\n")
			}
			i++
		}
	case "poloniex":
		url := "https://poloniex.com/public?command=returnTicker"
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer resp.Body.Close()
		var requestInput map[string]Poloniex
		body, err := ioutil.ReadAll(resp.Body)
		error1 := json.Unmarshal(body, &requestInput)
		if error1 != nil {
			fmt.Println(error1)
			return err
		}
		var i int
		for k, _ := range requestInput {
			fmt.Print(k + "   ")
			if i%15 == 0 {
				fmt.Printf("\n")
			}
			i++
		}
	}
	return nil
}

func init() {
	RootCmd.AddCommand(CoinListCmd)
	CoinListCmd.Flags().StringVar(&ApiFlag, "api", "poloniex", "specify api")
}

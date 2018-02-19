package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	// "net/http"
	"bytes"
	"sort"
)

var (
	CoinListCmd = &cobra.Command{
		Use:   "list",
		Short: "show coin list",
		RunE:  coinListCmd,
	}
)

func init() {
	RootCmd.AddCommand(CoinListCmd)
	CoinListCmd.Flags().StringVar(&ApiFlag, "api", "poloniex", "specify api")
}

func coinListCmd(cmd *cobra.Command, args []string) error {

	var outputVar bytes.Buffer

	switch ApiFlag {
	case "hitbtc":
		resp := getJson("https://api.hitbtc.com/api/2/public/ticker")
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

		formattedNameList := formatCoinList(nameListSlice)
		outputVar.WriteString(formattedNameList)

	case "poloniex":
		resp := getJson("https://poloniex.com/public?command=returnTicker")
		defer resp.Body.Close()
		var requestInput map[string]Poloniex
		body, err := ioutil.ReadAll(resp.Body)
		error1 := json.Unmarshal(body, &requestInput)
		if error1 != nil {
			fmt.Println(error1)
			return err
		}
		var nameListSlice []string
		for k, _ := range requestInput {
			nameListSlice = append(nameListSlice, k)
		}
		formattedNameList := formatCoinList(nameListSlice)
		outputVar.WriteString(formattedNameList)
	}

	fmt.Print(outputVar.String())
	return nil
}

func formatCoinList(nameList []string) string {
	var buffer bytes.Buffer
	var i int
	for _, name := range nameList {
		buffer.WriteString(name + "   ")
		if i%15 == 0 {
			buffer.WriteString("\n")
		}
		i++
	}
	return buffer.String()
}

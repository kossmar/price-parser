package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	// "strconv"
	"sort"
	"time"

	"github.com/fatih/structs"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Poloniex struct {
	Id            int    `json: "id"`
	Last          string `json: "last"`
	LowestAsk     string `json: "lowestAsk"`
	HighestBid    string `json: "highestBid"`
	PercentChange string `json: "percentChange"`
	BaseVolume    string `json: "baseVolume"`
	QuoteVolume   string `json: "quoteVolume"`
	IsFrozen      string `json: "isFrozen"`
	High24hr      string `json: "high24hr"`
	Low24hr       string `json: "low24hr"`
}

type HitBTC struct {
	Ask         string `json: "ask"`
	Bid         string `json: "bid"`
	Last        string `json: "last"`
	Open        string `json: "open"`
	Low         string `json: "low"`
	High        string `json: "high"`
	Volume      string `json: "volume"`
	VolumeQuote string `json: "volumeQuote"`
	Timestamp   string `json: "timestamp"`
	Symbol      string `json: "symbol"`
}

var (
	RootCmd = &cobra.Command{
		Use: "parser",
	}

	ParsePriceCmd = &cobra.Command{
		Use:   "parse",
		Short: "displays price information for various cryptocurrencies",
		RunE:  parsePriceCmd,
	}

	cfgFile      string
	coinString   string
	VerboseFlag  bool
	TimeFlag     bool
	JSONFlag     bool
	CoinNameFlag string
	ApiFlag      string
)

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.AddCommand(ParsePriceCmd)
	ParsePriceCmd.Flags().BoolVarP(&VerboseFlag, "verbose", "v", false, "verbose output")
	ParsePriceCmd.Flags().BoolVarP(&TimeFlag, "time", "T", false, "show the time between prints")
	ParsePriceCmd.Flags().BoolVarP(&JSONFlag, "json", "j", false, "print output in JSON format")
	ParsePriceCmd.Flags().StringVar(&CoinNameFlag, "coin", "USDT_BTC", "specify coin")
	ParsePriceCmd.Flags().StringVar(&ApiFlag, "api", "poloniex", "specify api")
	ParsePriceCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.price-parser.yaml)")
	ParsePriceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".price-parser")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	if ApiFlag == "hitbtc" {
		CoinNameFlag = "BTCUSD"
	}
}

func parsePriceCmd(cmd *cobra.Command, args []string) error {

	coinString = CoinNameFlag
	var currentCoin map[string]interface{}

	for {

		start := time.Now()
		time.Sleep(time.Second * 2)
		var outputVar bytes.Buffer
		outputVar.WriteString(coinString + "\n")

		_, err := setupOutputFile()
		if err != nil {
			return err
		}
		file, err := os.OpenFile("/Users/spinkringle/Documents/datazz", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		newWriter := bufio.NewWriter(file)

		switch ApiFlag {
		case "poloniex":
			resp, err := getJson("https://poloniex.com/public?command=returnTicker")
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			var requestInput map[string]Poloniex
			body, err := ioutil.ReadAll(resp.Body)
			error1 := json.Unmarshal(body, &requestInput)
			if error1 != nil {
				return err
			}
			s := structs.New(requestInput[coinString])
			m := s.Map()
			currentCoin = m

		case "hitbtc":
			resp, err := getJson("https://api.hitbtc.com/api/2/public/ticker")
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			var requestInput []HitBTC
			body, err := ioutil.ReadAll(resp.Body)
			error1 := json.Unmarshal(body, &requestInput)
			if error1 != nil {
				return err
			}
			coin := HitBTC{}
			for _, elem := range requestInput {
				if elem.Symbol == coinString {
					coin = elem
				}
			}
			s := structs.New(coin)
			m := s.Map()
			currentCoin = m
		}

		var jsonString string

		switch {
		case VerboseFlag && !JSONFlag:
			var keys []string
			for k := range currentCoin {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				jsonString += fmt.Sprint(k, ": ", currentCoin[k], "\n")
			}
		case VerboseFlag && JSONFlag:
			coin, err := json.Marshal(currentCoin)
			if err != nil {
				return err
			}
			jsonString += string(coin)
		case JSONFlag && !VerboseFlag:
			coin, err := json.Marshal(currentCoin["Last"])
			if err != nil {
				return err
			}
			jsonString += string(coin)
		default:
			jsonString += fmt.Sprint(currentCoin["Last"])
		}

		if TimeFlag == true {
			timeElapsed := (time.Since(start)).Seconds()
			jsonString += fmt.Sprintf("%.1f seconds\n", timeElapsed)
		}

		jsonString += ("\n")
		outputVar.WriteString(jsonString)
		fmt.Println(outputVar.String())
		err1 := outputToFile(outputVar.String(), newWriter)
		if err1 != nil {
			return err1
		}
	}
	return nil
}

func setupOutputFile() (*os.File, error) {
	file, err := os.OpenFile("/Users/spinkringle/Documents/datazz", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	return file, err
}

func outputToFile(output string, writer *bufio.Writer) error {
	_, err := writer.WriteString(output)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

func getJson(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, err
}

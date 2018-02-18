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
	ParsePriceCmd = &cobra.Command{
		Use:   "price-parser",
		Short: "displays price information for various cryptocurrencies",
		Long:  `...`,
		Run:   parsePriceCmd,
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
}

func parsePriceCmd(cmd *cobra.Command, args []string) {

	var url string
	coinString = CoinNameFlag
	var currentCoin map[string]interface{}

	for {

		start := time.Now()
		time.Sleep(time.Second * 2)
		var outputVar bytes.Buffer
		outputVar.WriteString(coinString + "\n")

		file, _ := os.OpenFile("/Users/spinkringle/Documents/datazz", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		newWriter := bufio.NewWriter(file)

		switch ApiFlag {
		case "poloniex":
			url = "https://poloniex.com/public?command=returnTicker"
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
				return
			}
			var requestInput map[string]Poloniex
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			error1 := json.Unmarshal(body, &requestInput)
			if error1 != nil {
				fmt.Println(error1)
				return
			}
			s := structs.New(requestInput[coinString])
			m := s.Map()
			currentCoin = m

		case "hitbtc":
			url = "https://api.hitbtc.com/api/2/public/ticker"
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
				return
			}
			var requestInput []HitBTC
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			error1 := json.Unmarshal(body, &requestInput)
			if error1 != nil {
				fmt.Println(error1)
				return
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

		switch {
		case VerboseFlag && !JSONFlag:
			var keys []string
			for k := range currentCoin {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				string := fmt.Sprint(k, ": ", currentCoin[k], "\n")
				outputVar.WriteString(string)
			}
		case VerboseFlag && JSONFlag:
			coin, err := json.Marshal(currentCoin)
			if err != nil {
				fmt.Println(err)
			}
			output := string(coin)
			outputVar.WriteString(output)
		case JSONFlag && !VerboseFlag:
			coin, err := json.Marshal(currentCoin["Last"])
			if err != nil {
				fmt.Println(err)
			}
			output := string(coin)
			outputVar.WriteString(output)
		default:
			output := fmt.Sprint(currentCoin["Last"])
			outputVar.WriteString(output)
		}

		if TimeFlag == true {
			timeElapsed := time.Since(start)
			time := timeElapsed.Seconds()
			output := fmt.Sprintf("%.1f seconds\n", time)
			outputVar.WriteString(output)
		}

		outputVar.WriteString("\n")

		fmt.Println(outputVar.String())
		outputToFile(outputVar.String(), newWriter)
	}
}

func setupOutputFile() (file *os.File) {
	file, err := os.OpenFile("/Users/spinkringle/Documents/datazz", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	return file
}

func outputToFile(output string, writer *bufio.Writer) {
	_, err := writer.WriteString(output)
	if err != nil {
		panic(err)
	}
	writer.Flush()
}

package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Coin struct {
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
)

func init() {
	cobra.OnInitialize(initConfig)
	ParsePriceCmd.Flags().BoolVarP(&VerboseFlag, "verbose", "v", false, "verbose output")
	ParsePriceCmd.Flags().BoolVarP(&TimeFlag, "time", "T", false, "show the time between prints")
	ParsePriceCmd.Flags().BoolVarP(&JSONFlag, "json", "j", false, "print output in JSON format")
	ParsePriceCmd.Flags().StringVar(&CoinNameFlag, "coin", "USDT_BTC", "specify coin")
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

	url := "https://poloniex.com/public?command=returnTicker"

	var requestInput map[string]Coin

	file, _ := os.OpenFile("/Users/spinkringle/Documents/datazz", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	newWriter := bufio.NewWriter(file)

	for {

		var outputVar bytes.Buffer

		start := time.Now()
		time.Sleep(time.Second * 5)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			return
		}

		unmarshalJSON(resp, &requestInput)
		coinString = CoinNameFlag
		outputVar.WriteString(coinString + "\n")

		switch {
		case VerboseFlag && !JSONFlag:
			verboseVal, verboseValues := verboseInfo(coinString, requestInput)
			for i := 0; i < verboseVal.NumField(); i++ {
				output := fmt.Sprint(verboseVal.Type().Field(i).Name, ": ", verboseValues[i], "\n")
				outputVar.WriteString(output)
			}
		case JSONFlag && !VerboseFlag:
			jsonString := jsonInfo(coinString, requestInput) + "\n"
			outputVar.WriteString(jsonString)
		case VerboseFlag && JSONFlag:
			verboseJSON := verboseJSONInfo(coinString, requestInput) + "\n"
			outputVar.WriteString(verboseJSON)
		default:
			defaultInfo := defaultInfo(coinString, requestInput)
			defaultInfoString := strconv.FormatFloat(defaultInfo, 'f', 6, 64) + "\n"
			outputVar.WriteString(defaultInfoString)
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

func unmarshalJSON(resp *http.Response, input *map[string]Coin) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	error1 := json.Unmarshal(body, &input)
	if error1 != nil {
		fmt.Println(err)
		return
	}
}

func defaultInfo(coin string, input map[string]Coin) float64 {
	coinName := input[coin]
	price, _ := strconv.ParseFloat(coinName.Last, 64)
	return price
}

func jsonInfo(coin string, input map[string]Coin) string {
	json, err := json.Marshal(input[coin].Last)
	if err != nil {
		fmt.Println(err)
	}
	jsonString := string(json)
	return jsonString
}

func verboseJSONInfo(coin string, input map[string]Coin) string {
	json, err := json.Marshal(input[coin])
	if err != nil {
		fmt.Println(err)
	}
	jsonString := string(json)
	return jsonString
}

func verboseInfo(coin string, input map[string]Coin) (reflect.Value, []interface{}) {

	coinName := input[coin]
	val := reflect.ValueOf(coinName)
	values := make([]interface{}, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		values[i] = val.Field(i).Interface()
	}
	return val, values
}

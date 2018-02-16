package cmd

import (
	"bufio"
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

var num = 5

var cfgFile string
var coinString string
var Verbose bool
var Time bool
var JSON bool
var CoinName string

var rootCmd = &cobra.Command{
	Use:   "price-parser",
	Short: "displays price information for various cryptocurrencies",
	Long:  `...`,
	Run: func(cmd *cobra.Command, args []string) {

		url := "https://poloniex.com/public?command=returnTicker"
		var requestInput map[string]Coin

		file, _ := os.OpenFile("/Users/spinkringle/Documents/datazz", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		// file, err := os.OpenFile("/Users/spinkringle/Documents/datazz", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		// if err != nil {
		// 	panic(err)
		// }
    //
		// defer file.Close()
    //
		// file.Sync()

		newWriter := bufio.NewWriter(file)

		for {
			start := time.Now()
			time.Sleep(time.Second * 5)

			resp, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
				return
			}

			UnmarshalJSON(resp, &requestInput)
			coinString = CoinName
			fmt.Println(coinString)
			OutputToFile(coinString + "\n", newWriter)
			
			switch {
			case Verbose && !JSON:
				verboseVal, verboseValues := VerboseInfo(coinString, requestInput)
				for i := 0; i < verboseVal.NumField(); i++ {
					output := fmt.Sprint(verboseVal.Type().Field(i).Name, ": ", verboseValues[i], "\n")
					fmt.Printf(output)
					OutputToFile(output, newWriter)
				}
			case JSON && !Verbose:
				jsonString := JSONInfo(coinString, requestInput) + "\n"
				fmt.Println(jsonString)
				OutputToFile(jsonString, newWriter)
			case Verbose && JSON:
				verboseJSON := VerboseJSONInfo(coinString, requestInput)
				fmt.Println(verboseJSON)
				OutputToFile(verboseJSON + "\n", newWriter)
			default:
				defaultInfo := DefaultInfo(coinString, requestInput)
				fmt.Printf("%.5f\n", defaultInfo)
				stringThing := (FloatToString(defaultInfo)) + "\n"
				OutputToFile(stringThing, newWriter)
			}

			if Time == true {
				time := ElapsedTime(start)
				output := fmt.Sprintf("%.1f seconds\n", time)
				fmt.Printf(output)
				OutputToFile(output, newWriter)
			}

			fmt.Print("\n")
			OutputToFile("\n", newWriter)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.Flags().BoolVarP(&Time, "time", "T", false, "show the time between prints")
	rootCmd.Flags().BoolVarP(&JSON, "json", "j", false, "print output in JSON format")
	rootCmd.Flags().StringVar(&CoinName, "coin", "USDT_BTC", "specify coin")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.price-parser.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// CUSTOM FUNCTIONS

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

func OutputToFile(output string, writer *bufio.Writer) {
	_, err := writer.WriteString(output)
	if err != nil {
		panic(err)
	}
	writer.Flush()
}

func UnmarshalJSON(resp *http.Response, input *map[string]Coin) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	error1 := json.Unmarshal(body, &input)
	if error1 != nil {
		fmt.Println(err)
		return
	}
}

func DefaultInfo(coin string, input map[string]Coin) float64 {
	coinName := input[coin]
	price, _ := strconv.ParseFloat(coinName.Last, 64)
	return price
}

func JSONInfo(coin string, input map[string]Coin) string {
	json, err := json.Marshal(input[coin].Last)
	if err != nil {
		fmt.Println(err)
	}
	jsonString := string(json)
	return jsonString
}

func VerboseJSONInfo(coin string, input map[string]Coin) string {
	json, err := json.Marshal(input[coin])
	if err != nil {
		fmt.Println(err)
	}
	jsonString := string(json)
	return jsonString
}

func VerboseInfo(coin string, input map[string]Coin) (reflect.Value, []interface{}) {

	coinName := input[coin]
	val := reflect.ValueOf(coinName)

	values := make([]interface{}, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		values[i] = val.Field(i).Interface()
	}
	// fmt.Println(val)
	// fmt.Println(values)

	return val, values
}

func ElapsedTime(start time.Time) float64 {
	timeElapsed := time.Since(start)
	time := timeElapsed.Seconds()
	return time
}

func FloatToString(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

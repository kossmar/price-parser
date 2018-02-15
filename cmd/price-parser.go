// Copyright © 2018 Nicholas Koss kossmar2@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"time"
	"net/http"
	"io/ioutil"
	"reflect"
	"encoding/json"
	"strconv"
	"bufio"
	"log"


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
	Long: `...`,
		Run: func(cmd *cobra.Command, args []string) {

			url := "https://poloniex.com/public?command=returnTicker"
			var requestInput map[string]Coin
			for {
				start := time.Now()
				time.Sleep(time.Second * 5)

				resp, err := http.Get(url)
				if err != nil {
					fmt.Println(err)
					return
				}
				// fmt.Println(resp)
				UnmarshalJSON(resp, &requestInput)
				coinString = CoinName
				fmt.Println(coinString)


				// f, err := os.Create("/Users/spinkringle/Documents/datazz")
				// if err != nil {
				// 	fmt.Println("you fucked up")
				// }

				f, err := os.OpenFile("/Users/spinkringle/Documents/datazz", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Fatal(err)
				}


				defer f.Close()

				f.Sync()
				w := bufio.NewWriter(f)

				_, err1 := w.WriteString(coinString + "\n")
				if err1 != nil {
					fmt.Println("you fucked up")
				}
				w.Flush()


				switch {
				case Verbose && !JSON:
					verboseVal, verboseValues := VerboseInfo(coinString, requestInput)
					for i := 0; i < verboseVal.NumField(); i++ {
					fmt.Print(verboseVal.Type().Field(i).Name, ": ", verboseValues[i], "\n")
					}
					fmt.Printf("\n\n")

				case JSON && !Verbose:
					jsonString := JSONInfo(coinString, requestInput)
					fmt.Println(jsonString)

				case Verbose && JSON:
					verboseJSON := VerboseJSONInfo(coinString, requestInput)
					fmt.Println(verboseJSON)

				default:
					defaultInfo := DefaultInfo(coinString, requestInput)
					fmt.Printf("%.5f\n", defaultInfo)

					stringThing := (FloatToString(defaultInfo)) + "\n\n"
					_, err := w.WriteString(stringThing)
					if err != nil {
						fmt.Println("you fucked up")
					}
					w.Flush()
				}

				if Time == true {
				time :=	ElapsedTime(start)
				fmt.Printf("%.1f seconds\n\n\n", time)
				}

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

func UnmarshalJSON(resp *http.Response, input *map[string]Coin) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	error1 := json.Unmarshal(body, &input)
	if error1 != nil {
		fmt.Println(err)
		return
	}
}

func DefaultInfo(coin string, input map[string]Coin) float64{
	coinName := input[coin]
	price, _ := strconv.ParseFloat(coinName.Last, 64)
	return price
}

func JSONInfo(coin string, input map[string]Coin) string{
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

func ElapsedTime(start time.Time) float64{
	timeElapsed := time.Since(start)
	time := timeElapsed.Seconds()
	return time
}

func FloatToString(input_num float64) string {
    return strconv.FormatFloat(input_num, 'f', 6, 64)
}

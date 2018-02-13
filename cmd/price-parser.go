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

var requestInput map[string]Coin
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

			for {
				start := time.Now()
				time.Sleep(time.Second * 5)

				unmarshalJSON(url)

				coinString = CoinName
				fmt.Println(coinString)

				switch {
				case Verbose && !JSON:
					verboseInfo(coinString)
				case JSON && !Verbose:
					JSONInfo(coinString)
				case Verbose && JSON:
					verboseJSONInfo(coinString)
				default:
					defaultInfo(coinString)
				}

				if Time == true {
					elapsedTime(start)
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

func unmarshalJSON(url string){
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	error1 := json.Unmarshal(body, &requestInput)
	if error1 != nil {
		fmt.Println(err)
		return
	}
}

func defaultInfo(coin string) {
	coinName := requestInput[coin]
	price, _ := strconv.ParseFloat(coinName.Last, 64)
	fmt.Printf("%.5f\n", price)
}

func JSONInfo(coin string) {
	jsonString, err := json.Marshal(requestInput[coin].Last)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jsonString))
}

func verboseJSONInfo(coin string) {
	jsonString, err := json.Marshal(requestInput[coin])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jsonString))
}

func verboseInfo(coin string) {

	coinName := requestInput[coin]
	val := reflect.ValueOf(coinName)

	values := make([]interface{}, val.NumField())

	fmt.Println("------ ", coin, " ------")

	for i := 0; i < val.NumField(); i++ {
		values[i] = val.Field(i).Interface()
		fmt.Print(val.Type().Field(i).Name, ": ", values[i], "\n")
	}

	fmt.Println("------ $$$$$$$$ ------ \n")
}

func elapsedTime(start time.Time) {
	timeElapsed := time.Since(start)
	time := timeElapsed.Seconds()
	fmt.Printf("%.1f seconds\n\n\n", time)
}

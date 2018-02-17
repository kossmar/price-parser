package cmd

import (
  "strconv"
  "encoding/json"
  "testing"
  // "github.com/kossmar/price-parser/cmd"
)

var input = map[string]Coin{"ETH_GNO":{Id:188, Last:"0.16920234", LowestAsk:"0.16920207", HighestBid:"0.16772051", PercentChange:"0.00395089", BaseVolume:"26.74316656", QuoteVolume:"158.52710200", IsFrozen:"0", High24hr:"0.17135812", Low24hr:"0.16565462"}, "BTC_NXT":{Id:69, Last:"0.00002284", LowestAsk:"0.00002286", HighestBid:"0.00002281", PercentChange:"0.01511111", BaseVolume:"94.57932181", QuoteVolume:"4133924.15000634", IsFrozen:"0", High24hr:"0.00002359", Low24hr:"0.00002209"}}

func TestDefaultInfo(t *testing.T) {
  defaultInfo := defaultInfo("ETH_GNO", input)
  exp := input["ETH_GNO"].Last
  expFloat, _ := strconv.ParseFloat(exp, 64)
  if defaultInfo != 0.16920234 {
    t.Errorf("Info incorrect, expected: %f got: %f", expFloat, defaultInfo)
  }
}

func TestVerboseInfo(t* testing.T) {
  verboseVal, verboseValues := VerboseInfo("ETH_GNO", input)
  expVal := []string{"Id", "Last", "LowestAsk", "HighestBid", "PercentChange", "BaseVolume", "QuoteVolume", "IsFrozen", "High24hr", "Low24hr"}
  expValues := []interface{}{188, "0.16920234", "0.16920207", "0.16772051", "0.00395089", "26.74316656", "158.52710200", "0", "0.17135812", "0.16565462"}

  for i := 0; i < len(expVal); i++ {
  if expVal[i] != verboseVal.Type().Field(i).Name {
    t.Errorf("Info incorrect, expected: %v got: %v ", expVal[i], verboseVal.Type().Field(i).Name)
    }
  }

  for i := 0; i < len(expValues); i++ {
      if expValues[i] != verboseValues[i] {
        t.Errorf("Info incorrect, expected: %v got: %v ", expValues[i], verboseValues[i])
    }
  }
}

func TestJSONInfo(t *testing.T) {
  JSONInfo := JSONInfo("ETH_GNO", input)
  exp, _ := json.Marshal(input["ETH_GNO"].Last)
  expString := string(exp)
  if JSONInfo != expString {
    t.Errorf("Info incorrect, expected: %s got: %s", exp, JSONInfo)
  }
}

func TestVerboseJSONInfo(t *testing.T) {
  verboseJSONInfo := VerboseJSONInfo("ETH_GNO", input)
  exp, _ := json.Marshal(input["ETH_GNO"])
  expString := string(exp)
  if verboseJSONInfo != expString {
    t.Errorf("Info incorrect, expected: %s got: %s", exp, verboseJSONInfo)
  }
}

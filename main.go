package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"collectd.org/api"
	"collectd.org/exec"
	"collectd.org/plugin"
)

var (
	progname = "ticker"
	conffile = "/etc/collectd/ticker.json"
	clicall  = strings.Contains(os.Args[0], progname)
	nerr     = -1.0
)

func fatalErrHandle(errMsg error) {
	if errMsg != nil {
		log.Fatal(errMsg)
	}
}

func tickerFetch(exchange string, url string, pricekey string) float64 {
	resp, getErr := http.Get(url)
	if getErr != nil {
		log.Print(getErr)
		return nerr
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Print(readErr)
		return nerr
	}

	res := make([]map[string]interface{}, 1)
	var jsonErr error

	switch exchange {
	case "coinmarketcap":
		jsonErr = json.Unmarshal(body, &res)
	default:
		jsonErr = json.Unmarshal(body, &res[0])
	}

	if jsonErr != nil {
		log.Print(jsonErr)
		return nerr
	}

	switch exchange {
	case "bitstamp", "hitbtc", "bitfinex", "binance", "coinmarketcap":
		last := res[0][pricekey]
		if last != nil {
			l, errConv := strconv.ParseFloat(last.(string), 64)
			if errConv == nil {
				return l
			}
			log.Print(errConv)
		}
	case "bittrex":
		if res[0]["result"] != nil {
			result := res[0]["result"].(map[string]interface{})
			if result != nil && result[pricekey] != nil {
				return result[pricekey].(float64)
			}
		}
	default:
		log.Fatal("Unsupported exchange")
	}

	return nerr
}

type Ticker struct{}

func (Ticker) Read() error {

	tickercf, readErr := ioutil.ReadFile(conffile)
	fatalErrHandle(readErr)

	j := make(map[string]interface{})
	fatalErrHandle(json.Unmarshal(tickercf, &j))

	// iterate through exchanges
	for k := range j {
		entry := j[k].(map[string]interface{})
		baseurl := entry["url"].(string)
		pairs := entry["pairs"]
		pricekey := entry["pricekey"].(string)

		convert, doconv := entry["convert"].(string)
		factor := 1.0
		if doconv {
			factor = tickerFetch(k, baseurl+convert, pricekey)
			if factor <= 0.0 {
				continue
			}
		}

		// iterate through pairs
		for _, c := range pairs.([]interface{}) {
			p := c.(string)
			url := baseurl + p

			l := tickerFetch(k, url, pricekey) * factor

			if l <= 0.0 {
				continue
			}

			p = strings.ToLower(strings.Replace(p, "-", "", -1))

			vl := api.ValueList{
				Identifier: api.Identifier{
					Host:           exec.Hostname(),
					Plugin:         progname,
					PluginInstance: p,
					Type:           "gauge",
				},
				Time:     time.Now(),
				Interval: 60 * time.Second,
				Values:   []api.Value{api.Gauge(l)},
			}

			if clicall {
				exec.Putval.Write(context.Background(), &vl)
			} else {
				if err := plugin.Write(&vl); err != nil {
					plugin.Error(err)
				}
			}
		}
	}

	return nil
}

func init() {

	switch clicall {
	case false:
		plugin.RegisterRead(progname, &Ticker{})
	default:
		if len(os.Args) > 1 {
			conffile = os.Args[1]
		}
		Ticker{}.Read()
	}
}

func main() {}

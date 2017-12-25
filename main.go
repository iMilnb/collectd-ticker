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

var progname = "ticker"
var conffile = "/etc/collectd/ticker.json"
var clicall = strings.Contains(os.Args[0], progname)

func errHandle(errMsg error) {
	if errMsg != nil {
		log.Fatal(errMsg)
	}
}

func tickerFetch(exchange string, url string) float64 {
	resp, getErr := http.Get(url)
	errHandle(getErr)

	body, readErr := ioutil.ReadAll(resp.Body)
	errHandle(readErr)

	res := make(map[string]interface{})

	errHandle(json.Unmarshal(body, &res))

	switch exchange {
	case "bitstamp", "hitbtc":
		if res["last"] != nil {
			l, errConv := strconv.ParseFloat(res["last"].(string), 64)
			errHandle(errConv)
			return l
		}
	case "bittrex":
		if res["result"] != nil {
			result := res["result"].(map[string]interface{})
			return result["Last"].(float64)
		}
	default:
		log.Fatal("Unsupported exchange")
	}

	return -1.0
}

type Ticker struct{}

func (Ticker) Read() error {

	tickercf, readErr := ioutil.ReadFile(conffile)
	errHandle(readErr)

	j := make(map[string]interface{})
	errHandle(json.Unmarshal(tickercf, &j))

	// iterate through exchanges
	for k := range j {
		entry := j[k].(map[string]interface{})
		baseurl := entry["url"].(string)
		pairs := entry["pairs"]

		convert, doconv := entry["convert"].(string)
		factor := 1.0
		if doconv {
			factor = tickerFetch(k, baseurl+convert)
		}

		// iterate through pairs
		for _, c := range pairs.([]interface{}) {
			p := c.(string)
			url := baseurl + p

			l := tickerFetch(k, url) * factor

			if l < 0.0 {
				return nil
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

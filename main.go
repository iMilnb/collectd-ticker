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
var clicall = strings.Contains(os.Args[0], progname)

func errHandle(errMsg error) {
	if errMsg != nil {
		log.Fatal(errMsg)
	}

}

func tickerFetch(url string) string {
	resp, getErr := http.Get(url)
	errHandle(getErr)

	body, readErr := ioutil.ReadAll(resp.Body)
	errHandle(readErr)

	res := make(map[string]string)

	errHandle(json.Unmarshal(body, &res))

	return res["last"]
}

type Ticker struct{}

func (Ticker) Read() error {

	var url string
	var v string
	var l float64
	var convErr error

	pairs := []string{"ethusd", "xrpusd", "ltcusd", "btcusd"}

	for _, v = range pairs {
		url = "https://www.bitstamp.net/api/v2/ticker/" + v + "/"
		l, convErr = strconv.ParseFloat(tickerFetch(url), 64)
		errHandle(convErr)

		vl := api.ValueList{
			Identifier: api.Identifier{
				Host:           exec.Hostname(),
				Plugin:         progname,
				PluginInstance: v,
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

	return nil
}

func init() {
	switch clicall {
	case false:
		plugin.RegisterRead(progname, &Ticker{})
	default:
		Ticker{}.Read()
	}
}

func main() {}

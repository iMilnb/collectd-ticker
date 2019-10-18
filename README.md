### ticker: a collectd plugin to watch your assets grow (hopefully)

Ticker is a [collectd][1] plugin written in [golang][2] that monitors cryptocurrencies pairs values on various exchanges in order to graph them using systems like [InfluxDB][7] and [Grafana][8].

![ticker screenshot](https://imil.net/stuff/ticker_screenshot.png "ticker in action")

#### Build

First, install [golang][2] version 1.7 or superior.

Then install [go-collectd][5].

Finally, install the `collectd-dev` package from your distribution. On _Debian_ and derivatives, a wrong patch is applied to `/usr/include/collectd/core/daemon/configfile.h` as mentioned in [this][3] bug report. In this file, change:

```
#include "collectd/core/config.h"
```

to

```
#include "collectd/liboconfig/oconfig.h"
```

In order to build the `ticker` [collectd][1] plugin, type

```
$ make plugin
```

To build the [Exec][4] version, type

```
$ make exec
```

To build both of them, simply type

```
$ make
```
 
You might need to modify [collectd][1] includes path in the `COLLECTD_SRC` variable from the `Makefile` depending on your distribution.

#### Installation

Use

```
$ make install
```

To install the plugin in [collectd][1] plugin directory as defined in the `Makefile`, or simply copy the `ticker.so` plugin to [collectd][1] plugin directory, for example in a _Debian_ based system

```
$ cp ticker.so /usr/lib/collectd/
```

And add the following to `collectd.conf`

```
Interval 60
LoadPlugin ticker
```

[collectd-go][5] does not support [collectd][1] configuration framework, so `ticker` has its own configuration file, which is expected to be located in `/etc/collectd` and named `ticker.json`. It is a `JSON` formatted file composed of

* Exchange name
* Exchange API base URL
* Cryptocurrencies pairs you'd like to monitor. The pairs must be in the format required by the associated exchange
* Name of the JSON key for currency "last price" in this exchange
* Optional: a conversion pair. For example, [Bittrex][6] does not offer `USD` pairs, but has a `BTC/USDT` pair, adding a conversion pair will multiply listed pairs with this factor.

Example

```
{
  "bitstamp":
    {
      "url": "https://www.bitstamp.net/api/v2/ticker/",
      "pairs": ["ethusd", "xrpusd", "ltcusd", "btcusd"],
      "pricekey": "last"
    },
  "bittrex":
    {
      "url": "https://bittrex.com/api/v1.1/public/getticker?market=",
      "pairs": ["BTC-XEM", "BTC-FUN", "BTC-XVG"],
      "convert": "USDT-BTC",
      "pricekey": "Last"
    },
  "hitbtc":
    {
      "url": "https://api.hitbtc.com/api/2/public/ticker/",
      "pairs": ["COSSBTC"],
      "convert": "BTCUSD",
      "pricekey": "last"
    },
  "bitfinex":
    {
      "url": "https://api.bitfinex.com/v1/pubticker/",
      "pairs": ["neousd", "btgusd"],
      "pricekey": "last_price"
    },
  "binance":
    {
      "url": "https://api.binance.com/api/v1/ticker/price?symbol=",
      "pairs": ["REQBTC", "LINKBTC"],
      "convert": "BTCUSDT",
      "pricekey": "price"
    },
  "coinmarketcap":
    {
      "url": "https://api.coinmarketcap.com/v1/ticker/",
      "pairs": ["kin", "electroneum"],
      "pricekey": "price_usd"
    }
}
```

You might want to test the plugin using its `Exec` version to see if everything works as expected.

```
$ ./ticker
PUTVAL "tatooine/ticker-ethusd/gauge" interval=60.000 1514193796.393:731.99
PUTVAL "tatooine/ticker-xrpusd/gauge" interval=60.000 1514193796.437:1.01117
PUTVAL "tatooine/ticker-ltcusd/gauge" interval=60.000 1514193796.488:278.78
PUTVAL "tatooine/ticker-btcusd/gauge" interval=60.000 1514193796.533:14212.55
PUTVAL "tatooine/ticker-btcxem/gauge" interval=60.000 1514193797.154:1.02831251000075
PUTVAL "tatooine/ticker-btcfun/gauge" interval=60.000 1514193797.361:0.0598963400000434
PUTVAL "tatooine/ticker-btcxvg/gauge" interval=60.000 1514193797.562:0.217641770000158
```

#### Exchanges

For now, `ticker` supports [Coinmarketcap][13], [Bitstamp][9], [Bittrex][6], [Bitfinex][11], [Binance][12] and [HitBTC][10].


[1]: https://collectd.org/
[2]: https://golang.org/
[3]: https://bugs.debian.org/cgi-bin/bugreport.cgi?bug=827759
[4]: https://collectd.org/wiki/index.php/Plugin:Exec
[5]: https://github.com/collectd/go-collectd
[6]: https://bittrex.com/
[7]: https://www.influxdata.com/
[8]: https://grafana.com/
[9]: https://www.bitstamp.net/
[10]: https://hitbtc.com/
[11]: https://www.bitfinex.com/
[12]: https://www.binance.com/
[13]: https://coinmarketcap.com/

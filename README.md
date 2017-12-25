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

[collectd-go][5] does not support [collectd][1] configuration framework, so it has its own configuration file, which is expected to be located in `/etc/collectd`. It is a `JSON` formatted file composed of

* The exchange name
* The exchange API base URL
* The cryptocurrencies pairs you'd like to monitor
* Optional: a conversion pair. For example, [Bittrex][6] does not offer `USD` pairs, but has a `BTC/USDT` pair, adding a conversion pair will multiply listed pairs with this factor.

Example

```
$ cat ticker.json
{
  "bitstamp":
    {
      "url": "https://www.bitstamp.net/api/v2/ticker/",
      "pairs": ["ethusd", "xrpusd", "ltcusd", "btcusd"]
    },
  "bittrex":
    {
      "url": "https://bittrex.com/api/v1.1/public/getticker?market=",
      "pairs": ["BTC-XEM", "BTC-FUN", "BTC-XVG"],
      "convert": "USDT-BTC"
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

For now, `ticker` supports [Bitstamp][9] and [Bittrex][6].

#### Show your appreciation

If you like the software, feel free to send a tip ;)

Bitcoin: ***REMOVED***
Ethereum: ***REMOVED***
Litecoin: ***REMOVED***

[1]: https://collectd.org/
[2]: https://golang.org/
[3]: https://bugs.debian.org/cgi-bin/bugreport.cgi?bug=827759
[4]: https://collectd.org/wiki/index.php/Plugin:Exec
[5]: https://github.com/collectd/go-collectd
[6]: https://bittrex.com/
[7]: https://www.influxdata.com/
[8]: https://grafana.com/
[9]: https://www.bitstamp.net/

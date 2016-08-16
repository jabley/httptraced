[![Build Status](https://travis-ci.org/jabley/go-latency.svg?branch=master)](https://travis-ci.org/jabley/go-latency)

Tool to monitor connection latency.

We want a tool that will happily run for long periods, as well as single-shot.

Hence this, rather than curl and cron (which only gives us minute-by-minute granularity).

It needs the new net/http/httptrace facility in Go 1.7.

## Usage

```shell
$ go-latency -c 4 -i 4 https://github.com/
{"timestamp":"2016-08-16T15:32:55.227148206Z","url":"https://github.com/","getConn":5.0107e-05,"gotConn":0.526888001,"ttfb":0.697935482,"dnsStart":8.283200000000001e-05,"dnsDone":0.001572922,"connectStart":0.0015815970000000001,"connectDone":0.004335918,"wroteRequest":0.52702092,"total":0.6981806340000001}
{"timestamp":"2016-08-16T15:32:59.231211464Z","url":"https://github.com/","getConn":4.1358e-05,"gotConn":0.833581494,"ttfb":0.9899488780000001,"dnsStart":9.874800000000001e-05,"dnsDone":0.0014137910000000002,"connectStart":0.001421156,"connectDone":0.0031775970000000003,"wroteRequest":0.8336668190000001,"total":1.127901335}
{"timestamp":"2016-08-16T15:33:03.227800115Z","url":"https://github.com/","getConn":8.032e-05,"gotConn":0.481622415,"ttfb":0.660823398,"dnsStart":0.00011373100000000001,"dnsDone":0.001378965,"connectStart":0.001386019,"connectDone":0.003737101,"wroteRequest":0.4816906,"total":0.661262624}
{"timestamp":"2016-08-16T15:33:07.22843019Z","url":"https://github.com/","getConn":6.982e-05,"gotConn":1.080641793,"ttfb":1.297710369,"dnsStart":0.00010004600000000001,"dnsDone":0.406959927,"connectStart":0.406969866,"connectDone":0.40870535,"wroteRequest":1.08070248,"total":1.297956954}
```

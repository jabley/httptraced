[![Build Status](https://travis-ci.org/jabley/go-latency.svg?branch=master)](https://travis-ci.org/jabley/go-latency)

Tool to monitor connection latency.

We want a tool that will happily run for long periods, as well as single-shot.

Hence this, rather than curl and cron (which only gives us minute-by-minute granularity).

It needs the new net/http/httptrace facility in Go 1.7.

## Usage

TODO update output using a version of Go that actually contains `net/http/httptrace`.

```
$ go-latency -i 4 https://github.com/
{"timestamp":"2016-05-14T13:47:55.041794271Z","url":"https://github.com/","getConn":0,"gotConn":0,"ttfb":0,"dnsStart":0,"dnsDone":0,"connectStart":0,"connectDone":0,"wroteRequest":0,"total":1.2073199510000001}
{"timestamp":"2016-05-14T13:47:59.046622346Z","url":"https://github.com/","getConn":0,"gotConn":0,"ttfb":0,"dnsStart":0,"dnsDone":0,"connectStart":0,"connectDone":0,"wroteRequest":0,"total":0.300005106}
{"timestamp":"2016-05-14T13:48:03.043108335Z","url":"https://github.com/","getConn":0,"gotConn":0,"ttfb":0,"dnsStart":0,"dnsDone":0,"connectStart":0,"connectDone":0,"wroteRequest":0,"total":0.30343416}
{"timestamp":"2016-05-14T13:48:07.045484402Z","url":"https://github.com/","getConn":0,"gotConn":0,"ttfb":0,"dnsStart":0,"dnsDone":0,"connectStart":0,"connectDone":0,"wroteRequest":0,"total":0.313531331}
{"timestamp":"2016-05-14T13:48:11.044279722Z","url":"https://github.com/","getConn":0,"gotConn":0,"ttfb":0,"dnsStart":0,"dnsDone":0,"connectStart":0,"connectDone":0,"wroteRequest":0,"total":0.299610692}
```

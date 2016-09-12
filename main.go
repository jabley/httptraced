package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"

	flag "github.com/ogier/pflag"

	_ "net/http/pprof"
)

var (
	transport = &http.Transport{DisableKeepAlives: true}
	client    = &http.Client{
		Transport: transport,
		Timeout:   time.Duration(10) * time.Second,
	}

	// CLI flags
	help     bool
	interval int64
	count    int64
)

const (
	usage = `
httptraced [options...] url

httptraced will make a GET request to a URL and report on the timings.

`

	helpUsage   = "Display this help message"
	helpDefault = false

	intervalUsage   = "positive number of seconds to wait between making requests"
	intervalDefault = 2

	countUsage   = "positive number of requests to make. If set, then interval is assumed to be set"
	countDefault = -1
)

func main() {
	// parse flags
	debug := flag.Bool("debug", false, "If true, you can debug this process at http://localhost:6060/debug/pprof/")

	flag.BoolVarP(&help, "help", "h", helpDefault, helpUsage)

	flag.Int64VarP(&interval, "interval", "i", intervalDefault, intervalUsage)

	flag.Int64VarP(&count, "count", "c", countDefault, countUsage)

	flag.Parse()

	if *debug {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	if help {
		showUsage()
		os.Exit(2)
	}

	// Enforce sensible defaults
	if count < 0 {
		count = countDefault
	}
	if interval < 0 {
		interval = intervalDefault
	}

	URL := flag.Arg(0)

	if URL == "" {
		showUsage()
		os.Exit(2)
	}

	poll(URL)
}

func showUsage() {
	fmt.Fprintf(os.Stderr, usage)
	flag.PrintDefaults()
}

type JSONTimestamp time.Time

func (j *JSONTimestamp) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(*j).UTC().Format(time.RFC3339Nano))
	return []byte(stamp), nil
}

type JSONError struct {
	Detail string `json:detail`
}

type JSONOutput struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []JSONError `json:"errors,omitempty"`
}

type TimingContext struct {
	StartTime            JSONTimestamp `json:"timestamp"`
	URL                  string        `json:"url"`
	GetConn              float64       `json:"getConn"`
	GotConn              float64       `json:"gotConn"`
	GotFirstResponseByte float64       `json:"ttfb"`
	DNSStart             float64       `json:"dnsStart"`
	DNSDone              float64       `json:"dnsDone"`
	ConnectStart         float64       `json:"connectStart"`
	ConnectDone          float64       `json:"connectDone"`
	WroteRequest         float64       `json:"wroteRequest"`
	Total                float64       `json:"total"`
}

func New(URL string) *TimingContext {
	t := TimingContext{}
	t.StartTime = JSONTimestamp(time.Now())
	t.URL = URL
	return &t
}

func (tc *TimingContext) Elapsed() float64 {
	return time.Since(time.Time(tc.StartTime)).Seconds()
}

func poll(URL string) {
	encoder := json.NewEncoder(os.Stdout)

	tickInterval := time.Duration(interval) * time.Second

	// channel used to do the initial poll
	start := make(chan struct{})

	// channel used to signal that we've done the required count of polls
	done := make(chan struct{})

	t := time.NewTicker(tickInterval)

	if count != countDefault {
		// stop the timer after count * interval seconds
		go func() {
			<-time.After(time.Duration(count-1) * tickInterval)
			close(done)
		}()
	}

	// This one weird trick to do the initial poll
	go func() {
		start <- struct{}{}
	}()

	for {
		select {
		case <-start:
			doPoll(URL, encoder)
		case <-done:
			t.Stop()
			return
		case <-t.C:
			doPoll(URL, encoder)
		}
	}
}

func doPoll(URL string, encoder *json.Encoder) {
	tc, err := doIt(URL)

	if err != nil {
		write(encoder,
			JSONOutput{
				Errors: []JSONError{
					JSONError{Detail: err.Error()},
				},
			})
		return
	}

	write(encoder, JSONOutput{Data: tc})
}

func doIt(URL string) (*TimingContext, error) {
	req, err := http.NewRequest("GET", URL, nil)

	if err != nil {
		return nil, err
	}

	timingContext := New(URL)

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), &httptrace.ClientTrace{
		GetConn:              func(hostPort string) { timingContext.GetConn = timingContext.Elapsed() },
		GotConn:              func(ci httptrace.GotConnInfo) { timingContext.GotConn = timingContext.Elapsed() },
		GotFirstResponseByte: func() { timingContext.GotFirstResponseByte = timingContext.Elapsed() },
		DNSStart:             func(e httptrace.DNSStartInfo) { timingContext.DNSStart = timingContext.Elapsed() },
		DNSDone:              func(e httptrace.DNSDoneInfo) { timingContext.DNSDone = timingContext.Elapsed() },
		ConnectStart:         func(network, addr string) { timingContext.ConnectStart = timingContext.Elapsed() },
		ConnectDone:          func(network, addr string, err error) { timingContext.ConnectDone = timingContext.Elapsed() },
		WroteRequest:         func(e httptrace.WroteRequestInfo) { timingContext.WroteRequest = timingContext.Elapsed() },
	}))

	res, err := client.Do(req)

	timingContext.Total = timingContext.Elapsed()

	if err != nil {
		return nil, err
	}

	res.Body.Close()

	return timingContext, nil
}

func write(encoder *json.Encoder, jo JSONOutput) {
	if err := encoder.Encode(jo); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

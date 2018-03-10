package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/skycoin-karl/otc/pkg/api/public"
	"github.com/skycoin-karl/otc/pkg/currencies"
	"github.com/skycoin-karl/otc/pkg/currencies/btc"
	"github.com/skycoin-karl/otc/pkg/currencies/sky"
	"github.com/skycoin-karl/otc/pkg/model"
	"github.com/skycoin-karl/otc/pkg/otc"
)

var CURRENCIES = currencies.New()

func init() {
	conf, err := otc.NewConfig("config.toml")
	if err != nil {
		panic(err)
	}

	SKY, err := sky.New(conf)
	if err != nil {
		panic(err)
	}
	CURRENCIES.Add(otc.SKY, SKY)

	BTC, err := btc.New(conf)
	if err != nil {
		panic(err)
	}
	CURRENCIES.Add(otc.BTC, BTC)
}

func main() {
	// for graceful shutdown / cleanup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	modl, err := model.New(CURRENCIES)
	if err != nil {
		panic(err)
	}

	public := public.New(CURRENCIES, modl)
	println("listening")
	go http.ListenAndServe(":8080", public)

	<-stop
	println("stopping")
	modl.Stop()
}

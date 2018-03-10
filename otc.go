package main

import (
	"net/http"

	"github.com/skycoin/services/otc/pkg/api/public"
	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/currencies/btc"
	"github.com/skycoin/services/otc/pkg/currencies/sky"
	"github.com/skycoin/services/otc/pkg/model"
	"github.com/skycoin/services/otc/pkg/otc"
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
	modl := model.New(CURRENCIES)

	public := public.New(CURRENCIES, modl)

	println("listening")
	http.ListenAndServe(":8080", public)
}

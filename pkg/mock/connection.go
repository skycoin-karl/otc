package mock

import (
	"github.com/skycoin/services/otc/pkg/otc"
)

type Connection struct {
	Currency otc.Currency
}

func NewConnection(curr otc.Currency) *Connection {
	return &Connection{
		Currency: curr,
	}
}

func (c *Connection) Balance(addr string) (uint64, error) {
	if c.Currency == otc.BTC {
		return 100000000, nil
	}

	return 0, nil
}

func (c *Connection) Confirmed(txid string) (bool, error) {
	return true, nil
}

func (c *Connection) Send(addr string, amount uint64) (string, error) {
	return "sky_txid", nil
}

func (c *Connection) Address() (string, error) {
	return "", nil
}

func (c *Connection) Connected() (bool, error) {
	return true, nil
}

func (c *Connection) Stop() error {
	return nil
}

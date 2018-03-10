package model

import (
	"encoding/json"
	"os"

	"github.com/skycoin/services/otc/pkg/otc"
)

const (
	PATH string = ".otc/"
	REQS string = "reqs/"
	LOGS string = "logs/"
)

func Save(req *otc.Request, res *otc.Result) error {
	file, err := os.OpenFile(
		PATH+REQS+req.Id()+".json",
		os.O_CREATE|os.O_RDWR,
		0644,
	)
	if err != nil {
		return err
	}

	// empty file
	file.Truncate(0)
	file.Seek(0, 0)

	// indent json
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	// write json to file
	if err = enc.Encode(req); err != nil {
		return err
	}

	// sync to disk
	if err = file.Sync(); err != nil {
		return err
	}

	// close file
	if err = file.Close(); err != nil {
		return err
	}

	return Log(req, res)
}

func Log(req *otc.Request, res *otc.Result) error {
	file, err := os.OpenFile(
		PATH+LOGS+"log.json",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return err
	}

	event := &otc.Event{
		Id:       req.Id(),
		Status:   req.Status,
		Finished: res.Finished,
	}
	if res.Err != nil {
		event.Err = res.Err.Error()
	}

	if err = json.NewEncoder(file).Encode(&event); err != nil {
		return err
	}

	if err = file.Sync(); err != nil {
		return err
	}

	return file.Close()
}

func Load(path string) ([]*otc.Request, error) {
	return nil, nil
}

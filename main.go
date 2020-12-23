package main

import "log"

var (
	cdmonResolvers = []string{"46.16.60.166", "46.16.60.159", "35.156.85.88"}

	updates = []*Updatable{
		&Updatable{
			User:    "dyndnsupdater",
			PassMD5: "cc2c3715b70c11f62c9ac6c70389e957",
			Email:   "jllopis@gimlab.net",
			Host:    "gimlab.net",
		},
		&Updatable{
			User:    "srv1gimlabupdater",
			PassMD5: "56f1c60558111852a08654d49b04d3ac",
			Email:   "jllopis@gimlab.net",
			Host:    "mx.gimlab.net",
		},
		&Updatable{
			User:    "mxgimlabupdater",
			PassMD5: "48834d0b40252d0960e35b624efac8c7",
			Email:   "jllopis@gimlab.net",
			Host:    "srv1.gimlab.net",
		},
	}
)

func main() {
	for _, updt := range updates {
		log.Printf("%#+v\n", updt)
	}
}

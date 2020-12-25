package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
)

var (
	// cdmonResolvers = []string{"46.16.60.166", "46.16.60.159", "35.156.85.88"}

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
			Host:    "srv1.gimlab.net",
		},
		&Updatable{
			User:    "mxgimlabupdater",
			PassMD5: "48834d0b40252d0960e35b624efac8c7",
			Email:   "jllopis@gimlab.net",
			Host:    "mx.gimlab.net",
		},
	}
)

func main() {
	var wg sync.WaitGroup

	externalIP, err := getExternalIP()
	if err != nil {
		log.Fatalf("Can not get external IP: %v", err)
	}
	for _, updt := range updates {
		updt.CurrentIP = externalIP
		wg.Add(1)
		go updt.execute(&wg)
	}
	wg.Wait()
}

func getExternalIP() (net.IP, error) {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		ipAddr := string(bodyBytes)
		log.Printf("Got external address: %s", ipAddr)
		return net.ParseIP(ipAddr), nil
	}
	return nil, nil
}

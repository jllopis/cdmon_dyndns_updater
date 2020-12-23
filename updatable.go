package main

import (
	"net"
	"sync"
)

// Updatable holds the values for an updatable dynamic dns entry
type Updatable struct {
	sync.WaitGroup

	User         string
	PassMD5      string
	Email        string
	Host         string
	CurrentIP    net.IP
	RegisteredIP net.IP
	shouldUpdate bool
}

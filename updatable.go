package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/smtp"
	"net/url"
	"sync"
)

var (
	cdmonUrl       = "https://dinamico.cdmon.org/onlineService.php?enctype=MD5&n=%s&p=%s"
	updateCdmonUrl = `https://dinamico.cdmon.org/onlineService.php?enctype=MD5&n=%s&p=%s&cip=%s`
)

// Updatable holds the values for an updatable dynamic dns entry
// Here we hold the cdmon user data
// User is the user name to authenticate into cdmon.com
// PassMD5 is the md5 sum of the user password
// Email is the email address we want to send notifications to
// Host is the domain/subdomain we want to update
type Updatable struct {
	sync.WaitGroup
	mu sync.Mutex

	User         string
	PassMD5      string
	Email        string
	Host         string
	CurrentIP    net.IP
	RegisteredIP net.IP
	updatedIP    bool
}

func (u *Updatable) execute(wg *sync.WaitGroup) {
	defer wg.Done()
	// get registered IP with cdmon
	regIP, err := u.getRegisteredIP()
	if err != nil {
		log.Printf("[ERROR] can't get registered address: %v", err)
		return
	}
	u.RegisteredIP = regIP
	log.Printf("Got registered IP: %s", regIP.String())
	// If IPs differ, update registered IP
	if u.CurrentIP.Equal(u.RegisteredIP) {
		log.Printf("(%s) %s IP hasn't changed", u.User, u.Host)
		// TODO: No t'oblidis de descomentar-ho quant el notify funcioni ok
		// return
	}
	// Change detected. Update
	if err := u.updateRegisteredIP(); err != nil {
		log.Printf("%v", err)
		return
	}

	// Notify change
	if u.updatedIP {
		err = u.notify()
		if err != nil {
			log.Printf("Can not send notification email: %v", err)
			return
		}
	}
}

func (u *Updatable) updateRegisteredIP() error {
	resp, err := http.Get(fmt.Sprintf(cdmonUrl+"&cip=%s", u.User, u.PassMD5, u.CurrentIP))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		b := string(bodyBytes)
		ipAddr, err := url.ParseQuery(b)
		if err != nil {
			return err
		}
		if status := ipAddr.Get("resultat"); status != "customok" {
			if status == "" {
				status = "unknown error"
			}
			return fmt.Errorf("error updating ip to cdmon: %s", status)
		}
	}
	log.Printf("Updated %s IP for %s", u.User, u.Host)
	u.updatedIP = true
	return nil
}

func (u *Updatable) getRegisteredIP() (net.IP, error) {
	resp, err := http.Get(fmt.Sprintf(cdmonUrl, u.User, u.PassMD5))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		b := string(bodyBytes)
		ipAddr, err := url.ParseQuery(b)
		if err != nil {
			return nil, err
		}
		return net.ParseIP(ipAddr.Get("newip")), nil
	}
	return nil, nil
}

func (u *Updatable) notify() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	// Sender data.
	from := "jllopis@gimlab.net"
	password := "Nua3aino"

	// Receiver email address.
	to := []string{
		"jllopis@gimlab.net",
	}

	// smtp server configuration.
	smtpHost := "mx.gimlab.net"
	smtpPort := "587"

	// Message.
	msg := []byte("From: jllopis@gmail.net\r\n" +
		"To: " + u.Email + "\r\n" +
		"Subject: Cdmon Dyn DNS Updater\r\n" +
		"\r\n" +
		"Your IP for " + u.Host + " has been updated\r\n\n" +
		"Old IP: " + u.RegisteredIP.String() + "\r\n" +
		"New IP: " + u.CurrentIP.String() + "\r\n" +
		"User " + u.User + "\r\n")

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		return err
	}
	fmt.Printf("Email for %s Sent Successfully!", u.Host)

	return nil
}

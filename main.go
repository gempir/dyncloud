package main

import (
	"log"
	"github.com/cloudflare/cloudflare-go"
	"os"
	"io/ioutil"
	"net/http"
)

var (
	config Config
)

func main() {
	config = ReadConfig("/etc/dyncloud")

	// Construct a new API object
	api, err := cloudflare.New(config.ApiKey, config.Email)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch the zone ID
	zoneId, err := api.ZoneIDByName(config.Domain) // Assuming example.com exists in your Cloudflare account already
	if err != nil {
		log.Fatal(err)
	}

	whyDoINeedThisRecord := new(cloudflare.DNSRecord)
	records, err := api.DNSRecords(zoneId, *whyDoINeedThisRecord)
	if err != nil {
		log.Fatal(err)
	}
	var dnsRecord cloudflare.DNSRecord
	for _, record := range records {
		if record.Name == config.DNSRecord {
			dnsRecord = record
		}
	}
	if dnsRecord.ID == "" {
		os.Exit(1)
	}

	resp, _ := httpRequest("https://api.ipify.org")
	publicIP := string(resp[:])

	dnsRecord.Content = publicIP

	log.Println("Updating DNS Record", dnsRecord.Name, "to IP", publicIP)
	err = api.UpdateDNSRecord(zoneId, dnsRecord.ID, dnsRecord)
	if err != nil {
		log.Fatal("Failed to update DNSRecord", dnsRecord.Name)
	}

	log.Println("Sucessfully updated DNSRecord", dnsRecord.Name)
}

func httpRequest(url string) ([]byte, error) {
	log.Printf("httpRequest %s", url)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return contents, nil
}
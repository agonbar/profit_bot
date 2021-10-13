package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type RestSno struct {
	NodeID         string      `json:"nodeID"`
	Wallet         string      `json:"wallet"`
	WalletFeatures interface{} `json:"walletFeatures"`
	Satellites     []struct {
		ID                 string      `json:"id"`
		URL                string      `json:"url"`
		Disqualified       interface{} `json:"disqualified"`
		Suspended          interface{} `json:"suspended"`
		CurrentStorageUsed int64       `json:"currentStorageUsed"`
	} `json:"satellites"`
	DiskSpace struct {
		Used      int64 `json:"used"`
		Available int64 `json:"available"`
		Trash     int64 `json:"trash"`
		Overused  int   `json:"overused"`
	} `json:"diskSpace"`
	Bandwidth struct {
		Used      int64 `json:"used"`
		Available int   `json:"available"`
	} `json:"bandwidth"`
	LastPinged     time.Time `json:"lastPinged"`
	Version        string    `json:"version"`
	AllowedVersion string    `json:"allowedVersion"`
	UpToDate       bool      `json:"upToDate"`
	StartedAt      time.Time `json:"startedAt"`
}

func getSpace(url string) [2]string {

	resp, err := http.Get(url + "/api/sno/")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	//Create a variable of the same type as our model
	var cResp RestSno

	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		log.Fatal("ooopsss! an error occurred, please try again")
	}

	var toret [2]string

	toret[0] = fmt.Sprintf("%.2f", (float32(cResp.DiskSpace.Used) / 1000000000000))
	toret[1] = fmt.Sprintf("%.2f", (float32(cResp.DiskSpace.Available) / 1000000000000))

	return toret
}

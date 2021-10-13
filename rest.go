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

type RestEstimatedPayout struct {
	CurrentMonth struct {
		EgressBandwidth         int64   `json:"egressBandwidth"`
		EgressBandwidthPayout   float64 `json:"egressBandwidthPayout"`
		EgressRepairAudit       int64   `json:"egressRepairAudit"`
		EgressRepairAuditPayout float64 `json:"egressRepairAuditPayout"`
		DiskSpace               float64 `json:"diskSpace"`
		DiskSpacePayout         float64 `json:"diskSpacePayout"`
		HeldRate                int     `json:"heldRate"`
		Payout                  float64 `json:"payout"`
		Held                    float64 `json:"held"`
	} `json:"currentMonth"`
	PreviousMonth struct {
		EgressBandwidth         int64   `json:"egressBandwidth"`
		EgressBandwidthPayout   float64 `json:"egressBandwidthPayout"`
		EgressRepairAudit       int64   `json:"egressRepairAudit"`
		EgressRepairAuditPayout float64 `json:"egressRepairAuditPayout"`
		DiskSpace               float64 `json:"diskSpace"`
		DiskSpacePayout         float64 `json:"diskSpacePayout"`
		HeldRate                int     `json:"heldRate"`
		Payout                  float64 `json:"payout"`
		Held                    float64 `json:"held"`
	} `json:"previousMonth"`
	CurrentMonthExpectations int `json:"currentMonthExpectations"`
}

func getPrice(url string) string {
	lastChar := url[len(url)-1:]
	if lastChar != "/" {
		url = url + "/"
	}
	resp, err := http.Get(url + "api/sno/estimated-payout")
	if err != nil {
		return err.Error()
	}

	defer resp.Body.Close()
	//Create a variable of the same type as our model
	var cResp RestEstimatedPayout

	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%.2f", (float32(cResp.CurrentMonth.Payout) / 100))
}

func getSpace(url string) [2]string {
	lastChar := url[len(url)-1:]
	if lastChar != "/" {
		url = url + "/"
	}
	resp, err := http.Get(url + "api/sno/")
	if err != nil {
		return [2]string{err.Error(), ""}
	}

	defer resp.Body.Close()
	//Create a variable of the same type as our model
	var cResp RestSno

	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		return [2]string{err.Error(), ""}
	}

	var toret [2]string

	toret[0] = fmt.Sprintf("%.2f", (float32(cResp.DiskSpace.Used) / 1000000000000))
	toret[1] = fmt.Sprintf("%.2f", (float32(cResp.DiskSpace.Available) / 1000000000000))

	return toret
}

func getStatus(url string) [3]string {

	space := getSpace(url)

	price := getPrice(url)

	log.Println([2]string{space[0] + "/" + space[1], price})

	return [3]string{space[0], space[1], price}
}

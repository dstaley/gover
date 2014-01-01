package main

import (
	"bytes"
	"encoding/json"
	"github.com/dstaley/cgminerapi"
	"log"
	"net/http"
	"strconv"
	"time"
)

type MobileminerAPIJSON struct {
	AcceptedShares        int     `json:"AcceptedShares,omitempty"`
	Algorithm             string  `json:"Algorithm,omitempty"`
	AverageHashrate       float64 `json:"AverageHashrate,omitempty"`
	CoinName              string  `json:"CoinName,omitempty"`
	CoinSymbol            string  `json:"CoinSymbol,omitempty"`
	CurrentHashrate       float64 `json:"CurrentHashrate,omitempty"`
	DeviceID              int     `json:"DeviceID,omitempty"`
	Enabled               bool    `json:"Enabled,omitempty"`
	FanPercent            float64 `json:"FanPercent,omitempty"`
	FanSpeed              int     `json:"FanSpeed,omitempty"`
	FullName              string  `json:"FullName,omitempty"`
	GpuActivity           int     `json:"GpuActivity,omitempty"`
	GpuClock              int     `json:"GpuClock,omitempty"`
	GpuVoltage            float64 `json:"GpuVoltage,omitempty"`
	HardwareErrors        int     `json:"HardwareErrors,omitempty"`
	HardwareErrorsPercent float64 `json:"HardwareErrorsPercent,omitempty"`
	Index                 int     `json:"Index,omitempty"`
	Intensity             string  `json:"Intensity,omitempty"`
	Kind                  string  `json:"Kind,omitempty"`
	MemoryClock           int     `json:"MemoryClock,omitempty"`
	MinerName             string  `json:"MinerName,omitempty"`
	Name                  string  `json:"Name,omitempty"`
	PoolIndex             float64 `json:"PoolIndex,omitempty"`
	PoolName              string  `json:"PoolName,omitempty"`
	PowerTune             int     `json:"PowerTune,omitempty"`
	RejectedShares        int     `json:"RejectedShares,omitempty"`
	RejectedSharesPercent float64 `json:"RejectedSharesPercent,omitempty"`
	Status                string  `json:"Status,omitempty"`
	Temperature           float64 `json:"Temperature,omitempty"`
	Utility               float64 `json:"Utility,omitempty"`
}

func updateMobileminer() {
	api := cgminerapi.APIClient{config.CGMiner.Hostname, config.CGMiner.Port}
	for {
		c := cgminerapi.APICommand{Method: "devs"}
		res, err := api.Send(&c)
		if err != nil {
			println("52")
		}
		var arr []MobileminerAPIJSON
		for _, v := range res.Devs {
			var data MobileminerAPIJSON
			data.AcceptedShares = v.Accepted
			data.Algorithm = "Scrypt"
			data.AverageHashrate = v.MHSav
			data.CurrentHashrate = v.MHS5s
			data.DeviceID = *v.GPU
			data.Enabled = true
			data.FanPercent = v.FanPercent
			data.FanSpeed = v.FanSpeed
			data.FullName = "GPU " + strconv.Itoa(*v.GPU)
			data.GpuActivity = v.GPUActivity
			data.GpuClock = v.GPUClock
			data.GpuVoltage = v.GPUVoltage
			data.HardwareErrors = 1
			data.HardwareErrorsPercent = 1
			data.Index = *v.GPU
			data.Intensity = v.Intensity
			data.Kind = "GPU"
			data.MemoryClock = v.MemoryClock
			data.MinerName = "Gover"
			data.Name = "GPU " + strconv.Itoa(*v.GPU)
			data.PowerTune = *v.Powertune
			data.RejectedShares = v.Rejected
			data.RejectedSharesPercent = 1
			data.Status = v.Status
			data.Temperature = v.Temperature
			data.Utility = v.Utility
			arr = append(arr, data)
		}
		blob, err := json.Marshal(arr)
		if err != nil {
			println("80")
		}
		r := bytes.NewReader(blob)
		url := "https://mobileminer.azurewebsites.net/api/MiningStatisticsInput"
		// rb := "http://requestb.in/17ppwlx1?emailAddress="
		_, err = http.Post(url+"?emailAddress="+config.Mobileminer.Email+"&applicationKey="+config.Mobileminer.APIKey+"&machineName="+config.Mobileminer.MachineName+"&apiKey=zYZJ35h8YRuE2W", "application/json", r)
		if err != nil {
			log.Print(err)
		}
		time.Sleep(time.Duration(config.Mobileminer.UpdateInterval) * time.Second)
	}
}

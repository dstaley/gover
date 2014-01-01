package main

import (
	"database/sql"
	"github.com/dstaley/cgminerapi"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
	"time"
)

type HistoricalGPUData struct {
	Hashrate    [][]int `json:"hashrate,omitempty"`
	Id          string  `json:"id,omitempty"`
	Temperature [][]int `json:"temperature,omitempty"`
}

type HistoricalGPUDatas []HistoricalGPUData

type GPUStats struct {
	Time        int
	GPU_ID      int
	Temperature int
	Fan_Speed   int
	Fan_Load    float64
	Hashrate    float64
	Core        int
	Memory      int
	Power       float64
}

func AddGPUStats(database *sql.DB, v cgminerapi.Devs) error {
	tx, err := database.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO GPUStats VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(time.Now().Unix(), *v.GPU, v.Temperature, int(v.FanSpeed), v.FanPercent, v.MHS5s, int(v.GPUClock), int(v.MemoryClock), v.GPUVoltage)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func updateGPUStats(db *sql.DB) {
	api := cgminerapi.APIClient{config.CGMiner.Hostname, config.CGMiner.Port}
	for {
		c := cgminerapi.APICommand{Method: "devs"}
		res, err := api.Send(&c)
		if err != nil {
			println(err)
		}
		for _, v := range res.Devs {
			err = AddGPUStats(db, v)
			if err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(time.Duration(config.Database.UpdateInterval) * time.Second)
	}
}

func GetHistoricalGPUData(db *sql.DB, id int) HistoricalGPUData {
	rows, err := db.Query("SELECT * FROM GPUStats WHERE GPU_ID = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	gpu := HistoricalGPUData{Id: "GPU" + strconv.Itoa(id)}
	for rows.Next() {
		var r GPUStats
		if err := rows.Scan(&r.Time, &r.GPU_ID, &r.Temperature, &r.Fan_Speed, &r.Fan_Load, &r.Hashrate, &r.Core, &r.Memory, &r.Power); err != nil {
			log.Fatal(err)
		}
		gpu.Hashrate = append(gpu.Hashrate, []int{r.Time, int(r.Hashrate * 1000)})
		gpu.Temperature = append(gpu.Temperature, []int{r.Time, r.Temperature})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return gpu
}

func GetHistoricalGPUDatas(db *sql.DB) HistoricalGPUDatas {
	rows, err := db.Query("SELECT DISTINCT GPU_ID FROM GPUStats")
	if err != nil {
		log.Fatal(err)
	}
	var gpus HistoricalGPUDatas
	for rows.Next() {
		var i int
		if err := rows.Scan(&i); err != nil {
			log.Fatal(err)
		}
		gpus = append(gpus, GetHistoricalGPUData(db, i))
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return gpus
}

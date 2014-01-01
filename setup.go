package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dstaley/cgminerapi"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

type CGMinerSettings struct {
	Hostname string `json:"hostname,omitempty"`
	Port     string `json:"port, omitempty"`
}

type DatabaseSettings struct {
	UpdateInterval int `json:"update_interval,omitempty"`
}

type MobileminerSettings struct {
	Email          string `json:"email,omitempty"`
	APIKey         string `json:"application_key,omitempty"`
	MachineName    string `json:"machine_name,omitempty"`
	UpdateInterval int    `json:"update_interval,omitempty"`
}

type Settings struct {
	CGMiner     CGMinerSettings      `json:"CGMiner,omitempty"`
	Database    DatabaseSettings     `json:"database,omitempty"`
	Mobileminer *MobileminerSettings `json:"mobileminer,omitempty"`
}

var api *cgminerapi.APIClient

func createDatabase() {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sql := `CREATE TABLE GPUStats (
	Time INTEGER,
	GPU_ID INTEGER,
	Temperature REAL,
	Fan_Speed INTEGER,
	Fan_Load REAL,
	Hashrate REAL,
	Core INTEGER,
	Memory INTEGER,
	Power REAL
);`
	_, err = db.Exec(sql)
	if err != nil {
		log.Printf("%q: %s\n", err, sql)
		return
	}
}

func setup() {
	var a string
	if _, err := os.Stat("./db.db"); os.IsNotExist(err) {
		println("Looks like you don't have a database yet? Shall I create one for you?")
		_, err := fmt.Scanln(&a)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for {
			if a == "Y" || a == "y" || a == "Yes" || a == "yes" {
				println("Okay! Creating a database!")
				createDatabase()
				break
			} else if a == "N" || a == "n" || a == "No" || a == "no" {
				println("Okay, I won't create a database for you.")
				break
			} else {
				println("Hmm, not sure what you mean. Shall I create a database for you? Yes or no?")
				_, err := fmt.Scanln(&a)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		}
		var updateinterval int
		println("How often would you like to store data in the database (in seconds)? [120]")
		_, _ = fmt.Scanln(&updateinterval)
		if updateinterval == 0 {
			updateinterval = 120
		}
		config.Database.UpdateInterval = updateinterval
		var hostname, port string
		println("Now let's setup your configuration.")
		println("What's the hostname for your CGMiner instance? [localhost]")
		_, _ = fmt.Scanln(&hostname)
		if hostname == "" {
			hostname = "localhost"
		}
		println("What port is the API listening in on? [4028]")
		_, _ = fmt.Scanln(&port)
		if port == "" {
			port = "4028"
		}
		config.CGMiner.Hostname = hostname
		config.CGMiner.Port = port

		println("Would you like to setup a Mobileminerapp connection?")
		_, err = fmt.Scanln(&a)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for {
			if a == "Y" || a == "y" || a == "Yes" || a == "yes" {
				println("Please enter your email address.")
				var email string
				_, err = fmt.Scanln(&email)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				println("Please enter your application key.")
				var appkey string
				_, err = fmt.Scanln(&appkey)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				println("What is the name of your rig?")
				var rigname string
				_, err = fmt.Scanln(&rigname)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				println("How often would you like to update Mobileminer (in seconds)? [120]")
				var mmupdate int
				_, _ = fmt.Scanln(&mmupdate)
				if mmupdate == 0 {
					mmupdate = 120
				}
				var m MobileminerSettings
				m.Email = email
				m.APIKey = appkey
				m.MachineName = rigname
				m.UpdateInterval = mmupdate
				config.Mobileminer = &m
				break
			} else if a == "N" || a == "n" || a == "No" || a == "no" {
				break
			} else {
				println("Hmm, not sure what you mean. Would you like to setup a Mobileminerapp connection? Yes or no?")
				_, err := fmt.Scanln(&a)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		}

		b, err := json.MarshalIndent(config, "", "	")
		if err != nil {
			log.Fatal("Problem marshalling configuration to json")
		}

		f, err := os.Create("config.json")
		if err != nil {
			log.Fatal("Problem writing config.json")
		}

		f.Write(b)
		f.Close()
	} else {
		log.Fatal("Hmm, looks like you already have a database. Please delete it and run setup again.")
	}
}

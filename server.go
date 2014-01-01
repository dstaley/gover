package main

import (
	"database/sql"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/gzip"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/dstaley/cgminerapi"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"strconv"
)

var config Settings

func server() {
	configure()
	m := martini.Classic()

	database, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		log.Fatal(err)
	}

	api := cgminerapi.APIClient{config.CGMiner.Hostname, config.CGMiner.Port}

	if _, err := os.Stat("./db.db"); os.IsNotExist(err) {
		log.Fatal("Missing a database")
	} else {
		go updateGPUStats(database)
	}

	if config.Mobileminer != nil {
		go updateMobileminer()
	}

	m.Use(func(w http.ResponseWriter) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Cache-Control", fmt.Sprintf("private, max-age=%d", config.Database.UpdateInterval))
	})

	m.Use(gzip.All())

	m.Use(render.Renderer())

	m.Get("/api/summary", func(w http.ResponseWriter, r render.Render) {
		c := cgminerapi.APICommand{Method: "summary"}
		resp, _ := api.Send(&c)
		r.JSON(200, resp)
	})

	m.Get("/api/devs", func(w http.ResponseWriter, r render.Render) {
		c := cgminerapi.APICommand{Method: "devs"}
		resp, _ := api.Send(&c)
		r.JSON(200, resp)
	})

	m.Get("/api/gpu/:id", func(w http.ResponseWriter, r render.Render, params martini.Params) {
		c := cgminerapi.APICommand{Method: "gpu", Parameter: params["id"]}
		resp, _ := api.Send(&c)
		r.JSON(200, resp)
	})

	m.Get("/api/gpu/:id/historical", func(w http.ResponseWriter, r render.Render, params martini.Params) {
		i, _ := strconv.Atoi(params["id"])
		r.JSON(200, GetHistoricalGPUData(database, i))
	})

	m.Get("/api/historical/gpu", func(w http.ResponseWriter, r render.Render) {
		r.JSON(200, GetHistoricalGPUDatas(database))
	})

	http.ListenAndServe(":8080", m)
}

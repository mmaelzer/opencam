package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/mmaelzer/opencam/model/migrations"
	"github.com/mmaelzer/opencam/pipeline/manager"
	"github.com/mmaelzer/opencam/server"
	"github.com/mmaelzer/opencam/settings"
	"github.com/mmaelzer/opencam/store"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var settingsfile = flag.String("settings", "", "settings file")
var sqldURL = flag.String("sqld", "http://localhost:9090", "the url for sqld")

func handleFlags() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		go func() {
			time.Sleep(time.Second * 10)
			pprof.StopCPUProfile()
			os.Exit(0)
		}()
	}

	if *settingsfile != "" {
		file, err := os.Open(*settingsfile)
		if err != nil {
			panic(err)
		}
		err = settings.Load(file)
		if err != nil {
			panic(err)
		}
	}

	settings.Set("sqldURL", *sqldURL)
	store.Initialize(*sqldURL)
	if err := migrations.Migrate("sqlite3"); err != nil {
		log.Fatal(`Problems connecting to a sqld instance.
For more information on setting up sqld, see https://github.com/mmaelzer/sqld.`, err)
	}
}

func main() {
	log.SetOutput(os.Stdout)
	log.Printf("=== opencam started ===")

	handleFlags()
	manager.Start()
	store.InitTrim()
	server.Serve(manager.Cameras())
}

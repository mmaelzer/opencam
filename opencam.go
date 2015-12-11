package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/mmaelzer/cam"
	"github.com/mmaelzer/opencam/client"
	"github.com/mmaelzer/opencam/pipeline"
	"github.com/mmaelzer/opencam/settings"
	"github.com/mmaelzer/opencam/store"
	"github.com/mmaelzer/opencam/transform"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var settingsfile = flag.String("settings", "", "settings file")

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
}

func writeToDisk(block pipeline.Block) {
	if len(block.Frames) == 0 {
		return
	}

	ff := block.Frames[0]
	name := ff.Timestamp.Format("15_04_05")
	out := settings.GetString("output")

	for _, frame := range block.Frames {
		store.Save(frame.Bytes, frame.Timestamp, out, name, "jpg")
	}
}

func main() {
	log.SetOutput(os.Stdout)
	log.Printf("=== opencam started ===")

	handleFlags()

	camera := cam.Camera{
		Name:     settings.GetString("camera.name"),
		URL:      settings.GetString("camera.url"),
		Username: settings.GetString("camera.username"),
		Password: settings.GetString("camera.password"),
		Log:      true,
	}

	pipeline.New().
		AddCamera(&camera).
		AddTransform(transform.Motion()).
		AddWriter(writeToDisk).
		Start()

	client.Serve([]cam.Camera{camera})
}

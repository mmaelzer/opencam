package store

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/mmaelzer/opencam/settings"
)

func InitTrim() {
	go start()
}

func start() {
	for {
		log.Print("Trim begin")
		trimmed := trim()
		log.Printf("Trim end. %d %s removed.", trimmed, naivePlural("folder", trimmed))
		time.Sleep(time.Hour)
	}
}

func naivePlural(s string, count int) string {
	if count == 1 {
		return s
	}
	return fmt.Sprintf("%ss", s)
}

func trim() (deleted int) {
	o := settings.GetString("output")
	files, err := ioutil.ReadDir(o)
	if err != nil {
		log.Fatalf("Unable to read from output folder '%s': %s", o, err.Error())
	}

	r := settings.Get("retention")
	retention, ok := r.(time.Duration)
	if !ok {
		log.Fatal("Invalid retention setting value: %v", retention)
	}

	trimDate := time.Now().Add(-retention).Format(DATE_FORMAT)
	var inspectDirs []string
	var deleteDirs []string
	for _, file := range files {
		if file.IsDir() {
			dir := path.Join(o, file.Name())
			if file.Name() < trimDate {
				deleteDirs = append(deleteDirs, dir)
			} else if file.Name() == trimDate {
				inspectDirs = append(inspectDirs, dir)
			}
		}
	}

	deleted += handleDirs(deleteDirs, deleteDir)
	deleted += handleDirs(inspectDirs, inspectDir)
	return
}

func handleDirs(dirs []string, handler func(string) (int, error)) (deleted int) {
	for _, d := range dirs {
		ed, err := handler(d)
		if err != nil {
			log.Printf("Error deleting %s: %s", d, err.Error())
		} else {
			deleted += ed
		}
	}
	return
}

func inspectDir(dir string) (deleted int, err error) {
	var files []os.FileInfo
	files, err = ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	now := time.Now().Format(TIME_FORMAT)
	for _, f := range files {
		fmt.Println(f.Name())
		if !f.IsDir() {
			continue
		}
		splitEvent := strings.Split(f.Name(), "-")
		// Unknown folder name
		if len(splitEvent) < 2 {
			continue
		}
		eventTime := splitEvent[1]
		if now > eventTime {
			d := path.Join(dir, f.Name())
			if err = deleteEventDir(d); err != nil {
				return
			} else {
				deleted++
			}
		}
	}
	return
}

func deleteDir(dir string) (deleted int, err error) {
	var files []os.FileInfo
	files, err = ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, f := range files {
		if f.IsDir() {
			d := path.Join(dir, f.Name())
			if err = deleteEventDir(d); err != nil {
				return
			} else {
				deleted++
			}
		}
	}

	err = os.RemoveAll(dir)
	return
}

func deleteEventDir(d string) error {
	var err error
	if err = os.RemoveAll(d); err != nil {
		return err
	}
	o := settings.GetString("output")
	fp := strings.TrimPrefix(strings.TrimPrefix(d, o), "/")
	_, err = Delete("event", map[string]string{
		"filepath": fp,
	})
	if err == nil {
		log.Printf("Trimmed event %s", d)
	}
	return err
}

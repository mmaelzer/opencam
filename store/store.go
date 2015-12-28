package store

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

func Save(data []byte, date time.Time, base, name, ext string) error {
	dateDir := date.Format("2006-01-02")
	fullPath := path.Join(base, dateDir, name)
	err := os.MkdirAll(fullPath, 0755)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("%d.%s", date.UnixNano(), ext)
	file := path.Join(fullPath, filename)
	log.Printf("saving %s", file)
	return ioutil.WriteFile(file, data, 0755)
}

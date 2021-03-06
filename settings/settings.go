package settings

import (
	"io"
	"time"

	"github.com/spf13/viper"
)

var s *viper.Viper

func init() {
	s = viper.New()
	s.SetConfigType("json")
	s.SetDefault("http.port", 8080)
	s.SetDefault("output", "videos")
	s.SetDefault("retention", time.Hour*24*7)
}

func Load(r io.Reader) error {
	return s.ReadConfig(r)
}

func Get(key string) interface{} {
	return s.Get(key)
}

func GetString(key string) string {
	return s.GetString(key)
}

func Set(key string, val interface{}) {
	s.Set(key, val)
}

func GetInt(key string) int {
	return s.GetInt(key)
}

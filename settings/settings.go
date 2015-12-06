package settings

import (
	"io"

	"github.com/spf13/viper"
)

var s *viper.Viper

func init() {
	s = viper.New()
	s.SetConfigType("json")
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

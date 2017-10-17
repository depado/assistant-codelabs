package conf

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// C is the main exported configuration struct
var C *Conf

// Conf is a struct used to configure the program
type Conf struct {
	Host  string `yaml:"host"`
	Port  int    `yaml:"port"`
	Debug bool   `yaml:"debug"`
	DB    string `yaml:"db"`
	Token string `yaml:"token"`
}

// ListenAddress returns the listen address of the server
func (c Conf) ListenAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// NewDefaults creates a new Conf struct with default values already filled
func NewDefaults() *Conf {
	return &Conf{
		Host:  "127.0.0.1",
		Port:  8080,
		Debug: true,
		DB:    "cars.db",
	}
}

// Load loads the configuration file into C
func Load(fp string) error {
	var err error
	var c []byte

	C = NewDefaults()

	if c, err = ioutil.ReadFile(fp); err != nil {
		return err
	}
	return yaml.Unmarshal(c, &C)
}

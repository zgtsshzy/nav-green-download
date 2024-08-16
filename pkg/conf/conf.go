package conf

import (
	"fmt"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var config = new(Conf)

func Get() *Conf {
	return config
}

type Conf struct {
	SMOCDir   string `mapstructure:"SMOC_DIR"`
	MFWAMDir  string `mapstructure:"MFWAM_DIR"`
	SeaIceDir string `mapstructure:"SEA_ICE_DIR"`
	ECDir     string `mapstructure:"EC_DIR"`
	GFSDir    string `mapstructure:"GFS_DIR"`
	Log       Log    `mapstructure:"LOG"`
}

func New() *Conf {
	viper.AutomaticEnv()

	config.SMOCDir = viper.GetString("SMOC_DIR")
	config.MFWAMDir = viper.GetString("MFWAM_DIR")
	config.SeaIceDir = viper.GetString("SEA_ICE_DIR")
	config.ECDir = viper.GetString("EC_DIR")
	config.GFSDir = viper.GetString("GFS_DIR")

	config.Log = NewLog()

	config.Log.File = viper.GetString("LOG_FILE")
	config.Log.Level = viper.GetString("LOG_LEVEL")

	config.Log.InitLog()

	return config
}

func (c *Conf) Show() {
	if b, err := yaml.Marshal(c); err != nil {
		return
	} else {
		fmt.Printf(`
-----------------------------------------------------------------------------------------
%s
-----------------------------------------------------------------------------------------
`, string(b))
	}
}

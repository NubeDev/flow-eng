package rubixos

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
)

var config *Config

func GetConfig() (*Config, error) {
	if config != nil {
		return config, nil
	} else {
		yamlFile, err := ioutil.ReadFile("/data/rubix-os/config/config.yml")
		if err != nil {
			return nil, errors.New(fmt.Sprintf("yamlFile.Get err   #%v ", err))
		}
		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("yamlFile.Get err   #%v ", err))
		}
		return config, nil
	}
}

type Config struct {
	Server struct {
		KeepAlivePeriodSeconds int
		ListenAddr             string `default:"0.0.0.0"`
		Port                   int
		ResponseHeaders        map[string]string
	}
	Database struct {
		Dialect    string `default:"sqlite3"`
		Connection string `default:"data.db"`
		LogLevel   string `default:"WARN"`
	}
	UnitsImperial bool `default:"false"` // metric or imperial units
	LogLevel      string
	Location      struct {
		GlobalDir string `default:"./"`
		ConfigDir string `default:"config"`
		DataDir   string `default:"data"`
		Data      struct {
			PluginsDir        string `default:"plugins"`
			ModulesDir        string `default:"modules"`
			UploadedImagesDir string `default:"images"`
			ModulesDataDir    string `default:"modules-data"`
		}
	} // leave it as default; don't include in config.eg.yml
	Prod         bool  `default:"false"` // set from commandline; don't include in config.eg.yml
	Auth         *bool `default:"true"`  // set from commandline; don't include in config.eg.yml
	PointHistory struct {
		Enable  *bool `default:"true"`
		Cleaner struct {
			Enable              *bool `default:"true"`
			Frequency           int   `default:"600"`
			DataPersistingHours int   `default:"24"`
		}
		IntervalHistoryCreator struct {
			Enable    *bool `default:"true"`
			Frequency int   `default:"10"`
		}
	}
	MQTT struct {
		Enable                *bool  `default:"true"`
		Address               string `default:"localhost"`
		Port                  int    `default:"1883"`
		Username              string `default:""`
		Password              string `default:""`
		AutoReconnect         *bool  `default:"true"`
		ConnectRetry          *bool  `default:"true"`
		QOS                   int    `default:"1"`
		Retain                *bool  `default:"true"`
		GlobalBroadcast       *bool  `default:"false"` // if set to true will include the plat details in the topic
		PublishPointCOV       *bool  `default:"true"`
		PublishPointList      *bool  `default:"false"`
		PointWriteListener    *bool  `default:"true"`
		PublishScheduleCOV    *bool  `default:"true"`
		PublishScheduleList   *bool  `default:"true"`
		ScheduleWriteListener *bool  `default:"true"`
	}
	Notification struct {
		Enable         *bool         `default:"false"`
		Frequency      time.Duration `default:"1m"`
		ResendDuration time.Duration `default:"1h"`
	}
}

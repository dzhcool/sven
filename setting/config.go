package setting

import (
	"log"
	"os"

	"github.com/safeie/goconfig"
)

type cfg struct {
	config *goconfig.ConfigFile
	env    string
}

var Config = &cfg{}

func (c *cfg) Set(key, value string) {
	c.config.AddOption(c.env, key, value)
}

func (c *cfg) Remove(key string) {
	c.config.RemoveOption(c.env, key)
}

func (c *cfg) Handle() *goconfig.ConfigFile {
	return c.config
}

func (c *cfg) GetString(key string) (string, error) {
	val, err := c.config.GetString(c.env, key)
	if err != nil {
		val, err = c.config.GetString("default", key)
	}
	return val, err
}

func (c *cfg) GetInt(key string) (int, error) {
	val, err := c.config.GetInt(c.env, key)
	if err != nil {
		val, err = c.config.GetInt("default", key)
	}
	return val, err
}

func (c *cfg) GetInt64(key string) (int64, error) {
	val, err := c.config.GetInt64(c.env, key)
	if err != nil {
		val, err = c.config.GetInt64("default", key)
	}
	return val, err
}

func (c *cfg) GetFloat(key string) (float64, error) {
	val, err := c.config.GetFloat(c.env, key)
	if err != nil {
		val, err = c.config.GetFloat("default", key)
	}
	return val, err
}

func (c *cfg) GetBool(key string) (bool, error) {
	val, err := c.config.GetBool(c.env, key)
	if err != nil {
		val, err = c.config.GetBool("default", key)
	}
	return val, err
}

func (c *cfg) MustString(key string, value string) string {
	val, err := c.config.GetString(c.env, key)
	if err != nil || val == "" {
		val = c.config.MustString("default", key, value)
	}
	return val
}

func (c *cfg) MustInt(key string, value int) int {
	val, err := c.config.GetInt(c.env, key)
	if err != nil || val == 0 {
		val = c.config.MustInt("default", key, value)
	}
	return val
}

func (c *cfg) MustInt64(key string, value int64) int64 {
	val, err := c.config.GetInt64(c.env, key)
	if err != nil || val == 0 {
		val = c.config.MustInt64("default", key, value)
	}
	return val
}

func (c *cfg) MustFloat(key string, value float64) float64 {
	val, err := c.config.GetFloat(c.env, key)
	if err != nil || val == 0.0 {
		val = c.config.MustFloat("default", key, value)
	}
	return val
}

func (c *cfg) MustBool(key string, value bool) bool {
	val, err := c.config.GetBool(c.env, key)
	if err != nil {
		val = c.config.MustBool("default", key, value)
	}
	return val
}

func findDefaultFile() string {
	root, err := os.Getwd()
	if err != nil {
		return ""
	}

	lookup := []string{"app.ini", "conf/app.ini", "../conf/app.ini"}

	for _, name := range lookup {
		filename := root + "/" + name
		if IsExist(filename) {
			return filename
		}
	}
	return ""
}

func InitSetting(filename, appEnv string) {
	Config.env = appEnv
	if Config.env == "" {
		Config.env = "default"
	}

	var err error
	if IsExist(filename) == false {
		log.Printf(" '%s' not exists and find the default directory \n", filename)
		filename = findDefaultFile()
		if filename == "" {
			Config.config = goconfig.NewConfigFile()
			return
		} else {
			log.Printf("use default directory file:%s \n", filename)
		}
	}

	Config.config, err = goconfig.ReadConfigFile(filename)
	if err != nil {
		log.Fatalf("can not load filename:%s \n", err)
		// init an empty config
		Config.config = goconfig.NewConfigFile()
	}

	initSetting(appEnv)
}

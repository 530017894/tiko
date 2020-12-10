package libs

import (
	"fmt"
	"path/filepath"
	"strings"

	logger "github.com/sirupsen/logrus"

	"github.com/jinzhu/configor"
)

var Config = struct {
	LogLevel string `default:"info" env:"Loglevel"`
	Bindata  bool   `default:"true" env:"Bindata"`
	Debug    bool   `default:"true" env:"Debug"`
	HTTPS    bool   `default:"false" env:"HTTPS"`
	Certpath string `default:"" env:"Certpath"`
	Certkey  string `default:"" env:"Certkey"`
	Port     int    `default:"8085" env:"PORT"`
	Host     string `default:"127.0.0.1" env:"Host"`
	DB       struct {
		Prefix   string `env:"DBPrefix" default:"iris_"`
		Name     string `env:"DBName" default:"goirisapi"`
		Adapter  string `env:"DBAdapter" default:"mysql"`
		Host     string `env:"DBHost" default:"localhost"`
		Port     string `env:"DBPort" default:"3306"`
		User     string `env:"DBUser" default:"root"`
		Password string `env:"DBPassword" default:""`
	}
	Redis struct {
		Host string `env:"RedisHost" default:"localhost"`
		Port string `env:"RedisPort" default:"6379"`
		Pwd  string `env:"RedisPwd" default:""`
	}
	Limit struct {
		Disable bool    `env:"LimitDisable" default:"true"`
		Limit   float64 `env:"LimitLimit" default:"1"`
		Burst   int     `env:"LimitBurst" default:"5"`
	}
	Qiniu struct {
		Enable    bool   `env:"QiniuEnable" default:"false"`
		Host      string `env:"QiniuHost" default:""`
		Accesskey string `env:"QiniuAccesskey" default:""`
		Secretkey string `env:"QiniuSecretkey" default:""`
		Bucket    string `env:"QiniuBucket" default:""`
	}
}{}

func InitConfig(p string) {

	configPath := filepath.Join(CWD(), "application.yml")
	if p != "" {
		configPath = p
	}

	if err := configor.Load(&Config, configPath); err != nil {
		logger.Println(fmt.Sprintf("Config Path:%s ,Error:%s", configPath, err.Error()))
		return
	}

	if Config.Debug {
		fmt.Println(fmt.Sprintf("配置项：%v", Config))
	}
}

func GetRedisUris() []string {
	addrs := make([]string, 0, 0)
	hosts := strings.Split(Config.Redis.Host, ";")
	ports := strings.Split(Config.Redis.Port, ";")
	for _, h := range hosts {
		for _, p := range ports {
			addrs = append(addrs, fmt.Sprintf("%s:%s", h, p))
		}
	}
	return addrs
}

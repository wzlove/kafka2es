package cmd

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"path"
	"whoops/kafka2es/src/model"
)

//Run 运行项目
func Run() {
	//获取配置文件
	if confFilePath, err := getConfPath(); err != nil {
		panic(err)
	} else {
		//解析配置文件并赋值
		gConf := new(model.GlobalConfig)
		if err = parseConfig(confFilePath, gConf); err != nil {
			panic(err)
		}
		//初始化服务
		InitService(gConf)
	}
}

//获取具体配置文件
func getConfPath() (string, error) {
	//配置文件 flag
	pflag.StringVar(&cfgFileInfo, "config", "",
		"config path, format:'--config /etc/config.yaml'")
	pflag.Parse()

	//从命令行中获取到配置文件
	if cfgFileInfo != "" {
		return cfgFileInfo, nil
	}

	//使用默认位置的配置文件
	if cfgFileInfo == "" {
		pwd, err := os.Getwd()
		if err != nil {
			os.Exit(1)
			return "", err
		}
		filePath := fmt.Sprintf("%s/src/etc/", pwd)
		fileName := "config.yaml"
		cfgFileInfo = path.Join(filePath, fileName)
	}
	return cfgFileInfo, nil
}

// 配置业务服务
// 读取配置文件配置service
func parseConfig(cfgFile string, gConfig *model.GlobalConfig) error {
	viper.SetConfigFile(cfgFile)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		pwd, _ := os.Getwd()
		return fmt.Errorf("|parseConfig| pwd:%s, err:%s", pwd, err)
	}
	return viper.Unmarshal(&gConfig)
}

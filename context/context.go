package context

import "gopkg.in/ini.v1"

type Config struct {
	DbHost string `ini:"mysql_host"`
	DbUserName string `ini:"mysql_user"`
	DbPassword string `ini:"mysql_psw"`
	DbSchema string `ini:"mysql_schema"`
	DbCharset string `ini:"mysql_charset"`
	DbMaxConn int `ini:"mysql_max_conn"`
	DbIdleConn int `ini:"mysql_idle_conn"`
	GrabMaxConcurrent int `ini:"grab_max_concurrent"`
}

//读取配置文件并转成结构体
func ReadConfig(path string) (Config, error) {
	var config Config
	conf, err := ini.Load(path)   //加载配置文件
	if err != nil {
		//log.Fatal("load config file fail!")
		return config, err
	}
	conf.BlockMode = false
	err = conf.MapTo(&config)   //解析成结构体
	if err != nil {
		//log.Fatal("mapto config file fail!")
		return config, err
	}
	return config, nil
}
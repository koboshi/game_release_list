package context

import (
	"gopkg.in/ini.v1"
	"path/filepath"
	"os"
	"errors"
	"github.com/koboshi/mole/database"
	"log"
	"sync"
	"github.com/koboshi/mole/work"
	"io/ioutil"
)

type Config struct {
	DbHost string `ini:"mysql_host"`
	DbUserName string `ini:"mysql_user"`
	DbPassword string `ini:"mysql_psw"`
	DbSchema string `ini:"mysql_schema"`
	DbCharset string `ini:"mysql_charset"`
	DbMaxConn int `ini:"mysql_max_conn"`
	DbIdleConn int `ini:"mysql_idle_conn"`
	GrabMaxConcurrent int `ini:"grab_max_concurrent"`
	LogOn int `ini:"log_on"`
}

var ErrLoadConf = errors.New("load gel.conf fail")
var ErrDbConnect = errors.New("connect database fail")
var ErrGoroutinePool = errors.New("create goroutine pool fail")

var appConfig Config
var db *database.Database
var logger *log.Logger

func init() {
	log.SetPrefix("LOG:")
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	var err error
	appConfig, err = load(filepath.Dir(os.Args[0]) + "/conf/gel.conf")
	if err != nil {
		panic(ErrLoadConf)
	}
}

func load(path string) (Config, error) {
	var config Config
	conf, err := ini.Load(path)   //加载配置文件
	if err != nil {
		return config, err
	}
	conf.BlockMode = false
	err = conf.MapTo(&config)   //解析成结构体
	if err != nil {
		return config, err
	}
	return config, nil
}

func GetConfig() (Config) {
	return appConfig
}

func GetGoroutinePool() (*work.Pool) {
	pool, err := work.New(appConfig.GrabMaxConcurrent)
	if err != nil {
		panic(ErrGoroutinePool)
	}
	return pool
}

func Logger() (*log.Logger) {
	l := new(sync.Mutex)
	l.Lock()
	defer l.Unlock()

	if logger != nil {
		return logger
	}

	if appConfig.LogOn != 0 {
		logger = log.New(os.Stdout, "[LOG]", log.Ldate | log.Ltime | log.Llongfile)
	}else {
		logger = log.New(ioutil.Discard, "[LOG]", log.Ldate | log.Ltime | log.Llongfile)
	}
	return logger
}

func GetDatabase() (*database.Database) {
	l := new(sync.Mutex)
	l.Lock()
	defer l.Unlock()

	if db != nil {
		return db
	}

	host := appConfig.DbHost
	username := appConfig.DbUserName
	password := appConfig.DbPassword
	schema := appConfig.DbSchema
	charset := appConfig.DbCharset
	conn, err := database.New(host, username, password, schema, charset)
	if err != nil {
		log.Println(err)
		panic(ErrDbConnect)
	}
	return conn
}
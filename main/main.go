package main

import (
	"fmt"
	"os"
	"time"
	"log"
	"path/filepath"
	"flag"
	"github.com/koboshi/game_release_list/context"
	"github.com/koboshi/game_release_list/engine"
	"github.com/Jeffail/tunny"
	"github.com/koboshi/go-tool"
)

var config context.Config

func init() {
	//读取配置文件并加载
	var err error
	config, err = context.ReadConfig(filepath.Dir(os.Args[0]) + "/conf/gel.conf")
	if err != nil {
		panic(err)
	}
}

func initRoutinePool(n int) (*tunny.Pool) {
	pool := tunny.NewFunc(n, func (arg interface{}) interface{} {
		rArg, ok := arg.(engine.GrabArg)
		if ok {
			engine.GrabReleaseList(rArg.Database, rArg.Year, rArg.Month)
		}else {
			panic("tunny pool arg error")
		}
		return true
	})
	log.Println(fmt.Sprintf("Max concurrent: %d", n))
	return pool
}

func initDatabase(config context.Config) (*tool.Database) {
	host := config.DbHost
	username := config.DbUserName
	password := config.DbPassword
	dbname := config.DbSchema
	charset := config.DbCharset
	customParams := make(map[string]string)
	customParams["readTimeout"] = "10s"
	customParams["writeTimeout"] = "10s"
	database := new(tool.Database)
	database.Connect(host, username, password, dbname, charset, customParams)
	database.SetPool(config.DbMaxConn, config.DbIdleConn, 0)
	return database
}

//执行游戏发售数据爬取
func grab(argYear *int, argMonth *int, argAll *bool) {
	//检查参数
	now := time.Now()
	if *argYear == 0 {
		//没有从外部取得正确年份，则为当前年份
		*argYear, _, _ = now.Date()
	}
	if *argMonth < 1 || *argMonth > 12 {
		//没有从外部取得正确月份，则为当前月
		_, tmp, _ := now.Date()
		*argMonth = int(tmp)
	}

	//创建协程池
	pool := initRoutinePool(config.GrabMaxConcurrent)
	defer pool.Close()

	//创建数据库连接池
	database := initDatabase(config)
	defer database.Close()

	//执行游戏发售日期抓取
	//var n sync.WaitGroup
	if *argAll {
		log.Println(fmt.Sprintf("Start all"))
		//抓取2010年至下一年的每一个月
		endYear, _, _ := now.Date()
		endYear++
		for i := 2010; i <= endYear; i++ {
			for j := 1; j <= 12; j ++ {
				//fmt.Println(fmt.Sprintf("multi_dest:%d %d", i, j))
				pool.Process(engine.GrabArg{database, i, j})
			}
		}
	}else {
		//只抓取指定年月
		log.Println(fmt.Sprintf("Start %d %d", *argYear, *argMonth))
		//fmt.Println(fmt.Sprintf("single_dest:%d %d", year, month))
		pool.Process(engine.GrabArg{database, *argYear, *argMonth})
	}
	//n.Wait()
	log.Println(fmt.Sprintf("Done"))
}

//执行游戏发售数据结构化整理
func cron() {
	//创建数据库连接池
	database := initDatabase(config)
	defer database.Close()

	//循环读取整个爬虫表

	//整理数据
}

func main() {
	//获取参数
	argCron := flag.Bool("cron", false, "cron")
	argYear := flag.Int("year", 0, "specified year")
	argMonth := flag.Int("month", 0, "specified month")
	argAll := flag.Bool("all", false, "grab all")
	flag.Parse()
	if *argCron {
		//执行定时任务
		cron()
	}else {
		//执行爬取
		grab(argYear, argMonth, argAll)
	}

}
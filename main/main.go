package main

import (
	"fmt"
	"os"
	"time"
	"sync"
	"log"
	"path/filepath"
	"flag"
	"github.com/koboshi/game_release_list/context"
	"github.com/koboshi/game_release_list/engine"
)

func main() {
	var err error
	now := time.Now()
	//获取参数
	argYear := flag.Int("year", 0, "specified year")
	argMonth := flag.Int("month", 0, "specified month")
	argAll := flag.Bool("all", false, "grab all")
	flag.Parse()
	//检查参数
	if *argYear == 0 {
		//没有从外部取得正确年份，则为当前年份
		*argYear, _, _ = now.Date()
	}
	if *argMonth < 1 || *argMonth > 12 {
		//没有从外部取得正确月份，则为当前月
		_, tmp, _ := now.Date()
		*argMonth = int(tmp)
	}

	//读取配置文件
	config, err := context.ReadConfig(filepath.Dir(os.Args[0]) + "/conf/gel.conf")
	if err != nil {
		panic(err)
	}

	var n sync.WaitGroup
	if *argAll {
		log.Println(fmt.Sprintf("Start all"))
		//抓取2010年至下一年的每一个月
		endYear, _, _ := now.Date()
		endYear++
		for i := 2010; i <= endYear; i++ {
			for j := 1; j <= 12; j ++ {
				//fmt.Println(fmt.Sprintf("multi_dest:%d %d", i, j))
				go engine.GrabReleaseList(&config, i, j, &n)
				n.Add(1)
			}
		}
	}else {
		//只抓取指定年月
		log.Println(fmt.Sprintf("Start %d %d", *argYear, *argMonth))
		//fmt.Println(fmt.Sprintf("single_dest:%d %d", year, month))
		go engine.GrabReleaseList(&config, *argYear, *argMonth, &n)
		n.Add(1)
	}
	n.Wait()
	log.Println(fmt.Sprintf("Done"))
}
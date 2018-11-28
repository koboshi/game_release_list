package main

import (
	"fmt"
	"time"
	"log"
	"flag"
	"github.com/koboshi/game_release_list/context"
	"github.com/koboshi/game_release_list/engine"
)

var NowTime time.Time

func init() {
	NowTime = time.Now()
}

//执行游戏发售数据爬取
func exec(argYear *int, argMonth *int, argAll *bool) {
	//创建goroutine池
	pool := context.GetGoroutinePool()
	defer pool.Shutdown()

	//执行游戏发售日期抓取
	if *argAll {
		log.Println(fmt.Sprintf("Start all"))
		//抓取2010年至下一年的每一个月
		endYear, _, _ := NowTime.Date()
		endYear++
		for i := 2010; i <= endYear; i++ {
			for j := 1; j <= 12; j ++ {
				pool.Run(func() {
					engine.GrabReleaseList(i, j)
				})
			}
		}
		defer log.Println(fmt.Sprintf("Done all"))
	}else {
		//只抓取指定年月
		log.Println(fmt.Sprintf("Start %d %d", *argYear, *argMonth))
		engine.GrabReleaseList(*argYear, *argMonth)
		log.Println(fmt.Sprintf("Done %d %d", *argYear, *argMonth))
	}
}

func main() {
	//获取参数
	argYear := flag.Int("year", 0, "specified year")
	argMonth := flag.Int("month", 0, "specified month")
	argAll := flag.Bool("all", false, "exec all")
	flag.Parse()
	//修正输入参数
	if *argYear == 0 {
		//没有从外部取得正确年份，则为当前年份
		*argYear, _, _ = NowTime.Date()
	}
	if *argMonth < 1 || *argMonth > 12 {
		//没有从外部取得正确月份，则为当前月
		_, tmp, _ := NowTime.Date()
		*argMonth = int(tmp)
	}
	//执行
	exec(argYear, argMonth, argAll)
}
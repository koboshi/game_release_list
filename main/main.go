package main

import (
	"fmt"
	"time"
	"log"
	"github.com/koboshi/game_release_list/context"
	"github.com/koboshi/game_release_list/engine"
	"flag"
)

var NowTime time.Time
var ArgYear *int
var ArgMonth *int
var ArgAll *bool

func init() {
	NowTime = time.Now()
}

func main() {
	//获取参数
	ArgYear = flag.Int("year", 0, "specified year")
	ArgMonth = flag.Int("month", 0, "specified month")
	ArgAll = flag.Bool("all", false, "exec all")
	flag.Parse()
	//修正输入参数
	if *ArgYear == 0 {
		//没有从外部取得正确年份，则为当前年份
		*ArgYear, _, _ = NowTime.Date()
	}
	if *ArgMonth < 1 || *ArgMonth > 12 {
		//没有从外部取得正确月份，则为当前月
		_, tmp, _ := NowTime.Date()
		*ArgMonth = int(tmp)
	}
	//执行
	exec(ArgYear, ArgMonth, ArgAll)
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
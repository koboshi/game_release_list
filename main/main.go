package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	var year int
	var month int
	var err error
	//获取外部参数
	if len(os.Args) > 1 {
		year, err = strconv.Atoi(os.Args[1])
		if err != nil {
			year = 0//转换失败，设置年份为0
		}
	}
	if len(os.Args) > 2 {
		month, err = strconv.Atoi(os.Args[2])
		if err != nil {
			month = 0//转换失败，设置月份为0
		}
	}

	now := time.Now()
	if year == 0 {
		//没有从外部取得年份，以当前年份开始
		year, _, _ = now.Date()
	}
	if month < 1 || month > 12 {
		//没有从外部取得月份，以1月开始
		_, tmp, _ := now.Date()
		month = int(tmp)
	}

	//爬取游戏发售表
	//http://www.a9vg.com/game/release?genres=&region=&platform=&year={year}&month={month}&quarter=
	for i := month; i <= 12; i++ {
		url := fmt.Sprintf("http://www.a9vg.com/game/release?genres=&region=&platform=&year=%d&month=%d&quarter=", year, i)
		fmt.Println(url)
		//发起请求
		//获取数据
		//分析组装
		//写入数据至mysql
	}

	//输出近期游戏发售时间
}
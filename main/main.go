package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"github.com/koboshi/go-tool"
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

	//链接数据库
	host := "127.0.0.1:3306"
	username := "admin"
	password := "802927"
	dbname := "life"
	charset := "utf8mb4,utf8"
	customParams := make(map[string]string)
	customParams["readTimeout"] = "10s"
	customParams["writeTimeout"] = "10s"
	database := new(tool.Database)
	database.Connect(host, username, password, dbname, charset, customParams)
	defer database.Close()

	//爬取游戏发售表
	//http://www.a9vg.com/game/release?genres=&region=&platform=&year={year}&month={month}&quarter=
	for i := month; i <= 12; i++ {
		//构造url
		url := fmt.Sprintf("http://www.a9vg.com/game/release?genres=&region=&platform=&year=%d&month=%d&quarter=", year, i)
		//fmt.Println(fmt.Sprintf(" process:%d %d", year, i))
		//fmt.Println("request:" + url)
		//发起请求
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}
		req.Header.Add("Cookie", "WxSS_a648_saltkey=gBKd9Ob3; WxSS_a648_lastvisit=1531750974; taihe=6bdc704084fdc7e7406e6ce106bb5050; Hm_lvt_68e4f3f877acf23e052991a583acf43e=1531754575,1532009703,1532266194; Hm_lpvt_68e4f3f877acf23e052991a583acf43e=1532266194; taihe_session=efcdb2fd68abb580885ad125bb9517e2")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
		client := &http.Client{}
		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != 200 {
			fmt.Println(fmt.Sprintf("status code error: %d %s", resp.StatusCode, resp.Status))
			continue
		}
		//分析html
 		doc, err := goquery.NewDocumentFromReader(resp.Body)
 		if err != nil {
			fmt.Println(fmt.Sprintf("goquery error: %s", err))
 			continue
		}
		doc.Find(".area_column_left .saletimebox .saleinfobox").Each(func (i int, s *goquery.Selection) {
			releaseDate := s.Find("h6").Text()
			releaseDate = strings.Replace(releaseDate, "发行时间：", "", -1)
			//fmt.Println(releaseDate)
			s.Find("dl").Each(func (j int, subS *goquery.Selection) {
				idHref, exist := subS.Find(".ddwz1 a").Attr("href")
				if !exist {
					return
				}
				idHref = strings.Replace(idHref,"/game/", "", -1)
				id, err := strconv.Atoi(idHref)
				if err != nil {
					return
				}
				title := subS.Find(".ddwz1 a").Text()
				platform, exist := subS.Find(".ddwz1 span em").Attr("class")
				if !exist {
					platform = "unknown"
				}
				platform = strings.Replace(platform, "icon_", "", -1)
				gameType := subS.Find(".ddwz2 span").Eq(0).Find("em").Text()
				lang := subS.Find(".ddwz2 span").Eq(1).Find("em").Text()
				company := subS.Find(".ddwz2 span").Eq(2).Find("em").Text()

				fmt.Println(year)
				fmt.Println(i)
				fmt.Println(fmt.Sprintf("游戏编号: %d", id))
				fmt.Println(fmt.Sprintf("发售日期: %s", releaseDate))
				fmt.Println(fmt.Sprintf("游戏名称: %s", title))
				fmt.Println(fmt.Sprintf("游戏类型: %s", gameType))
				fmt.Println(fmt.Sprintf("支持语种: %s", lang))
				fmt.Println(fmt.Sprintf("发行公司: %s", company))
				fmt.Println(fmt.Sprintf("运行平台: %s", platform))
				fmt.Println()

				//写入数据至mysql
				data := make(map[string]interface{})
				data["out_id"] = id
				data["name"] = title
				data["release_date"] = releaseDate
				data["type"] = gameType
				data["language"] = lang
				data["company"] = company
				data["platform"] = platform
				data["add_time"] = time.Now().Format("2006-01-02 15:04:05")
				data["edit_time"] = time.Now().Format("2006-01-02 15:04:05")
				database.Ignore(data, "game")
			})
		})
	}
}
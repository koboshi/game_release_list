package engine

import (
	"fmt"
	"net/http"
	"log"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
	"time"
	"github.com/koboshi/mole/database"
)

//抓取指定年月的发售列表
func GrabReleaseList(database *database.Database, year int, month int) {
	//爬取游戏发售表
	//http://www.a9vg.com/game/release?genres=&region=&platform=&year={year}&month={month}&quarter=
	//构造url
	url := fmt.Sprintf("http://www.a9vg.com/game/release?genres=&region=&platform=&year=%d&month=%d&quarter=", year, month)
	//构造请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	//构造请求头
	req.Header.Add("Cookie", "WxSS_a648_saltkey=gBKd9Ob3; WxSS_a648_lastvisit=1531750974; taihe=6bdc704084fdc7e7406e6ce106bb5050; Hm_lvt_68e4f3f877acf23e052991a583acf43e=1531754575,1532009703,1532266194; Hm_lpvt_68e4f3f877acf23e052991a583acf43e=1532266194; taihe_session=efcdb2fd68abb580885ad125bb9517e2")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
	client := &http.Client{}
	resp, err := client.Do(req)//发起http请求
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal(fmt.Sprintf("status code error: %d %s", resp.StatusCode, resp.Status))
		return
	}
	//分析html
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(fmt.Sprintf("goquery error: %s", err))
		return
	}
	doc.Find(".area_column_left .saletimebox .saleinfobox").Each(func (i int, s *goquery.Selection) {
		releaseDate := s.Find("h6").Text()
		releaseDate = strings.Replace(releaseDate, "发行时间：", "", -1)
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

			//写入数据
			addReleaseInfo(database, id, platform, title, releaseDate, gameType, lang, company)
		})
	})
}

//写入发售数据至数据库
func addReleaseInfo(database *database.Database, id int, platform string, title string, releaseDate string, gameType string, lang string, company string) {
	//检查是否存在
	sql := "SELECT COUNT(*) AS `count` FROM game WHERE out_id = ? AND platform = ?"
	var count int
	database.QueryOne(sql, id, platform).Scan(&count)
	if count > 0 {
		//存在则更新发售日期和语种
		data := make(map[string]interface{})
		data["release_date"] = releaseDate
		data["type"] = gameType
		data["language"] = lang
		data["company"] = company
		whereStr := "out_id = ? AND platform = ?"
		database.Update(data, "game", whereStr, id, platform)
	}else {
		//不存在则写入数据至mysql
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
	}
}
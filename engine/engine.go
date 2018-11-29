package engine

import (
	"fmt"
	"log"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
	"github.com/koboshi/game_release_list/model"
	"github.com/koboshi/game_release_list/utils"
)

//抓取指定年月的发售列表
func GrabReleaseList(year int, month int) {
	resp, err := utils.RequestA9VGReleaseList(year, month)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

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
			model.AddReleaseInfo(id, platform, title, releaseDate, gameType, lang, company)
		})
	})
}
package model

import (
	"github.com/koboshi/game_release_list/context"
	"time"
)

//写入发售数据至数据库
func AddReleaseInfo(id int, platform string, title string, releaseDate string, gameType string, lang string, company string) {
	db := context.GetDatabase()
	//检查是否存在
	sql := "SELECT COUNT(*) AS `count` FROM game WHERE out_id = ? AND platform = ?"
	var count int
	db.QueryOne(sql, id, platform).Scan(&count)
	if count > 0 {
		//存在则更新发售日期和语种
		data := make(map[string]interface{})
		data["release_date"] = releaseDate
		data["type"] = gameType
		data["language"] = lang
		data["company"] = company
		whereStr := "out_id = ? AND platform = ?"
		db.Update(data, "game", whereStr, id, platform)
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
		db.Ignore(data, "game")
	}
}
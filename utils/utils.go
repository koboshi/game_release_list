package utils

import (
	"net/http"
	"fmt"
)

func RequestA9VGReleaseList(year int, month int) (*http.Response, error) {
	//爬取游戏发售表
	//http://www.a9vg.com/game/release?genres=&region=&platform=&year={year}&month={month}&quarter=
	//构造url
	url := fmt.Sprintf("http://www.a9vg.com/game/release?genres=&region=&platform=&year=%d&month=%d&quarter=", year, month)
	//构造请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	//构造请求头
	req.Header.Add("Cookie", "WxSS_a648_saltkey=gBKd9Ob3; WxSS_a648_lastvisit=1531750974; taihe=6bdc704084fdc7e7406e6ce106bb5050; Hm_lvt_68e4f3f877acf23e052991a583acf43e=1531754575,1532009703,1532266194; Hm_lpvt_68e4f3f877acf23e052991a583acf43e=1532266194; taihe_session=efcdb2fd68abb580885ad125bb9517e2")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
	client := &http.Client{}
	resp, err := client.Do(req)//发起http请求
	if err != nil {
		return nil, err
	}
	return resp, nil
}
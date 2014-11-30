package main

import (
	"environmentspider/airspider"
	"log"
	"time"
)

func main() {
	count := 1
	url := "http://datacenter.mep.gov.cn/report/air_daily/airDairyCityHourMain.jsp"
	env := airspider.NewAirSpider(url)
	for {
		log.Println("正在进行第", count, "次抓取")
		env.Crawl()
		count++
		time.Sleep(20 * time.Minute)
	}
}

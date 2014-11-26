package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type CityInfo struct {
	Name      string
	Time      string
	AQI       int
	Level     string
	Pollution string
}

func (city *CityInfo) ToString() string {
	return city.Name + "	" + city.Time + "	" + strconv.Itoa(city.AQI) + "	" + city.Level
}

func (city *CityInfo) SaveDataToFile() {
	log.Println("save to file", city.Name)
	f, err := os.OpenFile(city.Name+".txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(city.ToString())
	f.WriteString("\n")
}

func ParseData(s *goquery.Selection) {
	log.Println("开始抓取城市列表和地址")
	s.Find("a").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		GetUrl(s.Text(), url)
	})
}

func GetUrl(city, url string) {
	name, _ := iconv.ConvertString(city, "GB18030", "UTF-8")
	log.Println("正在获取", name, "的数据")
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}
	table := doc.Find("table")
	for i := range table.Nodes {
		if single := table.Eq(i); single.HasClass("font12") {
			if frame := single.Has("IFRAME"); frame != nil {
				ParseFrame(frame)
			}
		}
	}
}

func GetData(url string) {
	if !strings.HasPrefix(url, "http://") {
		url = "http://datacenter.mep.gov.cn/report/air_daily/" + url
	}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}

	tr := doc.Find("table").Last().Find("tr")
	length := tr.Length()
	for i := 1; i < length; i++ {
		data := tr.Eq(i).Find("td")
		convert, _ := iconv.NewConverter("GB18030", "UTF-8")
		date := data.Eq(1).Text()
		addr, _ := convert.ConvertString(data.Eq(2).Text())
		aqi_s, _ := convert.ConvertString(data.Eq(3).Text())
		aqi, _ := strconv.Atoi(aqi_s)
		level, _ := convert.ConvertString(data.Eq(4).Text())
		pollution, _ := convert.ConvertString(data.Eq(5).Text())
		m := &CityInfo{
			Name:      addr,
			Time:      date,
			AQI:       aqi,
			Level:     level,
			Pollution: pollution,
		}
		m.SaveDataToFile()
	}
}

func ParseFrame(frames *goquery.Selection) {
	frame := frames.Find("IFRAME")
	for i := range frame.Nodes {
		single := frame.Eq(i)
		if src, ok := single.Attr("src"); ok {
			if src == "" {
				continue
			}
			GetData(src)
		}
	}
}

func main() {
	count := 1
	url := "http://datacenter.mep.gov.cn/report/air_daily/airDairyCityHourMain.jsp"
	for {
		log.Println("正在进行第", count, "次抓取")
		doc, err := goquery.NewDocument(url)
		if err != nil {
			panic(err)
		}
		table := doc.Find("table")
		for i := range table.Nodes {
			if single := table.Eq(i); single.HasClass("font12") {
				if single.Find("table").Text() != "" {
					continue
				}
				ParseData(single)
			}
		}
		log.Println("第", count, "次抓取完毕")
		count++
		time.Sleep(10 * time.Minute)
	}
}

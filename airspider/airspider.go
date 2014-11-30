package airspider

import (
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
)

type AirSpider struct {
	Url string
}

func NewAirSpider(url string) *AirSpider {
	return &AirSpider{
		Url: url,
	}
}

func (env *AirSpider) Crawl() {
	cities, urls := GetAllCityUrl(env.Url)
	for i := range urls {
		cityUTF8, _ := iconv.ConvertString(cities[i], "GB18030", "UTF-8")
		log.Println("正在获取", cityUTF8, "的环境数据")
		GetCityInfo(urls[i])
	}
}
func GetAllCityUrl(url string) ([]string, []string) {
	urls := make([]string, 0, 200)
	cities := make([]string, 0, 200)
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
			single.Find("a").Each(func(i int, s *goquery.Selection) {
				url, ok := s.Attr("href")
				if ok {
					urls = append(urls, url)
					cities = append(cities, s.Text())
				}
			})
		}
	}
	return cities, urls
}

func GetRealUrl(url string) string {
	result := ""
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}
	table := doc.Find("table")
	for i := range table.Nodes {
		single := table.Eq(i)
		if !single.HasClass("font12") {
			continue
		}
		frames := single.Has("IFRAME")
		if frames == nil {
			continue
		}
		frame := frames.Find("IFRAME")
		for i := range frame.Nodes {
			single := frame.Eq(i)
			if src, ok := single.Attr("src"); ok {
				if src != "" {
					result = src
					return result
				}
			}
		}
	}
	return result
}

func GetCityData(url string) []*AirInfo {
	cities := make([]*AirInfo, 0, 200)
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
		m := &AirInfo{
			Name:      addr,
			Time:      date,
			AQI:       aqi,
			Level:     level,
			Pollution: pollution,
		}
		cities = append(cities, m)
	}
	return cities
}

func GetCityInfo(url string) {
	convert, _ := iconv.NewConverter("GB18030", "UTF-8")
	urlUTF8, _ := convert.ConvertString(url)
	log.Println("源地址：", urlUTF8)
	realUrl := GetRealUrl(url)
	realUrlUTF8, _ := convert.ConvertString(realUrl)
	log.Println("真实地址：", realUrlUTF8)
	cities := GetCityData(realUrl)
	for _, value := range cities {
		value.SaveDataToFile()
	}
}

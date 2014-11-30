package airspider

import (
	"os"
	"strconv"
)

type AirInfo struct {
	Name      string
	Time      string
	AQI       int
	Level     string
	Pollution string
}

func (city *AirInfo) ToString() string {
	return city.Name + "	" + city.Time + "	" + strconv.Itoa(city.AQI) + "	" + city.Level + "	" + city.Pollution
}

func (city *AirInfo) SaveDataToFile() {
	//log.Println("save to file", city.Name)
	f, err := os.OpenFile(city.Name+".txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(city.ToString())
	f.WriteString("\n")
}

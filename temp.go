package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var TempList *GlobalTempList

type TempMeasurePoint struct {
	Name     string
	Lat, Lon float64
	CdkId    string  `json:"cdk_id"`
	Temp     float64 `json:"layers"."rws.temp"."data"."waarde"`
}

type GlobalTempList struct {
	Points []TempMeasurePoint
	Mutex  sync.RWMutex
}

type TempParse struct {
	Status  string `json:"status"`
	Results []struct {
		CdkId string `json:"cdk_id"`
		Name  string `json:"name"`
		Geom  struct {
			Coordinates [2]float64 `json:"coordinates"`
		} `json:"geom"`
		Layers struct {
			RWSTemp struct {
				Data struct {
					Waarde   string `json:"waarde"`
					MeetTijd string `json:"meettijd"`
				} `json:"data"`
			} `json:"rws.temp"`
		} `json:"Layers"`
	} `json:"results"`
}

func get_temperatures() []TempMeasurePoint {
	points := []TempMeasurePoint{}
	resp, err := http.Get("http://api.citysdk.waag.org/nodes?layer=rws.temp&geom&per_page=1000")
	if err != nil {
		log.Println(err)
		return points
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return points
	}
	var data = new(TempParse)
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err)
		return points
	}
	if data.Status != "success" {
		log.Println(data.Status)
		return points
	}
	for _, point := range data.Results {
		p := new(TempMeasurePoint)
		p.Name = point.Name
		p.CdkId = point.CdkId
		f, err := strconv.ParseFloat(point.Layers.RWSTemp.Data.Waarde, 64)
		if err != nil {
			log.Println(err)
			continue
		}
		p.Temp = f
		p.Lat = point.Geom.Coordinates[0]
		p.Lon = point.Geom.Coordinates[1]
		points = append(points, *p)
	}
	return points
}

func BackgroundTempUpdate() {
	for {
		points := get_temperatures()
		TempList.Mutex.Lock()
		TempList.Points = points
		TempList.Mutex.Unlock()
		time.Sleep(1 * time.Minute)
	}
}

func FindClosestTemp(Lat, Lon float64) float64 {
	closest := 0
	closestdist := math.MaxFloat64
	for i, point := range TempList.Points {
		dist := math.Sqrt(math.Pow(Lat-point.Lat, 2.) + math.Pow(Lon-point.Lon, 2.))
		if dist < closestdist {
			closest = i
			closestdist = dist
		}
	}
	return TempList.Points[closest].Temp
}

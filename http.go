package main

import (
	"net/http"
	"strconv"
    "html/template"
)

type Page struct {
	Temp float64
	Regen string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "template/index.html")
}

func DataHandler(w http.ResponseWriter, r *http.Request) {
    TempList.Mutex.RLock()
	defer TempList.Mutex.RUnlock()
	lat, err := strconv.ParseFloat(r.FormValue("Lat"), 64)
	if err != nil {
		http.Error(w, "Latitude not given", 417)
		return
	}
	lon, err := strconv.ParseFloat(r.FormValue("Lon"), 64)
	if err != nil {
		http.Error(w, "Longitude not given", 417)
		return
	}
	temp := FindClosestTemp(lat, lon)
	p := new(Page)
	p.Temp = temp
	p.Regen = getRain(lat, lon)
	t, err := template.ParseFiles("template/page.html")
	if err != nil {
	    http.Error(w, "AARG", 500)
	}
	t.Execute(w, p)
}

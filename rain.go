package main

import (
    "net/http"
    "fmt"
    "strings"
    "log"
    "io/ioutil"
)


func getRain(lat, lon float64) string {
    //return "Geen regen de komende twee uur."
    resp, err := http.Get(fmt.Sprintf("http://gps.buienradar.nl/getrr.php?lat=%f&lon=%f", lat, lon))
	if err != nil {
		log.Println(err)
        return "Geen regen de komende twee uur."
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
        return "Geen regen de komende twee uur."
	}
	for i, ammount := range strings.Split(string(body), " ") {
	    fmt.Println(i, ammount)
	}
    return "Geen regen de komende twee uur."
}

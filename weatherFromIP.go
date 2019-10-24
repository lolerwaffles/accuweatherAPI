package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", returnText)
	http.ListenAndServe(":8086", nil)
}

func returnText(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	remoteIP := ReadUserIP(r)
	fmt.Println("Accepted connection from " + remoteIP)
	w.Write([]byte(getWeather(remoteIP)))

}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	if idx := strings.IndexByte(IPAddress, ':'); idx >= 0 {
		IPAddress = IPAddress[:idx]
	}
	return IPAddress
}

func getWeather(ip string) string {
	accuweatherAPI := ""
	wURL := "http://dataservice.accuweather.com/currentconditions/v1/" + getLocationKey(ip, accuweatherAPI) + "?apikey=" + accuweatherAPI
	resp, err := http.Get(wURL)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body)
}

type locationKey struct {
	Key string `json:"Key"`
}

func getLocationKey(ip string, accuweatherAPI string) string {
	keyURL := "http://dataservice.accuweather.com/locations/v1/cities/geoposition/search?apikey=" + accuweatherAPI + "&q=" + getLoc(ip)
	resp, err := http.Get(keyURL)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	response := locationKey{}
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return response.Key
}

func getLoc(ip string) string {
	url := "https://ipapi.co/" + ip + "/latlong/"
	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body)
}

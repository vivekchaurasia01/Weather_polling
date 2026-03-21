package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	pollInterval = time.Second * 10
)
const (
	endpoint = "https://api.open-meteo.com/v1/forecast"//?latitude=52.52&longitude=13.41&hourly=temperature_2m"    //A string variable storing API URL
)

type WeatherData struct {
	Elevation float64 `json:"elevation"`
	Hourly map[string]any `json:"hourly"`
}

type WPoller struct {
	closech chan struct{}
}
func NewWPoller() *WPoller{
	return &WPoller{
		closech: make(chan struct{}),
	}
}

func(wp WPoller) close() {
	// now we can return Empty string to the channel or just close it.
	close(wp.closech)
}

func (wp WPoller) start() {
	fmt.Println("Startiing wpoller")

	ticker := time.NewTicker(pollInterval)  //ticker that sends a signal on a channel (ticker.C) at fixed intervals.
	free:
	for {
		select {
		case <-ticker.C:
			data,err := GetWeatherResults(52.52, 13.41)
			if err != nil {
				log.Fatal(err)
			}
			if err = wp.handleData(data); err != nil {
				log.Fatal(err)
			}
		case <- wp.closech:
			break free
		}
	}
}

func (wp WPoller) handleData(data *WeatherData) error {
	fmt.Println(data)
	return nil
}

func main(){
	wpoller := NewWPoller()
	wpoller.start()
}

func GetWeatherResults(lat, long float64) (*WeatherData, error) {
	// we shorted the Uri So we did Sprintf to write formatted data into a string.
	uri := fmt.Sprintf("%s?latitude=%.2f&longitude=%.2f&hourly=temperature_2m",endpoint,lat,long)

	req, err := http.NewRequest("GET",uri,nil)
	if err != nil{
		log.Fatal(err)
	}
	// Now we can actually create a client and use that client to Perforn Request.. We did it because we want to control TimeOut and other things..
	client := &http.Client{}
	resp, err := client.Do(req)  //Do() sends an HTTP request and returns an HTTP response, following policy (such as redirects, cookies, auth) as configured on the client.
	if err != nil{
		log.Fatal(err)
	}

	// resp, err := http.Get(uri)  //Sends an HTTP GET request,Return: resp → pointer to http.Response and err → error if request fails
	// if err != nil{
	// 	// log.Fatal(err)
	// 	return nil,err
	// }

	// Lets Explore The Response...
	// var data map[string]map[string]any

	defer resp.Body.Close() // important to close..

	var data WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil{
		// log.Fatal(err)
		return nil,err
	}
	// fmt.Println(data["hourly"]["temperature_2m"])
	return &data,nil
}
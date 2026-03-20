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

type Data struct {
	Elevation float64 `json:"elevation"`
	Hourly map[string]any `json:"hourly"`
}

func main(){
	ticker := time.NewTicker(pollInterval)  //Creates a ticker that sends a signal on a channel (ticker.C) at fixed intervals.
	for {
		<-ticker.C  //This is a blocking receive from the ticker channel.
		data, err := GetWeatherResults(52.52,13.41)
		if err != nil{
			log.Fatal(err)
		}
		fmt.Println(data)
	}
	// data, err := GetWeatherResults(52.52,13.41)
	// if err != nil{
	// 	log.Fatal(err)
	// }
	// fmt.Println(data)
}

func GetWeatherResults(lat, long float64) (*Data, error) {
	// we shorted the Uri So we did it to write formatted data into a string.
	uri := fmt.Sprintf("%s?latitude=%.2f&longitude=%.2f&hourly=temperature_2m",endpoint,lat,long)


	// Rather than calling http.Get lets first made request ....
	req, err := http.NewRequest("GET",uri,nil)
	if err != nil{
		log.Fatal(err)
	}
	// Now we can actually create a client and use that client to Perforn Request.. We did it because we want to control TimeOut and other things..
	client := &http.Client{}
	resp, err := client.Do(req)  //Do sends an HTTP request and returns an HTTP response, following policy (such as redirects, cookies, auth) as configured on the client.
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

	defer resp.Body.Close()

	var data Data
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil{
		// log.Fatal(err)
		return nil,err
	}
	// fmt.Println(data["hourly"]["temperature_2m"])
	return &data,nil
}
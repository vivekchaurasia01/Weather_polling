package main

import (
	"fmt"
	"log"
	"net/http"
)

func main(){
	endpoint := "https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&hourly=temperature_2m"    //A string variable storing API URL
	resp,err := http.Get(endpoint)  //Sends an HTTP GET request,Return: resp → pointer to http.Response and err → error if request fails
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(resp)
}
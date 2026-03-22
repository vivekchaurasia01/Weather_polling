package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"time"
)

var (
	pollInterval = time.Second * (60 * 60 * 4)
)

const (
	endpoint = "https://api.open-meteo.com/v1/forecast" //?latitude=52.52&longitude=13.41&hourly=temperature_2m"    //A string variable storing API URL
)

type Sender interface {
	Send(*WeatherData) error
}
type SMSSender struct {
	number string
}

func NewSMSSender(number string) *SMSSender {
	return &SMSSender{
		number: number,
	}
}

func (s *SMSSender) Send(data *WeatherData) error {
	fmt.Println("Sending weather to the number: ", s.number)
	return nil
}

type EmailSender struct {
	from     string
	password string
	to       []string // slice so we can send email to multiple users.
}

func NewEmailSender(from, password string, to []string) *EmailSender {
	return &EmailSender{
		from:     from,
		password: password,
		to:       to,
	}
}
func (e *EmailSender) Send(data *WeatherData) error {
	auth := smtp.PlainAuth("", e.from, e.password, "smtp.gmail.com")

	temps, ok := data.Hourly["temperature_2m"].([]any)
	if !ok || len(temps) == 0 {
		return fmt.Errorf("invalid temperature data")
	}

	latestTemp := temps[0].(float64)

	msg := fmt.Sprintf(
		"Subject: Weather Update\r\n\r\nCity: Jammu\nElevation: %.2f m\nTemperature: %.2f°C\n",
		data.Elevation,
		latestTemp,
	)
	return smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		e.from,
		e.to,
		[]byte(msg),
	)
}

type WeatherData struct {
	Elevation float64        `json:"elevation"`
	Hourly    map[string]any `json:"hourly"`
}

type WPoller struct {
	closech chan struct{}
	senders []Sender
}

func NewWPoller(senders ...Sender) *WPoller {
	return &WPoller{
		closech: make(chan struct{}),
		senders: senders,
	}
}

func (wp *WPoller) close() {
	// now we can return Empty string to the channel or just close it.
	close(wp.closech)
}

func (wp *WPoller) start() {
	fmt.Println("Startiing wpoller")

	ticker := time.NewTicker(pollInterval) //ticker that sends a signal on a channel (ticker.C) at fixed intervals.
outer:
	for {
		select {
		case <-ticker.C:
			data, err := GetWeatherResults(52.52, 13.41)
			if err != nil {
				log.Fatal(err)
			}
			if err = wp.handleData(data); err != nil {
				log.Fatal(err)
			}
		case <-wp.closech:
			break outer
		}
	}
	fmt.Println("wpoller stop gracefully")
}

func (wp *WPoller) handleData(data *WeatherData) error {
	for _, s := range wp.senders {
		if err := s.Send(data); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
func GetWeatherResults(lat, long float64) (*WeatherData, error) {
	// we shorted the Uri So we did Sprintf to write formatted data into a string.
	uri := fmt.Sprintf("%s?latitude=%.2f&longitude=%.2f&hourly=temperature_2m", endpoint, lat, long)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Now we can actually create a client and use that client to Perforn Request.. We did it because we want to control TimeOut and other things..
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req) //Do() sends an HTTP request and returns an HTTP response, following policy (such as redirects, cookies, auth) as configured on the client.
	if err != nil {
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
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		// log.Fatal(err)
		return nil, err
	}
	// fmt.Println(data["hourly"]["temperature_2m"])
	return &data, nil
}

func main() {
	smsSender := NewSMSSender("9596******")

	emailSender := &EmailSender{
		from:     "vivekkumarr221106@gmail.com",
		password: "hyfjmmnpsvugbznf",
		to: []string{
			"raymorrow003@gmail.com",
			"chaurasiavivek840@gmail.com",
			"amansamraj@gmail.com",
			"Luthergraham2007@gmail.com",
			"Zenaniverseamv@gmail.com",
			"yadavg7970@gmail.com",
			"Shivamrathore1234sk@gmail.com",
		},
	}

	wpoller := NewWPoller(smsSender, emailSender)
	wpoller.start()
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

// json response from weather api struct
type WeatherAPIResponse struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	}
	Forecast struct {
		ForecastDay []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	q := "bahia_blanca"

	if len(os.Args) >= 2 {
		q = os.Args[1]
	}

	res, err := http.Get("https://api.weatherapi.com/v1/forecast.json?key=f38f1a779bbe499bb9b132430240704&q=" + q + "&days=1&aqi=no&alerts=no")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close() //kind of an await on js

	if res.StatusCode != 200 {
		panic("Weather API not available")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var weather WeatherAPIResponse
	//json.Unmarshal takes the body and converts it to the type you give it
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.ForecastDay[0].Hour

	fmt.Printf(
		"%s, %s: %s, %.0fC\n",
		location.Name,
		location.Country,
		current.Condition.Text,
		current.TempC,
	)

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0) //converting the TimeEpoch in an actual date object

		if date.Before(time.Now()) {
			continue
		}

		msg := fmt.Sprintf(
			"%s - %.0fC, Rain chance: %.0f%%, %s\n", // %.0f is for floats
			date.Format("15:04"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)

		if hour.ChanceOfRain < 40 {
			fmt.Print(msg)
		} else {
			color.Red(msg)
		}
	}
}

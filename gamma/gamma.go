package gamma

import (
	"fmt"
	"github.com/kelvins/sunrisesunset"
	"github.com/FedeDP/golight/clightd"
	"github.com/FedeDP/golight/conf"
	"github.com/FedeDP/golight/day"
	"github.com/FedeDP/golight/state"
	"time"
)

var api, _ = clightd.NewGammaApi()

func Subscribe() <-chan time.Time {
	Next()

	now := time.Now()
	if state.NextSunset.Sub(now) > 0 {
		state.DayTime = day.Day
		fmt.Println("Sunset timer elapsing in", state.NextSunset.Sub(now).Truncate(time.Second))
		return time.After(state.NextSunset.Sub(now))
	}
	state.DayTime = day.Night
	fmt.Println("Sunrise timer elapsing in", state.NextSunrise.Sub(now).Truncate(time.Second))
	return time.After(state.NextSunrise.Sub(now))
}

func Next() {
	latitude, err := state.Location.GetLatitude()
	if err != nil {
		panic(err)
	}

	longitude, err := state.Location.GetLongitude()
	if err != nil {
		panic(err)
	}

	t := time.Now()
	_, offset := t.Zone()
	p := sunrisesunset.Parameters {
		Latitude: 	latitude,
		Longitude: 	longitude,
		UtcOffset: 	float64(offset) / 3600,
		Date:      	t,
	}

	state.NextSunrise, state.NextSunset, err = p.GetSunriseSunset()
	if err != nil {
		panic(err)
	}

	/* Library does not use today */
	state.NextSunrise = state.NextSunrise.AddDate(t.Year() - 1, int(t.Month()) - 1, t.Day() - 1)
	state.NextSunset = state.NextSunset.AddDate(t.Year() - 1, int(t.Month()) - 1, t.Day() - 1)
}

func Update() {
	err := api.GoSetTemp(conf.Temps[state.DayTime], &conf.GSmooth, nil)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Set %d gamma temp.\n", conf.Temps[state.DayTime])
	}
}

func Close() {
	err := api.Destroy()
	if err != nil {
		fmt.Println(err)
	}
}
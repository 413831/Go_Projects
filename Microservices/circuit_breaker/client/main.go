package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sony/gobreaker"
)

const (
	requestLimit      = 10
	failureRatioLimit = 0.6
)

var cb *gobreaker.CircuitBreaker

func init() {
	cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:          "HTTP GET",
		ReadyToTrip:   setLimit,
		Timeout:       time.Millisecond,
		OnStateChange: stateChangeHandler,
	})
}

func setLimit(counts gobreaker.Counts) bool {
	// circuit breaker will trip when 60% of requests failed and at least 10 request were made
	failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)

	return counts.Requests >= requestLimit && failureRatio >= failureRatioLimit
}

func stateChangeHandler(name string, from gobreaker.State, to gobreaker.State) {
	if to == gobreaker.StateOpen {
		log.Error().Msg("State Open")
	}
	if from == gobreaker.StateOpen && to == gobreaker.StateHalfOpen {
		log.Info().Msg("Open State to Half Open State")
	}
	if from == gobreaker.StateHalfOpen && to == gobreaker.StateClosed {
		log.Info().Msg("Half Open State to Closed State")
	}
}

func Get(url string) ([]byte, error) {
	body, err := cb.Execute(func() (interface{}, error) {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("http GET request failed")

			return nil, err
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	})

	if err != nil {
		return nil, err
	}

	return body.([]byte), nil
}

func main() {
	var (
		body []byte
		err  error
	)

	urlIncorrect := "http://localhost:8091"
	urlCorrect := "http://localhost:8090"

	for i := 0; i < 20; i++ {
		body, err = Get(urlIncorrect)
		if err != nil {
			log.Error().Err(err).Msg("Error")
		}

		fmt.Println(string(body))
		if i > 15 {
			urlIncorrect = urlCorrect
		}

		time.Sleep(time.Millisecond)
	}
}

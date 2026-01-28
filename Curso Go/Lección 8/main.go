package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func FetchAll(ctx context.Context, ids []string) ([]string, error) {
	// Define channels for results and errors
	var (
		wg         sync.WaitGroup
		results    = make(chan string, len(ids))
		errorsChan = make(chan error, len(ids))
	)

	// Request IDs iteration
	for _, id := range ids {
		wg.Add(1)
		// Go routine declaration to fetch each ID data in parallel
		go func(id string) {
			defer wg.Done()
			data, err := fetchService(id)
			if err != nil {
				// Send error through error channel
				errorsChan <- err
				return
			}
			// Send data through results channel
			results <- data
		}(id)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var out []string
	for {
		select {
		case err := <-errorsChan:
			return nil, err
		case res, ok := <-results:
			if !ok {
				return out, nil
			}
			out = append(out, res)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func fetchService(id string) (string, error) {
	if id == "" {
		return "", errors.New("invalid id")
	}

	// Simulate time processing
	delay := time.Duration(rand.Intn(300)+100) * time.Millisecond
	time.Sleep(delay)

	return "data-for-" + id, nil
}

func main() {
	// To simulate
	rand.New(rand.NewSource(time.Now().UnixNano()))

	ids := []string{"1", "2", "3"}

	fmt.Println("Fetching data...")

	results, err := FetchAll(context.Background(), ids)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Results:", results)
	for _, result := range results {
		fmt.Println("-", result)
	}
}

/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package load

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var (
	duration   string
	workersNum int
	method     string
	url        string
)

// loadCmd represents the load command
var LoadCmd = &cobra.Command{
	Use:   "load",
	Short: "load sends requests for given period of time to given url",
	Long:  `load sends requests for given period of time to given url`,
	Run: func(cmd *cobra.Command, args []string) {
		load()
	},
}

func worker(requests chan *http.Request, results chan string, client *http.Client) {
	var (
		res *http.Response
		err error
	)

	for r := range requests {
		res, err = client.Do(r)
		if err != nil {
			panic(err)
		}
		results <- res.Status
	}
}

func load() {
	requests := make(chan *http.Request, workersNum)
	results := make(chan string)
	defer close(requests)
	defer close(results)

	client := http.Client{}

	// this loop shoud generate requests
	for i := 0; i < cap(requests); i++ {
		go worker(requests, results, &client)
	}

	parsedDuration, err := time.ParseDuration(duration)
	if err != nil {
		panic(err)
	}

	go func() {
		for start := time.Now(); time.Since(start) < parsedDuration; {
			r, err := http.NewRequest(method, url, nil)
			if err != nil {
				fmt.Println(err.Error())
			}

			requests <- r
		}
	}()

	for start := time.Now(); time.Since(start) < parsedDuration; {
		fmt.Println(start)
		res := <-results
		fmt.Println(res)
	}
}

func init() {
	LoadCmd.Flags().StringVarP(&duration, "duration", "d", "", "Please provide load duration with string. e.g. 1s, 3m, 1h")
	if err := LoadCmd.MarkFlagRequired("duration"); err != nil {
		fmt.Println(err)
	}

	LoadCmd.Flags().IntVarP(&workersNum, "workersnum", "w", 1, "Please provide number of workers. e.g. 1, 5 , 100")
	if err := LoadCmd.MarkFlagRequired("workersnum"); err != nil {
		fmt.Println(err)
	}

	LoadCmd.Flags().StringVarP(&method, "method", "m", "", "Please provide http Method for requests. e.g. GET, HEAD, POST")
	if err := LoadCmd.MarkFlagRequired("method"); err != nil {
		fmt.Println(err)
	}

	LoadCmd.Flags().StringVarP(&url, "url", "", "", "Please provide url")
	if err := LoadCmd.MarkFlagRequired("url"); err != nil {
		fmt.Println(err)
	}
}

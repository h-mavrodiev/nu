/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package pscan

import (
	"fmt"
	"net"
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	urlPath string
)

// pscanCmd represents the pscan command
var PscanCmd = &cobra.Command{
	Use:   "pscan",
	Short: "Scans URL for available ports",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		pscan()
	},
}

func worker(ports, results chan int, domain string) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", domain, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func pscan() {
	ports := make(chan int, 100)
	results := make(chan int)

	fmt.Println(viper.GetString("test"))
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results, urlPath)
	}

	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
	}()

	for i := 1; i <= 1024; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}
	close(ports)
	close(results)
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}

func init() {
	PscanCmd.Flags().StringVarP(&urlPath, "url", "u", "", "Please provide url to scan for open ports")

	if err := PscanCmd.MarkFlagRequired("url"); err != nil {
		fmt.Println(err)
	}
}

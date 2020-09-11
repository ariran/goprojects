package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	unixDateValue, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err == nil {
		unixTimeUTC := time.Unix(unixDateValue, 0)
		custom := unixTimeUTC.UTC().Format("2006-01-02-15.4.5")
		fmt.Println(custom)
	} else {
		timeValue, err := time.Parse("2006-1-2-15.4.5", os.Args[1])
		if err == nil {
			fmt.Println(timeValue.Unix())
		} else {
			fmt.Println("Not parseable!")
		}
	}
}

func usage() {
	fmt.Println("Converts unixtime to datetime format and vice versa.")
	fmt.Println("Usage: unixtime {1599657322 | 2020-12-03-13.15.22}")
}

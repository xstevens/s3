package main

import (
	"fmt"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	region = kingpin.Flag("region", "S3 region").Default("us-east-1").String()
)

func main() {
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Version("0.0.1")

	switch kingpin.Parse() {
	case "upload":
		runUpload()
	default:
		fmt.Println("No command specified")
	}
}

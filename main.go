package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("s3", "A general command line client for S3.")
	app.HelpFlag.Short('h')
	app.Version("0.1.0")
	configureCatCommand(app)
	configureUploadCommand(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}

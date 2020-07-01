package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/ysmood/kit"

	"github.com/go-rod/wayang"
)

var (
	headless   = flag.Bool("headless", true, "decide between whether to run chrome in windowed mode or not")
	filePath   = flag.String("file", "", "*the location of the file which will be executed")
	verbose    = flag.Bool("verbose", false, "verbose logging information")
	output     = flag.Bool("output", true, "print JSON output to stdout")
	outputFile = flag.String("outputFile", "", "the file location of the output json")
	timeout    = flag.Int("timeout", 30, "timeout for program")
)

func main() {
	flag.Parse()

	if *filePath == "" {
		log.Fatal("You must provide a file path to a JSON file.")
	}

	url := launcher.New().Headless(*headless).Launch()
	timeout := time.Duration(*timeout)

	browser := rod.New().ControlURL(url).Timeout(timeout * time.Second).Connect()
	defer browser.Close()

	page := browser.Page("")

	runner := &wayang.Runner{
		B:      browser,
		P:      page,
		ENV:    map[string]interface{}{},
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}

	var program wayang.Program
	readRes := kit.ReadJSON(*filePath, &program)
	if readRes != nil {
		log.Fatal("Error while reading to input file:", readRes)
	}

	res, err := runner.RunProgram(program)
	if err != nil {
		log.Print(err.Error())
		if *outputFile != "" {
			writeRes := kit.OutputFile(*outputFile, err, nil)
			if writeRes != nil {
				log.Fatal("Error while writing to output file:", writeRes)
			}
		}
		if *verbose {
			fmt.Println(err.Dump())
		}
		return
	}

	if *output {
		var bin []byte

		switch t := res.(type) {
		case []byte:
			bin = t
		case string:
			bin = []byte(t)
		default:
			var err error
			bin, err = json.MarshalIndent(res, "", "    ")

			if err != nil {
				log.Fatal("Error while converting result to binary representation")
			}
		}
		fmt.Println(string(bin))
	}
	if *outputFile != "" {
		writeRes := kit.OutputFile(*outputFile, res, nil)
		if writeRes != nil {
			log.Fatal("Error while writing to output file:", writeRes)
		}
	}
}

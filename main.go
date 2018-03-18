package main

import (
	"github.com/jessevdk/go-flags"
	overseer "github.com/m-mizutani/overseer/lib"
	"log"
	"os"
)

type Options struct {
	FileName string `short:"r" description:"A pcap file" value-name:"FILE"`
	DevName  string `short:"i" description:"Interface name" value-name:"DEV"`
}

func main() {
	var opts Options

	_, ParseErr := flags.ParseArgs(&opts, os.Args)
	if ParseErr != nil {
		os.Exit(1)
	}

	seer := overseer.New()

	if opts.DevName != "" {
		err := seer.SetPcapDev(opts.DevName)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}

	if opts.FileName != "" {
		err := seer.SetPcapFile(opts.FileName)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}

}

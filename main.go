package main

import (
	"flag"
	"strings"
	"fmt"
	"log"
	"log/syslog"
	"os"
)

var excludesvar string
var modevar string
var forcevar bool
var privatevar bool
var modes map[string]bool
var LogWriter *syslog.Writer
var SyslogError error

func init() {
	modeNames := []string{"ending", "another", "any", "random", "popular"}
	modes = make(map[string]bool)
	for _, modeName := range modeNames {
		modes[modeName] = true
	}
	modeUsage := fmt.Sprintf("Spoof mode - %s", strings.Join(modeNames, ", "))
	flag.Usage = func() {
		usage := fmt.Sprintf(
			"%s -m mode [-f] [-e \"iface,iface,iface\"]\n\nThis tool spoofs all network interfaces at once, logging results and errors to syslog\n\n",
			os.Args[0])
		fmt.Fprintf(os.Stderr, usage)
		flag.PrintDefaults()
	}
	LogWriter, SyslogError = syslog.New(syslog.LOG_INFO, "macouflage-multi")
	if SyslogError != nil {
		log.Fatal(SyslogError)
	}
	flag.StringVar(&excludesvar, "e", "",
		"List of interface names to exclude")
	flag.StringVar(&modevar, "m", "",
		modeUsage)
	flag.BoolVar(&forcevar, "f", false,
		"Force all intefaces to be spoofed each time the program is run, the default is to spoof interfaces that are not already spoofed")
	flag.BoolVar(&privatevar, "p", true,
		"Do not include MAC addresses in the logs")
}

func main() {
	flag.Parse()
	if modevar != "" {
		if modes[modevar] {
			spoofErrors := spoofAll(modevar)
			if len(spoofErrors) > 0 {
				for _, err := range spoofErrors {
					LogWriter.Err(err.Error())
				}
			}
		} else {
			log.Fatalf("Cannot set MACs, Invalid mode: %s", modevar)
		}
	} else {
		flag.Usage()
	}

}

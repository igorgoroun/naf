package ftnlib

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/teris-io/cli"
)

var (
	logger *log.Logger
)

// InitApp - Creates all needed environment variables
func InitApp() (app cli.App, err error) {
	// Init logger
	var logFile, logErr = os.OpenFile("naf.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if logErr != nil {
		err = fmt.Errorf("Can't open log file: %v", logErr)
		return
	}
	logger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
	logger.Printf("Init %s", "app")

	// Only one setup - FIDOPATH environment variable
	DBPath, hasKey := os.LookupEnv("FIDOPATH")
	if !hasKey {
		err = fmt.Errorf("Environment variable FIDOPATH not exists, please define it")
		return
	}

	// Tosser
	toss := cli.NewCommand("toss", "Toss inbound packets").
		WithShortcut("t").
		WithOption(cli.NewOption("dry-run", "do only test run, do not really toss").WithChar('d').WithType(cli.TypeBool)).
		WithAction(func(args []string, options map[string]string) int {
			// do something
			logger.Print(options["verbose"])
			logger.Println("Tosser")
			return 0
		})

	// Nodelist
	nodelistCompile := cli.NewCommand("nodelist-compile", "Nodelist operations").WithShortcut("nlc").
		//WithOption(cli.NewOption("create", "Compile new nodelist index from file").WithChar('c').WithType(cli.TypeBool)).
		WithOption(cli.NewOption("drop", "Clear index before compilation").WithChar('d').WithType(cli.TypeBool)).
		WithArg(cli.NewArg("nodelistFile", "Path to nodelist file").WithType(cli.TypeString)).
		WithAction(func(args []string, options map[string]string) int {
			if _, err := os.Stat(args[0]); os.IsNotExist(err) {
				logger.Fatal(args[0], err)
				return 1
			}
			// Check to drop old index
			if _, ok := options["drop"]; ok {
				drop, _ := strconv.ParseBool(options["drop"])
				if drop {
					logger.Print("Deleting index")
					_, err := DropNodelistIndex(DBPath)
					if err != nil {
						logger.Fatal(err)
					}
				}
			}
			// Create index
			logger.Print("Parsing nodelist file")
			nlc, _ := ParseNodelist(args[0])
			logger.Print("Creating nodelist index")
			res, err := nlc.UpdateIndex(DBPath)
			if !res {
				if err != nil {
					logger.Fatal("Nodelist error", err)
				} else {
					logger.Fatal("Undefined nodelist error")
				}
			}

			// return res
			return 0
		})

	app = cli.New("New Age Fidonet Tosser & Editor, v.0.1.1").
		WithOption(cli.NewOption("verbose", "Verbose execution").WithChar('v').WithType(cli.TypeBool)).
		WithCommand(toss).
		WithCommand(nodelistCompile)

	return
}

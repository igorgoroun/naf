package ftnlib

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/teris-io/cli"
)

var (
	appenvs = []string{"toss", "scan", "pack", "read", "nl"}
)

// App - Main application runner
type App interface {
	Run()
	Check() (result bool)
	Type() (runtype string)
}

// AppEnv - object that contains all needed variables
type AppEnv struct {
	ToRun  string
	Argv   InputArgument
	DBPath string
	Force  bool
}

// InputArgument - arguments structure
type InputArgument struct {
	runEditor  bool
	runTosser  bool
	runScanner bool
	runPacker  bool
	runSender  bool
}

// InitApp - Creates all needed environment variables
func InitApp() (app cli.App, err error) {
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
			fmt.Println("Tosser")
			return 0
		})

	// Nodelist
	nodelistCompile := cli.NewCommand("nodelist-compile", "Nodelist operations").
		WithShortcut("nlc").
		//WithOption(cli.NewOption("create", "Compile new nodelist index from file").WithChar('c').WithType(cli.TypeBool)).
		WithOption(cli.NewOption("drop", "Clear index before compilation").WithChar('d').WithType(cli.TypeBool)).
		WithArg(cli.NewArg("nodelistFile", "Path to nodelist file").WithType(cli.TypeString)).
		WithAction(func(args []string, options map[string]string) int {
			if _, err := os.Stat(args[0]); os.IsNotExist(err) {
				log.Fatal(args[0], err)
				return 1
			}
			// Check to drop old index
			if _, ok := options["drop"]; ok {
				drop, _ := strconv.ParseBool(options["drop"])
				if drop {
					log.Print("Deleting index")
					_, err := DropNodelistIndex(DBPath)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
			// Create index
			log.Print("Parsing nodelist file")
			nlc, _ := ParseNodelist(args[0])
			log.Print("Creating nodelist index")
			res, err := nlc.UpdateIndex(DBPath)
			if !res {
				if err != nil {
					log.Fatal("Nodelist error", err)
				} else {
					log.Fatal("Undefined nodelist error")
				}
			}

			// return res
			return 0
		})

	app = cli.New("New Age Fidonet Tosser & Editor, v.0.1.1").
		WithCommand(toss).
		WithCommand(nodelistCompile)
	return
}

// Run - run app
func (argv *AppEnv) Run() {
	fmt.Println(appenvs)
	fmt.Println(argv)
}

// Check - check arguments
func (argv *InputArgument) Check() (result bool) {
	result = true

	return
}

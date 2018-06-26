package ftnlib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	hostTypes = [...]string{"Hub", "Pvt", "Down"}
)

type NodelistFile interface {
	GetJson(nfl os.File) (jsonNodelist json.Token)
}

type Nodelist struct {
	list []nodelistHost
}

type nodelistHost struct {
	address   FtnAddress
	hostType  string
	hostName  string
	sysopName string
	zone      string
	country   string
	region    string
	city      string
	phone     string
	flags     []nodelistLineFlags
}

type nodelistLineFlags struct {
	flag  string
	value string
}

// currentPlace - For machine cycle
type currentPlace struct {
	zone   int
	region int
	net    int
	node   int
}

// ParseNodelist - the main function
func ParseNodelist(nflPath string) (nl Nodelist, err error) {
	file, err := os.Open(nflPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var current = currentPlace{
		zone:   0,
		region: 0,
		net:    0,
		node:   0,
	}

	for scanner.Scan() {
		line := strings.Trim(strings.TrimSpace(scanner.Text()), "\n")
		// skip commented line
		if line[:1] == ";" || len(line) <= 1 {
			continue
		}

		lineData := strings.Split(line, ",")

		current.node = 0
		switch indicator := strings.ToLower(lineData[0]); indicator {
		case "zone":
			newZone, _ := strconv.Atoi(lineData[1])
			if current.zone != newZone {
				current.zone = newZone
				current.region = 0
				current.net = 0
			}
			continue
		case "region":
			newRegion, _ := strconv.Atoi(lineData[1])
			if current.region != newRegion {
				current.region = newRegion
				current.net = 0
			}
			continue
		case "host":
			newHost, _ := strconv.Atoi(lineData[1])
			if current.net != newHost {
				current.net = newHost
			}
		default:
			nodeNum, _ := strconv.Atoi(lineData[1])
			current.node = nodeNum
		}

		fmt.Println(line)
		lineAddress, _ := FormFtnAddress(current.zone, current.region, current.net, current.node, 0, "")
		fmt.Println(lineAddress)

		var hostline = nodelistHost{
			address: lineAddress,
		}
		fmt.Println(hostline)
		nl.list = append(nl.list, hostline)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

func (nodeList *Nodelist) GetJson() (jsonNodelist json.Token) {

	return
}

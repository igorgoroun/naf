package ftnlib

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	hostTypes        = []string{"Hub", "Pvt", "Down"}
	sqNodelistSearch = `
		SELECT n.id, n.address, n.hostName, n.sysop, n.country, n.city
		FROM nodelist n
		LEFT JOIN nodelist_flags f ON f.node=n.id
		WHERE 1=1			
			%s
		GROUP BY n.id
	`
	// Nodelist line insert
	sqNodelistInsertOne = `
		INSERT INTO nodelist (zone, region, net, node, point, domain, address, hostType, hostName, sysop, country, city, phone)
		VALUES (%d, %d, %d, %d, %d, "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s");
	`
	// Flags
	sqlNodelistFlagInsertOne = `
		INSERT INTO nodelist_flags (node, flag, value)
		VALUES (%d, "%s", "%s");
	`
	sqlNodelistFlagInsertList = `
		INSERT INTO nodelist_flags (node, flag, value)
		VALUES %s;
	`
	// get record by zone/net/node/point
	sqNodelistGetByAddress = `
		SELECT id FROM nodelist
		WHERE 1=1
			AND zone=%d
			AND net=%d
			AND node=%d
			AND point=%d
		LIMIT 1;
	`
	// drop nodelist
	sqNodelistDropIndex = []string{`DROP TABLE nodelist;`, `DROP TABLE nodelist_flags;`}
	// check tables exists
	sqNodelistCheckTables = `
		SELECT name FROM sqlite_master WHERE type='table' AND name IN ('nodelist', 'nodelist_flags');
	`
	// create modelist table
	sqNodelistCreateTables = []string{`
		CREATE TABLE IF NOT EXISTS nodelist (
			id INTEGER NOT NULL PRIMARY KEY,
			zone INTEGER NOT NULL,
			region INTEGER NOT NULL,
			net INTEGER NOT NULL,
			node INTEGER NOT NULL,
			point INTEGER DEFAULT 0,
			domain TEXT,
			address TEXT NOT NULL,
			hostType TEXT,
			hostName TEXT,
			sysop TEXT,
			country TEXT,
			city TEXT,
			phone TEXT,
			UNIQUE(zone, net, node, point) ON CONFLICT REPLACE
		);
		`,
		`DELETE FROM nodelist;`,
		`
		CREATE TABLE IF NOT EXISTS nodelist_flags (
			id INTEGER NOT NULL PRIMARY KEY,
			node INTEGER NOT NULL,
			flag TEXT,
			value TEXT
		);`,
		`DELETE FROM nodelist_flags;`}
)

// Nodelist - a list of parsed nodelist lines
type Nodelist struct {
	List []nodelistHost
}

type nodelistHost struct {
	Address   FtnAddress
	HostType  string
	HostName  string
	SysopName string
	Country   string
	City      string
	Phone     string
	Flags     []nodelistLineFlags
}

type nodelistLineFlags struct {
	Flag  string
	Value string
}

// currentPlace - For machine cycle
type currentPlace struct {
	zone    int
	region  int
	country string
	net     int
	node    int
	append  bool
}

// ParseNodelist - the main function
func ParseNodelist(nflPath string) (nl Nodelist, err error) {
	file, err := os.Open(nflPath)
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var current = currentPlace{
		zone:   0,
		region: 0,
		net:    0,
		node:   0,
		append: false,
	}

	for scanner.Scan() {
		line := strings.Trim(strings.TrimSpace(scanner.Text()), "\n")
		// skip commented line
		if line[:1] == ";" || len(line) <= 1 {
			continue
		}

		lineData := strings.Split(line, ",")
		current.node = 0
		current.append = false

		switch indicator := strings.ToLower(lineData[0]); indicator {
		case "zone":
			newZone, _ := strconv.Atoi(lineData[1])
			if current.zone != newZone {
				current.zone = newZone
				current.region = 0
				current.country = lineData[2]
				current.net = 0
			}
			continue
		case "region":
			newRegion, _ := strconv.Atoi(lineData[1])
			if current.region != newRegion {
				current.region = newRegion
				current.country = lineData[2]
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
			current.append = true
		}

		// Append only real addresses to nodelist index
		if current.append {
			// Address object
			lineAddress, _ := FormFtnAddress(current.zone, current.region, current.net, current.node, 0, "")
			// Line data
			var hostline = nodelistHost{
				Address:   lineAddress,
				HostName:  strings.Replace(lineData[2], "_", " ", -1),
				Country:   strings.Replace(current.country, "_", " ", -1),
				City:      strings.Replace(lineData[3], "_", " ", -1),
				SysopName: strings.Replace(lineData[4], "_", " ", -1),
				Phone:     lineData[5],
			}
			// Host type
			if ContainsString(lineData[0], hostTypes) {
				hostline.HostType = lineData[0]
			} else {
				hostline.HostType = "Node"
			}
			// Node flags
			for _, fl := range lineData[6:] {
				flagData := strings.Split(fl, ":")
				var flag = nodelistLineFlags{
					Flag:  flagData[0],
					Value: "",
				}
				if len(flagData) == 2 {
					flag.Value = flagData[1]
				}
				hostline.Flags = append(hostline.Flags, flag)
			}
			// Add to index list
			nl.List = append(nl.List, hostline)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

// DropNodelistIndex - Delete index from db
func DropNodelistIndex(DBPath string) (res bool, err error) {
	db, err := GetStorage(DBPath, "nodelist")
	if err != nil {
		return false, fmt.Errorf("Cannot get nodelist storage")
	}
	for _, q := range sqNodelistDropIndex {
		_, err := db.Exec(q)
		if err != nil {
			return false, fmt.Errorf("Cannot drop nodelist tables")
		}
	}
	return
}

// UpdateIndex - Check nodelist and update nodelist address book
// This operaton can take a time to run
func (nodeList *Nodelist) UpdateIndex(DBPath string) (res bool, err error) {
	store, err := GetStorage(DBPath, "nodelist")
	for _, line := range nodeList.List {
		sqlQuery := fmt.Sprintf(sqNodelistGetByAddress, line.Address.zone, line.Address.net, line.Address.node, line.Address.point)
		rows, err := store.Query(sqlQuery)
		if err != nil {
			return false, fmt.Errorf("Error querying nodelist: %s", err)
		}

		rowsCnt := 0
		for rows.Next() {
			var id int
			err = rows.Scan(&id)
			rowsCnt++
			if err != nil {
				return false, err
			}
			// TODO: update row if force specified

		}
		if err = rows.Err(); err != nil {
			return false, err
		}
		if rowsCnt == 0 {
			// insert new row
			sqlQ := fmt.Sprintf(sqNodelistInsertOne,
				line.Address.zone,
				line.Address.region,
				line.Address.net,
				line.Address.node,
				line.Address.point,
				line.Address.domain,
				line.Address.Dump3D(),
				line.HostType,
				line.HostName,
				line.SysopName,
				line.Country,
				line.City,
				line.Phone,
			)
			sqlRes, err := store.Exec(sqlQ)
			if err != nil {
				return false, fmt.Errorf("Can't create nodelist record: %s", err)
			}
			id, err := sqlRes.LastInsertId()
			if err != nil {
				return false, fmt.Errorf("Can't get inserted id: %s", err)
			}
			flagsList := []string{}
			for _, flag := range line.Flags {
				strFlag := strings.Join([]string{"(", fmt.Sprintf("%d", id), ",\"", flag.Flag, "\",NULL)"}, "")
				if flag.Value != "" {
					strFlag = strings.Join([]string{"(", fmt.Sprintf("%d", id), ",\"", flag.Flag, "\",\"", flag.Value, "\")"}, "")
				}
				flagsList = append(flagsList, strFlag)
			}
			// create flags query
			sqlF := fmt.Sprintf(sqlNodelistFlagInsertList, strings.Join(flagsList, ", "))
			sqlRes, err = store.Exec(sqlF)
			if err != nil {
				log.Output(0, "Nodelist flag wasn't inserted")
			}
		}

	}
	store.Close()
	return true, nil
}

// SearchNodelist - Search for records in nodelist
func SearchNodelist(DBPath string) (err error) {
	return
}

// NodeInfo - Get full node info by 3D-address
func NodeInfo(DBPath string) (err error) {
	return
}

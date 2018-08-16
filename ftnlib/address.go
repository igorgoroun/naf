package ftnlib

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	defaultDomain = "fidonet"
)

// Address - address operations interface
type Address interface {
	Dump2D() (addressDump string)
	Dump3D() (addressDump string)
	Dump5D() (addressDump string)
}

// FtnAddress - address object structure
type FtnAddress struct {
	zone   int
	region int
	net    int
	node   int
	point  int
	domain string
	parsed bool
}

// ParseFtnAddress - parse ftn-address string to object FtnAddress
func ParseFtnAddress(addr string) (address FtnAddress, err error) {
	pattern := `(\d{1,})(?:\:)(\d{1,})(?:\/)(\d{1,})(?:(?:\.)?(\d{1,})|)(?:(?:\@)?([a-zA-Z\._]+)|)`
	reg, err := regexp.Compile(pattern)
	match := reg.FindStringSubmatch(addr)
	address.parsed = true
	if len(match) < 5 {
		address.parsed = false
		return
	}
	parsed := make(map[int]int)
	for i := 1; i <= 4; i++ {
		parsed[i], err = strconv.Atoi(match[i])
		if err != nil && i != 4 {
			address.parsed = false
			return
		}
	}
	address.zone, address.net, address.node, address.point, address.domain = parsed[1], parsed[2], parsed[3], parsed[4], match[5]

	if address.domain == "" {
		address.domain = defaultDomain
	}

	// TODO: incorrect region setter, need to be based on nodelist index
	regSlice := 2
	if len(match[2]) < regSlice {
		regSlice = len(match[2])
	}
	address.region, err = strconv.Atoi(match[2][:regSlice])
	return
}

// FormFtnAddress - ints to address
func FormFtnAddress(zone int, region int, net int, node int, point int, domain string) (address FtnAddress, err error) {
	address.zone, address.region, address.net, address.node, address.point, address.domain = zone, region, net, node, point, domain
	if address.domain == "" {
		address.domain = defaultDomain
	}
	address.parsed = true
	return
}

// Dump2D - returns z:n, string
func (address *FtnAddress) Dump2D() (addressDump string) {
	/*
		if address.parsed && address.zone > 0 && address.net > 0 {
			addressDump = fmt.Sprintf("%d:%d", address.zone, address.net)
		} else if address.parsed && address.Special() && address.zone > 0 {
			addressDump = fmt.Sprintf("%d:%d", address.zone, address.region)
		}
	*/
	addressDump = fmt.Sprintf("%d:%d", address.zone, address.net)
	return
}

// Dump3D - returns z:n/f(.p), string
func (address *FtnAddress) Dump3D() (addressDump string) {
	//if address.parsed && address.zone > 0 && (address.net > 0 || address.region > 0) && address.node >= 0 {
	addressDump = fmt.Sprintf("%s/%d", address.Dump2D(), address.node)
	if address.point > 0 {
		addressDump = fmt.Sprintf("%s.%d", addressDump, address.point)
	}
	//}
	return
}

// Dump5D - returns z:n/f.p@domain, string
func (address *FtnAddress) Dump5D() (addressDump string) {
	if address.parsed && address.zone > 0 && address.net > 0 && address.node >= 0 && len(address.domain) > 0 {
		addressDump = fmt.Sprintf("%s@%s", address.Dump3D(), address.domain)
	}
	return
}

// GetType - returns type of parsed address (string) and special flag (bool)
func (address *FtnAddress) GetType() (addressType string, special bool) {
	addressType = "invalid"
	special = false
	if !address.parsed {
		return
	}
	if address.zone > 0 && address.net > 0 && address.node > 0 && address.point > 0 {
		// point
		addressType = "point"
	} else if address.zone > 0 && address.net > 0 && address.node >= 0 {
		// node
		addressType = "node"
		// region coord
		if address.zone == address.net && address.node == 0 {
			addressType = "zoneCoordinator"
			special = true
		} else if address.zone == address.net && address.node > 0 {
			addressType = "zoneGate"
			special = true
		} else if address.region == address.net && address.node == 0 {
			addressType = "regionCoordinator"
			special = true
		} else if address.zone > 0 && address.net > 0 && address.node == 0 {
			addressType = "networkHost"
			special = true
		}
	} else if address.zone > 0 && address.region > 0 && address.net == 0 && address.node > 0 {
		addressType = "zoneGate"
		special = true
	}
	return
}

// Type - returns address type, string
func (address *FtnAddress) Type() (res string) {
	res, _ = address.GetType()
	return
}

// Special - returns special flag, bool
func (address *FtnAddress) Special() (res bool) {
	_, res = address.GetType()
	return
}

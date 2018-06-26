package ftnlib

type ftnPacketBinary struct {
	fromZone   uint32
	fromNet    uint32
	fromNode   uint32
	toZone     uint32
	toNet      uint32
	toNode     uint32
	dateYear   uint32
	dateMonth  uint32
	dateDay    uint32
	dateHour   uint32
	dateMinute uint32
	dateSecond uint32
}

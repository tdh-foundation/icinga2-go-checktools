package icinga2_go_checktools

import (
	"fmt"
)

//noinspection ALL
const (
	// icinga2 Constant

	// Ok SwStatus
	OkMsg  = "OK"
	OkExit = 0

	// Warning SwStatus
	WarMsg  = "WARNING"
	WarExit = 1

	// Critial SwStatus
	CriMsg  = "CRITICAL"
	CriExit = 2

	// Unknown SwStatus
	UnkMsg  = "UNKNOWN"
	UnkExit = 3
)

type Icinga struct {
	Message string
	Exit    int
	Metric  string
}

//Start Definition for switch monitoring
type SwitchInterfaceStatus struct {
	Port   string
	Name   string
	Status string
	Vlan   string
	Duplex string
	Speed  string
	Type   string
}

type SwitchStatus struct {
	Name     string
	SwStatus []SwitchInterfaceStatus
}

type SwitchInterface interface {
	ParseInterfaceStatus(response string) error
	ReturnIcingaResult() Icinga
	CheckInterfaceStatus(host string, username string, password string, identity string, port int) (Icinga, error)
	Status() []SwitchInterfaceStatus
}

//End Definition for switch monitoring

// Stringer to display/print Icinga struct
func (ict Icinga) String() string {
	outString := ""

	switch ict.Exit {
	case OkExit:
		outString = OkMsg + ":"
	case WarExit:
		outString = WarMsg + ":"
	case CriExit:
		outString = CriMsg + ":"
	default:
		outString = UnkMsg + ":"
	}

	if ict.Message != "" {
		outString += " " + ict.Message
	}

	if ict.Metric != "" {
		outString += " |" + ict.Metric
	}

	return fmt.Sprintf("%s", outString)
}

// String to display/print Switch
func (sws SwitchInterfaceStatus) String() string {
	return fmt.Sprintf("Port: %s\tName: %s\tStatus: %s\tVlan: %s\tDuplex: %s\tSpeed: %s\tType: %s\t",
		sws.Port, sws.Name, sws.Status, sws.Vlan, sws.Duplex, sws.Speed, sws.Type)
}

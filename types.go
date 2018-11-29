package main

import "time"

type connector struct {
	isConnected	bool
	host		string
	port		string
	timeout		time.Duration
}

type Connection interface {
	connect()
}

//general json type of data for mdm checkin request
type Request struct {
     Request     requestType  //timestamp and type of request where send to server
     Hardware    requestHardware // hardware type included serial number device model
     Network     requestWIFI //only wifi,  included name
     User        Account // account type: short and full name
     OS          OSInfo // information aboout device OS, version, build, name
     Versions    Version // version of application and api, static
     UID		string
}

type requestWIFI struct {
    InterfaceName   string
}

type requestHardware struct {
    Model   string
    CPU		string
    Memory	string
    Disk	string
}


// part of checkin request
type Account struct {
    Account  string
    FullName string
}

//part of checking request
type OSInfo struct {
    OSName   string
    OVersion string
    OSBuild  string
}

type requestType struct { 
    Date    string
    Type    string
    Serial	string
}

//part of request
type Version struct {
    Application int 
    API         int 
}

type Authorize struct {
	Request		requestAuth
    Uri string
    isAuthorize bool
}

type requestAuth struct {
    Date    string
    Type    string
    Serial  string
}

type AuthAnswer struct {
        Serial  string  `json:"Serial"`
        Time    string  `json:"Status"`
        UID     string  `json:"UID"`
        ReceivedStatus	bool	`json:"ReceivedStatus"`
}



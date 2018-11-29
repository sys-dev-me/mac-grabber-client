package main

import "log"

import "os"
import "os/exec"
import "encoding/json"
import "errors"

import "net"

import "time"
import "strings"
import "bytes"
import "strconv"

//gocron
import "github.com/jasonlvhit/gocron"

func (a connector) connect ( host string, port string) (bool, connector, error) {

	a.host = host
	a.port = port
	a.timeout = time.Second
	
	var b bytes.Buffer
	b.WriteString ( a.host )
	b.WriteString ( ":" )
	b.WriteString ( a.port )

	conn, err := net.Dial("tcp", b.String())
	
	if err != nil {
	
		return false, a, errors.New( "Something went wrong" )
	
	}

	if conn != nil {
		a.isConnected = true
		return true, a, nil
	}


	return false, a, errors.New("Host is down or port is closed")
}

func Path( uri string, auth *Authorize){
	idx := len(uri)
	if uri[len(uri)-1] != '/' {
		uri = uri[:idx] + "/"
	}
	if auth.isAuthorize {
		auth.Uri = uri + "serve/"
	} 

	if !auth.isAuthorize {
		auth.Uri = uri + "serve/"
	}
}

func start() Authorize {
	t := time.Now()
	a := Authorize{ requestAuth { t.String(), "auth", getSerialNumber() }, "", false}
	return a
}

func collect() (Request){
	md, _ := exec.Command("bash", "-c", "/usr/sbin/sysctl hw.model | awk -F ':' '{print $2}' | tr -d ' \n'").Output()
	accName, _ := exec.Command("bash", "-c", "id -F").Output()
	wifi, _ := exec.Command("bash", "-c", "networksetup -listnetworkserviceorder | grep -oE '^(.*)Wi-Fi,(.*)'").Output()
	os, _ := exec.Command("bash", "-c", "sw_vers").Output()
	s := strings.Split(string(os), "\n")
	m := make([]string, 3, 3)
	for idx, _ := range m {
		ss := strings.Split(string(s[idx]), ":\t")
		m[idx] = ss[1]
	}
	tmp1 := strings.Trim(string(wifi), "()\n")
	tmp2 := strings.Split(string(tmp1), ": ")
	t := time.Now()
	a := Request{
		requestType{t.String(), "checkin", getSerialNumber()}, 
		requestHardware{string(md), getCPU(), getMemory(), getDisk()}, 
		requestWIFI{string(tmp2[2])},
		Account{getAccount(), string(accName)}, 
		OSInfo{string(m[0]), string(m[1]), string(m[2])}, 
		Version{1, 1},
		getKey(),
		}
	return a
}

func send( conn connector, uri string ){
	a := collect()
	b, _ := json.Marshal(a)

	sendData( conn.host, conn.port, string(b), false)

}

func sendAuthorize (a connector, uri string ) bool {
	t := time.Now()

	request := Authorize{ requestAuth{t.String(), "auth", getSerialNumber()}, "", false }
	b,_ := json.Marshal (request) 

	sendData ( a.host, a.port, string(b), true )

	var ac AuthAnswer
	ac.ReceivedStatus = true;
	log.Println( "Response status is: ", strconv.FormatBool (ac.ReceivedStatus) )
	
	return ac.ReceivedStatus
}

func ask(conn connector,  a *Authorize ) *Authorize{
	log.Println("Current authorization status: ", a.isAuthorize)
	if a.isAuthorize {
		log.Println( "Device authorized, send details." )
        send(conn, a.Uri)
    }

    if !a.isAuthorize {
    	log.Println("Device not authorized, send auth request.")
        a.isAuthorize = sendAuthorize(conn, a.Uri)
    }
	return a
}

var conn = connector{}

func main() {
	if len(os.Args) < 2 {
	    log.Println ( "Not enough params, usage: grabber [host] [port]" )
	    os.Exit(1)
	}
	isConnected, conn, err := conn.connect( os.Args[1], os.Args[2] )

	if err != nil {
		log.Println ( "Connection status:", isConnected )
		os.Exit(0);
	}

	a := start()
	log.Println( "Initiate application." )

	if !a.isAuthorize {
    	    log.Println("Device not authorized, send auth request.")
    	    log.Println( "Will use: ", conn )
    	    a.isAuthorize = sendAuthorize(conn, a.Uri)
	}
	
	    log.Println("Will used url: ", a.Uri)	

	
	gocron.Every(1).Minute().Do( ask, conn, &a )

	<- gocron.Start()
}


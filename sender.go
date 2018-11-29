package main

import "log"
import "net"
import "strings"
import "encoding/json"
import "github.com/keybase/go-keychain"

const (
	StopCharacter = "\r\n\r\n"
)

var key string = ""

func sendData(ip string, port string, data string, isAuth bool) {
	log.Println("Used host: ", ip)
	log.Println("Used port: ", port)
	addr := strings.Join([]string{ip, port}, ":")
	log.Println("Used: ", addr)
	conn, err := net.Dial("tcp", addr)
	if err !=nil {
		log.Println( "Сервер не отвечает" )
		}

	defer conn.Close()

	if err != nil {
		log.Fatalln(err)
	}

	conn.Write([]byte(data))
	conn.Write([]byte(StopCharacter))

	log.Println("Send: ", data)

	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	var answer AuthAnswer
	//json.NewDecoder(strings.NewReader(string(buff))).Decode(&answer)

	if isAuth {
    	json.NewDecoder(strings.NewReader(string(buff))).Decode(&answer)
		log.Println ( "Received: ", answer )
    	item := keychain.NewItem()
    	item.SetSecClass( keychain.SecClassGenericPassword )
    	item.SetService ( "mdm_service" )
    	item.SetAccount ( getAccount() )
    	item.SetLabel( "grabber" )
    	item.SetAccessGroup ( "mdm.group" )
    	item.SetData ( []byte(answer.UID))
    	key = answer.UID
    	item.SetSynchronizable(keychain.SynchronizableNo)
		item.SetAccessible(keychain.AccessibleWhenUnlocked)
		err := keychain.AddItem(item)
		if err != nil {
			log.Println ( "Something went wrong: ", err )
		}
		if err == keychain.ErrorDuplicateItem {
			log.Println ( "Duplicate found" )
		}
   	}

	log.Println ( "Answer size: ", n )
	log.Println ( "Receive: ", answer )

}


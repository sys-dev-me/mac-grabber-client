package main

import "log"
import "github.com/keybase/go-keychain"


func getKey () string {
query := keychain.NewItem()
query.SetSecClass(keychain.SecClassGenericPassword)
query.SetService("mdm_service")
query.SetAccount(getAccount())
query.SetAccessGroup("mdm.group")
query.SetMatchLimit ( keychain.MatchLimitOne )

// single attribue
query.SetReturnData(true)


results, err := keychain.QueryItem(query)
if err != nil {
  log.Println ( "Unable to retrive key" )
}
return string( results[0].Data )	
}


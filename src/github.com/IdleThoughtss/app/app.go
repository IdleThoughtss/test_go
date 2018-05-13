package app

import (
	"net/http"
	"sync"
	"github.com/IdleThoughtss/dataStruct"
)

type App struct {
	 baseRequest  map[string]string
	 passTicket  string
	 syncKey
	 httpClient http.Client
	 fileLocker  sync.Mutex{}
	 contactList = map[string]Contact
	 user User
}



func (app *App) init(){

}
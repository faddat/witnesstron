package witnesstron

import (
	"fmt"
	"os"
	"time"
	"net/http"
    "encoding/json"
    "github.com/faddat/steemconnect"
)


func main(){
    

        address := "https://steem.yt"
    client := Steemconnect(address)
    
    for {
        time.Sleep()
    }


    
}


func cryptonator(){
        type cryptonator struct {
        price `json:"price"`
    }
        jsonreply, err := http.Get("https://api.cryptonator.com/api/ticker/steem-usd")
if err != nil {
	// handle err
}
defer resp.Body.Close()

}

func coinmarketcap(){
    type coinmarketcap struct {
        price `json:"price_usd"`
    }
    
resp, err := http.Get("https://api.coinmarketcap.com/v1/ticker/steem")
if err != nil {

    log.fatal
}
defer resp.Body.Close()

    
}

func cryptocompare(){
    type cryptocompare struct {
        price `json:"usd"`
    }
    var cooter cryptocompare
    resp, err := http.Get("https://api.coinmarketcap.com/v1/ticker/steem")
if err != nil {
    json.unmarshal
	// handle err
}
defer resp.Body.Close()
}
package witnesstron

import (
	"encoding/json"
	"fmt"
	"github.com/bluele/gforms"
	"github.com/faddat/steemconnect"
	"log"
	"net/http"
	"time"
)

type Witnessupdate struct {
	witnessname                string `gforms:"witnessname"`
	steem                      int
	wifactive                  string `gforms:"wifactive"`
	steem_account_creation_fee int    `gforms:"account_creation_fee"`
	steem_maximum_block_size   int    `gforms:"steem_maximum_block_size"`
	steem_sbd_interest_rate    int    `gforms:"steem_sbd_interest_rate"`
	steem_witness_url          string `gforms:"steem_witness_url"`
}

func main() {
	type priceProvider interface {
		price() (int, error)
	}
	go gatherdata()
	address := "https://steem.yt"
	client := steemconnect.Steemconnect(address)
	witnessschedule := client.Database.GetWitnessScheduleRaw()
	for {
		one := make(chan int)
		two := make(chan int)
		three := make(chan int)
		go cryptonator(one)
		go coinmarketcap(two)
		go cryptocompare(three)
		on := <-one
		tw := <-two
		th := <-three
		steem := (on + tw + th) / 3
		fmt.Print("Average of cryptonator, coinmarketcap, and cryptocompare is:  ", steem, "\n")
		time.Sleep(90 * time.Minute)
		go witness()
	}
}

func cryptonator(one chan int) {
	var k struct {
		price int `json:"price"`
	}
	resp, err := http.Get("https://api.cryptonator.com/api/ticker/steem-usd")
	if err != nil {
		log.Fatal(err)

	}
	if err := json.NewDecoder(resp.Body).Decode(&k); err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Print("CRPYTONATOR API: ", k.price)
	one <- k.price
}

func coinmarketcap(two chan int) {
	var k struct {
		price int `json:"price_usd"`
	}
	resp, err := http.Get("https://api.coinmarketcap.com/v1/ticker/steem")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(resp.Body).Decode(&k); err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Print("COINMARKETCAP API: ", k.price)
	two <- k.price

}

func cryptocompare(three chan int) {
	var k struct {
		price int `json:"usd"`
	}
	resp, err := http.Get("https://min-api.cryptocompare.com/data/price?fsym=STEEM&tsyms=USD")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(resp.Body).Decode(&k); err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Print("CRYPTOCOMPARE API", k.price)
	three <- k.price
}

func gatherdata() {
	witnessupdateform := gforms.DefineModelForm(witnessupdate{}, gforms.NewFields())
}

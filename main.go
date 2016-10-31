package witnesstron

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bluele/gforms"
	"github.com/faddat/steemconnect"
)

//Witnessupdate is used to build the form where update info is entered.
type Witnessupdate struct {
	witnessname             string `gforms:"witnessname"`
	steem                   int
	wifactive               string `gforms:"wifactive"`
	steemAccountCreationFee int    `gforms:"accountCreationFee"`
	steemMaximumBlockSize   int    `gforms:"steemMaximumBlockSize"`
	steemSBDInterestRate    int    `gforms:"steemSBDInterestRate"`
	steemWitnessURLss       string `gforms:"steemWitnessURL"`
}

func main() {
	type priceProvider interface {
		price() (int, error)
	}
	go gatherdata()
	address := "https://steem.yt"
	client := steemconnect.Steemconnect(address)
	witnessschedule, err := client.Database.GetWitnessScheduleRaw()
	if err != nil {
		log.Fatal(err)
	}
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
		fmt.Print("the witness schedule is", witnessschedule)
		time.Sleep(90 * time.Minute)
		updatewitness(steem)
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

	tplText := `<form method="post">
{{range $i, $field := .Fields}}
  <label>{{$field.GetName}}: </label>{{$field.Html}}
  {{range $ei, $err := $field.Errors}}<label class="error">{{$err}}</label>{{end}}<br />
{{end}}<input type="submit">
</form>
  `
	tpl := template.Must(template.New("tpl").Parse(tplText))

	witnessUpdateForm := gforms.DefineModelForm(Witnessupdate{}, gforms.NewFields())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		form := witnessUpdateForm(r)
		if r.Method != "POST" {
			tpl.Execute(w, form)
			return
		}
		if !form.IsValid() {
			tpl.Execute(w, form)
			return
		}
		user := Witnessupdate{}
		form.MapTo(&user)
		fmt.Fprintf(w, "ok: %v", user)
	})
	http.ListenAndServe(":9000", nil)
}

func updatewitness(steem int) {
	body := strings.NewReader(`{"jsonrpc": "2.0", "method": "call", "params": ["witness_api", "witness_update", owner=officialfuzzy, url=https://steemit.com/witness-category/@anyx/witness-application-anyx, block_signing_key=STM5ha3wiAZX1PL1RkBH8tsm4vuMw1mGbCVLcaFjDe21FbBxMrj48, props={ account_creation_fee=26.900 STEEM, maximum_block_size=65536, sbd_interest_rate=1000}, fee=0.000 STEEM}`)
	req, err := http.NewRequest("POST", "http://127.0.0.1:8090", body)
	if err != nil {
		// handle err

	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err

	}
	defer resp.Body.Close()
}

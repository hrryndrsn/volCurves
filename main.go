package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetInstramentsResponse struct {
	UsOut   int64    `json:"usOut"`
	UsIn    int64    `json:"usIn"`
	UsDiff  int      `json:"usDiff"`
	Testnet bool     `json:"testnet"`
	Success bool     `json:"success"`
	Result  []Result `json:"result"`
	Message string   `json:"message"`
}

type Result struct {
	Kind           string  `json:"kind"`
	BaseCurrency   string  `json:"baseCurrency"`
	Currency       string  `json:"currency"`
	MinTradeSize   float64 `json:"minTradeSize"`
	InstrumentName string  `json:"instrumentName"`
	IsActive       bool    `json:"isActive"`
	Settlement     string  `json:"settlement"`
	Created        string  `json:"created"`
	TickSize       float64 `json:"tickSize"`
	PricePrecision int     `json:"pricePrecision"`
	Expiration     string  `json:"expiration"`
	Strike         float64 `json:"strike,omitempty"`
	OptionType     string  `json:"optionType,omitempty"`
	ContractSize   float64 `json:"contractSize,omitempty"`
}

type OrderBookResponse struct {
	UsOut   int64     `json:"usOut"`
	UsIn    int64     `json:"usIn"`
	UsDiff  int       `json:"usDiff"`
	Testnet bool      `json:"testnet"`
	Success bool      `json:"success"`
	Result  OrderBook `json:"result"`
	Message string    `json:"message"`
}

type OrderBook struct {
	State           string  `json:"state"`
	SettlementPrice float64 `json:"settlementPrice"`
	Instrument      string  `json:"instrument"`
	Bids            []Bid   `json:"bids"`
	Asks            []Ask   `json:"asks"`
	Tstamp          int64   `json:"tstamp"`
	Last            float64 `json:"last"`
	Low             float64 `json:"low"`
	High            float64 `json:"high"`
	Mark            float64 `json:"mark"`
	UPx             float64 `json:"uPx"`
	UIx             string  `json:"uIx"`
	IR              float64 `json:"iR"`
	MarkIv          float64 `json:"markIv"`
	AskIv           float64 `json:"askIv"`
	BidIv           float64 `json:"bidIv"`
	Delta           float64 `json:"delta"`
	Gamma           float64 `json:"gamma"`
	Vega            float64 `json:"vega"`
	Theta           float64 `json:"theta"`
}

type Bid struct {
	Quantity float64 `json:"quantity"`
	Amount   float64 `json:"amount"`
	Price    float64 `json:"price"`
	Cm       float64 `json:"cm"`
	CmAmount float64 `json:"cm_amount"`
}

type Ask struct {
	Quantity float64 `json:"quantity"`
	Amount   float64 `json:"amount"`
	Price    float64 `json:"price"`
	Cm       float64 `json:"cm"`
	CmAmount float64 `json:"cm_amount"`
}

func main() {

	//routing
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	// Route
	router.GET("/instruments", getInstruments)
	router.GET("/orderbooks", getOrderBooks)

	router.Run(":8080")
}

func getInstruments(c *gin.Context) {
	var data GetInstramentsResponse

	//http get request
	res, err := http.Get("https://www.deribit.com/api/v1/public/getinstruments")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not find entry"})
		log.Fatal(err)
		return
	}
	defer res.Body.Close()

	//read body
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not find entry"})
		log.Fatal(err)
		return
	}

	//unmarshall into the reponse struct
	err = json.Unmarshal([]byte(contents), &data)
	if err != nil {
		log.Fatalln(err)
	}

	//send the unmarshalled data back to the caller
	c.JSON(http.StatusOK, gin.H{"instruments": data})
}

func getOrderBooks(c *gin.Context) {
	var responseData []OrderBook
	// instruments := []string{"BTC-28DEC18-5250-C", "BTC-30NOV18-4000-P", "BTC-28JUN19-20000-C"}

	// for i, v := range instruments {

	// }
	//get an option orderbook
	data1, err := handleGetOrderBook("BTC-28JUN19-20000-C")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed getting orderbook"})
		log.Fatal(err)
		return
	}

	//get an option orderbook
	data2, err := handleGetOrderBook("BTC-30NOV18-4000-P")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed getting orderbook"})
		log.Fatal(err)
		return
	}

	//get an option orderbook
	data3, err := handleGetOrderBook("BTC-28JUN19-20000-C")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed getting orderbook"})
		log.Fatal(err)
		return
	}

	responseData = append(responseData, data1.Result)
	responseData = append(responseData, data2.Result)
	responseData = append(responseData, data3.Result)
	c.JSON(http.StatusOK, gin.H{"orderbooks": responseData})
}

//handles get request, returns the orderbook response
func handleGetOrderBook(name string) (OrderBookResponse, error) {
	var data OrderBookResponse

	//http get request
	res, err := http.Get("https://www.deribit.com/api/v1/public/getorderbook?instrument=" + "BTC-28DEC18-5250-C")
	if err != nil {
		log.Fatal(err)
		return data, errors.New("failed to get Url")
	}
	defer res.Body.Close()

	//read body
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return data, errors.New("failed reading body")
	}
	// fmt.Println(string(contents))
	//unmarshall into the reponse struct
	err = json.Unmarshal([]byte(contents), &data)
	if err != nil {
		log.Fatalln(err)
	}
	return data, nil
}

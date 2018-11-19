package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

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
	// Low             float64 `json:"low,omitempty"`
	// High            float64 `json:"high,omitempty"`
	Mark   float64 `json:"mark"`
	UPx    float64 `json:"uPx"`
	UIx    string  `json:"uIx"`
	IR     float64 `json:"iR"`
	MarkIv float64 `json:"markIv"`
	AskIv  float64 `json:"askIv"`
	BidIv  float64 `json:"bidIv"`
	Delta  float64 `json:"delta"`
	Gamma  float64 `json:"gamma"`
	Vega   float64 `json:"vega"`
	Theta  float64 `json:"theta"`
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
	// router.GET("/concurrent", concurrentGetOB)

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

	//datais the main object
	// we also want data.result array of Results
	results := data.Result
	var names []string

	for _, v := range results {
		names = append(names, v.InstrumentName)
	}

	books, err := concurrentGetOB(names)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not find entry"})
		log.Fatal(err)
		return
	}

	//send the unmarshalled data back to the caller
	c.JSON(http.StatusOK, gin.H{"instruments": data.Result, "names": names, "count": len(names), "orderbooks": books})
}

func getOrderBooks(c *gin.Context) {
	var responseData []OrderBook

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

	//unmarshall into the reponse struct
	err = json.Unmarshal([]byte(contents), &data)
	if err != nil {
		log.Fatalln(err)
	}
	return data, nil
}

func concurrentGetOB(names []string) ([]OrderBook, error) {
	//change this function from a a handler to _in_ another handler and feed in the []string of instrument names

	baseUrl := "https://www.deribit.com/api/v1/public/getorderbook?instrument="
	start := time.Now()
	// var list []OrderBook

	//////
	var result []OrderBook
	var wg sync.WaitGroup

	//loop over the fragments
	for _, v := range names {
		// Increment the WaitGroup counter.
		wg.Add(1)
		// Launch a goroutine to fetch the URL.
		go func(v string) {
			var data OrderBookResponse
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
			// Fetch the URL.
			url := baseUrl + v
			fmt.Println("getting:", url)
			res, err := http.Get(url)

			//decode response`
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Println("Routine exitng early: failed to read resp.body:", err)
				return
			}
			defer res.Body.Close()
			// bx := []byte(body)
			// fmt.Println("body:", string(bx))

			//unmarshall into the reponse struct
			err = json.Unmarshal([]byte(body), &data)
			if err != nil {
				log.Println("Routine exitng early: unmarshal resp.body into Orderbook:", err)
				return
			}
			result = append(result, data.Result)
		}(v)
	}
	// Wait for all HTTP fetches to complete.
	wg.Wait()

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
	return result, nil

}

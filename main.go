package main

import (
	"encoding/json"
	"fmt"
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

func main() {

	//routing
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})
	//test Route
	router.GET("/instruments", getInstruments)

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
	//unmarshall
	err = json.Unmarshal([]byte(contents), &data)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s\n", data)

	// xs := []string{"hello", "You", "cunt"}

	c.JSON(http.StatusOK, gin.H{"instruments": data})
}

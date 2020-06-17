package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

var db *sql.DB

type link struct {
	LongURL  string `json:"longURL,omitempty"`
	ShortURL string `json:"shortURL,omitempty"`
}

// updates the table to include shortURL key for a given id
func updateURL(id int, shortURL, longURL string) string {
	fmt.Println(longURL, "------------updateURL------------", shortURL)
	queryString := fmt.Sprintf("UPDATE urls SET shortURL = '%s' WHERE id = '%d' returning shortURL;", shortURL, id)
	fmt.Println(queryString)
	result, err := db.Query(queryString)
	if err != nil {
		log.Println("error executing updateURL query")
		panic(err)
	}
	defer result.Close()
	var resultStr string

	result.Next()
	result.Scan(&resultStr)
	fmt.Println(resultStr)
	return resultStr
}

// inserts long url into db
// - returns id
func addURL(url string) int {
	fmt.Println("------------addURL------------")
	queryString := fmt.Sprintf("insert into urls(longURL) values('%s') returning id;", url)
	fmt.Println(queryString)
	result, err := db.Query(queryString)
	if err != nil {
		log.Println("error in addURL")
		panic(err)
	}
	defer result.Close()

	//fmt.Println(err)
	var resultStr string
	resultStr = ""
	result.Next()
	err = result.Scan(&resultStr)
	if resultStr != "" {
		fmt.Println(err, "-", resultStr)

	}
	id, err := strconv.Atoi(resultStr)
	fmt.Print(err)
	return id
}

// curl -X PUT -H "Content-Type: application/json" -d '{"longURL":"www.swag.com"}' http://localhost:8080/ -i
func putURL(w http.ResponseWriter, r *http.Request) {
	log.Println("In putURL() function")
	var tempLink link
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal([]byte(body), &tempLink)
	if err != nil {
		fmt.Println(err)
	}
	log.Println("BODY", tempLink.LongURL)
	id := addURL(tempLink.LongURL)
	log.Println("id for link:", id)
	shortURL := idToShortURL(id)
	updateURL(id, shortURL, tempLink.LongURL)
	var resp link
	resp = link{"", shortURL}
	log.Println(resp)
	jsonData, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}

func lookupLongURL(shortURL string) string {
	queryString := fmt.Sprintf("SELECT longURL from urls WHERE shortURL='%s';", shortURL)
	result, _ := db.Query(queryString)
	var resultStr string
	result.Next()
	result.Scan(&resultStr)
	fmt.Println(resultStr)
	return resultStr
}

// curl -X GET -H "Content-Type: application/json" -d '{"shortURL":"beZ"}' http://localhost:8080/ -i
func getURL(w http.ResponseWriter, r *http.Request) {
	log.Println("in getURL() function")
	//parse JSON
	var tempLink link
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal([]byte(body), &tempLink)
	if err != nil {
		fmt.Println(err)
	}
	log.Println("BODY", tempLink.ShortURL)
	longURL := lookupLongURL(tempLink.ShortURL)
	type resp struct {
		URL string `json:"URL"`
	}
	respT := resp{longURL}
	jsonData, _ := json.Marshal(respT)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}

func apiHelper(w http.ResponseWriter, r *http.Request) {
	log.Println("in apiHelper()")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT, GET")
	//Access - Control - Allow - Headers
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	if r.Method == http.MethodPut {
		log.Println("Put request")
		putURL(w, r)
	} else if r.Method == http.MethodGet {
		log.Println("Get request")
		getURL(w, r)
	} else {
		log.Println("unsupported request type!")
	}

}

func main() {
	fmt.Println("hello server")
	var err error
	db, err = sql.Open("postgres", "postgres://twdcnlmu:Hd7RXw1kL22RCi6Qbn0rldKHJMfGcSXp@hansken.db.elephantsql.com:5432/twdcnlmu?connect_timeout=5")
	// postgres://twdcnlmu:Hd7RXw1kL22RCi6Qbn0rldKHJMfGcSXp@hansken.db.elephantsql.com:5432/twdcnlmu
	//Make sure you setup the ELEPHANTSQL_URL to be a uri, e.g. 'postgres://user:pass@host/db?options'
	fmt.Println(err)
	err = db.Ping()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("db connected")
	}

	//note USE OF GO FUNCTION. listenandserve blocks execution!
	go http.HandleFunc("/", apiHelper)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func idToShortURL(id int) string {
	strmap := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortURL := ""
	for id > 0 {
		shortURL += string(strmap[id%62])
		id /= 62
	}
	//reverse url
	return reverse(shortURL)
}

func shortURLtoID(shortURL string) int {
	id := 0
	strmap := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i, j := range reverse(shortURL) {
		id += strings.Index(strmap, string(j)) * int(math.Pow(float64(62), float64(i)))
	}
	fmt.Println()
	return id
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

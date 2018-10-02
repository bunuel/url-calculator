
package main
import (
	"net/http"
	"encoding/json"
	"time"
	"sync"
	"strconv"
	"fmt"
)

type TelData struct {
	Action string `json:"action"`
	X float64 `json:"x,number"`
	Y float64 `json:"y,number"`
	Answer float64 `json:"answer,number"`
	Cached bool `json:"cached,bool"`
	TimeAdded time.Time `json:"-"`
	URL string `json:"-"`
}

type cacheType struct {
	myData   map[string]TelData
	mux sync.Mutex
}

var myCache = cacheType{myData: make(map[string]TelData)}

func main() {
	
	http.HandleFunc("/add", performAction)
	http.HandleFunc("/subtract", performAction)
	http.HandleFunc("/multiply", performAction)
	http.HandleFunc("/divide", performAction)

	http.ListenAndServe(":8080",nil)
 }

 

 func performAction (w http.ResponseWriter, r *http.Request) {
	
	myData := make(map[string]TelData)
	
	myCache.mux.Lock()
	myData = myCache.myData
	myCache.mux.Unlock()

	myURL := r.URL.String()

	cached := isItCached(w, myURL, myData)

	if cached {
		return
	}

	x, err  := strconv.ParseFloat(r.FormValue("x"), 64)
	if err != nil {
		returnErrorPage(w, "Please include a value for x in your URL")
		return
	}

	y, err  := strconv.ParseFloat(r.FormValue("y"), 64)
	if err != nil {
		returnErrorPage(w, "Please include a value for y in your URL")
		return
	}

	currentTime := time.Now()
	action := r.URL.Path[1:]
	answer := float64(0)

	switch {
    case action == "add":
		answer = x + y
	case action == "subtract":
		answer = x - y
	case action == "multiply":
		answer = x * y
	case action == "divide":
		if (y == 0) {
			returnErrorPage(w, "Please don't divide by zero : (")
			return
			answer = 0
		} else {
			answer = x / y
		}
    default:
		returnErrorPage(w, "Please request a valid URL... /add, /subtract, /multiply, /divide")
		return
	}

	teldata := TelData{action, x, y, answer, cached, currentTime, myURL}
	writeResponse(w, myURL, teldata, myData)
	
 }

 func returnErrorPage(w http.ResponseWriter, myError string) {
	fmt.Fprint(w, myError)
 }

 func isItCached(w http.ResponseWriter, myURL string, myData map[string]TelData) bool {
	  
	if val, ok := myData[myURL]; ok {

		duration := time.Since(val.TimeAdded)

		if ( duration > time.Minute ) {
			delete(myData, myURL)
			myCache.mux.Lock()
			myCache.myData = myData
			myCache.mux.Unlock()
		} else {
			val.Cached = true
			teldata := val
			writeResponse(w, myURL, teldata, myData)
			return true
		}

	}
	return false
}

func writeResponse(w http.ResponseWriter, myURL string, teldata TelData, myData map[string]TelData) {

	myData[myURL] = teldata

	myCache.mux.Lock()
	myCache.myData = myData
	myCache.mux.Unlock()

	js, err := json.Marshal(teldata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

 }

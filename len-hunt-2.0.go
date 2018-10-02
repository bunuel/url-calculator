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
	timeAdded time.Time
	//URL string `json:"-"`
}


/*
type MutexMap struct {
    m map[string]string
    *sync.RWMutex
}
*/

type cacheType struct {
	myData   map[string]TelData
	*sync.RWMutex
}

var myCache = cacheType{myData: make(map[string]TelData)}
/*
type myData struct {
	map[string]TelData
}
*/

//The correct way to use a cache or anything global is to use a sync.Mutex or a sync.RWMutex. This allows you to "lock" something so you prevent race conditions.

//var myCache = cacheType{myData: make(map[string]TelData)}

func main() {
	//myCache.mux.Lock()
	/*
	s := &Stat{
		counters: make(map[string]int64),
	}
	
	myCache := &cacheType{
		myData: make(map[string]TelData),
	}
	*/

	// initialize the mutex lock here
	// then pass it to each function?zz

	http.HandleFunc("/add", performAction)
	http.HandleFunc("/subtract", performAction)
	http.HandleFunc("/multiply", performAction)
	http.HandleFunc("/divide", performAction)

	http.ListenAndServe(":8080",nil)
}



func getTeldataFromCache(myURL string, myData map[string]TelData) {
	fmt.Println(myData)
	//for k, v := range myData { 
		//fmt.Println(myURL, k, ":", v)
		//if (k == myURL) {
			//return v[k]
			//fmt.Println(k,v)
			//return true
		
	//}
	//fmt.Println(k, ":", v)
}


func ReadSomething(myData map[string]TelData, myURL string, w http.ResponseWriter) TelData {
	teldata := TelData{"action", 0, 0, 0, false, time.Now()}
    myCache.RLock() // locking for writing here because we may have to delete something
	defer myCache.RUnlock() // make SURE you do this, else it will be locked permanently
	//if val, ok := o.myData[myURL]; ok {
	//if _, present := o.myData[myURL] {
		//fmt.Printf("%v\nis cached", myCache.myData[myURL])
		//return myCache.myData[myURL]
		//teldata := 
		
		getTeldataFromCache(myURL, myCache.myData)

		if teldata.Cached {
			// do something
			teldata = myCache.myData[myURL]
			teldata.Cached = true
			
		} else {
			//teldata.Cached = false
			teldata.Cached = false

		}
		return teldata
	}
		

		
func WriteSomething(myCache cacheType, myURL string, value TelData) {
    myCache.Lock() // lock for writing, blocks until the Mutex is ready
    defer myCache.Unlock() // again, make SURE you do this, else it will be locked permanently
	myCache.myData[myURL] = value
	// set expiration here?!!
}

func deleteSomething(myCache cacheType, myURL string) {
    myCache.Lock() // lock for writing, blocks until the Mutex is ready
    defer myCache.Unlock() // again, make SURE you do this, else it will be locked permanently
	delete(myCache.myData, myURL)
	// set expiration here?!!
}

func performAction (w http.ResponseWriter, r *http.Request) {
	
	//myData := make(map[string]TelData)
	
	//myCache.mux.Unlock()
	//myData = myCache.myData
	//myCache.mux.Lock()

	myURL := r.URL.String()



	//cached := isItCached(w, myURL, myData)
	teldata := ReadSomething(myCache.myData, myURL, w)

	if teldata.Cached  {
		//return
		// if it's cached already just write it and return
		writeResponse(w, myURL, teldata)
	}

	x, err  := strconv.ParseFloat(r.FormValue("x"), 64)
	y, err  := strconv.ParseFloat(r.FormValue("y"), 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
			returnErrorPage(w, "Please don't divide by zero")
			answer = 0
		} else {
			answer = x / y
		}
    default:
        returnErrorPage(w, "Please request a valid URL... /add, /subtract, /multiply, /divide")
	}

	teldata = TelData{action, x, y, answer, false, currentTime}
	writeResponse(w, myURL, teldata)
	
 }

 func returnErrorPage(w http.ResponseWriter, myError string) {
	fmt.Fprint(w, myError)
 }

 /*
 func isItCached(w http.ResponseWriter, myURL string, myData map[string]TelData) bool {
	  
	if val, ok := myData[myURL]; ok {

		duration := time.Since(val.timeAdded)

		if ( duration > time.Minute ) {
			delete(myData, myURL)
			myCache.mux.Unlock()
			myCache.myData = myData
			myCache.mux.Lock()
		} else {
			val.Cached = true
			teldata := val
			writeResponse(w, myURL, teldata, myData)
			return true
		}

	}
	return false
}
*/

func writeResponse(w http.ResponseWriter, myURL string, teldata TelData) {

	//myData[myURL] = teldata
/*
	myCache.mux.Unlock()
	myCache.myData = myData
	myCache.mux.Lock()
	*/
	WriteSomething(myCache, myURL, teldata)

	js, err := json.Marshal(teldata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

 }

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
)

type idLocation struct {
	Id   int     `json:"id"`
	Lat  float64 `json:"latitude"`
	Long float64 `json:"longitude"`
}

// wrapping this in a closure for scope purposes?
// coming back to this 4 months after I wrote it. Did I need to wrap this in a seperate function?
func sendResponse(locations []idLocation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//testing http headers
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("X-Custom-Header", "HelloWorld")

		param1 := r.URL.Query().Get("lat")
		param2 := r.URL.Query().Get("long")

		lat, err1 := strconv.ParseFloat(param1, 64)
		long, err2 := strconv.ParseFloat(param2, 64)

		//handling conversion errors
		if err1 != nil || err2 != nil {
			fmt.Println(param1 + " " + param2)
			http.Error(w, "Invalid float values", http.StatusBadRequest)
			return
		}

		resp, err := getNearby(1, lat, long, locations)
		if err != nil {
			fmt.Println("RAHHHH GETNEARBY ERROR")
		}
		fmt.Fprint(w, resp)
	}

}

func getNearby(threshhold float64, lat float64, long float64, compArray []idLocation) (string, error) {
	// start := time.Now()

	//CANNOT RETURN THE IDLOCATION WITHOUT STRUCTURING IT AS JSON.
	//THERE IS DEF A HEURISTIC FOR THIS BY STARTING FROM THE ORIGIN OF THE WORLD AND SORTING FROM THERE BUT THAT'S FOR A LEETCODE PROBLEM NOT A PERSONAL PROJECT
	nearby := []idLocation{}
	for _, store := range compArray {
		if math.Sqrt(math.Abs(store.Lat-lat)+math.Abs(store.Long-long)) < threshhold {
			nearby = append(nearby, store)
		}
	}
	// elapsed := time.Since((start))
	// fmt.Println("lat: " + strconv.FormatFloat(lat, 'f', -1, 64) + " long: " + strconv.FormatFloat(long, 'f', -1, 64) + " " + elapsed.String())

	fmt.Println(nearby)
	jsonData, err := json.Marshal(nearby)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// creating an http handler with cors rules enabled on it
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {

	//friggin CORS
	// corsMiddleware := cors.New(cors.Options{
	// 	// Specify allowed origins explicitly instead of using * for security
	// 	AllowedOrigins: []string{"http://localhost:3000", "https://yourdomain.com"},
	// 	// Only allow specific HTTP methods needed by your API
	// 	AllowedMethods: []string{"GET", "POST", "OPTIONS"},
	// 	// Only allow needed headers
	// 	AllowedHeaders: []string{"Content-Type", "Authorization"},
	// 	// Don't allow credentials by default unless you specifically need them
	// 	AllowCredentials: false,
	// 	// Cache preflight results to improve performance
	// 	MaxAge: 86400, // 24 hours
	// })

	//initialize list of locations
	file, err := os.Open("updated_output.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	//read and decode
	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)

	// Read the JSON data into a slice of Location
	var locations []idLocation
	if err := decoder.Decode(&locations); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	//fmt.Println(locations[5])
	// getNearby(.4, 40.4387, -79.9972, locations)
	//handle requests for locations
	mux := http.NewServeMux()
	mux.HandleFunc("/givechipotles", sendResponse(locations))

	handler := corsMiddleware(mux)

	fmt.Println("Server is running on :8080...")
	http.ListenAndServe(":8080", handler)
}

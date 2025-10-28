package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/joho/godotenv"
)

type idLocation struct {
	Id   int     `json:"id"`
	Lat  float64 `json:"latitude"`
	Long float64 `json:"longitude"`
}

// type idMenu struct {
// 	Id   int            `json:"id"`
// 	Menu map[string]int `json:"menu`
// }

// Combined struct that has to be main
type Chipotle struct {
	Id   int            `json:"id"`
	Menu map[string]int `json:"menu"`
	Lat  float64        `json:"latitude"`
	Long float64        `json:"longitude"`
}

// package level variables
// var LATS []idLocation
// var longs []idLocation
// var menus []map[string]int

//DATA HANDLING

// return -> nested array with lat and long sorted array, id map
func buildArrays(path string) ([]idLocation, []idLocation, []Chipotle, error) {
	//initialize list of locations
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error on file open")
		return nil, nil, nil, err
	}
	defer file.Close()

	//read and decode
	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)

	// Read the JSON data into a slice of stores
	var chipotles []Chipotle

	// I don't like this but it works
	if err := decoder.Decode(&chipotles); err != nil {
		fmt.Println("Error on lats")
		return nil, nil, nil, err
	}

	//test fill
	// fmt.Println("First item: ", chipotles[0])

	//DO NOT FORGET ITDS START ON 1

	//TODO: Decode into new arrays
	//lats is copy of chips with idlocation
	var lats []idLocation
	for _, b := range chipotles {
		lats = append(lats, idLocation{
			//MOVING ID TO 0 INDEXING. BE CAREFUL WITH THIS
			Id:   b.Id - 1,
			Lat:  b.Lat,
			Long: b.Long,
		})
	}

	sort.Slice(lats, func(i, j int) bool {
		return lats[i].Lat < lats[j].Lat
	})

	//long is resorted copy of lats
	longs := make([]idLocation, len(lats))
	copy(longs, lats)

	sort.Slice(longs, func(i, j int) bool {
		return longs[i].Long < longs[j].Long
	})

	// //building menu arrays
	// var menus []map[string]int
	// for _, chip := range chipotles {
	// 	menus = append(menus, chip.Menu)
	// }
	return lats, longs, chipotles, nil
}

func getNearby(threshold float64, lat float64, long float64, latList []idLocation, longList []idLocation, chipotles []Chipotle) ([]Chipotle, error) {
	// start := time.Now()

	//search for upper and lower boundary of value + threshold
	//value: array to search, array: array of idlocations(sorted by lat or long), searchLong: whether to search long(0) or lat(1)
	boundarySearch := func(value float64, array []idLocation, searchLat bool) []idLocation {
		if len(array) == 0 {
			fmt.Println("Empty array")
			return []idLocation{}
		}
		//find upper boundary
		if searchLat {
			tgt := value + threshold
			lowi, upi := 0, len(array)-1

			// Find upper bound (last index <= tgt)
			for lowi <= upi {
				mid := lowi + (upi-lowi)/2
				if array[mid].Lat <= tgt {
					lowi = mid + 1
				} else {
					upi = mid - 1
				}
			}
			upperBound := upi // inclusive

			// Find lower bound (first index >= value - threshold)
			tgt = value - threshold
			lowi, upi = 0, len(array)-1
			for lowi <= upi {
				mid := lowi + (upi-lowi)/2
				if array[mid].Lat < tgt {
					lowi = mid + 1
				} else {
					upi = mid - 1
				}
			}
			lowerBound := lowi

			// Slice range: inclusive bounds
			return array[lowerBound : upperBound+1]
		} else {
			tgt := value + threshold
			lowi, upi := 0, len(array)-1

			// Find upper bound (last index <= tgt)
			for lowi <= upi {
				mid := lowi + (upi-lowi)/2
				if array[mid].Long <= tgt {
					lowi = mid + 1
				} else {
					upi = mid - 1
				}
			}
			upperBound := upi // inclusive

			// Find lower bound (first index >= value - threshold)
			tgt = value - threshold
			lowi, upi = 0, len(array)-1
			for lowi <= upi {
				mid := lowi + (upi-lowi)/2
				if array[mid].Long < tgt {
					lowi = mid + 1
				} else {
					upi = mid - 1
				}
			}
			lowerBound := lowi
			// fmt.Println("Lower Bound: ", lowerBound, "Upper Bound: ", upperBound)
			// fmt.Println("Bounds: ", array[lowerBound:upperBound+1])
			return array[lowerBound : upperBound+1]
		}
	}
	//CANNOT RETURN THE IDLOCATION WITHOUT STRUCTURING IT AS JSON.
	//THERE IS DEF A HEURISTIC FOR THIS BY STARTING FROM THE ORIGIN OF THE WORLD AND SORTING FROM THERE BUT THAT'S FOR A LEETCODE PROBLEM NOT A PERSONAL PROJECT
	//USING RTREE DATA STRUCTURE... PLAN B IS TO USE 2 BINARY SEARCHES AND FIND WHAT IS IN RANGE

	//TODO: LOCATION IDS:
	// double b search returning array of ids
	latRange := boundarySearch(lat, latList, true)
	longRange := boundarySearch(long, longList, false)

	// fmt.Println("latrange", latRange)
	//Intersection these
	m := make(map[int]idLocation)
	for _, v := range latRange {
		m[v.Id] = v
	}

	var common []idLocation
	for _, v := range longRange {
		if _, ok := m[v.Id]; ok {
			common = append(common, v)
		}
	}

	//return menus of all these from the indexes
	sending := []Chipotle{}
	for _, item := range common {
		sending = append(sending, chipotles[item.Id])
	}
	fmt.Println(sending)
	return sending, nil
}

//WEB CODE

// wrapping this in a closure for scope purposes?
// coming back to this 4 months after I wrote it. Did I need to wrap this in a seperate function?
func sendResponse(lats []idLocation, longs []idLocation, chips []Chipotle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//testing http headers
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("X-Custom-Header", "HelloWorld")

		param1 := r.URL.Query().Get("lat")
		param2 := r.URL.Query().Get("long")

		lat, err1 := strconv.ParseFloat(param1, 64)
		long, err2 := strconv.ParseFloat(param2, 64)
		fmt.Println("lat", lat, "long", long)
		//handling conversion errors
		if err1 != nil || err2 != nil {
			fmt.Println("params", param1+" "+param2)
			http.Error(w, "Invalid float values", http.StatusBadRequest)
			return
		}

		//
		resp, err := getNearby(1, lat, long, lats, longs, chips)
		if err != nil {
			fmt.Println("RAHHHH GETNEARBY ERROR")
		}

		//convert to json
		bytes, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		fmt.Fprint(w, string(bytes))
	}

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

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//NEW INTERPRETATIONS: MENUS ARE STORED IN A MAP SORTED BY ID, LOCATIONS ARE STORED IN 2 SORTED ARRAYS
	path := os.Getenv("json_path")
	lats, longs, menuMap, err := buildArrays(path)
	// fmt.Println("lats", lats)
	if err != nil {
		fmt.Println("Error building arrays:", err)
	}

	//fmt.Println(locations[5])
	// getNearby(.4, 40.4387, -79.9972, menuMap)
	//handle requests for locations
	mux := http.NewServeMux()
	mux.HandleFunc("/givechipotles", sendResponse(lats, longs, menuMap))

	handler := corsMiddleware(mux)

	fmt.Println("Server is running on :8080...")
	http.ListenAndServe(":8080", handler)
}

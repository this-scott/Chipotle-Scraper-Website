package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
)

type idLocation struct {
	Id   int     `json:"id"`
	Lat  float64 `json:"latitude"`
	Long float64 `json:"longitude"`
}

type idMenu struct {
	Id   int            `json:"id"`
	Menu map[string]int `json:"menu`
}

// Combined struct that has to be main
type LocationWithMenu struct {
	Menu map[string]idMenu `json:"menu"`
	Id   int               `json:"id"`
	Lat  float64           `json:"latitude"`
	Long float64           `json:"longitude"`
}

// package level variables
// var latList []idLocation
// var longlist []idLocation
// var menuMap map[string]idMenu

//DATA HANDLING

// return -> nested array with lat and long sorted array, id map
func buildArrays(path string) ([]idLocation, []idLocation, map[int]idMenu, error) {
	//initialize list of locations
	file, err := os.Open("updated_output.json")
	if err != nil {
		fmt.Println("Error on file open")
		return nil, nil, nil, err
	}
	defer file.Close()

	//read and decode
	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)

	// Read the JSON data into a slice of Location
	var lats []idLocation
	var menus map[int]idMenu

	// I don't like this but it works
	if err := decoder.Decode(&lats); err != nil {
		fmt.Println("Error on lats")
		return nil, nil, nil, err
	}
	if err := decoder.Decode(&menus); err != nil {
		fmt.Println("Error on menus")
		return nil, nil, nil, err
	}

	//sort to get lats
	sort.Slice(lats, func(i, j int) bool {
		return lats[i].Lat < lats[j].Lat
	})

	longs := make([]idLocation, len(lats))
	copy(longs, lats) // Sort by latitude (ascending)

	sort.Slice(longs, func(i, j int) bool {
		return longs[i].Long < longs[j].Long
	})
	return lats, longs, menus, nil
}

// parse
func getNearby(threshhold float64, lat float64, long float64, latList []idLocation, longList []idLocation, menuMap map[int]idMenu) ([]idLocation, error) {
	// start := time.Now()
	fmt.Println(menuMap)

	//search for upper and lower boundary of value + threshold
	//value: array to search, array: array of idlocations(sorted by lat or long), searchLong: whether to search long(0) or lat(1)
	boundarySearch := func(value float64, array []idLocation, searchField bool) []idLocation {
		//this trick is so funny
		//find upper boundary

		fmt.Println(array)
		if searchField {
			tgt := value + threshhold
			lowi, upi := 0, len(array)-1

			for lowi <= upi {
				mid := lowi + (upi-lowi)/2
				if array[mid].Lat <= tgt {
					lowi = mid + 1 // Search in the right half
				} else {
					upi = mid - 1 // Search in the left half
				}
			}

			upperBound := upi

			//find lower boundary
			tgt = value - threshhold
			lowi, upi = 0, len(array)-1
			for lowi <= upi {
				mid := lowi + (upi-lowi)/2
				if array[mid].Lat < tgt {
					lowi = mid + 1 // Search in the right half
				} else {
					upi = mid - 1 // Search in the left half
				}
			}

			lowerBound := lowi

			return array[lowerBound : upperBound+1]
		} else {
			tgt := value + threshhold
			lowi, upi := 0, len(array)-1

			for lowi <= upi {
				mid := lowi + (upi-lowi)/2
				if array[mid].Long <= tgt {
					lowi = mid + 1 // Search in the right half
				} else {
					upi = mid - 1 // Search in the left half
				}
			}

			upperBound := upi

			//find lower boundary
			tgt = value - threshhold
			lowi, upi = 0, len(array)-1
			for lowi <= upi {
				mid := lowi + (upi-lowi)/2
				if array[mid].Long < tgt {
					lowi = mid + 1 // Search in the right half
				} else {
					upi = mid - 1 // Search in the left half
				}
			}

			lowerBound := lowi
			return array[lowerBound : upperBound+1]
		}
	}
	//CANNOT RETURN THE IDLOCATION WITHOUT STRUCTURING IT AS JSON.
	//THERE IS DEF A HEURISTIC FOR THIS BY STARTING FROM THE ORIGIN OF THE WORLD AND SORTING FROM THERE BUT THAT'S FOR A LEETCODE PROBLEM NOT A PERSONAL PROJECT
	//USING RTREE DATA STRUCTURE... PLAN B IS TO USE 2 BINARY SEARCHES AND FIND WHAT IS IN RANGE
	//I LIED, DOUBLE B SEARCH

	//TODO: LOCATION IDS:
	// double b search returning array of ids

	//TODO: GET MENUS FROM IDS:

	// nearby := []idLocation{}
	// for _, store := range compArray {
	// 	if math.Sqrt(math.Abs(store.Lat-lat)+math.Abs(store.Long-long)) < threshhold {
	// 		nearby = append(nearby, store)
	// 	}
	// }
	// elapsed := time.Since((start))
	// fmt.Println("lat: " + strconv.FormatFloat(lat, 'f', -1, 64) + " long: " + strconv.FormatFloat(long, 'f', -1, 64) + " " + elapsed.String())

	latRange := boundarySearch(lat, latList, false)
	longRange := boundarySearch(long, longList, true)
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
	sending := []idMenu{}
	for _, item := range common {
		sending = append(sending, menuMap[item.Id])
	}

	return common, nil
}

//WEB CODE

// wrapping this in a closure for scope purposes?
// coming back to this 4 months after I wrote it. Did I need to wrap this in a seperate function?
func sendResponse(lats []idLocation, longs []idLocation, menuMap map[int]idMenu) http.HandlerFunc {
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

		//
		resp, err := getNearby(1, lat, long, lats, longs, menuMap)
		if err != nil {
			fmt.Println("RAHHHH GETNEARBY ERROR")
		}
		fmt.Fprint(w, resp)
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

	//NEW INTERPRETATIONS: MENUS ARE STORED IN A MAP SORTED BY ID, LOCATIONS ARE STORED IN 2 SORTED ARRAYS
	lats, longs, menuMap, err := buildArrays(os.Getenv("json_path"))

	if err != nil {
		fmt.Println("Error building arrays:", err)
	}
	//fmt.Println(locations[5])
	// getNearby(.4, 40.4387, -79.9972, locations)
	//handle requests for locations
	mux := http.NewServeMux()
	mux.HandleFunc("/givechipotles", sendResponse(lats, longs, menuMap))

	handler := corsMiddleware(mux)

	fmt.Println("Server is running on :8080...")
	http.ListenAndServe(":8080", handler)
}

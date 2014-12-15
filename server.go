package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	//	"net/url"
	"flag"
	"github.com/gorilla/mux"
	"github.com/kellydunn/golang-geo"
	"log"
	"net/http"
	//	"github.com/olivere/elastic"
	"os"
)

type ElasticSearchHit struct {
	Index    string         `json:"_index"`
	Type     string         `json:"_type"`
	Score    float64        `json:"_score"`
	Provider HealthProvider `json:"_source"`
}

type ElasticSearchHits struct {
	Total    uint32             `json:"total"`
	MaxScore float64            `json:"max_score"`
	Hits     []ElasticSearchHit `json:"hits"`
}

type ElasticSearchResponse struct {
	Timedout bool              `json:"timed_out"`
	Took     uint32            `json:"took"`
	Hits     ElasticSearchHits `json:"hits"`
}

type HealthProvider struct {
	APC                              string  `json:"apc"` // The service
	ProviderId                       string  `json:"providerid"`
	ProviderName                     string  `json:"providername"`
	ProviderStreetAddress            string  `json:"providerstreetaddress"`
	ProviderCity                     string  `json:"providercity"`
	ProviderState                    string  `json:"providerstate"`
	ProviderZipCode                  string  `json:"providerzipcode"`
	ProviderHRR                      string  `json:"providerhrr"` // HRR = Hospital Referral Region
	OutpatientServices               int64   `json:"outpatientservices"`
	AverageEstimatedSubmittedCharges float64 `json:"averageestimatedsubmittedcharges"`
	AverageTotalPayments             float64 `json:"averagetotalpayments"`
	GpsLocation                      string  `json:"gpslocation"`
}

// var searcher elastic.Client
var restsearcher *http.Client
var geocoder *geo.GoogleGeocoder

// var hostnameAndPort string

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryParams := r.URL.Query()
	address := queryParams.Get("address")
	miles := queryParams.Get("miles")
	procedure := queryParams.Get("procedure")
	sortBy := queryParams.Get("sortby")
	if address == "" {
		http.Error(w, "The required 'address' query parameter was not present", http.StatusBadRequest)
		return
	}
	if procedure == "" {
		http.Error(w, "The required 'procedure' query parameter was not present", http.StatusBadRequest)
		return
	}
	if miles == "" {
		http.Error(w, "The required 'miles' query parameter was not present", http.StatusBadRequest)
		return
	}
	if sortBy == "" {
		http.Error(w, "The required 'sortby' query parameter was not present", http.StatusBadRequest)
		return
	}
	log.Printf("address: %s, procedure: %s, miles: %s", address, procedure, miles)
	// Fetch the GPS for that address
	geocode, err := geocoder.Geocode(address)
	if err != nil {
		log.Printf("Could not find geo code for %s", address)
		w.Write([]byte(`{ "healthproviders": []}`))
		return
	}
	log.Printf("The geocode for '%s' is '%f,%f'", address, geocode.Lat(), geocode.Lng())

	// Search for the given values
	query := fmt.Sprintf(`{"from": 0, "size": 100, "query":{"filtered":{"query":{"match":{"apc":"%s"}}, "filter":{"geo_distance":{"distance":"%smi", "service.gpslocation":"%f, %f"}}}},
	"sort":[{"averageestimatedsubmittedcharges" : { "order" : "asc" } }]}`, procedure, miles, geocode.Lat(), geocode.Lng())
	log.Printf("query: %s", query)
	results, err := restsearcher.Post("http://localhost:9200/healthadvisor/service/_search", "application/x-www-form-urlencoded", bytes.NewBufferString(query))
	if err != nil {
		http.Error(w, "There was an error talking to the search engine", http.StatusInternalServerError)
		return
	}
	defer results.Body.Close()
	body, err := ioutil.ReadAll(results.Body)
	if err != nil {
		http.Error(w, "There was an error deserializing the data from ES", http.StatusInternalServerError)
		return
	}
	log.Printf("Got a response of size %d from ES", len(body))
	decoder := json.NewDecoder(bytes.NewReader(body))
	var esResponse ElasticSearchResponse
	decoder.Decode(&esResponse)
	providers := make([]HealthProvider, len(esResponse.Hits.Hits))
	for i, hit := range esResponse.Hits.Hits {
		providers[i] = hit.Provider
	}
	serializedProviders, err := json.Marshal(providers)
	if err != nil {
		http.Error(w, "There was an error marshaling the providers to JSON", http.StatusInternalServerError)
		return
	}
	w.Write(serializedProviders)
}

func main() {
	// Read configuration parameters
	searchHostname := os.Getenv("ES_HOSTNAME")
	searchPort := os.Getenv("ES_PORT")
	// Validate configuration parameters
	if searchHostname == "" {
		log.Println("No searchHostname provided, will exit")
		return
	}
	if searchPort == "" {
		log.Println("No searchPort provided, will exit")
		return
	}
	// Setting up the HTTP server
	log.Println("Starting HTTP server")
	//	hostnameAndPort := fmt.Sprintf("%s:%s", searchHostname, searchPort)

	/*
		searcher, err := elastic.NewClient(http.DefaultClient, hostnameAndPort)
		if err != nil {
			log.Printf("Could not connect to ES: %s", err)
			return
		}
	*/
	restsearcher = &http.Client{}
	geocoder = &geo.GoogleGeocoder{}
	r := mux.NewRouter()
	r.HandleFunc("/@api/healthproviders", SearchHandler).Methods("GET")
	http.Handle("/@api/", r)

	// Server public assets
	http.Handle("/", http.FileServer(http.Dir("./public/")))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	http.ListenAndServe(":"+port, nil)
}

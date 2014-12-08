package main

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"bytes"
//	"net/url"
	"net/http"
	"flag"
	"log"
	"github.com/gorilla/mux"
	"github.com/kellydunn/golang-geo"
//	"github.com/olivere/elastic"
)

type ElasticSearchHit struct {
	Index string `json:"_index"`
	Type string `json:"_type"`
	Score float64 `json:"_score"`
	Provider HealthProvider `json:"_source"`
}

type ElasticSearchHits struct {
	Total uint32 `json:"total"`
	MaxScore float64 `json:"max_score"`
	Hits []ElasticSearchHit `json:"hits"`
}

type ElasticSearchResponse struct {
	Timedout bool `json:"timed_out"`
	Took uint32 `json:"took"`
	Hits ElasticSearchHits `json:"hits"`
}

type HealthProvider struct {
	APC string `json:"apc"` // The service
	ProviderId string `json:"providerid"`
	ProviderName string `json:"providername"`
	ProviderStreetAddress string `json:"providerstreetaddress"`
	ProviderCity string `json:"providercity"`
	ProviderState string `json:"providerstate"`
	ProviderZipCode string `json:"providerzipcode"`
	ProviderHRR string `json:"providerhrr"` // HRR = Hospital Referral Region
	OutpatientServices int64 `json:"outpatientservices"`
	AverageEstimatedSubmittedCharges float64 `json:"averageestimatedsubmittedcharges"`
	AverageTotalPayments float64 `json:"averagetotalpayments"`
	GpsLocation string `json:"gpslocation"`
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
	// TODO: transform miles to km
	query := fmt.Sprintf(`{"query":{"filtered":{"query":{"match":{"apc":"%s"}}, "filter":{"geo_distance":{"distance":"100km", "service.gpslocation":"%f, %f"}}}}}`, procedure, geocode.Lat(), geocode.Lng())
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
	providersHits := make([]HealthProvider, len(esResponse.Hits.Hits))
	for i, hit := range esResponse.Hits.Hits {
		providerHits[i] = hit
	}
	hits, err := json.Marshal(esResponse.Hits.Hits)
	if err != nil {
		http.Error(w, "There was an error marshaling the hits array", http.StatusInternalServerError)
		return
	}
	w.Write(providerHits)
}

func main() {
	// Read configuration parameters
	var searchHostname string
	var searchPort string
	flag.StringVar(&searchHostname, "search-hostname", "", "The search engine's hostname")
	flag.StringVar(&searchPort, "search-port", "", "The search engine's port")
	flag.Parse()
	// Validate configuration parameters
	if searchHostname == "" {
		log.Println("No searchHostname provided, will exit")
		return
	}
	if searchPort ==  "" {
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
	restsearcher = &http.Client {}
    geocoder = &geo.GoogleGeocoder {}
	r := mux.NewRouter()
	r.HandleFunc("/@api/healthproviders", SearchHandler).Methods("GET")
	http.Handle("/@api/", r)

	// Server public assets
	http.Handle("/", http.FileServer(http.Dir("./public/")))

	log.Println("Listing on port 8080")
	http.ListenAndServe(":8080", nil)
}

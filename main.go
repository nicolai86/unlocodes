package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	unlocodes = make(map[string]string)
	dataPath  string
	listen    string
)

func init() {
	flag.StringVar(&listen, "listen", ":3030", "http listen address")
	flag.StringVar(&dataPath, "data", "combined.csv", "combined.csv path")

	flag.Parse()

	f, err := os.Open(dataPath)
	if err != nil {
		log.Fatalf("Failed to open csv: %s", err)
	}

	r := csv.NewReader(f)

	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		if record[2] != "" {
			unlocodes[fmt.Sprintf("%s%s", record[1], record[2])] = record[3]
		}
	}
}

func UnLocodeLookup(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	codes, ok := req.URL.Query()["code"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	names := make(map[string]string)
	for _, code := range codes {
		if translation, ok := unlocodes[code]; ok {
			names[code] = translation
		}
	}
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(names)
	w.Write(b)
}

func Handler() http.Handler {
	mux := http.DefaultServeMux
	mux.HandleFunc("/", UnLocodeLookup)

	return http.Handler(mux)
}

func main() {
	fmt.Printf("Listening on %s\n", listen)
	http.ListenAndServe(listen, Handler())
}

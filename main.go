package main

import (
	"archive/zip"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var (
	unlocodes = make(map[string]string)
	dataPath  string
	listen    string
)

func parseLocode(f io.Reader) {
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

func parseLocodes() {
	r, err := zip.OpenReader(dataPath)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.Index(f.Name, "CodeListPart") != -1 {
			rc, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}
			parseLocode(rc)
			rc.Close()
		}
	}
}

func init() {
	flag.StringVar(&listen, "listen", ":3030", "http listen address")
	flag.StringVar(&dataPath, "data", "loc152csv.zip", "path to UNECE locodes")

	flag.Parse()

	parseLocodes()
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

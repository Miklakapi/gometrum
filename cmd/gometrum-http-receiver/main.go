package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	addr := "0.0.0.0:8088"

	http.HandleFunc("/ingest", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		codec := r.URL.Query().Get("codec")

		switch codec {

		case "text":
			fmt.Print(string(body))

		case "event_json":
			fmt.Println(string(body))

		case "ndjson":
			sc := bufio.NewScanner(strings.NewReader(string(body)))
			for sc.Scan() {
				fmt.Println(sc.Text())
			}

		case "loki":
			var p struct {
				Streams []struct {
					Stream map[string]string `json:"stream"`
					Values [][]string        `json:"values"`
				} `json:"streams"`
			}
			if json.Unmarshal(body, &p) == nil {
				for _, s := range p.Streams {
					for _, v := range s.Values {
						if len(v) == 2 {
							fmt.Println(v[1])
						}
					}
				}
			}

		default:
			fmt.Println(string(body))
		}

		w.WriteHeader(http.StatusNoContent)
	})

	log.Printf("HTTP receiver listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

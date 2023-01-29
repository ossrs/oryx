package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const stable = "v1.0.279"
const latest = "v1.0.279"
const api = "v1.0.362";

func main() {
	listen := os.Getenv("PORT")
	if listen == "" {
		listen = ":2023"
	}

	fmt.Println(fmt.Sprintf("Run with listen=%v", listen))

	ep := "/terraform/v1/releases"
	fmt.Println(fmt.Sprintf("Serve at %v", ep))
	http.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		res := &struct {
			Stable string `json:"stable"`
			Latest string `json:"latest"`
			API string `json:"api"`
		}{
			Stable: stable,
			Latest: latest,
			API: api,
		}

		w.Header().Set("Server", fmt.Sprintf("srs-cloud/%v", api))

		b,err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(b)
	})

	if err := http.ListenAndServe(listen, nil); err != nil {
		panic(err)
	}
}

//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: MIT
//
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// We should keep the stable version as 256, because for new architecture, we don't support automatically upgrade, so
// this feature is actually not used, but we should keep a specified version for compatibility.
const stable = "v1.0.256";
const latest = "v1.0.281";
const api = "v1.0.367";

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

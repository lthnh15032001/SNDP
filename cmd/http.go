package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/spf13/cobra"
)

func newHttpCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "http",
		Short: "Iot Streaming http Log",
		Long:  "Iot Streaming http Log",
		Run: func(cmd *cobra.Command, args []string) {
			httpServe()
		},
	}
}

func httpServe() {
	// Define the local server to be exposed
	localServerURL, _ := url.Parse("http://localhost:8080")

	// Create a reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(localServerURL)

	// Handle incoming requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Modify the request to point to the local server
		r.URL = localServerURL
		r.Host = localServerURL.Host

		// Serve the request using the reverse proxy
		proxy.ServeHTTP(w, r)
	})

	// Start the proxy server on port 8081
	go func() {
		fmt.Println("Proxy server is running on :8081")
		if err := http.ListenAndServe(":8081", nil); err != nil {
			fmt.Println(err)
		}
	}()

	select {}
}

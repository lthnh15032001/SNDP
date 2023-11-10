package iot

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/jpillora/requestlog"
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
	localServerURL, _ := url.Parse("http://192.168.1.97:9000")

	// Create a reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(localServerURL)

	// Handle incoming requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Modify the request to point to the local server
		r.URL = localServerURL
		r.Host = localServerURL.Host

	})

	// Start the proxy server on port 8081
	go func() {
		fmt.Println("Proxy server is running on :8081")
		if err := http.ListenAndServe("103.191.147.139:8082", requestlog.Wrap(logRequest(proxy))); err != nil {
			fmt.Println(err)
		}
	}()

	select {}
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
}

package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Proxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

// receive host URL
// Check its subdomain
// if subdomain exists,
//
//	then we can route to the correct(LB dns) service
//
// else
//
//	route to the default(main) service
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Printf("Complete Host URL: %s", r.Host)
	host := strings.Split(r.Host, ".")
	log.Printf("Host: %+v", host)
	if len(host) > 1 {
		log.Printf("Subdomain: %s", host[0])
		// get the LB DNS from the database
		// and update the target URL
		dns, err := GetLB_DNS(host[0])
		if err != nil {
			log.Printf("Error getting LB DNS: %v", err)
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}

		p.target, err = url.Parse("http://" + dns)
		if err != nil {
			log.Printf("Error parsing DNS: %v", err)
			http.Error(w, "Invalid service URL", http.StatusInternalServerError)
			return
		}
		log.Printf("Updated target URL: %s", p.target.String())
		p.proxy = httputil.NewSingleHostReverseProxy(p.target)
		p.proxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

	} else {
		log.Printf("No subdomain found")

		// Update the target URL
		var err error
		p.target, err = url.Parse("http://localhost:2001/")
		if err != nil {
			log.Printf("Error parsing target URL: %v", err)
			http.Error(w, "Invalid target URL", http.StatusInternalServerError)
			return
		}
		log.Printf("Updated target URL: %s", p.target.String())
	}

	r.URL.Scheme = p.target.Scheme

	r.URL.Host = p.target.Host

	// Log the updated URL
	log.Printf("Updated URL: %s", r.URL.String())

	p.proxy.ServeHTTP(w, r)
}

func main() {

	// URL of the server to proxy to
	target, err := url.Parse("http://localhost:2001/")
	if err != nil {
		panic(err)
	}

	// Create a new ReverseProxy instance
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Configure the reverse proxy to use HTTPS and HTTPS
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Create a new Proxy instance
	p := &Proxy{target: target, proxy: proxy}

	// Start the HTTP server and register the Proxy instance as the handler
	err = http.ListenAndServe(":3001", p)
	if err != nil {
		panic(err)
	}
}

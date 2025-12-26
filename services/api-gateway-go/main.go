package main

import (
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "strings"
)

const (
    userServiceURL    = "http://user-service-python:8001"
    catalogServiceURL = "http://catalog-service-go:8002"
    cartServiceURL    = "http://cart-service-python:8003"
    orderServiceURL   = "http://order-service-go:8005"
)

func newReverseProxy(target string) *httputil.ReverseProxy {
    targetURL, err := url.Parse(target)
    if err != nil {
        log.Fatalf("invalid target URL %s: %v", target, err)
    }
    return httputil.NewSingleHostReverseProxy(targetURL)
}

func loginHandler(proxy *httputil.ReverseProxy) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // /api/login -> user-service /login
        r.URL.Path = "/login"
        proxy.ServeHTTP(w, r)
    }
}

func productsHandler(proxy *httputil.ReverseProxy) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // /api/products -> catalog-service /products
        r.URL.Path = "/products"
        proxy.ServeHTTP(w, r)
    }
}

func cartHandler(proxy *httputil.ReverseProxy) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // /api/cart/* -> cart-service /cart/*
        // Strip the /api prefix
        if strings.HasPrefix(r.URL.Path, "/api") {
            r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
        }
        proxy.ServeHTTP(w, r)
    }
}

func orderHandler(proxy *httputil.ReverseProxy) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // /api/order -> order-service /order/checkout
        r.URL.Path = "/order/checkout"
        proxy.ServeHTTP(w, r)
    }
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"status":"ok"}`))
}

func main() {
    userProxy := newReverseProxy(userServiceURL)
    catalogProxy := newReverseProxy(catalogServiceURL)
    cartProxy := newReverseProxy(cartServiceURL)
    orderProxy := newReverseProxy(orderServiceURL)

    mux := http.NewServeMux()
    mux.HandleFunc("/api/login", loginHandler(userProxy))
    mux.HandleFunc("/api/products", productsHandler(catalogProxy))
    mux.HandleFunc("/api/cart/", cartHandler(cartProxy))
    mux.HandleFunc("/api/order", orderHandler(orderProxy))
    mux.HandleFunc("/health", healthHandler)

    addr := ":8000"
    log.Println("api-gateway-go listening on", addr)
    if err := http.ListenAndServe(addr, mux); err != nil {
        log.Fatal(err)
    }
}

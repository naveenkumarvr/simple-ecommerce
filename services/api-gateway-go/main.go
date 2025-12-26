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

// orderCheckoutHandler handles POST /api/order -> /order/checkout
func orderCheckoutHandler(proxy *httputil.ReverseProxy) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // /api/order -> order-service /order/checkout
        r.URL.Path = "/order/checkout"
        proxy.ServeHTTP(w, r)
    }
}

// orderGetHandler handles GET /api/order/{user_id} -> /order/{user_id}
func orderGetHandler(proxy *httputil.ReverseProxy) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Expect path like /api/order/{user_id}
        if strings.HasPrefix(r.URL.Path, "/api") {
            r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
        }
        // Now path is /order/{user_id}, which order-service understands
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

func logRequest(r *http.Request) {
    log.Printf("Received request: %s %s", r.Method, r.URL.Path)
}

// withCORS wraps a handler and adds very simple CORS headers so a static frontend
// loaded from a different origin (e.g. file:// or another port) can call the API.
func withCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func main() {
    userProxy := newReverseProxy(userServiceURL)
    catalogProxy := newReverseProxy(catalogServiceURL)
    cartProxy := newReverseProxy(cartServiceURL)
    orderProxy := newReverseProxy(orderServiceURL)

    mux := http.NewServeMux()
    mux.HandleFunc("/api/login", loginHandler(userProxy))
    mux.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
        logRequest(r)
        productsHandler(catalogProxy)(w, r)
    })
    mux.HandleFunc("/api/cart/", cartHandler(cartProxy))

    // Order routes:
    // POST /api/order -> /order/checkout
    // GET  /api/order/{user_id} -> /order/{user_id}
    mux.HandleFunc("/api/order", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            orderCheckoutHandler(orderProxy)(w, r)
            return
        }
        if r.Method == http.MethodGet {
            // Optional: could return 405 or some summary; for now, 405
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }
        w.WriteHeader(http.StatusMethodNotAllowed)
    })
    mux.HandleFunc("/api/order/", orderGetHandler(orderProxy))

    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        logRequest(r)
        healthHandler(w, r)
    })

    addr := ":8000"
    log.Println("api-gateway-go listening on", addr)
    if err := http.ListenAndServe(addr, withCORS(mux)); err != nil {
        log.Fatal(err)
    }
}

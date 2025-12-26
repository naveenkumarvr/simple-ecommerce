package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strings"
)

type Product struct {
    ID    string  `json:"id"`
    Price float64 `json:"price"`
}

var products = map[string]Product{
    "car":   {ID: "car", Price: 100},
    "bike":  {ID: "bike", Price: 40},
    "bus":   {ID: "bus", Price: 150},
    "truck": {ID: "truck", Price: 200},
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(data); err != nil {
        log.Println("error encoding response:", err)
    }
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    list := make([]Product, 0, len(products))
    for _, p := range products {
        list = append(list, p)
    }
    writeJSON(w, http.StatusOK, list)
}

func productByIDHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    // Expected path: /products/{id}
    path := strings.TrimPrefix(r.URL.Path, "/products/")
    if path == "" || strings.Contains(path, "/") {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    id := path
    p, ok := products[id]
    if !ok {
        http.NotFound(w, r)
        return
    }

    writeJSON(w, http.StatusOK, p)
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/products", productsHandler)
    mux.HandleFunc("/products/", productByIDHandler)

    addr := ":8002"
    log.Println("catalog-service-go listening on", addr)
    if err := http.ListenAndServe(addr, mux); err != nil {
        log.Fatal(err)
    }
}

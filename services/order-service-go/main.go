package main

import (
    "bytes"
    "encoding/json"
    "io"
    "log"
    "math/rand"
    "net/http"
    "strconv"
    "strings"
    "sync"
    "time"
)

const (
    cartServiceURL    = "http://cart-service-python:8003"    // Kubernetes service name or docker-compose service name
    paymentServiceURL = "http://payment-service-python:8004" // Kubernetes service name or docker-compose service name
)

type checkoutRequest struct {
    UserID string `json:"user_id"`
}

type cartItem struct {
    ProductID string `json:"product_id"`
    Quantity  int    `json:"quantity"`
}

type cartResponse struct {
    UserID string     `json:"user_id"`
    Items  []cartItem `json:"items"`
}

type paymentRequest struct {
    UserID   string  `json:"user_id"`
    Amount   float64 `json:"amount"`
    Currency string  `json:"currency"`
    Source   string  `json:"source"`
}

type paymentResponse struct {
    PaymentID string `json:"payment_id"`
    Status    string `json:"status"`
}

// Order represents a very simple in-memory order record.
type Order struct {
    OrderID   string     `json:"order_id"`
    UserID    string     `json:"user_id"`
    Items     []cartItem `json:"items"`
    Total     float64    `json:"total"`
    PaymentID string     `json:"payment_id"`
    Status    string     `json:"status"`
    CreatedAt time.Time  `json:"created_at"`
}

// In-memory store: user_id -> list of orders.
var (
    ordersMu    sync.Mutex
    ordersByUser = make(map[string][]Order)
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(data); err != nil {
        log.Println("error encoding response:", err)
    }
}

func fetchCart(userID string) (*cartResponse, error) {
    url := cartServiceURL + "/cart/" + userID
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        log.Printf("cart-service returned %d: %s", resp.StatusCode, string(body))
        return nil, err
    }

    var cr cartResponse
    if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
        return nil, err
    }
    return &cr, nil
}

func callPayment(userID string, amount float64) (*paymentResponse, error) {
    reqBody := paymentRequest{
        UserID:   userID,
        Amount:   amount,
        Currency: "USD",
        Source:   "demo-source",
    }
    buf, err := json.Marshal(reqBody)
    if err != nil {
        return nil, err
    }

    resp, err := http.Post(paymentServiceURL+"/pay", "application/json", bytes.NewReader(buf))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        log.Printf("payment-service returned %d: %s", resp.StatusCode, string(body))
        return nil, err
    }

    var pr paymentResponse
    if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
        return nil, err
    }
    return &pr, nil
}

// generateOrderID creates a simple pseudo-unique order ID for demo purposes.
func generateOrderID() string {
    return "order-" + strconv.FormatInt(time.Now().UnixNano(), 10) + "-" + strconv.Itoa(rand.Intn(1000))
}

// storeOrder appends an order for the given user into the in-memory store.
func storeOrder(o Order) {
    ordersMu.Lock()
    defer ordersMu.Unlock()
    ordersByUser[o.UserID] = append(ordersByUser[o.UserID], o)
}

// getOrdersForUser returns all orders for a given user.
func getOrdersForUser(userID string) []Order {
    ordersMu.Lock()
    defer ordersMu.Unlock()
    return append([]Order(nil), ordersByUser[userID]...)
}

func checkoutHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    var req checkoutRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body; user_id required"})
        return
    }

    // 1) Get cart for the user
    cart, err := fetchCart(req.UserID)
    if err != nil {
        log.Println("error calling cart-service:", err)
        writeJSON(w, http.StatusBadGateway, map[string]string{"error": "failed to fetch cart"})
        return
    }

    // 2) Compute a very naive total amount: just sum quantity (pretend each item costs 1)
    var total float64
    for _, item := range cart.Items {
        total += float64(item.Quantity)
    }

    if total == 0 {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "cart is empty"})
        return
    }

    // 3) Call payment-service
    payment, err := callPayment(req.UserID, total)
    if err != nil {
        log.Println("error calling payment-service:", err)
        writeJSON(w, http.StatusBadGateway, map[string]string{"error": "payment failed"})
        return
    }

    if payment.Status != "success" {
        writeJSON(w, http.StatusBadGateway, map[string]string{"error": "payment not successful"})
        return
    }

    // 4) Build and store order record in memory
    order := Order{
        OrderID:   generateOrderID(),
        UserID:    req.UserID,
        Items:     cart.Items,
        Total:     total,
        PaymentID: payment.PaymentID,
        Status:    "completed",
        CreatedAt: time.Now().UTC(),
    }
    storeOrder(order)

    // 5) Return success response including order details
    writeJSON(w, http.StatusOK, map[string]interface{}{
        "message":    "order successful",
        "payment_id": payment.PaymentID,
        "order_id":   order.OrderID,
        "total":      order.Total,
    })
}

// getOrdersHandler returns all orders for a given user.
// URL pattern: /order/{user_id}
func getOrdersHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    // Expect path like /order/{user_id}
    path := r.URL.Path
    if !strings.HasPrefix(path, "/order/") {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid path"})
        return
    }
    userID := strings.TrimPrefix(path, "/order/")
    if userID == "" {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user_id required"})
        return
    }

    orders := getOrdersForUser(userID)
    writeJSON(w, http.StatusOK, orders)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
    writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/order/checkout", checkoutHandler)
    mux.HandleFunc("/order/", getOrdersHandler) // GET /order/{user_id}
    mux.HandleFunc("/health", healthHandler)

    addr := ":8005"
    log.Println("order-service-go listening on", addr)
    if err := http.ListenAndServe(addr, mux); err != nil {
        log.Fatal(err)
    }
}

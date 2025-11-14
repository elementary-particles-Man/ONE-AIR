package main
import (
    "log"
    "net/http"
)
func main() {
    http.Handle("/", http.FileServer(http.Dir("./frontend")))
    log.Println("ONE-AIR server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

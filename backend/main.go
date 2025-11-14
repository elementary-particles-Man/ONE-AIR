package main

import (
    "database/sql"
    "log"
    "net/http"
    "os"

    _ "modernc.org/sqlite"
)

func main() {
    db, err := sql.Open("sqlite", "./db/airgate.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    initDB(db)

    fs := http.FileServer(http.Dir("../frontend"))
    http.Handle("/", fs)

    log.Println("ONE-AIR running at :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func initDB(db *sql.DB) {
    content, err := os.ReadFile("./db/init.sql")
    if err != nil {
        log.Fatal(err)
    }
    if _, err := db.Exec(string(content)); err != nil {
        log.Fatal(err)
    }
}

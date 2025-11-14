package main

import (
    "database/sql"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "os"
    "strings"

    "oneair/internal"

    _ "modernc.org/sqlite"
)

func main() {
    db, err := sql.Open("sqlite", "./db/airgate.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    initDB(db)

    ip := internal.LocalIP()
    port := internal.RandomPort()
    code := internal.ShortCode()
    redirectURL := fmt.Sprintf("http://%s:%d/", ip, port)
    redirectMap := map[string]string{
        code: redirectURL,
    }

    shutdownMDNS := internal.PublishMDNS(port)
    defer shutdownMDNS()

    backendMux := http.NewServeMux()
    backendMux.Handle("/", http.FileServer(http.Dir("../frontend")))

    go func() {
        log.Printf("ONE-AIR backend running: %s", redirectURL)
        if err := http.ListenAndServe(fmt.Sprintf(":%d", port), backendMux); err != nil {
            log.Fatal(err)
        }
    }()

    launcherMux := http.NewServeMux()
    launcherMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        key := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/"))
        if dest, ok := redirectMap[key]; ok && key != "" {
            http.Redirect(w, r, dest, http.StatusFound)
            return
        }
        renderLaunch(w, ip, port, code)
    })

    log.Printf("ONE-AIR launcher ready: http://%s/ (smartphone: http://oneair.local/%s)", ip, code)
    log.Fatal(http.ListenAndServe(":80", launcherMux))
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

func renderLaunch(w http.ResponseWriter, ip string, port int, code string) {
    launch := `<!DOCTYPE html><html><body>
    <h1>ONE-AIR 起動完了</h1>
    <p>スマホで以下にアクセスしてください：</p>
    <pre>http://oneair.local/{{.Code}}</pre>
    <p>ローカルPC直アクセス：</p>
    <pre>http://{{.IP}}:{{.Port}}/</pre>
    </body></html>`
    t := template.Must(template.New("launch").Parse(launch))
    if err := t.Execute(w, map[string]any{
        "IP":   ip,
        "Port": port,
        "Code": code,
    }); err != nil {
        log.Println("render error:", err)
    }
}

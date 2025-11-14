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

const (
	launcherHost = "127.0.0.1"
	launcherPort = 8800
	mdnsHost     = "oneair.local"
)

func main() {
	db, err := sql.Open("sqlite", "./db/airgate.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	initDB(db)

	localIP := internal.LocalIP()
	clientPort := internal.RandomPort()
	code := internal.ShortCode()

	redirectURL := fmt.Sprintf("http://%s:%d/", localIP, clientPort)
	redirectMap := map[string]string{
		code: redirectURL,
	}

	shutdownMDNS := internal.PublishMDNS(mdnsHost, launcherPort, localIP)
	defer shutdownMDNS()

	backendMux := http.NewServeMux()
	backendMux.Handle("/", http.FileServer(http.Dir("../frontend")))

	go func() {
		log.Printf("ONE-AIR backend running: %s", redirectURL)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", clientPort), backendMux); err != nil {
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
		renderLaunch(w, launchTemplateData{
			LocalIP:      localIP,
			ClientPort:   clientPort,
			Code:         code,
			LauncherPort: launcherPort,
			Hostname:     mdnsHost,
		})
	})

	log.Printf("ONE-AIR launcher ready: http://%s:%d/ (smartphone: http://%s:%d/%s)", launcherHost, launcherPort, mdnsHost, launcherPort, code)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", launcherPort), launcherMux))
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

type launchTemplateData struct {
	LocalIP      string
	ClientPort   int
	Code         string
	LauncherPort int
	Hostname     string
}

func renderLaunch(w http.ResponseWriter, data launchTemplateData) {
	launch := `<!DOCTYPE html><html><body>
    <h1>ONE-AIR 起動完了</h1>
    <p>スマホで以下にアクセスしてください：</p>
    <pre>http://{{.Hostname}}:{{.LauncherPort}}/{{.Code}}</pre>
    <p>PCからランチャーを開く：</p>
    <pre>http://127.0.0.1:{{.LauncherPort}}/</pre>
    <p>ローカルPC直アクセス：</p>
    <pre>http://{{.LocalIP}}:{{.ClientPort}}/</pre>
    </body></html>`
	t := template.Must(template.New("launch").Parse(launch))
	if err := t.Execute(w, data); err != nil {
		log.Println("render error:", err)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

// Build-time variables injected during compilation
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

// StatusInfo holds the application status information
type StatusInfo struct {
	Version     string    `json:"version"`
	GitCommit   string    `json:"git_commit"`
	BuildTime   string    `json:"build_time"`
	Environment string    `json:"environment"`
	PodName     string    `json:"pod_name"`
	Namespace   string    `json:"namespace"`
	Status      string    `json:"status"`
	Uptime      string    `json:"uptime"`
	Timestamp   time.Time `json:"timestamp"`
}

var startTime = time.Now()

// getStatusInfo retrieves current application status
func getStatusInfo() StatusInfo {
	uptime := time.Since(startTime).Round(time.Second).String()

	return StatusInfo{
		Version:     Version,
		GitCommit:   GitCommit,
		BuildTime:   BuildTime,
		Environment: getEnv("ENVIRONMENT", "development"),
		PodName:     getEnv("POD_NAME", "localhost"),
		Namespace:   getEnv("POD_NAMESPACE", "default"),
		Status:      "healthy",
		Uptime:      uptime,
		Timestamp:   time.Now(),
	}
}

// getEnv retrieves environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// healthHandler handles /health endpoint for K8s probes
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// statusAPIHandler handles /api/status endpoint (JSON)
func statusAPIHandler(w http.ResponseWriter, r *http.Request) {
	status := getStatusInfo()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// metricsHandler handles /metrics endpoint for Prometheus
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(startTime).Seconds()

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "# HELP devops_demo_uptime_seconds Application uptime in seconds\n")
	fmt.Fprintf(w, "# TYPE devops_demo_uptime_seconds counter\n")
	fmt.Fprintf(w, "devops_demo_uptime_seconds %.0f\n", uptime)
	fmt.Fprintf(w, "# HELP devops_demo_info Application info\n")
	fmt.Fprintf(w, "# TYPE devops_demo_info gauge\n")
	fmt.Fprintf(w, "devops_demo_info{version=\"%s\",git_commit=\"%s\",environment=\"%s\"} 1\n",
		Version, GitCommit, getEnv("ENVIRONMENT", "development"))
}

// rootHandler handles / endpoint (HTML UI)
func rootHandler(w http.ResponseWriter, r *http.Request) {
	status := getStatusInfo()

	tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>DevOps Demo App</title>
    <style>
        body {
            font-family: 'Courier New', monospace;
            background: #1e1e1e;
            color: #00ff00;
            padding: 40px;
            line-height: 1.6;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
            border: 2px solid #00ff00;
            padding: 30px;
            border-radius: 5px;
        }
        h1 {
            text-align: center;
            border-bottom: 2px solid #00ff00;
            padding-bottom: 10px;
        }
        .status-grid {
            display: grid;
            grid-template-columns: 200px 1fr;
            gap: 15px;
            margin: 20px 0;
        }
        .label {
            font-weight: bold;
            color: #00aaff;
        }
        .value {
            color: #00ff00;
        }
        .healthy {
            color: #00ff00;
            font-weight: bold;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #00ff00;
            color: #888;
            font-size: 0.9em;
        }
        a {
            color: #00aaff;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 DevOps Demo Status Service</h1>

        <div class="status-grid">
            <div class="label">Version:</div>
            <div class="value">{{.Version}}</div>

            <div class="label">Git Commit:</div>
            <div class="value">{{.GitCommit}}</div>

            <div class="label">Built:</div>
            <div class="value">{{.BuildTime}}</div>

            <div class="label">Environment:</div>
            <div class="value">{{.Environment}}</div>

            <div class="label">Pod Name:</div>
            <div class="value">{{.PodName}}</div>

            <div class="label">Namespace:</div>
            <div class="value">{{.Namespace}}</div>

            <div class="label">Status:</div>
            <div class="value healthy">✓ {{.Status}}</div>

            <div class="label">Uptime:</div>
            <div class="value">{{.Uptime}}</div>

            <div class="label">Timestamp:</div>
            <div class="value">{{.Timestamp.Format "2006-01-02 15:04:05 MST"}}</div>
        </div>

        <div class="footer">
            <p>API Endpoints:
                <a href="/api/status">/api/status</a> |
                <a href="/health">/health</a> |
                <a href="/metrics">/metrics</a>
            </p>
            <p>CI/CD + GitOps Demo Project | Built with ❤️ in Go</p>
        </div>
    </div>
</body>
</html>`

	t, err := template.New("status").Parse(tmpl)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, status)
}

func main() {
	port := getEnv("PORT", "8080")

	// Register handlers
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/status", statusAPIHandler)
	http.HandleFunc("/metrics", metricsHandler)

	log.Printf("🚀 DevOps Demo App starting...")
	log.Printf("Version: %s", Version)
	log.Printf("Git Commit: %s", GitCommit)
	log.Printf("Build Time: %s", BuildTime)
	log.Printf("Listening on port %s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

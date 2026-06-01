package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"
)

// MonitorConfig holds configuration for the monitoring middleware.
type MonitorConfig struct {
	// URI is the endpoint path for the monitor dashboard.
	URI string
	// Title is the title shown on the dashboard.
	Title string
	// RefreshInterval is how often the dashboard refreshes (seconds).
	RefreshInterval int
}

// DefaultMonitorConfig returns default monitoring configuration.
func DefaultMonitorConfig() MonitorConfig {
	return MonitorConfig{
		URI:             "/monitor",
		Title:           "GoW Monitor",
		RefreshInterval: 3,
	}
}

// monitorStats holds real-time server statistics.
type monitorStats struct {
	CPUUsage        float64 `json:"cpu_usage"`
	MemoryUsage     uint64  `json:"memory_usage"`
	MemoryAlloc     uint64  `json:"memory_alloc"`
	MemorySys       uint64  `json:"memory_sys"`
	GoRoutines      int     `json:"go_routines"`
	OpenConnections int     `json:"open_connections"`
	ResponseTime    float64 `json:"response_time"`
	TotalRequests   int64   `json:"total_requests"`
	Uptime          string  `json:"uptime"`
	Timestamp       int64   `json:"timestamp"`
}

var (
	monitorStartTime = time.Now()
	monitorTotalReqs int64
	monitorConnCount int64
	monitorMu        sync.RWMutex
)

// MonitorMiddleware returns middleware that serves a monitoring dashboard.
func MonitorMiddleware(config ...MonitorConfig) func(http.Handler) http.Handler {
	cfg := DefaultMonitorConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == cfg.URI {
				if r.URL.Query().Get("format") == "json" {
					serveMonitorJSON(w, r)
					return
				}
				serveMonitorDashboard(w, r, cfg)
				return
			}

			// Track requests
			monitorMu.Lock()
			monitorTotalReqs++
			monitorConnCount++
			monitorMu.Unlock()

			start := time.Now()
			next.ServeHTTP(w, r)

			monitorMu.Lock()
			monitorConnCount--
			monitorMu.Unlock()

			_ = time.Since(start)
		})
	}
}

func getMonitorStats() monitorStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	monitorMu.RLock()
	totalReqs := monitorTotalReqs
	connCount := monitorConnCount
	monitorMu.RUnlock()

	uptime := time.Since(monitorStartTime).Round(time.Second)

	return monitorStats{
		MemoryUsage:     m.Sys,
		MemoryAlloc:     m.Alloc,
		MemorySys:       m.Sys,
		GoRoutines:      runtime.NumGoroutine(),
		OpenConnections: int(connCount),
		TotalRequests:   totalReqs,
		Uptime:          uptime.String(),
		Timestamp:       time.Now().Unix(),
	}
}

func serveMonitorJSON(w http.ResponseWriter, r *http.Request) {
	stats := getMonitorStats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func serveMonitorDashboard(w http.ResponseWriter, r *http.Request, cfg MonitorConfig) {
	stats := getMonitorStats()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <style>
        .chart-bar { transition: height 0.3s ease; }
        .stat-card { backdrop-filter: blur(10px); }
    </style>
</head>
<body class="bg-zinc-950 text-white min-h-screen antialiased">
    <nav class="bg-zinc-900 border-b border-zinc-800 px-6 py-4">
        <div class="max-w-6xl mx-auto flex items-center justify-between">
            <h1 class="text-xl font-bold text-emerald-400">%s</h1>
            <span class="text-sm text-zinc-500">Auto-refresh: %ds</span>
        </div>
    </nav>

    <div class="max-w-6xl mx-auto p-6">
        <!-- Stats Grid -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
            <div class="stat-card bg-zinc-900/50 border border-zinc-800 rounded-xl p-4">
                <div class="text-xs text-zinc-500 mb-1">CPU Cores</div>
                <div class="text-2xl font-bold text-emerald-400">%d</div>
            </div>
            <div class="stat-card bg-zinc-900/50 border border-zinc-800 rounded-xl p-4">
                <div class="text-xs text-zinc-500 mb-1">Memory Usage</div>
                <div class="text-2xl font-bold text-blue-400">%.1f MB</div>
            </div>
            <div class="stat-card bg-zinc-900/50 border border-zinc-800 rounded-xl p-4">
                <div class="text-xs text-zinc-500 mb-1">Go Routines</div>
                <div class="text-2xl font-bold text-purple-400">%d</div>
            </div>
            <div class="stat-card bg-zinc-900/50 border border-zinc-800 rounded-xl p-4">
                <div class="text-xs text-zinc-500 mb-1">Open Connections</div>
                <div class="text-2xl font-bold text-yellow-400">%d</div>
            </div>
        </div>

        <!-- Metrics -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Memory Chart -->
            <div class="bg-zinc-900/50 border border-zinc-800 rounded-xl p-6">
                <h3 class="text-sm font-medium text-zinc-400 mb-4">Memory Usage</h3>
                <div class="flex items-end gap-1 h-32">
                    <div class="chart-bar bg-emerald-500 rounded-t" style="width: 100%%; height: %.0f%%;"></div>
                </div>
                <div class="flex justify-between mt-2 text-xs text-zinc-500">
                    <span>Alloc: %.1f MB</span>
                    <span>Sys: %.1f MB</span>
                </div>
            </div>

            <!-- Response Time -->
            <div class="bg-zinc-900/50 border border-zinc-800 rounded-xl p-6">
                <h3 class="text-sm font-medium text-zinc-400 mb-4">Response Time</h3>
                <div class="flex items-center justify-center h-32">
                    <div class="text-center">
                        <div class="text-4xl font-bold text-emerald-400">0.0</div>
                        <div class="text-sm text-zinc-500">ms avg</div>
                    </div>
                </div>
            </div>

            <!-- Total Requests -->
            <div class="bg-zinc-900/50 border border-zinc-800 rounded-xl p-6">
                <h3 class="text-sm font-medium text-zinc-400 mb-4">Total Requests</h3>
                <div class="flex items-center justify-center h-32">
                    <div class="text-center">
                        <div class="text-4xl font-bold text-blue-400">%d</div>
                        <div class="text-sm text-zinc-500">since server start</div>
                    </div>
                </div>
            </div>

            <!-- Uptime -->
            <div class="bg-zinc-900/50 border border-zinc-800 rounded-xl p-6">
                <h3 class="text-sm font-medium text-zinc-400 mb-4">Uptime</h3>
                <div class="flex items-center justify-center h-32">
                    <div class="text-center">
                        <div class="text-4xl font-bold text-purple-400">%s</div>
                        <div class="text-sm text-zinc-500">server uptime</div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Server Info -->
        <div class="mt-6 bg-zinc-900/50 border border-zinc-800 rounded-xl p-6">
            <h3 class="text-sm font-medium text-zinc-400 mb-4">Server Information</h3>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                <div>
                    <span class="text-zinc-500">Framework:</span>
                    <span class="text-white ml-2">GoW</span>
                </div>
                <div>
                    <span class="text-zinc-500">Go Version:</span>
                    <span class="text-white ml-2">%s</span>
                </div>
                <div>
                    <span class="text-zinc-500">OS/Arch:</span>
                    <span class="text-white ml-2">%s/%s</span>
                </div>
                <div>
                    <span class="text-zinc-500">Last Updated:</span>
                    <span class="text-white ml-2">%s</span>
                </div>
            </div>
        </div>
    </div>

    <script>
        setTimeout(function() { location.reload(); }, %d000);
    </script>
</body>
</html>`,
		cfg.Title,
		cfg.Title,
		cfg.RefreshInterval,
		runtime.NumCPU(),
		float64(stats.MemoryAlloc)/1024/1024,
		stats.GoRoutines,
		stats.OpenConnections,
		float64(stats.MemoryAlloc)*100/float64(stats.MemorySys),
		float64(stats.MemoryAlloc)/1024/1024,
		float64(stats.MemorySys)/1024/1024,
		stats.TotalRequests,
		stats.Uptime,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
		time.Now().Format("15:04:05"),
		cfg.RefreshInterval,
	)
}

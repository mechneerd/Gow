package pulse

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"
)

// Metric represents a single metric data point
type Metric struct {
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
	Unit      string    `json:"unit"`
	Timestamp time.Time `json:"timestamp"`
}

// MetricResult holds aggregated metric results
type MetricResult struct {
	Name    string   `json:"name"`
	Values  []Metric `json:"values"`
	Average float64  `json:"average"`
	Min     float64  `json:"min"`
	Max     float64  `json:"max"`
}

// MetricProvider is the interface for collecting metrics
type MetricProvider interface {
	Name() string
	Collect(ctx context.Context) (float64, string, error)
}

// Pulse provides application health monitoring
type Pulse struct {
	mu        sync.RWMutex
	providers []MetricProvider
	metrics   map[string][]Metric
	maxPoints int
	router    *http.ServeMux
}

// New creates a new Pulse monitoring instance
func New() *Pulse {
	p := &Pulse{
		providers: make([]MetricProvider, 0),
		metrics:   make(map[string][]Metric),
		maxPoints: 60, // Keep 60 data points
		router:    http.NewServeMux(),
	}
	p.routes()

	// Register default providers
	p.RegisterProvider(&MemoryProvider{})
	p.RegisterProvider(&GoroutineProvider{})
	p.RegisterProvider(&UptimeProvider{})
	p.RegisterProvider(&GCProvider{})

	return p
}

// RegisterProvider adds a metric provider
func (p *Pulse) RegisterProvider(provider MetricProvider) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.providers = append(p.providers, provider)
}

// Handler returns the HTTP handler
func (p *Pulse) Handler() http.Handler {
	return p.router
}

func (p *Pulse) routes() {
	p.router.HandleFunc("/", p.indexHandler)
	p.router.HandleFunc("/api/metrics", p.metricsHandler)
	p.router.HandleFunc("/api/providers", p.providersHandler)
	p.router.HandleFunc("/api/collect", p.collectHandler)
}

func (p *Pulse) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, pulseHTML)
}

func (p *Pulse) metricsHandler(w http.ResponseWriter, r *http.Request) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	results := make(map[string]*MetricResult)
	for name, values := range p.metrics {
		if len(values) == 0 {
			continue
		}
		result := &MetricResult{
			Name:   name,
			Values: values,
		}

		var sum, min, max float64
		min = values[0].Value
		max = values[0].Value
		for _, v := range values {
			sum += v.Value
			if v.Value < min {
				min = v.Value
			}
			if v.Value > max {
				max = v.Value
			}
		}
		result.Average = sum / float64(len(values))
		result.Min = min
		result.Max = max
		results[name] = result
	}

	json.NewEncoder(w).Encode(results)
}

func (p *Pulse) providersHandler(w http.ResponseWriter, r *http.Request) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var names []string
	for _, provider := range p.providers {
		names = append(names, provider.Name())
	}
	json.NewEncoder(w).Encode(names)
}

func (p *Pulse) collectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	p.CollectAll()
	json.NewEncoder(w).Encode(map[string]string{"status": "collected"})
}

// CollectAll runs all metric providers and stores results
func (p *Pulse) CollectAll() {
	ctx := context.Background()
	for _, provider := range p.providers {
		value, unit, err := provider.Collect(ctx)
		if err != nil {
			continue
		}
		p.addMetric(provider.Name(), value, unit)
	}
}

func (p *Pulse) addMetric(name string, value float64, unit string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	metric := Metric{
		Name:      name,
		Value:     value,
		Unit:      unit,
		Timestamp: time.Now(),
	}

	p.metrics[name] = append(p.metrics[name], metric)
	if len(p.metrics[name]) > p.maxPoints {
		p.metrics[name] = p.metrics[name][1:]
	}
}

// --- Default Metric Providers ---

// MemoryProvider collects memory usage metrics
type MemoryProvider struct{}

func (m *MemoryProvider) Name() string {
	return "memory"
}

func (m *MemoryProvider) Collect(ctx context.Context) (float64, string, error) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return float64(mem.Alloc) / 1024 / 1024, "MB", nil
}

// GoroutineProvider collects goroutine count
type GoroutineProvider struct{}

func (g *GoroutineProvider) Name() string {
	return "goroutines"
}

func (g *GoroutineProvider) Collect(ctx context.Context) (float64, string, error) {
	return float64(runtime.NumGoroutine()), "", nil
}

// UptimeProvider tracks application uptime
type UptimeProvider struct {
	startTime time.Time
}

func (u *UptimeProvider) Name() string {
	return "uptime"
}

func (u *UptimeProvider) Collect(ctx context.Context) (float64, string, error) {
	if u.startTime.IsZero() {
		u.startTime = time.Now()
	}
	uptime := time.Since(u.startTime).Seconds()
	return uptime, "seconds", nil
}

// GCProvider collects garbage collection metrics
type GCProvider struct{}

func (g *GCProvider) Name() string {
	return "gc_cycles"
}

func (g *GCProvider) Collect(ctx context.Context) (float64, string, error) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return float64(mem.NumGC), "", nil
}

// CustomMetricProvider allows creating custom metrics
type CustomMetricProvider struct {
	name     string
	collectFn func(ctx context.Context) (float64, string, error)
}

// NewCustomProvider creates a custom metric provider
func NewCustomProvider(name string, fn func(ctx context.Context) (float64, string, error)) *CustomMetricProvider {
	return &CustomMetricProvider{
		name:      name,
		collectFn: fn,
	}
}

func (c *CustomMetricProvider) Name() string {
	return c.name
}

func (c *CustomMetricProvider) Collect(ctx context.Context) (float64, string, error) {
	return c.collectFn(ctx)
}

// Template for Pulse dashboard
const pulseHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GoW Pulse - Health Monitoring</title>
    <style>
        :root {
            --primary: #6366f1;
            --success: #22c55e;
            --warning: #f59e0b;
            --danger: #ef4444;
            --bg: #0f172a;
            --surface: #1e293b;
            --text: #f8fafc;
            --muted: #94a3b8;
        }
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: var(--bg);
            color: var(--text);
            line-height: 1.6;
        }
        .container { max-width: 1200px; margin: 0 auto; padding: 2rem; }
        header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 2rem;
            padding-bottom: 1rem;
            border-bottom: 1px solid rgba(255,255,255,0.1);
        }
        h1 { font-size: 1.5rem; font-weight: 600; }
        .subtitle { color: var(--muted); font-size: 0.875rem; }
        .status-badge {
            display: inline-flex;
            align-items: center;
            gap: 0.5rem;
            padding: 0.5rem 1rem;
            background: rgba(34,197,94,0.2);
            color: var(--success);
            border-radius: 9999px;
            font-weight: 500;
        }
        .status-dot {
            width: 8px;
            height: 8px;
            background: var(--success);
            border-radius: 50%;
            animation: pulse 2s infinite;
        }
        @keyframes pulse {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.5; }
        }
        .metrics-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 1.5rem;
            margin-bottom: 2rem;
        }
        .metric-card {
            background: var(--surface);
            border-radius: 0.75rem;
            border: 1px solid rgba(255,255,255,0.1);
            padding: 1.5rem;
        }
        .metric-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 1rem;
        }
        .metric-title {
            font-size: 0.875rem;
            color: var(--muted);
            text-transform: uppercase;
            letter-spacing: 0.05em;
        }
        .metric-value {
            font-size: 2.5rem;
            font-weight: 700;
        }
        .metric-unit {
            font-size: 1rem;
            color: var(--muted);
            margin-left: 0.5rem;
        }
        .metric-chart {
            height: 60px;
            display: flex;
            align-items: flex-end;
            gap: 2px;
            margin-top: 1rem;
        }
        .chart-bar {
            flex: 1;
            background: var(--primary);
            border-radius: 2px 2px 0 0;
            min-height: 4px;
            transition: height 0.3s;
        }
        .metric-stats {
            display: flex;
            gap: 1.5rem;
            margin-top: 1rem;
            padding-top: 1rem;
            border-top: 1px solid rgba(255,255,255,0.1);
        }
        .stat {
            font-size: 0.75rem;
        }
        .stat-label { color: var(--muted); }
        .stat-value { color: var(--text); font-weight: 600; }
        .panel {
            background: var(--surface);
            border-radius: 0.75rem;
            border: 1px solid rgba(255,255,255,0.1);
            margin-bottom: 1.5rem;
        }
        .panel-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 1rem 1.5rem;
            border-bottom: 1px solid rgba(255,255,255,0.1);
        }
        .panel-title { font-weight: 600; }
        .panel-body { padding: 1.5rem; }
        .btn {
            padding: 0.5rem 1rem;
            border-radius: 0.5rem;
            cursor: pointer;
            border: none;
            font-size: 0.875rem;
            font-weight: 500;
            background: var(--primary);
            color: white;
            transition: opacity 0.2s;
        }
        .btn:hover { opacity: 0.8; }
        .provider-list {
            display: flex;
            flex-wrap: wrap;
            gap: 0.5rem;
        }
        .provider-tag {
            padding: 0.25rem 0.75rem;
            background: rgba(99,102,241,0.2);
            color: var(--primary);
            border-radius: 9999px;
            font-size: 0.75rem;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <div>
                <h1>GoW Pulse</h1>
                <div class="subtitle">Health Monitoring Dashboard</div>
            </div>
            <div style="display: flex; gap: 1rem; align-items: center;">
                <span class="status-badge">
                    <span class="status-dot"></span>
                    Healthy
                </span>
                <button class="btn" onclick="collect()">Collect Metrics</button>
            </div>
        </header>

        <div class="metrics-grid" id="metrics"></div>

        <div class="panel">
            <div class="panel-header">
                <span class="panel-title">Metric Providers</span>
            </div>
            <div class="panel-body">
                <div class="provider-list" id="providers"></div>
            </div>
        </div>
    </div>

    <script>
        async function fetchJSON(url) {
            const res = await fetch(url);
            return res.json();
        }

        async function loadMetrics() {
            const data = await fetchJSON('/api/metrics');
            const container = document.getElementById('metrics');
            
            if (Object.keys(data).length === 0) {
                container.innerHTML = '<div class="metric-card"><div class="empty-state">No metrics collected yet. Click "Collect Metrics" to start.</div></div>';
                return;
            }

            let html = '';
            for (const [name, result] of Object.entries(data)) {
                const latest = result.values[result.values.length - 1];
                const chartBars = result.values.map(v => {
                    const height = Math.max(4, (v.value / result.max) * 60);
                    return '<div class="chart-bar" style="height: ' + height + 'px"></div>';
                }).join('');

                html += '<div class="metric-card">' +
                    '<div class="metric-header">' +
                        '<div class="metric-title">' + name + '</div>' +
                    '</div>' +
                    '<div>' +
                        '<span class="metric-value">' + latest.value.toFixed(1) + '</span>' +
                        '<span class="metric-unit">' + latest.unit + '</span>' +
                    '</div>' +
                    '<div class="metric-chart">' + chartBars + '</div>' +
                    '<div class="metric-stats">' +
                        '<div class="stat"><div class="stat-label">Min</div><div class="stat-value">' + result.min.toFixed(1) + '</div></div>' +
                        '<div class="stat"><div class="stat-label">Avg</div><div class="stat-value">' + result.average.toFixed(1) + '</div></div>' +
                        '<div class="stat"><div class="stat-label">Max</div><div class="stat-value">' + result.max.toFixed(1) + '</div></div>' +
                    '</div>' +
                '</div>';
            }
            container.innerHTML = html;
        }

        async function loadProviders() {
            const providers = await fetchJSON('/api/providers');
            const container = document.getElementById('providers');
            container.innerHTML = providers.map(p => '<span class="provider-tag">' + p + '</span>').join('');
        }

        async function collect() {
            await fetch('/api/collect', {method: 'POST'});
            loadMetrics();
        }

        loadMetrics();
        loadProviders();
        setInterval(loadMetrics, 5000);
    </script>
</body>
</html>`

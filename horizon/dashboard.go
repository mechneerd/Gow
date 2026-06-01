package horizon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

// Job represents a queued job for dashboard display
type Job struct {
	ID            string        `json:"id"`
	Queue         string        `json:"queue"`
	Payload       string        `json:"payload"`
	Status        string        `json:"status"` // pending, processing, completed, failed
	Attempts      int           `json:"attempts"`
	MaxAttempts   int           `json:"max_attempts"`
	CreatedAt     time.Time     `json:"created_at"`
	StartedAt     *time.Time    `json:"started_at,omitempty"`
	CompletedAt   *time.Time    `json:"completed_at,omitempty"`
	FailedAt      *time.Time    `json:"failed_at,omitempty"`
	Error         string        `json:"error,omitempty"`
	RetryCount    int           `json:"retry_count"`
	ReservedFor   string        `json:"reserved_for,omitempty"`
	AvailableAt   time.Time     `json:"available_at"`
}

// QueueStats represents statistics for a queue
type QueueStats struct {
	Name        string `json:"name"`
	Pending     int    `json:"pending"`
	Processing  int    `json:"processing"`
	Completed   int    `json:"completed"`
	Failed      int    `json:"failed"`
	Total       int    `json:"total"`
	AvgWaitTime string `json:"avg_wait_time"`
}

// Worker represents a queue worker
type Worker struct {
	ID         string     `json:"id"`
	Queue      string     `json:"queue"`
	Status     string     `json:"status"` // idle, busy, paused
	StartedAt  time.Time  `json:"started_at"`
	LastJobAt  *time.Time `json:"last_job_at,omitempty"`
	JobsProcessed int     `json:"jobs_processed"`
}

// Store is the interface for persisting dashboard data
type Store interface {
	GetJobs(status string, limit, offset int) ([]*Job, int, error)
	GetJobByID(id string) (*Job, error)
	GetQueueStats() ([]*QueueStats, error)
	GetWorkers() ([]*Worker, error)
	GetRecentJobs(limit int) ([]*Job, error)
}

// Dashboard provides a queue monitoring dashboard
type Dashboard struct {
	store  Store
	router *http.ServeMux
	mu     sync.RWMutex
}

// NewDashboard creates a new Horizon dashboard
func NewDashboard(store Store) *Dashboard {
	d := &Dashboard{
		store:  store,
		router: http.NewServeMux(),
	}
	d.routes()
	return d
}

func (d *Dashboard) routes() {
	d.router.HandleFunc("/", d.indexHandler)
	d.router.HandleFunc("/api/stats", d.statsHandler)
	d.router.HandleFunc("/api/jobs", d.jobsHandler)
	d.router.HandleFunc("/api/jobs/", d.jobDetailHandler)
	d.router.HandleFunc("/api/workers", d.workersHandler)
	d.router.HandleFunc("/api/recent", d.recentJobsHandler)
	d.router.HandleFunc("/api/retry/", d.retryHandler)
	d.router.HandleFunc("/api/kill/", d.killHandler)
}

// Handler returns the HTTP handler for the dashboard
func (d *Dashboard) Handler() http.Handler {
	return d.router
}

func (d *Dashboard) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, dashboardHTML)
}

func (d *Dashboard) statsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := d.store.GetQueueStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

func (d *Dashboard) jobsHandler(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "all"
	}
	limit := 50
	offset := 0
	fmt.Sscanf(r.URL.Query().Get("limit"), "%d", &limit)
	fmt.Sscanf(r.URL.Query().Get("offset"), "%d", &offset)

	jobs, total, err := d.store.GetJobs(status, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"jobs":   jobs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (d *Dashboard) jobDetailHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/jobs/"):]
	job, err := d.store.GetJobByID(id)
	if err != nil {
		http.Error(w, `{"error": "job not found"}`, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(job)
}

func (d *Dashboard) workersHandler(w http.ResponseWriter, r *http.Request) {
	workers, err := d.store.GetWorkers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(workers)
}

func (d *Dashboard) recentJobsHandler(w http.ResponseWriter, r *http.Request) {
	limit := 20
	fmt.Sscanf(r.URL.Query().Get("limit"), "%d", &limit)

	jobs, err := d.store.GetRecentJobs(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(jobs)
}

func (d *Dashboard) retryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Path[len("/api/retry/"):]
	job, err := d.store.GetJobByID(id)
	if err != nil {
		http.Error(w, `{"error": "job not found"}`, http.StatusNotFound)
		return
	}
	job.Status = "pending"
	job.Attempts = 0
	job.Error = ""
	json.NewEncoder(w).Encode(map[string]string{"status": "retried"})
}

func (d *Dashboard) killHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Path[len("/api/kill/"):]
	job, err := d.store.GetJobByID(id)
	if err != nil {
		http.Error(w, `{"error": "job not found"}`, http.StatusNotFound)
		return
	}
	now := time.Now()
	job.Status = "failed"
	job.FailedAt = &now
	job.Error = "Manually killed"
	json.NewEncoder(w).Encode(map[string]string{"status": "killed"})
}

// InMemoryStore is a simple in-memory store for testing
type InMemoryStore struct {
	mu      sync.RWMutex
	jobs    map[string]*Job
	workers []*Worker
}

// NewInMemoryStore creates a new in-memory store
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		jobs: make(map[string]*Job),
	}
}

func (s *InMemoryStore) AddJob(job *Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[job.ID] = job
}

func (s *InMemoryStore) SetWorkers(workers []*Worker) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.workers = workers
}

func (s *InMemoryStore) GetJobs(status string, limit, offset int) ([]*Job, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filtered []*Job
	for _, job := range s.jobs {
		if status == "all" || job.Status == status {
			filtered = append(filtered, job)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	total := len(filtered)
	start := offset
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}

	return filtered[start:end], total, nil
}

func (s *InMemoryStore) GetJobByID(id string) (*Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, ok := s.jobs[id]
	if !ok {
		return nil, fmt.Errorf("job not found: %s", id)
	}
	return job, nil
}

func (s *InMemoryStore) GetQueueStats() ([]*QueueStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	queueMap := make(map[string]*QueueStats)
	for _, job := range s.jobs {
		stats, ok := queueMap[job.Queue]
		if !ok {
			stats = &QueueStats{Name: job.Queue}
			queueMap[job.Queue] = stats
		}
		stats.Total++
		switch job.Status {
		case "pending":
			stats.Pending++
		case "processing":
			stats.Processing++
		case "completed":
			stats.Completed++
		case "failed":
			stats.Failed++
		}
	}

	var result []*QueueStats
	for _, stats := range queueMap {
		result = append(result, stats)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result, nil
}

func (s *InMemoryStore) GetWorkers() ([]*Worker, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*Worker, len(s.workers))
	copy(result, s.workers)
	return result, nil
}

func (s *InMemoryStore) GetRecentJobs(limit int) ([]*Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var all []*Job
	for _, job := range s.jobs {
		all = append(all, job)
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})

	if limit > len(all) {
		limit = len(all)
	}
	return all[:limit], nil
}

// Template for dashboard HTML
const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GoW Horizon - Queue Dashboard</title>
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
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 1rem;
            margin-bottom: 2rem;
        }
        .stat-card {
            background: var(--surface);
            padding: 1.5rem;
            border-radius: 0.75rem;
            border: 1px solid rgba(255,255,255,0.1);
        }
        .stat-label { color: var(--muted); font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.05em; }
        .stat-value { font-size: 2rem; font-weight: 700; margin-top: 0.25rem; }
        .stat-value.success { color: var(--success); }
        .stat-value.warning { color: var(--warning); }
        .stat-value.danger { color: var(--danger); }
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
        .tabs {
            display: flex;
            gap: 0.5rem;
            margin-bottom: 1rem;
        }
        .tab {
            padding: 0.5rem 1rem;
            border-radius: 0.5rem;
            cursor: pointer;
            background: transparent;
            color: var(--muted);
            border: 1px solid transparent;
            transition: all 0.2s;
        }
        .tab:hover { background: rgba(255,255,255,0.05); }
        .tab.active {
            background: var(--primary);
            color: white;
        }
        table {
            width: 100%;
            border-collapse: collapse;
        }
        th, td {
            padding: 0.75rem 1rem;
            text-align: left;
            border-bottom: 1px solid rgba(255,255,255,0.05);
        }
        th {
            color: var(--muted);
            font-size: 0.75rem;
            text-transform: uppercase;
            letter-spacing: 0.05em;
            font-weight: 500;
        }
        tr:hover { background: rgba(255,255,255,0.02); }
        .badge {
            display: inline-block;
            padding: 0.25rem 0.5rem;
            border-radius: 9999px;
            font-size: 0.75rem;
            font-weight: 500;
        }
        .badge-pending { background: rgba(245,158,11,0.2); color: var(--warning); }
        .badge-processing { background: rgba(99,102,241,0.2); color: var(--primary); }
        .badge-completed { background: rgba(34,197,94,0.2); color: var(--success); }
        .badge-failed { background: rgba(239,68,68,0.2); color: var(--danger); }
        .btn {
            padding: 0.5rem 1rem;
            border-radius: 0.5rem;
            cursor: pointer;
            border: none;
            font-size: 0.875rem;
            font-weight: 500;
            transition: opacity 0.2s;
        }
        .btn:hover { opacity: 0.8; }
        .btn-primary { background: var(--primary); color: white; }
        .btn-danger { background: var(--danger); color: white; }
        .btn-sm { padding: 0.25rem 0.5rem; font-size: 0.75rem; }
        .empty-state {
            text-align: center;
            padding: 3rem;
            color: var(--muted);
        }
        .refresh-btn {
            background: transparent;
            border: 1px solid rgba(255,255,255,0.2);
            color: var(--muted);
            padding: 0.5rem 1rem;
            border-radius: 0.5rem;
            cursor: pointer;
            font-size: 0.875rem;
        }
        .refresh-btn:hover { background: rgba(255,255,255,0.05); }
        .auto-refresh {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            color: var(--muted);
            font-size: 0.875rem;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <div>
                <h1>GoW Horizon</h1>
                <div class="subtitle">Queue Dashboard</div>
            </div>
            <div class="auto-refresh">
                <input type="checkbox" id="autoRefresh" checked>
                <label for="autoRefresh">Auto-refresh</label>
                <button class="refresh-btn" onclick="loadAll()">Refresh</button>
            </div>
        </header>

        <div class="stats-grid" id="stats"></div>

        <div class="panel">
            <div class="panel-header">
                <span class="panel-title">Queues</span>
            </div>
            <div class="panel-body">
                <div id="queues"></div>
            </div>
        </div>

        <div class="panel">
            <div class="panel-header">
                <span class="panel-title">Workers</span>
            </div>
            <div class="panel-body">
                <div id="workers"></div>
            </div>
        </div>

        <div class="panel">
            <div class="panel-header">
                <span class="panel-title">Recent Jobs</span>
                <div class="tabs" id="jobTabs">
                    <button class="tab active" data-status="all">All</button>
                    <button class="tab" data-status="pending">Pending</button>
                    <button class="tab" data-status="processing">Processing</button>
                    <button class="tab" data-status="completed">Completed</button>
                    <button class="tab" data-status="failed">Failed</button>
                </div>
            </div>
            <div class="panel-body">
                <div id="jobs"></div>
            </div>
        </div>
    </div>

    <script>
        let currentStatus = 'all';
        let refreshInterval;

        async function fetchJSON(url) {
            const res = await fetch(url);
            return res.json();
        }

        async function loadStats() {
            const stats = await fetchJSON('/api/stats');
            const totalPending = stats.reduce((sum, s) => sum + s.pending, 0);
            const totalProcessing = stats.reduce((sum, s) => sum + s.processing, 0);
            const totalCompleted = stats.reduce((sum, s) => sum + s.completed, 0);
            const totalFailed = stats.reduce((sum, s) => sum + s.failed, 0);

            document.getElementById('stats').innerHTML = 
                '<div class="stat-card"><div class="stat-label">Pending</div><div class="stat-value warning">' + totalPending + '</div></div>' +
                '<div class="stat-card"><div class="stat-label">Processing</div><div class="stat-value">' + totalProcessing + '</div></div>' +
                '<div class="stat-card"><div class="stat-label">Completed</div><div class="stat-value success">' + totalCompleted + '</div></div>' +
                '<div class="stat-card"><div class="stat-label">Failed</div><div class="stat-value danger">' + totalFailed + '</div></div>';

            let queuesHtml = '<table><thead><tr><th>Queue</th><th>Pending</th><th>Processing</th><th>Completed</th><th>Failed</th><th>Total</th></tr></thead><tbody>';
            for (const q of stats) {
                queuesHtml += '<tr><td>' + q.name + '</td><td>' + q.pending + '</td><td>' + q.processing + '</td><td>' + q.completed + '</td><td>' + q.failed + '</td><td>' + q.total + '</td></tr>';
            }
            queuesHtml += '</tbody></table>';
            if (stats.length === 0) queuesHtml = '<div class="empty-state">No queues configured</div>';
            document.getElementById('queues').innerHTML = queuesHtml;
        }

        async function loadWorkers() {
            const workers = await fetchJSON('/api/workers');
            let html = '<table><thead><tr><th>ID</th><th>Queue</th><th>Status</th><th>Jobs Processed</th><th>Last Job</th></tr></thead><tbody>';
            for (const w of workers) {
                html += '<tr><td>' + w.id + '</td><td>' + w.queue + '</td><td><span class="badge badge-' + w.status + '">' + w.status + '</span></td><td>' + w.jobs_processed + '</td><td>' + (w.last_job_at || 'Never') + '</td></tr>';
            }
            html += '</tbody></table>';
            if (workers.length === 0) html = '<div class="empty-state">No workers running</div>';
            document.getElementById('workers').innerHTML = html;
        }

        async function loadJobs() {
            const data = await fetchJSON('/api/jobs?status=' + currentStatus + '&limit=50');
            const jobs = data.jobs || [];
            let html = '<table><thead><tr><th>ID</th><th>Queue</th><th>Status</th><th>Attempts</th><th>Created</th><th>Actions</th></tr></thead><tbody>';
            for (const j of jobs) {
                const actions = j.status === 'failed' ? 
                    '<button class="btn btn-primary btn-sm" onclick="retryJob(\'' + j.id + '\')">Retry</button> ' +
                    '<button class="btn btn-danger btn-sm" onclick="killJob(\'' + j.id + '\')">Kill</button>' : '';
                html += '<tr><td>' + j.id + '</td><td>' + j.queue + '</td><td><span class="badge badge-' + j.status + '">' + j.status + '</span></td><td>' + j.attempts + '/' + j.max_attempts + '</td><td>' + j.created_at + '</td><td>' + actions + '</td></tr>';
            }
            html += '</tbody></table>';
            if (jobs.length === 0) html = '<div class="empty-state">No jobs found</div>';
            document.getElementById('jobs').innerHTML = html;
        }

        async function retryJob(id) {
            await fetch('/api/retry/' + id, {method: 'POST'});
            loadJobs();
        }

        async function killJob(id) {
            await fetch('/api/kill/' + id, {method: 'POST'});
            loadJobs();
        }

        function loadAll() {
            loadStats();
            loadWorkers();
            loadJobs();
        }

        document.getElementById('jobTabs').addEventListener('click', function(e) {
            if (e.target.classList.contains('tab')) {
                document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
                e.target.classList.add('active');
                currentStatus = e.target.dataset.status;
                loadJobs();
            }
        });

        document.getElementById('autoRefresh').addEventListener('change', function() {
            if (this.checked) {
                refreshInterval = setInterval(loadAll, 5000);
            } else {
                clearInterval(refreshInterval);
            }
        });

        loadAll();
        refreshInterval = setInterval(loadAll, 5000);
    </script>
</body>
</html>`

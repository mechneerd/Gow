package nova

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
	"sync"
)

// Resource represents a resource that can be managed in Nova
type Resource struct {
	Name       string
	TableName  string
	Fields     []Field
	GetAll     func() ([]map[string]any, error)
	GetByID    func(id string) (map[string]any, error)
	Create     func(data map[string]any) error
	Update     func(id string, data map[string]any) error
	Delete     func(id string) error
}

// Field represents a field in a resource
type Field struct {
	Name     string
	Label    string
	Type     string // text, number, email, password, select, textarea, boolean, date
	Options  []string // For select fields
	Required bool
}

// Panel represents a dashboard panel
type Panel struct {
	Title  string
	Stats  []Stat
}

// Stat represents a dashboard statistic
type Stat struct {
	Label string
	Value string
	Type  string // info, success, warning, danger
}

// Admin is the Nova admin panel
type Admin struct {
	mu        sync.RWMutex
	resources map[string]*Resource
	panels    []*Panel
	router    *http.ServeMux
	authFunc  func(r *http.Request) bool
}

// New creates a new Nova admin panel
func New() *Admin {
	a := &Admin{
		resources: make(map[string]*Resource),
		panels:    make([]*Panel, 0),
		router:    http.NewServeMux(),
		authFunc:  defaultAuth,
	}
	a.routes()
	return a
}

// SetAuthFunc sets the authentication function
func (a *Admin) SetAuthFunc(fn func(r *http.Request) bool) {
	a.authFunc = fn
}

// Resource registers a new resource
func (a *Admin) Resource(name string, resource *Resource) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.resources[strings.ToLower(name)] = resource
}

// Panel adds a dashboard panel
func (a *Admin) Panel(panel *Panel) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.panels = append(a.panels, panel)
}

// Handler returns the HTTP handler
func (a *Admin) Handler() http.Handler {
	return a.router
}

func (a *Admin) routes() {
	a.router.HandleFunc("/", a.dashboardHandler)
	a.router.HandleFunc("/resources/", a.resourceHandler)
	a.router.HandleFunc("/api/resources/", a.apiResourceHandler)
}

func defaultAuth(r *http.Request) bool {
	// Default: allow all (in production, check session/auth)
	return true
}

func (a *Admin) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if !a.authFunc(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	a.mu.RLock()
	defer a.mu.RUnlock()

	data := map[string]any{
		"Resources": a.resources,
		"Panels":    a.panels,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.New("dashboard").Parse(dashboardHTML))
	tmpl.Execute(w, data)
}

func (a *Admin) resourceHandler(w http.ResponseWriter, r *http.Request) {
	if !a.authFunc(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/resources/")
	parts := strings.SplitN(path, "/", 2)
	resourceName := strings.ToLower(parts[0])

	a.mu.RLock()
	resource, ok := a.resources[resourceName]
	a.mu.RUnlock()

	if !ok {
		http.Error(w, "Resource not found", http.StatusNotFound)
		return
	}

	data := map[string]any{
		"Resource": resource,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.New("resource").Parse(resourceHTML))
	tmpl.Execute(w, data)
}

func (a *Admin) apiResourceHandler(w http.ResponseWriter, r *http.Request) {
	if !a.authFunc(r) {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/resources/")
	parts := strings.SplitN(path, "/", 2)
	resourceName := strings.ToLower(parts[0])

	a.mu.RLock()
	resource, ok := a.resources[resourceName]
	a.mu.RUnlock()

	if !ok {
		http.Error(w, `{"error": "resource not found"}`, http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if len(parts) > 1 && parts[1] != "" {
			// Get by ID
			item, err := resource.GetByID(parts[1])
			if err != nil {
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(item)
		} else {
			// Get all
			items, err := resource.GetAll()
			if err != nil {
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(items)
		}

	case http.MethodPost:
		var data map[string]any
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
			return
		}
		if err := resource.Create(data); err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "created"})

	case http.MethodPut:
		if len(parts) < 2 || parts[1] == "" {
			http.Error(w, `{"error": "id required"}`, http.StatusBadRequest)
			return
		}
		var data map[string]any
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
			return
		}
		if err := resource.Update(parts[1], data); err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "updated"})

	case http.MethodDelete:
		if len(parts) < 2 || parts[1] == "" {
			http.Error(w, `{"error": "id required"}`, http.StatusBadRequest)
			return
		}
		if err := resource.Delete(parts[1]); err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})

	default:
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GoW Nova - Admin Panel</title>
    <style>
        :root {
            --primary: #6366f1;
            --success: #22c55e;
            --warning: #f59e0b;
            --danger: #ef4444;
            --bg: #f8fafc;
            --surface: #ffffff;
            --text: #0f172a;
            --muted: #64748b;
            --border: #e2e8f0;
        }
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: var(--bg);
            color: var(--text);
            line-height: 1.6;
        }
        .layout {
            display: flex;
            min-height: 100vh;
        }
        .sidebar {
            width: 250px;
            background: var(--surface);
            border-right: 1px solid var(--border);
            padding: 1.5rem;
        }
        .logo {
            font-size: 1.25rem;
            font-weight: 700;
            color: var(--primary);
            margin-bottom: 2rem;
        }
        .nav-item {
            display: block;
            padding: 0.75rem 1rem;
            color: var(--muted);
            text-decoration: none;
            border-radius: 0.5rem;
            margin-bottom: 0.25rem;
            transition: all 0.2s;
        }
        .nav-item:hover {
            background: var(--bg);
            color: var(--text);
        }
        .nav-item.active {
            background: rgba(99,102,241,0.1);
            color: var(--primary);
        }
        .main {
            flex: 1;
            padding: 2rem;
        }
        .header {
            margin-bottom: 2rem;
        }
        h1 { font-size: 1.5rem; font-weight: 600; }
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
            border: 1px solid var(--border);
        }
        .stat-label { color: var(--muted); font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.05em; }
        .stat-value { font-size: 2rem; font-weight: 700; margin-top: 0.25rem; }
        .stat-value.info { color: var(--primary); }
        .stat-value.success { color: var(--success); }
        .stat-value.warning { color: var(--warning); }
        .stat-value.danger { color: var(--danger); }
        .panel {
            background: var(--surface);
            border-radius: 0.75rem;
            border: 1px solid var(--border);
            margin-bottom: 1.5rem;
        }
        .panel-header {
            padding: 1rem 1.5rem;
            border-bottom: 1px solid var(--border);
            font-weight: 600;
        }
        .panel-body { padding: 1.5rem; }
        .resource-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 1rem;
        }
        .resource-card {
            background: var(--bg);
            padding: 1.5rem;
            border-radius: 0.75rem;
            text-decoration: none;
            color: var(--text);
            transition: all 0.2s;
        }
        .resource-card:hover {
            background: rgba(99,102,241,0.05);
            border-color: var(--primary);
        }
        .resource-name { font-weight: 600; margin-bottom: 0.25rem; }
        .resource-table { color: var(--muted); font-size: 0.875rem; }
    </style>
</head>
<body>
    <div class="layout">
        <div class="sidebar">
            <div class="logo">GoW Nova</div>
            <a href="/" class="nav-item active">Dashboard</a>
            {{range $name, $res := .Resources}}
            <a href="/resources/{{$name}}" class="nav-item">{{$res.Name}}</a>
            {{end}}
        </div>
        <div class="main">
            <div class="header">
                <h1>Dashboard</h1>
            </div>

            {{if .Panels}}
            <div class="stats-grid">
                {{range .Panels}}
                {{range .Stats}}
                <div class="stat-card">
                    <div class="stat-label">{{.Label}}</div>
                    <div class="stat-value {{.Type}}">{{.Value}}</div>
                </div>
                {{end}}
                {{end}}
            </div>
            {{end}}

            <div class="panel">
                <div class="panel-header">Resources</div>
                <div class="panel-body">
                    <div class="resource-grid">
                        {{range $name, $res := .Resources}}
                        <a href="/resources/{{$name}}" class="resource-card">
                            <div class="resource-name">{{$res.Name}}</div>
                            <div class="resource-table">{{$res.TableName}}</div>
                        </a>
                        {{end}}
                    </div>
                </div>
            </div>
        </div>
    </div>
</body>
</html>`

const resourceHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Resource.Name}} - GoW Nova</title>
    <style>
        :root {
            --primary: #6366f1;
            --success: #22c55e;
            --warning: #f59e0b;
            --danger: #ef4444;
            --bg: #f8fafc;
            --surface: #ffffff;
            --text: #0f172a;
            --muted: #64748b;
            --border: #e2e8f0;
        }
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: var(--bg);
            color: var(--text);
            line-height: 1.6;
        }
        .layout { display: flex; min-height: 100vh; }
        .sidebar {
            width: 250px;
            background: var(--surface);
            border-right: 1px solid var(--border);
            padding: 1.5rem;
        }
        .logo { font-size: 1.25rem; font-weight: 700; color: var(--primary); margin-bottom: 2rem; }
        .nav-item {
            display: block;
            padding: 0.75rem 1rem;
            color: var(--muted);
            text-decoration: none;
            border-radius: 0.5rem;
            margin-bottom: 0.25rem;
        }
        .nav-item:hover { background: var(--bg); color: var(--text); }
        .nav-item.active { background: rgba(99,102,241,0.1); color: var(--primary); }
        .main { flex: 1; padding: 2rem; }
        .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
        h1 { font-size: 1.5rem; font-weight: 600; }
        .btn {
            padding: 0.5rem 1rem;
            border-radius: 0.5rem;
            cursor: pointer;
            border: none;
            font-size: 0.875rem;
            font-weight: 500;
            background: var(--primary);
            color: white;
        }
        .panel {
            background: var(--surface);
            border-radius: 0.75rem;
            border: 1px solid var(--border);
        }
        table { width: 100%; border-collapse: collapse; }
        th, td { padding: 0.75rem 1rem; text-align: left; border-bottom: 1px solid var(--border); }
        th { color: var(--muted); font-size: 0.75rem; text-transform: uppercase; font-weight: 500; }
        tr:hover { background: var(--bg); }
        .badge {
            display: inline-block;
            padding: 0.25rem 0.5rem;
            border-radius: 9999px;
            font-size: 0.75rem;
        }
        .badge-success { background: rgba(34,197,94,0.1); color: var(--success); }
        .empty { text-align: center; padding: 3rem; color: var(--muted); }
    </style>
</head>
<body>
    <div class="layout">
        <div class="sidebar">
            <div class="logo">GoW Nova</div>
            <a href="/" class="nav-item">Dashboard</a>
            <a href="/resources/{{.Resource.Name | ToLower}}" class="nav-item active">{{.Resource.Name}}</a>
        </div>
        <div class="main">
            <div class="header">
                <h1>{{.Resource.Name}}</h1>
                <button class="btn" onclick="showCreate()">Create {{.Resource.Name}}</button>
            </div>
            <div class="panel">
                <table>
                    <thead>
                        <tr>
                            {{range .Resource.Fields}}
                            <th>{{.Label}}</th>
                            {{end}}
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody id="table-body">
                        <tr><td colspan="100" class="empty">Loading...</td></tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
    <script>
        async function loadData() {
            const res = await fetch('/api/resources/{{.Resource.Name | ToLower}}');
            const data = await res.json();
            const tbody = document.getElementById('table-body');
            if (!data || data.length === 0) {
                tbody.innerHTML = '<tr><td colspan="100" class="empty">No records found</td></tr>';
                return;
            }
            let html = '';
            for (const item of data) {
                html += '<tr>';
                {{range .Resource.Fields}}
                html += '<td>' + (item['{{.Name}}'] || '') + '</td>';
                {{end}}
                html += '<td><button class="btn" style="padding:0.25rem 0.5rem;font-size:0.75rem">Edit</button></td>';
                html += '</tr>';
            }
            tbody.innerHTML = html;
        }
        loadData();
    </script>
</body>
</html>`

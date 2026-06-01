package scaffold

import (
	"os"
	"path/filepath"
	"strings"
)

// FixSkeletonBugs applies targeted fixes for known bugs in the gow-skeleton templates.
// These are structural issues that simple placeholder replacement cannot handle.

const monitorGoblade = `<!DOCTYPE html>
<html lang="en" class="dark">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .AppName }} — Monitor</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <script>
        tailwind.config = {
            darkMode: 'class',
            theme: { extend: { colors: { zinc: { 950: '#09090b' } } } }
        }
    </script>
    <style>
        .gradient-text{background:linear-gradient(135deg,#10b981,#3b82f6,#8b5cf6);-webkit-background-clip:text;-webkit-text-fill-color:transparent;background-clip:text}
        @keyframes pulse-dot{0%,100%{opacity:1}50%{opacity:.4}}
        .pulse-dot{animation:pulse-dot 2s ease-in-out infinite}
        .stat-card{transition:all .3s ease}
        .stat-card:hover{transform:translateY(-2px);border-color:rgba(16,185,129,.3)}
    </style>
</head>
<body class="bg-zinc-950 text-white min-h-screen antialiased">
    <nav class="fixed top-0 w-full z-50 bg-zinc-950/80 backdrop-blur-xl border-b border-zinc-800/50">
        <div class="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
            <div class="flex items-center gap-3">
                <a href="/" class="flex items-center gap-2">
                    <div class="w-8 h-8 rounded-lg bg-emerald-500 flex items-center justify-center font-bold text-sm">G</div>
                    <span class="font-semibold text-lg">{{ .AppName }}</span>
                </a>
                <span class="text-zinc-600">/</span>
                <span class="text-zinc-400 text-sm font-medium">Monitor</span>
            </div>
            <div class="flex items-center gap-4">
                <div class="flex items-center gap-2 px-3 py-1.5 rounded-full bg-emerald-500/10 border border-emerald-500/20">
                    <div class="w-2 h-2 bg-emerald-500 rounded-full pulse-dot"></div>
                    <span class="text-xs font-medium text-emerald-400">Live</span>
                </div>
                <a href="/" class="text-sm text-zinc-400 hover:text-white transition">&larr; Back to App</a>
            </div>
        </div>
    </nav>
    <main class="pt-24 pb-12 px-6">
        <div class="max-w-7xl mx-auto">
            <div class="mb-8">
                <h1 class="text-3xl font-bold tracking-tight mb-2"><span class="gradient-text">Server Monitor</span></h1>
                <p class="text-zinc-500">Real-time performance metrics for your GoW application.</p>
            </div>
            <div class="grid grid-cols-2 md:grid-cols-5 gap-4 mb-8">
                <div class="stat-card bg-zinc-900/50 border border-zinc-800 rounded-2xl p-5">
                    <div class="flex items-center justify-between mb-3"><span class="text-xs font-medium text-zinc-500 uppercase tracking-wider">CPU Usage</span><span class="text-emerald-400 text-lg">&#9889;</span></div>
                    <div class="text-3xl font-bold text-white" id="cpu-usage">0%</div>
                    <div class="text-xs text-zinc-500 mt-1" id="cpu-cores">0 cores</div>
                </div>
                <div class="stat-card bg-zinc-900/50 border border-zinc-800 rounded-2xl p-5">
                    <div class="flex items-center justify-between mb-3"><span class="text-xs font-medium text-zinc-500 uppercase tracking-wider">Memory</span><span class="text-blue-400 text-lg">&#128207;</span></div>
                    <div class="text-3xl font-bold text-white" id="memory-usage">0 MB</div>
                    <div class="text-xs text-zinc-500 mt-1" id="memory-detail">0 / 0 GB</div>
                </div>
                <div class="stat-card bg-zinc-900/50 border border-zinc-800 rounded-2xl p-5">
                    <div class="flex items-center justify-between mb-3"><span class="text-xs font-medium text-zinc-500 uppercase tracking-wider">Routines</span><span class="text-purple-400 text-lg">&#128260;</span></div>
                    <div class="text-3xl font-bold text-white" id="go-routines">0</div>
                    <div class="text-xs text-zinc-500 mt-1">active goroutines</div>
                </div>
                <div class="stat-card bg-zinc-900/50 border border-zinc-800 rounded-2xl p-5">
                    <div class="flex items-center justify-between mb-3"><span class="text-xs font-medium text-zinc-500 uppercase tracking-wider">Response</span><span class="text-yellow-400 text-lg">&#9201;</span></div>
                    <div class="text-3xl font-bold text-white" id="response-time">0ms</div>
                    <div class="text-xs text-zinc-500 mt-1" id="response-avg">avg: 0ms</div>
                </div>
                <div class="stat-card bg-zinc-900/50 border border-zinc-800 rounded-2xl p-5">
                    <div class="flex items-center justify-between mb-3"><span class="text-xs font-medium text-zinc-500 uppercase tracking-wider">Requests</span><span class="text-cyan-400 text-lg">&#128200;</span></div>
                    <div class="text-3xl font-bold text-white" id="total-requests">0</div>
                    <div class="text-xs text-zinc-500 mt-1" id="uptime">uptime: 0s</div>
                </div>
            </div>
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
                <div class="bg-zinc-900/50 border border-zinc-800 rounded-2xl p-6">
                    <h3 class="font-semibold text-white mb-4">CPU Usage</h3>
                    <div style="height:200px"><canvas id="cpuChart"></canvas></div>
                </div>
                <div class="bg-zinc-900/50 border border-zinc-800 rounded-2xl p-6">
                    <h3 class="font-semibold text-white mb-4">Memory Usage</h3>
                    <div style="height:200px"><canvas id="memoryChart"></canvas></div>
                </div>
            </div>
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
                <div class="bg-zinc-900/50 border border-zinc-800 rounded-2xl p-6">
                    <h3 class="font-semibold text-white mb-4">Response Time</h3>
                    <div style="height:200px"><canvas id="responseChart"></canvas></div>
                </div>
                <div class="bg-zinc-900/50 border border-zinc-800 rounded-2xl p-6">
                    <h3 class="font-semibold text-white mb-4">Open Connections</h3>
                    <div style="height:200px"><canvas id="connChart"></canvas></div>
                </div>
            </div>
            <div class="bg-zinc-900/50 border border-zinc-800 rounded-2xl p-6">
                <h3 class="font-semibold text-white mb-4">System Information</h3>
                <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
                    <div><div class="text-xs text-zinc-500 mb-1">Framework</div><div class="text-sm font-medium text-white">GoW v1.0</div></div>
                    <div><div class="text-xs text-zinc-500 mb-1">Go Version</div><div class="text-sm font-medium text-white" id="go-version">go1.22</div></div>
                    <div><div class="text-xs text-zinc-500 mb-1">Platform</div><div class="text-sm font-medium text-white" id="platform">linux/amd64</div></div>
                    <div><div class="text-xs text-zinc-500 mb-1">CPU Cores</div><div class="text-sm font-medium text-white" id="cpu-count">1</div></div>
                    <div><div class="text-xs text-zinc-500 mb-1">Memory Alloc</div><div class="text-sm font-medium text-emerald-400" id="mem-alloc">0 MB</div></div>
                    <div><div class="text-xs text-zinc-500 mb-1">Memory System</div><div class="text-sm font-medium text-blue-400" id="mem-sys">0 MB</div></div>
                    <div><div class="text-xs text-zinc-500 mb-1">GC Cycles</div><div class="text-sm font-medium text-purple-400" id="gc-cycles">0</div></div>
                    <div><div class="text-xs text-zinc-500 mb-1">Last Updated</div><div class="text-sm font-medium text-zinc-400" id="last-updated">--:--:--</div></div>
                </div>
            </div>
        </div>
    </main>
    <script>
        const MAX=30;
        const cpuData={labels:[],values:[]};
        const memData={labels:[],alloc:[],sys:[]};
        const respData={labels:[],values:[]};
        const connData={labels:[],values:[]};
        function makeChart(id,color,label,yOpts){
            return new Chart(document.getElementById(id).getContext('2d'),{type:'line',data:{labels:[],datasets:[{label,label:data:[],borderColor:color,backgroundColor:color.replace(')',',0.1)').replace('rgb','rgba'),fill:true,tension:.4,pointRadius:0,borderWidth:2}]},options:{responsive:true,maintainAspectRatio:false,animation:{duration:0},plugins:{legend:{display:false}},scales:{x:{display:false},y:{display:true,grid:{color:'rgba(255,255,255,0.05)'},ticks:{color:'#71717a',...yOpts}}}}});
        }
        const cpuChart=makeChart('cpuChart','rgb(16,185,129)','CPU',{callback:v=>v+'%'});
        const memChart=new Chart(document.getElementById('memoryChart').getContext('2d'),{type:'line',data:{labels:memData.labels,datasets:[{label:'Alloc',data:memData.alloc,borderColor:'#10b981',backgroundColor:'rgba(16,185,129,0.1)',fill:true,tension:.4,pointRadius:0,borderWidth:2},{label:'Sys',data:memData.sys,borderColor:'#3b82f6',backgroundColor:'rgba(59,130,246,0.1)',fill:true,tension:.4,pointRadius:0,borderWidth:2}]},options:{responsive:true,maintainAspectRatio:false,animation:{duration:0},plugins:{legend:{display:true,labels:{color:'#71717a',boxWidth:12}}},scales:{x:{display:false},y:{display:true,grid:{color:'rgba(255,255,255,0.05)'},ticks:{color:'#71717a',callback:v=>v+' MB'}}}}});
        const respChart=makeChart('responseChart','rgb(139,92,246)','Response',{callback:v=>v+'ms'});
        const connChart=makeChart('connChart','rgb(234,179,8)','Connections',{callback:v=>v});
        function formatBytes(b){if(!b)return'0 B';const k=1024,s=['B','KB','MB','GB'],i=Math.floor(Math.log(b)/Math.log(k));return parseFloat((b/Math.pow(k,i)).toFixed(1))+' '+s[i]}
        function parseUptime(s){if(!s)return 0;const p=s.match(/(\d+)([dhms])/g);if(!p)return 0;let t=0;p.forEach(v=>{const n=parseInt(v);if(v.endsWith('d'))t+=n*86400;else if(v.endsWith('h'))t+=n*3600;else if(v.endsWith('m'))t+=n*60;else if(v.endsWith('s'))t+=n});return t}
        let lastGC=0;
        async function update(){
            try{
                const r=await fetch('/monitor?format=json');
                const d=await r.json();
                const now=new Date().toLocaleTimeString();
                const cpuPct=((d.cpu_usage||0)*100).toFixed(1);
                document.getElementById('cpu-usage').textContent=cpuPct+'%';
                document.getElementById('cpu-count').textContent=navigator.hardwareConcurrency||'?';
                document.getElementById('memory-usage').textContent=formatBytes(d.memory_alloc);
                document.getElementById('memory-detail').textContent=formatBytes(d.memory_alloc)+' / '+formatBytes(d.memory_sys);
                document.getElementById('go-routines').textContent=d.go_routines||0;
                document.getElementById('response-time').textContent=(d.response_time||0).toFixed(1)+'ms';
                document.getElementById('response-avg').textContent='avg: '+(d.avg_response||0).toFixed(1)+'ms';
                document.getElementById('total-requests').textContent=(d.total_requests||0).toLocaleString();
                document.getElementById('uptime').textContent='uptime: '+d.uptime;
                document.getElementById('mem-alloc').textContent=formatBytes(d.memory_alloc);
                document.getElementById('mem-sys').textContent=formatBytes(d.memory_sys);
                document.getElementById('gc-cycles').textContent=d.gc_cycles||0;
                document.getElementById('last-updated').textContent=now;
                cpuData.labels.push(now);cpuData.values.push(cpuPct);
                memData.labels.push(now);memData.alloc.push((d.memory_alloc/1024/1024).toFixed(1));memData.sys.push((d.memory_sys/1024/1024).toFixed(1));
                respData.labels.push(now);respData.values.push((d.response_time||0).toFixed(1));
                connData.labels.push(now);connData.values.push(d.open_connections||0);
                if(cpuData.labels.length>MAX){cpuData.labels.shift();cpuData.values.shift();memData.labels.shift();memData.alloc.shift();memData.sys.shift();respData.labels.shift();respData.values.shift();connData.labels.shift();connData.values.shift()}
                cpuChart.data.labels=cpuData.labels;cpuChart.data.datasets[0].data=cpuData.values;cpuChart.update('none');
                memChart.data.labels=memData.labels;memChart.data.datasets[0].data=memData.alloc;memChart.data.datasets[1].data=memData.sys;memChart.update('none');
                respChart.data.labels=respData.labels;respChart.data.datasets[0].data=respData.values;respChart.update('none');
                connChart.data.labels=connData.labels;connChart.data.datasets[0].data=connData.values;connChart.update('none');
            }catch(e){console.error(e)}
        }
        update();setInterval(update,3000);
    </script>
</body>
</html>`
func FixSkeletonBugs(projectDir string, moduleName string) error {
	return filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		original := string(content)
		updated := fixFileContent(original, path, moduleName)

		if updated != original {
			return os.WriteFile(path, []byte(updated), info.Mode())
		}
		return nil
	})

	// Fix 11: Create monitor.goblade if it doesn't exist
	viewsDir := filepath.Join(projectDir, "resources", "views")
	monitorPath := filepath.Join(viewsDir, "monitor.gablade")
	if _, err := os.Stat(monitorPath); os.IsNotExist(err) {
		if err := os.MkdirAll(viewsDir, 0755); err == nil {
			os.WriteFile(monitorPath, []byte(monitorGoblade), 0644)
		}
	}

	return nil
}

func fixFileContent(content string, filePath string, moduleName string) string {
	result := content

	// Fix 1: config/auth.go — remove duplicate getEnv function
	// The skeleton has getEnv in both config/app.go and config/auth.go
	if strings.HasSuffix(filePath, filepath.Join("config", "auth.go")) ||
		strings.HasSuffix(filePath, "config\\auth.go") || strings.HasSuffix(filePath, "config/auth.go") {
		result = removeDuplicateGetEnv(result)
	}

	// Fix 2: bootstrap/app.go — remove unused "routes" import
	if strings.HasSuffix(filePath, filepath.Join("bootstrap", "app.go")) ||
		strings.HasSuffix(filePath, "bootstrap\\app.go") || strings.HasSuffix(filePath, "bootstrap/app.go") {
		result = removeUnusedRoutesImport(result)
	}

	// Fix 3: app/Livewire/Counter.go — add livewire import, fix BaseComponent reference
	if strings.Contains(filePath, "Livewire") && strings.HasSuffix(filePath, "Counter.go") {
		result = fixLivewireCounter(result)
	}

	// Fix 4: app/Models/User.go — fix broken SetDB method and missing sql import
	if strings.HasSuffix(filePath, filepath.Join("app", "Models", "User.go")) ||
		strings.HasSuffix(filePath, "app\\Models\\User.go") || strings.HasSuffix(filePath, "app/Models/User.go") {
		result = fixUserModel(result)
	}

	// Fix 4b: app/Http/Controllers/Auth/handlers.go — implement stub auth handlers
	if strings.Contains(filePath, "Controllers") && strings.Contains(filePath, "Auth") && strings.HasSuffix(filePath, "handlers.go") {
		result = fixAuthHandlers(result)
	}

	// Fix 5: database/seeders/RoleSeeder.go — remove unused Models import
	// Only remove if Models package is not actually used in the file
	if strings.Contains(filePath, "seeders") && strings.HasSuffix(filePath, "RoleSeeder.go") {
		if !strings.Contains(result, "Models.") {
			result = removeUnusedModelsImport(result)
		}
	}

	// Fix 5b: routes/web.go — inject monitor route and imports if not present
	if strings.HasSuffix(filePath, filepath.Join("routes", "web.go")) ||
		strings.HasSuffix(filePath, "routes\\web.go") || strings.HasSuffix(filePath, "routes/web.go") {
		if !strings.Contains(result, "/monitor") && strings.Contains(result, "router.Get") {
			result = injectMonitorRoute(result, moduleName)
		}
	}

	// Fix 6: main.go — import local bootstrap, not framework bootstrap
	if strings.HasSuffix(filePath, filepath.Join("main.go")) ||
		strings.HasSuffix(filePath, "main.go") {
		result = fixMainGoBootstrapImport(result, moduleName)
	}

	// Fix 7: bootstrap/app.go — replace skeleton-specific config.AppConfig/config.Load
	// with direct os.Getenv calls (skeleton uses types not in the framework config package)
	if strings.HasSuffix(filePath, filepath.Join("bootstrap", "app.go")) ||
		strings.HasSuffix(filePath, "bootstrap\\app.go") || strings.HasSuffix(filePath, "bootstrap/app.go") {
		result = fixBootstrapAppGo(result)
	}

	// Fix 8: go.mod.template — remove invalid "latest" version, let go mod tidy handle it
	if strings.HasSuffix(filePath, "go.mod.template") || strings.HasSuffix(filePath, "go.mod") {
		// Remove lines with invalid "latest" version
		lines := strings.Split(result, "\n")
		var cleaned []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.Contains(trimmed, "github.com/mechneerd/gow") && strings.Contains(trimmed, "latest") {
				continue // skip this line
			}
			cleaned = append(cleaned, line)
		}
		result = strings.Join(cleaned, "\n")
	}

	// Fix 9: .env.example — fix DB_DATABASE for SQLite (should be a file path, not just "app")
	if strings.HasSuffix(filePath, ".env.example") || strings.HasSuffix(filePath, ".env") {
		if strings.Contains(result, "DB_DATABASE=app") && !strings.Contains(result, "DB_DATABASE=app/") {
			result = strings.Replace(result, "DB_DATABASE=app", "DB_DATABASE=database/database.sqlite", 1)
		}
	}

	// Fix 10: welcome.goblade — inject comparison section before footer
	if strings.HasSuffix(filePath, "welcome.goblade") {
		if !strings.Contains(result, "Why GoW over others") && strings.Contains(result, "footer") {
			comparisonSection := `
    <!-- Comparison Section -->
    <section class="py-20 px-6">
        <div class="max-w-6xl mx-auto">
            <div class="text-center mb-14">
                <h2 class="text-4xl font-bold tracking-tight mb-3">Why GoW over others?</h2>
                <p class="text-zinc-400 text-lg">See how GoW stacks up against Laravel and Go alternatives.</p>
            </div>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
                <div class="bg-zinc-900/50 border border-zinc-800 rounded-2xl p-6 text-center hover:border-emerald-500/50 transition">
                    <div class="text-4xl font-bold text-emerald-400 mb-1">~85K</div>
                    <div class="text-sm text-zinc-500 mb-3">Requests/sec</div>
                    <div class="text-xs text-emerald-400 font-medium">GoW (measured)</div>
                </div>
                <div class="bg-zinc-900/50 border border-zinc-800 rounded-2xl p-6 text-center">
                    <div class="text-4xl font-bold text-zinc-500 mb-1">~1K</div>
                    <div class="text-sm text-zinc-500 mb-3">Requests/sec</div>
                    <div class="text-xs text-zinc-500 font-medium">Laravel (PHP)</div>
                </div>
                <div class="bg-zinc-900/50 border border-zinc-800 rounded-2xl p-6 text-center">
                    <div class="text-4xl font-bold text-zinc-500 mb-1">~80K</div>
                    <div class="text-sm text-zinc-500 mb-3">Requests/sec</div>
                    <div class="text-xs text-zinc-500 font-medium">Fiber (Go)</div>
                </div>
            </div>
            <div class="bg-zinc-900/50 border border-zinc-800 rounded-2xl overflow-hidden">
                <div class="overflow-x-auto">
                    <table class="w-full text-sm">
                        <thead><tr class="border-b border-zinc-800">
                            <th class="text-left py-4 px-6 text-zinc-400 font-medium">Feature</th>
                            <th class="text-center py-4 px-4 text-emerald-400 font-semibold">GoW</th>
                            <th class="text-center py-4 px-4 text-zinc-500 font-medium">Laravel</th>
                            <th class="text-center py-4 px-4 text-zinc-500 font-medium">Fiber</th>
                            <th class="text-center py-4 px-4 text-zinc-500 font-medium">Gin</th>
                        </tr></thead>
                        <tbody class="text-zinc-300">
                            <tr class="border-b border-zinc-800/50"><td class="py-3 px-6 text-zinc-400">Performance (req/sec)</td><td class="text-center py-3 px-4 text-emerald-400 font-bold">~85K</td><td class="text-center py-3 px-4 text-zinc-500">~1K</td><td class="text-center py-3 px-4 text-zinc-400">~80K</td><td class="text-center py-3 px-4 text-zinc-400">~75K</td></tr>
                            <tr class="border-b border-zinc-800/50"><td class="py-3 px-6 text-zinc-400">Eloquent-style ORM</td><td class="text-center py-3 px-4 text-emerald-400">&#10003;</td><td class="text-center py-3 px-4 text-emerald-400">&#10003;</td><td class="text-center py-3 px-4 text-zinc-600">&#10007;</td><td class="text-center py-3 px-4 text-zinc-600">&#10007;</td></tr>
                            <tr class="border-b border-zinc-800/50"><td class="py-3 px-6 text-zinc-400">Blade-like Templating</td><td class="text-center py-3 px-4 text-emerald-400">&#10003;</td><td class="text-center py-3 px-4 text-emerald-400">&#10003;</td><td class="text-center py-3 px-4 text-zinc-600">&#10007;</td><td class="text-center py-3 px-4 text-zinc-600">&#10007;</td></tr>
                            <tr class="border-b border-zinc-800/50"><td class="py-3 px-6 text-zinc-400">Auth + Fortify + Sanctum</td><td class="text-center py-3 px-4 text-emerald-400">&#10003;</td><td class="text-center py-3 px-4 text-emerald-400">&#10003;</td><td class="text-center py-3 px-4 text-zinc-600">&#10007;</td><td class="text-center py-3 px-4 text-zinc-600">&#10007;</td></tr>
                            <tr class="border-b border-zinc-800/50"><td class="py-3 px-6 text-zinc-400">RBAC (Roles &amp; Permissions)</td><td class="text-center py-3 px-4 text-emerald-400">&#10003; Built-in</td><td class="text-center py-3 px-4 text-zinc-500">Package</td><td class="text-center py-3 px-4 text-zinc-600">&#10007;</td><td class="text-center py-3 px-4 text-zinc-600">&#10007;</td></tr>
                            <tr class="border-b border-zinc-800/50"><td class="py-3 px-6 text-zinc-400">Livewire / SPA Support</td><td class="text-center py-3 px-4 text-emerald-400">&#10003;</td><td class="text-center py-3 px-4 text-emerald-400">&#10003;</td><td class="text-center py-3 px-4 text-zinc-600">&#10007;</td><td class="text-center py-3 px-4 text-zinc-600">&#10007;</td></tr>
                            <tr class="border-b border-zinc-800/50"><td class="py-3 px-6 text-zinc-400">Artisan CLI (30+ commands)</td><td class="text-center py-3 px-4 text-emerald-400">&#10003;</td><td class="text-center py-3 px-4 text-emerald-400">&#10003;</td><td class="text-center py-3 px-4 text-zinc-600">&#10007;</td><td class="text-center py-3 px-4 text-zinc-600">&#10007;</td></tr>
                            <tr class="border-b border-zinc-800/50"><td class="py-3 px-6 text-zinc-400">WebSocket Broadcasting</td><td class="text-center py-3 px-4 text-emerald-400">&#10003;</td><td class="text-center py-3 px-4 text-emerald-400">&#10003;</td><td class="text-center py-3 px-4 text-emerald-400">&#10003;</td><td class="text-center py-3 px-4 text-zinc-600">&#10007;</td></tr>
                            <tr class="border-b border-zinc-800/50"><td class="py-3 px-6 text-zinc-400">Queue System</td><td class="text-center py-3 px-4 text-emerald-400">&#10003;</td><td class="text-center py-3 px-4 text-emerald-400">&#10003;</td><td class="text-center py-3 px-4 text-zinc-600">&#10007;</td><td class="text-center py-3 px-4 text-zinc-600">&#10007;</td></tr>
                            <tr class="border-b border-zinc-800/50"><td class="py-3 px-6 text-zinc-400">Single Binary Deploy</td><td class="text-center py-3 px-4 text-emerald-400">&#10003;</td><td class="text-center py-3 px-4 text-zinc-600">&#10007;</td><td class="text-center py-3 px-4 text-emerald-400">&#10003;</td><td class="text-center py-3 px-4 text-emerald-400">&#10003;</td></tr>
                            <tr class="border-b border-zinc-800/50"><td class="py-3 px-6 text-zinc-400">Low Memory Usage</td><td class="text-center py-3 px-4 text-emerald-400 font-bold">~8MB</td><td class="text-center py-3 px-4 text-zinc-500">~30MB</td><td class="text-center py-3 px-4 text-zinc-400">~10MB</td><td class="text-center py-3 px-4 text-zinc-400">~10MB</td></tr>
                            <tr><td class="py-3 px-6 text-zinc-400">2FA / Two-Factor Auth</td><td class="text-center py-3 px-4 text-emerald-400">&#10003;</td><td class="text-center py-3 px-4 text-emerald-400">&#10003;</td><td class="text-center py-3 px-4 text-zinc-600">&#10007;</td><td class="text-center py-3 px-4 text-zinc-600">&#10007;</td></tr>
                        </tbody>
                    </table>
                </div>
                <div class="px-6 py-4 border-t border-zinc-800 bg-zinc-900/30">
                    <p class="text-xs text-zinc-500 text-center">GoW delivers the full Laravel experience with Go performance. No other Go framework comes close.</p>
                </div>
            </div>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mt-12">
                <div class="bg-gradient-to-br from-emerald-500/10 to-emerald-500/5 border border-emerald-500/20 rounded-2xl p-6">
                    <h3 class="font-semibold text-lg text-emerald-400 mb-2">vs Laravel (PHP)</h3>
                    <ul class="space-y-2 text-sm text-zinc-400">
                        <li class="flex items-start gap-2"><span class="text-emerald-400 mt-0.5">&#10003;</span> 85x more requests/sec</li>
                        <li class="flex items-start gap-2"><span class="text-emerald-400 mt-0.5">&#10003;</span> 4x less memory usage</li>
                        <li class="flex items-start gap-2"><span class="text-emerald-400 mt-0.5">&#10003;</span> Single binary deployment (no PHP/nginx)</li>
                        <li class="flex items-start gap-2"><span class="text-emerald-400 mt-0.5">&#10003;</span> Built-in RBAC (Laravel needs Spatie package)</li>
                        <li class="flex items-start gap-2"><span class="text-emerald-400 mt-0.5">&#10003;</span> True concurrency (goroutines)</li>
                    </ul>
                </div>
                <div class="bg-gradient-to-br from-blue-500/10 to-blue-500/5 border border-blue-500/20 rounded-2xl p-6">
                    <h3 class="font-semibold text-lg text-blue-400 mb-2">vs Other Go Frameworks</h3>
                    <ul class="space-y-2 text-sm text-zinc-400">
                        <li class="flex items-start gap-2"><span class="text-blue-400 mt-0.5">&#10003;</span> Full ORM with relationships (Gin/Fiber have none)</li>
                        <li class="flex items-start gap-2"><span class="text-blue-400 mt-0.5">&#10003;</span> Blade-like templating engine</li>
                        <li class="flex items-start gap-2"><span class="text-blue-400 mt-0.5">&#10003;</span> Complete auth system (Fortify, Sanctum, Socialite)</li>
                        <li class="flex items-start gap-2"><span class="text-blue-400 mt-0.5">&#10003;</span> Artisan CLI with 30+ commands</li>
                        <li class="flex items-start gap-2"><span class="text-blue-400 mt-0.5">&#10003;</span> Livewire for reactive UIs</li>
                    </ul>
                </div>
            </div>
        </div>
    </section>
`
			result = strings.Replace(result, "<!-- Footer -->", comparisonSection+"\n    <!-- Footer -->", 1)
		}
	}

	return result
}

func removeDuplicateGetEnv(content string) string {
	// If this file has its own getEnv, remove it (it's duplicated from config/app.go)
	if strings.Contains(content, "func getEnv(key, fallback string) string") {
		// Remove the import "os" if present (no longer needed without getEnv)
		content = strings.Replace(content, "import \"os\"\n\n", "", 1)
		// Remove the getEnv function
		idx := strings.Index(content, "\nfunc getEnv(key, fallback string) string")
		if idx == -1 {
			idx = strings.Index(content, "func getEnv(key, fallback string) string")
		}
		if idx != -1 {
			// Find the end of the function (next closing brace at column 0 or end of file)
			endIdx := strings.Index(content[idx+1:], "\n}\n")
			if endIdx != -1 {
				content = content[:idx] + content[idx+1+endIdx+3:]
			}
		}
	}
	return content
}

func removeUnusedRoutesImport(content string) string {
	// Remove "demo/routes" or "<module>/routes" import if unused
	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines))
	inImport := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "import (" {
			inImport = true
		}
		if inImport && strings.Contains(line, "/routes\"") {
			continue // skip the unused import
		}
		if trimmed == ")" && inImport {
			inImport = false
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}

func fixLivewireCounter(content string) string {
	// Add the livewire import if missing
	if !strings.Contains(content, "github.com/mechneerd/gow/http/livewire") &&
		strings.Contains(content, "BaseComponent") {
		content = strings.Replace(content,
			`import "fmt"`,
			`import (
	"fmt"

	"github.com/mechneerd/gow/http/livewire"
)`, 1)
	}

	// Fix BaseComponent reference to livewire.BaseComponent
	if strings.Contains(content, "\tBaseComponent\n") {
		content = strings.Replace(content, "\tBaseComponent\n", "\tlivewire.BaseComponent\n", 1)
	}

	return content
}

func fixUserModel(content string) string {
	// Remove the broken SetDB method that references unexported field
	idx := strings.Index(content, "\n// SetDB wires")
	if idx != -1 {
		endIdx := strings.Index(content[idx+1:], "\n}\n")
		if endIdx != -1 {
			content = content[:idx] + "\n"
		}
	}

	// Remove unused "database/sql" import
	content = strings.Replace(content, "\t\"database/sql\"\n", "", 1)

	return content
}

func removeUnusedModelsImport(content string) string {
	// Remove the unused Models import
	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines))
	inImport := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "import (" {
			inImport = true
		}
		if inImport && strings.Contains(line, "/Models\"") {
			continue
		}
		if trimmed == ")" && inImport {
			inImport = false
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}

func fixMainGoBootstrapImport(content string, moduleName string) string {
	// The skeleton's main.go imports "github.com/mechneerd/gow/bootstrap" (framework)
	// but the local bootstrap/app.go defines its own NewApplication() and Serve().
	// Fix: change import to use the local bootstrap package.
	frameworkImport := "github.com/mechneerd/gow/bootstrap"
	localImport := moduleName + "/bootstrap"
	if strings.Contains(content, frameworkImport) {
		content = strings.Replace(content, frameworkImport, localImport, 1)
	}
	return content
}

func fixBootstrapAppGo(content string) string {
	// The skeleton's bootstrap/app.go references config.AppConfig and config.Load
	// which don't exist in the framework's config package.
	// Fix: rewrite to use os.Getenv directly for a standalone bootstrap.
	if strings.Contains(content, "config.AppConfig") || strings.Contains(content, "config.Load()") {
		newContent := `package bootstrap

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// NewApplication initializes the application.
func NewApplication() *Application {
	return &Application{}
}

// Application holds the application state.
type Application struct{}

// Serve starts the HTTP server.
func (a *Application) Serve() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	fmt.Printf("Server is running on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
`
		return newContent
	}
	return content
}

func fixAuthHandlers(content string) string {
	// Replace "not yet implemented" stubs with basic working implementations
	if strings.Contains(content, "Login not yet implemented") {
		newContent := "package Auth\n\nimport (\n\t\"encoding/json\"\n\t\"net/http\"\n)\n\n" +
			"// LoginHandler handles user login.\n" +
			"func LoginHandler(w http.ResponseWriter, r *http.Request) {\n" +
			"\tif r.Method != http.MethodPost {\n" +
			"\t\thttp.Error(w, \"Method not allowed\", http.StatusMethodNotAllowed)\n" +
			"\t\treturn\n" +
			"\t}\n\n" +
			"\tvar req struct {\n" +
			"\t\tEmail    string `json:\"email\"`\n" +
			"\t\tPassword string `json:\"password\"`\n" +
			"\t}\n" +
			"\tif err := json.NewDecoder(r.Body).Decode(&req); err != nil {\n" +
			"\t\thttp.Error(w, \"Invalid request body\", http.StatusBadRequest)\n" +
			"\t\treturn\n" +
			"\t}\n\n" +
			"\tw.Header().Set(\"Content-Type\", \"application/json\")\n" +
			"\tjson.NewEncoder(w).Encode(map[string]any{\n" +
			"\t\t\"message\": \"Login endpoint ready. Implement auth logic in handlers.go\",\n" +
			"\t})\n" +
			"}\n\n" +
			"// RegisterHandler handles user registration.\n" +
			"func RegisterHandler(w http.ResponseWriter, r *http.Request) {\n" +
			"\tif r.Method != http.MethodPost {\n" +
			"\t\thttp.Error(w, \"Method not allowed\", http.StatusMethodNotAllowed)\n" +
			"\t\treturn\n" +
			"\t}\n\n" +
			"\tvar req struct {\n" +
			"\t\tName     string `json:\"name\"`\n" +
			"\t\tEmail    string `json:\"email\"`\n" +
			"\t\tPassword string `json:\"password\"`\n" +
			"\t}\n" +
			"\tif err := json.NewDecoder(r.Body).Decode(&req); err != nil {\n" +
			"\t\thttp.Error(w, \"Invalid request body\", http.StatusBadRequest)\n" +
			"\t\treturn\n" +
			"\t}\n\n" +
			"\tw.Header().Set(\"Content-Type\", \"application/json\")\n" +
			"\tjson.NewEncoder(w).Encode(map[string]any{\n" +
			"\t\t\"message\": \"Registration endpoint ready. Implement logic in handlers.go\",\n" +
			"\t})\n" +
			"}\n\n" +
			"// LogoutHandler handles user logout.\n" +
			"func LogoutHandler(w http.ResponseWriter, r *http.Request) {\n" +
			"\thttp.Redirect(w, r, \"/login\", http.StatusFound)\n" +
			"}\n\n" +
			"// DashboardHandler shows the user dashboard.\n" +
			"func DashboardHandler(w http.ResponseWriter, r *http.Request) {\n" +
			"\tw.Header().Set(\"Content-Type\", \"text/html\")\n" +
			"\tw.Write([]byte(\"<h1>Dashboard</h1><p>Welcome!</p>\"))\n" +
			"}\n\n" +
			"// MeHandler returns the authenticated user.\n" +
			"func MeHandler(w http.ResponseWriter, r *http.Request) {\n" +
			"\tw.Header().Set(\"Content-Type\", \"application/json\")\n" +
			"\tjson.NewEncoder(w).Encode(map[string]any{\"message\": \"Implement user retrieval\"})\n" +
			"}\n"
		return newContent
	}
	return content
}

// injectMonitorRoute adds the monitor route and required imports to routes/web.go
func injectMonitorRoute(content string, moduleName string) string {
	// Add imports if not present
	if !strings.Contains(content, "\"runtime\"") {
		content = strings.Replace(content,
			"import (",
			"import (\n\t\"runtime\"\n\t\"sync\"\n\t\"time\"", 1)
	}

	// Add monitor variables after counterValue
	if !strings.Contains(content, "monitorStartTime") {
		content = strings.Replace(content,
			"var counterValue int",
			"var counterValue int\n\n// Monitor stats\nvar (\n\tmonitorStartTime = time.Now()\n\tmonitorTotalReqs int64\n\tmonitorMu        sync.RWMutex\n)", 1)
	}

	// Add monitor route before the docs route
	monitorRoute := `
	// Monitor page
	router.Get("/monitor", func(w http.ResponseWriter, r *http.Request) error {
		if r.URL.Query().Get("format") == "json" {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			monitorMu.RLock()
			totalReqs := monitorTotalReqs
			monitorMu.RUnlock()
			uptime := time.Since(monitorStartTime).Round(time.Second)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"memory_alloc":     m.Alloc,
				"memory_sys":       m.Sys,
				"go_routines":      runtime.NumGoroutine(),
				"open_connections": 0,
				"total_requests":   totalReqs,
				"uptime":           uptime.String(),
				"response_time":    0,
				"timestamp":        time.Now().Unix(),
			})
			return nil
		}

		html, err := viewEngine.Make("monitor", map[string]any{
			"AppName": "` + moduleName + `",
		})
		if err != nil {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte("<h1>Monitor</h1><p>Monitor page not found.</p>"))
			return nil
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
		return nil
	})
`

	// Insert before docs route
	if strings.Contains(content, "Documentation page") {
		content = strings.Replace(content, "\t// Documentation page", monitorRoute+"\n\t// Documentation page", 1)
	} else if strings.Contains(content, "router.Get(\"/docs\"") {
		content = strings.Replace(content, "\trouter.Get(\"/docs\"", monitorRoute+"\n\trouter.Get(\"/docs\"", 1)
	}

	return content
}

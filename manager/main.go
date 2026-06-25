// Manager es una consola web independiente para gestionar y monitorear
// el sistema Fixture D15. No tiene dependencias con el backend.
package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	BackendPort  = "8080"
	FrontendPort = "5173"
	ManagerPort  = "9099"
)

//go:embed dashboard.html
var dashboard embed.FS

var (
	rootDir     = findRoot()
	backendDir  = filepath.Join(rootDir, "backend")
	frontendDir = filepath.Join(rootDir, "frontend")
	backendLog  = "/tmp/fixture-backend.log"
	frontendLog = "/tmp/fixture-frontend.log"
)

// ─── Buscar raíz del proyecto ─────────────────────

func findRoot() string {
	exe, _ := os.Executable()
	dir := filepath.Dir(exe)
	for range 10 {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	wd, _ := os.Getwd()
	return wd
}

// ─── Tipos ─────────────────────────────────────────

type Status struct {
	Backend  *Process    `json:"backend"`
	Frontend *Process    `json:"frontend"`
	DB       *DBInfo     `json:"db"`
	System   *SysInfo    `json:"system"`
	Firewall *FWInfo     `json:"firewall"`
}

type Process struct {
	Running bool   `json:"running"`
	PID     int    `json:"pid,omitempty"`
	CPU     string `json:"cpu,omitempty"`
	Mem     string `json:"mem,omitempty"`
	Uptime  string `json:"uptime,omitempty"`
}

type DBInfo struct {
	Exists       bool   `json:"exists"`
	Size         string `json:"size"`
	Autos        int    `json:"autos"`
	Competencias int    `json:"competencias"`
	Archivos     int    `json:"archivos"`
}

type SysInfo struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	CPU      string `json:"cpu"`
	Mem      string `json:"mem"`
	Disk     string `json:"disk"`
	Arch     string `json:"arch"`
	Uptime   string `json:"uptime"`
}

type FWInfo struct {
	Port8080 bool `json:"8080"`
	Port5173 bool `json:"5173"`
}

type ActionResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Output  string `json:"output,omitempty"`
}

// ─── Funciones de sistema ─────────────────────────

func pidOf(port string) int {
	out, _ := sh("lsof -ti:%s 2>/dev/null | head -1", port)
	pid, _ := strconv.Atoi(strings.TrimSpace(out))
	return pid
}

func processInfo(port string) *Process {
	p := &Process{Running: false}
	pid := pidOf(port)
	if pid == 0 {
		return p
	}
	p.Running = true
	p.PID = pid
	p.CPU = trim(shs("ps -p %d -o %%cpu=", pid))
	p.Mem = trim(shs("ps -p %d -o %%mem=", pid))
	p.Uptime = trim(shs("ps -p %d -o etime=", pid))
	return p
}

func sh(cmd string, args ...any) (string, error) {
	c := exec.Command("sh", "-c", fmt.Sprintf(cmd, args...))
	c.Dir = rootDir
	out, err := c.CombinedOutput()
	return string(out), err
}

func shs(cmd string, args ...any) string {
	s, _ := sh(cmd, args...)
	return s
}

func trim(s string) string { return strings.TrimSpace(s) }

func getIP() string {
	if ip := trim(shs("ip route get 1 2>/dev/null | awk '{print $7; exit}'")); ip != "" {
		return ip
	}
	if ip := trim(shs("hostname -I 2>/dev/null | awk '{print $1}'")); ip != "" {
		return ip
	}
	return "localhost"
}

// ─── Handlers ──────────────────────────────────────

func handleStatus(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Status{
		Backend:  processInfo(BackendPort),
		Frontend: processInfo(FrontendPort),
		DB:       dbInfo(),
		System:   sysInfo(),
		Firewall: fwInfo(),
	})
}

func dbInfo() *DBInfo {
	d := &DBInfo{}
	path := filepath.Join(backendDir, "fixture.db")
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		return d
	}
	d.Exists = true
	d.Size = fmt.Sprintf("%.1fKB", float64(fi.Size())/1024)

	for _, q := range []struct {
		name string
		dest *int
	}{
		{"autos", &d.Autos},
		{"competencias", &d.Competencias},
		{"archivos_carrera", &d.Archivos},
	} {
		s := trim(shs("sqlite3 %q \"SELECT COUNT(*) FROM %s\"", path, q.name))
		*q.dest, _ = strconv.Atoi(s)
	}
	return d
}

func sysInfo() *SysInfo {
	s := &SysInfo{
		Hostname: trim(shs("hostname")),
		IP:       getIP(),
		Arch:     runtime.GOOS + "/" + runtime.GOARCH,
	}
	s.CPU = trim(shs("top -bn1 2>/dev/null | grep 'Cpu(s)' | awk '{print $2}'")) + "%"
	s.Mem = trim(shs("free -m 2>/dev/null | awk '/Mem:/ {printf \"%d%%\", $3/$2*100}'"))
	s.Disk = trim(shs("df -h / 2>/dev/null | awk 'NR==2 {print $5}'"))
	s.Uptime = trim(shs("uptime -p 2>/dev/null | sed 's/up //'"))
	return s
}

func fwInfo() *FWInfo {
	f := &FWInfo{}
	for _, p := range []struct {
		port string
		dest *bool
	}{
		{"8080", &f.Port8080},
		{"5173", &f.Port5173},
	} {
		out := trim(shs("firewall-cmd --query-port=%s/tcp 2>/dev/null", p.port))
		*p.dest = out == "yes"
	}
	return f
}

func handleAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	action := r.PathValue("action")

	var res ActionResult
	switch action {
	case "start-backend":
		res = startBackend()
	case "stop-backend":
		res = stopProcess(BackendPort, "Backend")
	case "restart-backend":
		res = stopProcess(BackendPort, "Backend")
		if res.Success {
			time.Sleep(1 * time.Second)
			res = startBackend()
		}
	case "start-frontend":
		res = startFrontend()
	case "stop-frontend":
		res = stopProcess(FrontendPort, "Frontend")
	case "restart-frontend":
		res = stopProcess(FrontendPort, "Frontend")
		if res.Success {
			time.Sleep(1 * time.Second)
			res = startFrontend()
		}
	case "build":
		res = buildBackend()
	case "db-reset":
		res = dbReset()
	case "ports":
		res = openPorts()
	default:
		res = ActionResult{Success: false, Message: "Acción desconocida: " + action}
	}
	json.NewEncoder(w).Encode(res)
}

func stopProcess(port, name string) ActionResult {
	pid := pidOf(port)
	if pid == 0 {
		return ActionResult{Success: true, Message: name + " no estaba corriendo"}
	}
	sh("kill %d", pid)
	return ActionResult{Success: true, Message: name + " detenido"}
}

func startBackend() ActionResult {
	if pidOf(BackendPort) != 0 {
		return ActionResult{Success: true, Message: "Backend ya está corriendo"}
	}
	binary := filepath.Join(backendDir, "tmp", "main")
	if _, err := os.Stat(binary); os.IsNotExist(err) {
		if out, err := sh("cd %s && go build -o ./tmp/main ./cmd/api", backendDir); err != nil {
			return ActionResult{Success: false, Message: "Error al compilar", Output: out}
		}
	}
	cmd := exec.Command(binary)
	cmd.Dir = backendDir
	logF, _ := os.Create(backendLog)
	cmd.Stdout = logF
	cmd.Stderr = logF
	if err := cmd.Start(); err != nil {
		return ActionResult{Success: false, Message: "Error al iniciar: " + err.Error()}
	}
	for range 10 {
		if pidOf(BackendPort) != 0 {
			return ActionResult{Success: true, Message: fmt.Sprintf("Backend iniciado (PID %d)", cmd.Process.Pid)}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return ActionResult{Success: false, Message: "Backend no respondió"}
}

func startFrontend() ActionResult {
	if pidOf(FrontendPort) != 0 {
		return ActionResult{Success: true, Message: "Frontend ya está corriendo"}
	}
	cmd := exec.Command("npm", "run", "dev", "--", "--host")
	cmd.Dir = frontendDir
	logF, _ := os.Create(frontendLog)
	cmd.Stdout = logF
	cmd.Stderr = logF
	if err := cmd.Start(); err != nil {
		return ActionResult{Success: false, Message: "Error al iniciar: " + err.Error()}
	}
	for range 15 {
		if pidOf(FrontendPort) != 0 {
			return ActionResult{Success: true, Message: fmt.Sprintf("Frontend iniciado (PID %d)", cmd.Process.Pid)}
		}
		time.Sleep(1 * time.Second)
	}
	return ActionResult{Success: false, Message: "Frontend no respondió"}
}

func buildBackend() ActionResult {
	out, err := sh("cd %s && go build -o ./tmp/main ./cmd/api", backendDir)
	if err != nil {
		return ActionResult{Success: false, Message: "Error de compilación", Output: out}
	}
	return ActionResult{Success: true, Message: "Backend compilado correctamente"}
}

func dbReset() ActionResult {
	pidOf(BackendPort)
	sh("kill %d", pidOf(BackendPort))
	for _, f := range []string{"fixture.db", "fixture.db-wal", "fixture.db-shm"} {
		os.Remove(filepath.Join(backendDir, f))
	}
	return ActionResult{Success: true, Message: "Base de datos eliminada. Usá Start Backend para recrearla."}
}

func openPorts() ActionResult {
	var msgs []string
	for _, p := range []string{"8080", "5173"} {
		if out, err := sh("firewall-cmd --permanent --add-port=%s/tcp 2>/dev/null && firewall-cmd --reload 2>/dev/null", p); err != nil {
			msgs = append(msgs, fmt.Sprintf("Puerto %s: error (%s)", p, strings.TrimSpace(out)))
		} else {
			msgs = append(msgs, fmt.Sprintf("Puerto %s: abierto", p))
		}
	}
	return ActionResult{Success: true, Message: strings.Join(msgs, " | ")}
}

// ─── Logs SSE ──────────────────────────────────────

func handleLogs(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "no flush", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	logPath := backendLog
	if r.URL.Query().Get("source") == "frontend" {
		logPath = frontendLog
	}

	// Enviar logs existentes
	if data, _ := os.ReadFile(logPath); len(data) > 0 {
		for _, line := range strings.Split(string(data), "\n") {
			if line != "" {
				fmt.Fprintf(w, "data: %s\n\n", line)
			}
		}
		flusher.Flush()
	}

	f, err := os.Open(logPath)
	if err != nil {
		fmt.Fprintf(w, "data: [ERROR] %v\n\n", err)
		flusher.Flush()
		return
	}
	defer f.Close()
	f.Seek(0, io.SeekEnd)

	buf := make([]byte, 4096)
	for {
		select {
		case <-r.Context().Done():
			return
		default:
			n, err := f.Read(buf)
			if err != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			for _, line := range strings.Split(string(buf[:n]), "\n") {
				if line != "" {
					fmt.Fprintf(w, "data: %s\n\n", line)
					flusher.Flush()
				}
			}
		}
	}
}

// ─── Dashboard ─────────────────────────────────────

func handleDashboard(w http.ResponseWriter, _ *http.Request) {
	data, _ := dashboard.ReadFile("dashboard.html")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

// ─── Main ──────────────────────────────────────────

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/status", handleStatus)
	mux.HandleFunc("POST /api/action/{action}", handleAction)
	mux.HandleFunc("GET /api/logs", handleLogs)
	mux.HandleFunc("GET /", handleDashboard)

	ip := getIP()
	log.Printf("🏁 Consola de Gestión Fixture D15")
	log.Printf("   Local:  http://localhost:%s", ManagerPort)
	if ip != "localhost" {
		log.Printf("   Red:    http://%s:%s", ip, ManagerPort)
	}
	log.Printf("   Puertos: Backend %s | Frontend %s", BackendPort, FrontendPort)

	if err := http.ListenAndServe(":"+ManagerPort, mux); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

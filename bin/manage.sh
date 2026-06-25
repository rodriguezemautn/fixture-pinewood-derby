#!/usr/bin/env bash
# ═══════════════════════════════════════════════════
# Gestor del Fixture D15
# Dashboard de gestión para el sistema de carreras
# pinewood derby — Destacamento 15
# ═══════════════════════════════════════════════════

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
BACKEND_DIR="$ROOT_DIR/backend"
FRONTEND_DIR="$ROOT_DIR/frontend"
BACKEND_LOG="/tmp/fixture-backend.log"

# ─── Colores ────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

# ─── Funciones auxiliares ──────────────────────
info()  { echo -e "${CYAN}[INFO]${NC} $1"; }
ok()    { echo -e "${GREEN}[OK]${NC} $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
err()   { echo -e "${RED}[ERR]${NC} $1"; }
title() { echo -e "\n${BOLD}═══ $1 ═══${NC}\n"; }

# ─── Detectar puertos ───────────────────────────
BACKEND_PORT="8080"
FRONTEND_PORT="5173"

get_backend_pid() { lsof -ti:$BACKEND_PORT 2>/dev/null | head -1 || true; }
get_frontend_pid() { lsof -ti:$FRONTEND_PORT 2>/dev/null | head -1 || true; }

backend_running() { [ -n "$(get_backend_pid)" ]; }
frontend_running() { [ -n "$(get_frontend_pid)" ]; }

get_ip() {
	ip route get 1 2>/dev/null | awk '{print $7; exit}' || hostname -I 2>/dev/null | awk '{print $1}' || echo "localhost"
}

# ═══════════════════════════════════════════════════
# COMANDOS
# ═══════════════════════════════════════════════════

cmd_status() {
	title "ESTADO DEL SISTEMA"

	local ip
	ip=$(get_ip)

	# Backend
	if backend_running; then
		local pid
		pid=$(get_backend_pid)
		local cpu mem uptime
		cpu=$(ps -p "$pid" -o %cpu= 2>/dev/null | tr -d ' ')
		mem=$(ps -p "$pid" -o %mem= 2>/dev/null | tr -d ' ')
		uptime=$(ps -p "$pid" -o etime= 2>/dev/null | tr -d ' ')
		echo -e "  ${GREEN}●${NC} Backend  →  ${BOLD}PID $pid${NC}  |  CPU: ${cpu}%  |  RAM: ${mem}%  |  uptime: ${uptime}"
		echo -e "              http://$ip:$BACKEND_PORT/health"
	else
		echo -e "  ${RED}○${NC} Backend  →  ${RED}Detenido${NC}"
	fi

	# Frontend
	if frontend_running; then
		local pid
		pid=$(get_frontend_pid)
		echo -e "  ${GREEN}●${NC} Frontend →  ${BOLD}PID $pid${NC}"
		echo -e "              http://$ip:$FRONTEND_PORT"
	else
		echo -e "  ${RED}○${NC} Frontend →  ${RED}Detenido${NC}"
	fi

	# DB
	echo ""
	local db_size
	if [ -f "$BACKEND_DIR/fixture.db" ]; then
		db_size=$(du -h "$BACKEND_DIR/fixture.db" 2>/dev/null | cut -f1)
		local autos comps
		autos=$(sqlite3 "$BACKEND_DIR/fixture.db" "SELECT COUNT(*) FROM autos" 2>/dev/null || echo "?")
		comps=$(sqlite3 "$BACKEND_DIR/fixture.db" "SELECT COUNT(*) FROM competencias" 2>/dev/null || echo "?")
		echo -e "  ${CYAN}◇${NC} DB SQLite →  ${BOLD}${db_size}B${NC}  |  ${autos} autos  |  ${comps} competencias"
	else
		echo -e "  ${YELLOW}◇${NC} DB SQLite →  ${YELLOW}No existe${NC}"
	fi

	# IPs de red
	echo ""
	echo -e "  ${BLUE}🌐${NC} Acceso red local: ${BOLD}http://$ip:$FRONTEND_PORT${NC}"
}

cmd_start() {
	title "INICIAR SERVICIOS"

	# Backend
	if backend_running; then
		ok "Backend ya está corriendo (PID $(get_backend_pid))"
	else
		info "Iniciando backend..."
		cd "$BACKEND_DIR"
		if [ ! -f "./tmp/main" ]; then
			warn "Backend no compilado, compilando..."
			go build -o ./tmp/main ./cmd/api
		fi
		nohup ./tmp/main > "$BACKEND_LOG" 2>&1 &
		local pid=$!
		sleep 2
		if kill -0 "$pid" 2>/dev/null; then
			ok "Backend iniciado (PID $pid) → http://localhost:$BACKEND_PORT"
		else
			err "El backend falló al iniciar. Revisá: $BACKEND_LOG"
			tail -5 "$BACKEND_LOG"
		fi
	fi

	# Frontend
	if frontend_running; then
		ok "Frontend ya está corriendo (PID $(get_frontend_pid))"
	else
		info "Iniciando frontend..."
		cd "$FRONTEND_DIR"
		nohup npm run dev -- --host > /tmp/fixture-frontend.log 2>&1 &
		local pid=$!
		# Esperar a que vite arranque
		for i in $(seq 1 10); do
			if curl -s "http://localhost:$FRONTEND_PORT" >/dev/null 2>&1; then
				ok "Frontend iniciado (PID $pid) → http://localhost:$FRONTEND_PORT"
				break
			fi
			sleep 1
		done
		if ! curl -s "http://localhost:$FRONTEND_PORT" >/dev/null 2>&1; then
			warn "Frontend arrancando... revisá /tmp/fixture-frontend.log"
		fi
	fi

	echo ""
	local ip
	ip=$(get_ip)
	info "Acceso desde la red: ${BOLD}http://$ip:$FRONTEND_PORT${NC}"
}

cmd_stop() {
	title "DETENER SERVICIOS"

	local pid

	pid=$(get_backend_pid)
	if [ -n "$pid" ]; then
		kill "$pid" 2>/dev/null && ok "Backend detenido (PID $pid)" || warn "No se pudo detener backend"
	else
		info "Backend no estaba corriendo"
	fi

	pid=$(get_frontend_pid)
	if [ -n "$pid" ]; then
		kill "$pid" 2>/dev/null && ok "Frontend detenido (PID $pid)" || warn "No se pudo detener frontend"
	else
		info "Frontend no estaba corriendo"
	fi
}

cmd_restart() {
	cmd_stop
	sleep 1
	cmd_start
}

cmd_logs() {
	if ! backend_running; then
		warn "Backend no está corriendo. Mostrando últimos logs..."
		if [ -f "$BACKEND_LOG" ]; then
			tail -20 "$BACKEND_LOG"
		else
			err "No hay logs de backend"
		fi
		echo ""
		info "Para ver logs del frontend: ${BOLD}tail -f /tmp/fixture-frontend.log${NC}"
		return
	fi

	title "LOGS DEL BACKEND (Ctrl+C para salir)"
	tail -f "$BACKEND_LOG"
}

cmd_build() {
	title "CONSTRUIR BACKEND"

	cd "$BACKEND_DIR"
	info "Compilando backend..."
	if go build -o ./tmp/main ./cmd/api; then
		ok "Backend compilado correctamente"
		# Verificar si está corriendo para reiniciar
		if backend_running; then
			warn "Backend corriendo — reiniciando para aplicar cambios..."
			local pid
			pid=$(get_backend_pid)
			kill "$pid" 2>/dev/null
			sleep 1
			nohup ./tmp/main > "$BACKEND_LOG" 2>&1 &
			sleep 2
			if backend_running; then
				ok "Backend reiniciado (PID $(get_backend_pid))"
			else
				err "Error al reiniciar backend"
			fi
		fi
	else
		err "Error de compilación"
		exit 1
	fi

	# Frontend
	info "Compilando frontend..."
	cd "$FRONTEND_DIR"
	if npm run build 2>/dev/null; then
		ok "Frontend compilado correctamente"
	else
		warn "Error en compilación del frontend (revisá dependencias)"
	fi
}

cmd_db_reset() {
	title "RESETEAR BASE DE DATOS"

	warn "⚠️  Esto ELIMINA TODOS LOS DATOS (categorías, autos, carreras, historial)"
	echo ""
	read -r -p "¿Estás seguro? (escribí 'BORRAR' para confirmar): " confirm
	if [ "$confirm" != "BORRAR" ]; then
		info "Cancelado"
		return
	fi

	# Detener backend si está corriendo
	if backend_running; then
		warn "Deteniendo backend..."
		kill "$(get_backend_pid)" 2>/dev/null || true
		sleep 1
	fi

	cd "$BACKEND_DIR"
	rm -f fixture.db fixture.db-wal fixture.db-shm
	ok "Base de datos eliminada"

	info "Reiniciando backend para recrear DB..."
	cmd_start
}

cmd_ports() {
	title "PUERTOS DEL FIREWALL"

	if ! command -v firewall-cmd &>/dev/null; then
		warn "firewall-cmd no disponible. ¿Estás en Rocky/RHEL?"
		return
	fi

	if ! sudo -n true 2>/dev/null; then
		err "Se necesita sudo para abrir puertos en el firewall"
		echo ""
		echo "  Ejecutá manualmente:"
		echo "    sudo firewall-cmd --permanent --add-port=$BACKEND_PORT/tcp"
		echo "    sudo firewall-cmd --permanent --add-port=$FRONTEND_PORT/tcp"
		echo "    sudo firewall-cmd --reload"
		return
	fi

	for port in "$BACKEND_PORT" "$FRONTEND_PORT"; do
		if sudo firewall-cmd --query-port="$port/tcp" &>/dev/null; then
			ok "Puerto $port ya está abierto"
		else
			info "Abriendo puerto $port..."
			sudo firewall-cmd --permanent --add-port="$port/tcp"
			sudo firewall-cmd --reload
			ok "Puerto $port abierto"
		fi
	done

	local ip
	ip=$(get_ip)
	echo ""
	info "Ahora deberías poder acceder desde la red:"
	echo -e "  Frontend: ${BOLD}http://$ip:$FRONTEND_PORT${NC}"
	echo -e "  Backend:  ${BOLD}http://$ip:$BACKEND_PORT${NC}"
}

cmd_test() {
	title "TEST RÁPIDO DEL SISTEMA"

	if ! backend_running; then
		err "Backend no está corriendo. Ejecutá: ${BOLD}$0 start${NC}"
		return
	fi

	info "Verificando health endpoint..."
	if curl -sf "http://localhost:$BACKEND_PORT/health" >/dev/null; then
		ok "Health check OK"
	else
		err "Health check falló"
		return
	fi

	info "Probando login..."
	local login_res token
	login_res=$(curl -sf -X POST "http://localhost:$BACKEND_PORT/api/auth/login" \
		-H "Content-Type: application/json" \
		-d '{"pin":"1234"}' 2>/dev/null || true)

	if [ -z "$login_res" ]; then
		err "Login falló (sin respuesta del backend)"
		return
	fi

	token=$(echo "$login_res" | python3 -c "import sys,json;print(json.load(sys.stdin).get('token',''))" 2>/dev/null || true)

	if [ -n "$token" ]; then
		ok "Login OK (token: ${token:0:20}...)"
	else
		err "Login falló: $(echo "$login_res" | python3 -c "import sys,json;print(json.load(sys.stdin).get('error','desconocido'))" 2>/dev/null)"
		return
	fi

	info "Listando categorías..."
	if curl -sf "http://localhost:$BACKEND_PORT/api/categorias" >/dev/null; then
		ok "API categorías OK"
	else
		err "API categorías falló"
	fi

	info "Frontend..."
	if frontend_running; then
		ok "Frontend corriendo en http://localhost:$FRONTEND_PORT"
	else
		warn "Frontend no está corriendo"
	fi

	echo ""
	info "Sistema funcionando correctamente ✅"
}

cmd_help() {
	title "GESTOR DEL FIXTURE D15"
	echo "  Uso: $0 [comando]"
	echo ""
	echo "  Comandos:"
	echo "    status     Ver estado de backend y frontend"
	echo "    start      Iniciar backend y frontend"
	echo "    stop       Detener backend y frontend"
	echo "    restart    Reiniciar servicios"
	echo "    logs       Ver logs del backend en tiempo real"
	echo "    build      Compilar backend y frontend"
	echo "    db-reset   Eliminar y recrear la base de datos"
	echo "    ports      Abrir puertos en el firewall"
	echo "    test       Verificar que el sistema funciona"
	echo "    menu       Menú interactivo (predeterminado)"
	echo "    help       Mostrar esta ayuda"
}

# ═══════════════════════════════════════════════════
# MENÚ INTERACTIVO
# ═══════════════════════════════════════════════════

cmd_menu() {
	local whiptail_available=false
	command -v whiptail &>/dev/null && whiptail_available=true

	while true; do
		if $whiptail_available; then
			# ── Menú gráfico con whiptail ──────────
			local running_icon
			backend_running && running_icon="●" || running_icon="○"
			local f_running_icon
			frontend_running && f_running_icon="●" || f_running_icon="○"

			local choice
			choice=$(whiptail --title "🏁 Gestor Fixture D15" \
				--menu "
  Backend: $running_icon  |  Frontend: $f_running_icon
  ───────────────────────────────" \
				18 56 9 \
				"status"   "📊  Ver estado del sistema" \
				"start"    "▶️   Iniciar backend + frontend" \
				"stop"     "⏹️   Detener servicios" \
				"restart"  "🔄  Reiniciar servicios" \
				"logs"     "📋  Ver logs del backend" \
				"build"    "🔨  Compilar backend + frontend" \
				"db-reset" "🗑️   Resetear base de datos" \
				"ports"    "🌐  Abrir puertos en firewall" \
				"test"     "🧪  Test rápido del sistema" \
				"salir"    "🚪  Salir" \
				3>&1 1>&2 2>&3)

			case "$choice" in
				"") break ;;
				"salir") break ;;
				status|start|stop|restart|logs|build|db-reset|ports|test)
					clear
					cmd_$choice
					echo ""
					read -r -p "Presioná Enter para volver al menú..."
					;;
			esac
		else
			# ── Menú de texto (fallback) ──────────
			clear
			title "GESTOR DEL FIXTURE D15"

			backend_running && echo -e "  Backend:  ${GREEN}● Corriendo${NC}" || echo -e "  Backend:  ${RED}○ Detenido${NC}"
			frontend_running && echo -e "  Frontend: ${GREEN}● Corriendo${NC}" || echo -e "  Frontend: ${RED}○ Detenido${NC}"
			echo ""
			echo "  ${BOLD}1${NC}) 📊  Ver estado"
			echo "  ${BOLD}2${NC}) ▶️   Iniciar servicios"
			echo "  ${BOLD}3${NC}) ⏹️   Detener servicios"
			echo "  ${BOLD}4${NC}) 🔄  Reiniciar servicios"
			echo "  ${BOLD}5${NC}) 📋  Ver logs del backend"
			echo "  ${BOLD}6${NC}) 🔨  Compilar backend + frontend"
			echo "  ${BOLD}7${NC}) 🗑️   Resetear base de datos"
			echo "  ${BOLD}8${NC}) 🌐  Abrir puertos firewall"
			echo "  ${BOLD}9${NC}) 🧪  Test rápido"
			echo "  ${BOLD}0${NC}) 🚪  Salir"
			echo ""
			read -r -p "  Elegí una opción: " opt

			case "$opt" in
				0|"") break ;;
				1) clear; cmd_status; echo ""; read -r -p "Enter para continuar..." ;;
				2) clear; cmd_start; echo ""; read -r -p "Enter para continuar..." ;;
				3) clear; cmd_stop; echo ""; read -r -p "Enter para continuar..." ;;
				4) clear; cmd_restart; echo ""; read -r -p "Enter para continuar..." ;;
				5) clear; cmd_logs ;;
				6) clear; cmd_build; echo ""; read -r -p "Enter para continuar..." ;;
				7) clear; cmd_db_reset; echo ""; read -r -p "Enter para continuar..." ;;
				8) clear; cmd_ports; echo ""; read -r -p "Enter para continuar..." ;;
				9) clear; cmd_test; echo ""; read -r -p "Enter para continuar..." ;;
				*) warn "Opción inválida"; sleep 1 ;;
			esac
		fi
	done
	clear
	info "¡Hasta luego! 🏁"
}

# ═══════════════════════════════════════════════════
# ENTRY POINT
# ═══════════════════════════════════════════════════

# Verificar dependencias
for cmd in curl go npm; do
	if ! command -v "$cmd" &>/dev/null; then
		err "Requerido: $cmd no está instalado"
		exit 1
	fi
done

# Si no hay comando o es 'menu', abrir menú interactivo
case "${1:-menu}" in
	menu|interactive)  cmd_menu ;;
	status|start|stop|restart|logs|build|db-reset|ports|test|help)
		shift 0
		"cmd_${1:-help}"
		;;
	*)  cmd_help ;;
esac

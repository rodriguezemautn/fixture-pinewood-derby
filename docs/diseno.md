# Documento de Diseño — Fixture D15

> **Versión**: 1.0  
> **Última actualización**: 2026-06-25  
> **Tema**: Arcade Racing 8-bit

---

## 1. Sistema de Diseño

### 1.1 Filosofía Visual

La interfaz combina la estética de los **arcades de carreras de los 80s/90s** (OutRun, Daytona USA, Need for Speed Underground) con elementos **pixel art 8-bit** y una paleta **naranja retro** cálida.

```
Inspiración visual:
  🏁 Arcade racing (Daytona, NFS)
  🕹️ Pixel art 8-bit
  🌅 Atardecer retro
  🚄 Shinkansen bala
  🏭 Estética industrial
```

### 1.2 Paleta de Colores

#### Tema Oscuro (default)

| Variable | Color | Uso |
|----------|-------|-----|
| `--arcade-black` | `#141824` | Fondo principal |
| `--arcade-dark` | `#1C2233` | Nav, headers |
| `--arcade-surface` | `#242B3D` | Cards, tablas |
| `--orange` | `#E07B0A` | 🔶 Naranja principal (acción) |
| `--orange-glow` | `#FBBF24` | Brillo, hover |
| `--orange-border` | `#A85B0A` | Bordes pixel |
| `--text-primary` | `#F5F8FC` | Texto principal |
| `--text-body` | `#D8E0EA` | Texto secundario |
| `--text-dim` | `#94A0B5` | Texto tenue |
| `--border-color` | `#3E4860` | Bordes de cards |
| `--red-race` | `#EF4444` | Peligro, eliminar |
| `--green-go` | `#34D399` | Éxito, completado |
| `--blue-neon` | `#60A5FA` | Acento NFS |

#### Tema Claro

| Variable | Color | Uso |
|----------|-------|-----|
| `--arcade-black` | `#FFF5EB` | Fondo crema anaranjado |
| `--arcade-dark` | `#FFE8D0` | Nav, headers |
| `--arcade-surface` | `#FFFFFF` | Cards blancas |
| `--orange` | `#D97706` | Naranja principal |
| `--text-primary` | `#1A1F2E` | Texto oscuro |
| `--border-color` | `#D4C5B0` | Bordes suaves |

### 1.3 Tipografía

| Uso | Fuente | Estilo |
|-----|--------|--------|
| **Logo + títulos grandes** | `Press Start 2P` | Pixel 8-bit, monospace |
| **Subtítulos** | `Black Ops One` | Bold racing, mayúsculas |
| **UI / botones / labels** | `VT323` | Pixel compacto, monospace |
| **Cuerpo de texto** | `Inter` | Sans-serif legible |

```css
.font-pixel { font-family: 'Press Start 2P', monospace; }
.font-terminal { font-family: 'VT323', monospace; }
```

### 1.4 Efectos Visuales

#### Scanlines CRT
```css
body::after {
    background: repeating-linear-gradient(
        0deg, transparent, transparent 2px,
        rgba(0,0,0,0.04) 2px, rgba(0,0,0,0.04) 4px
    );
}
```
Capa fija sobre toda la pantalla simulando líneas de monitor CRT.

#### Pixel Border
```css
.pixel-border {
    border: 2px solid var(--orange-border);
    box-shadow: 3px 3px 0 0 var(--orange-border);
}
```
Bordes con sombra descentrada simulando pixels.

#### Glow
```css
.glow-orange {
    box-shadow: 0 0 12px rgba(217,119,6,0.25), 0 0 25px rgba(217,119,6,0.08);
}
```

### 1.5 Animaciones

| Animación | Duración | Propósito |
|-----------|:--------:|-----------|
| `pulse-orange` | 2s | Latido en elementos importantes |
| `slide-up` | 0.5s | Entrada de secciones |
| `zoom-in` | 0.3s | Apertura de modales |
| `pixel-flicker` | 4s | Parpadeo retro en títulos |
| `crt-on` | 0.6s | Efecto de encendido CRT |
| `speed-lines` | 0.5s | Líneas de velocidad en hero |
| `confetti-fall` | 4s | Confeti de celebración |
| `wheel-spin` | 0.8s | Ruedas del tren animado |

---

## 2. Sistema de Temas

El proyecto soporta **tema oscuro** (default) y **tema claro**, seleccionables mediante un botón ☀️/🌙 en el header.

### Mecanismo

```html
<!-- +layout.svelte -->
<html data-theme="light">  <!-- o "dark" -->
```

```css
/* layout.css */
:root { /* vars del tema oscuro */ }
[data-theme="light"] { /* vars del tema claro */ }
```

La preferencia se guarda en `localStorage('arcade-theme')` y persiste entre sesiones.

---

## 3. Componentes

### 3.1 Árbol de Componentes

```
routes/
├── +layout.svelte           ← Header arcade + footer
├── +page.svelte             ← Landing con escenario CSS
├── login/+page.svelte       ← Login PIN
├── admin/
│   ├── +layout.svelte       ← Nav admin + breadcrumbs
│   └── categorias/
│       ├── +page.svelte     ← CRUD categorías
│       └── [id]/
│           ├── autos/+page.svelte   ← CRUD autos + competencias
│           └── fixture/+page.svelte ← Fixture + resultados + podio
├── carreras/
│   ├── +page.svelte         ← Lista de categorías
│   └── [id]/+page.svelte    ← Vista pública de carrera
└── autos/+page.svelte       ← Todos los autos

lib/components/
├── AutoForm.svelte          ← Modal crear/editar auto
├── CategoriaForm.svelte     ← Modal crear/editar categoría
├── Podium.svelte            ← Podio animado (🥇🥈🥉)
└── Celebration.svelte       ← Pantalla de celebración con confeti
```

### 3.2 Patrón Modal

Todos los modales siguen el mismo patrón:

```html
<div class="modal-overlay" onclick={onclose}>
    <div class="modal" onclick={(e) => e.stopPropagation()}>
        <h2>Título</h2>
        <form onsubmit={handleSubmit}>
            <!-- campos -->
            <div class="buttons">
                <button onclick={onclose}>Cancelar</button>
                <button type="submit">Guardar</button>
            </div>
        </form>
    </div>
</div>
```

Características:
- `max-height: 85vh; overflow-y: auto` (scroll si el contenido excede la pantalla)
- Cierra con Escape o click fuera del modal
- Transición `zoom-in` (0.3s)

---

## 4. Escenario Arcade Racing (Landing)

La landing page incluye una escena dibujada completamente con **CSS puro**:

```
🎨 Capas del escenario (de atrás hacia adelante):

1. Cielo degradado (púrpura → naranja atardecer → verde)
2. Estrellas (radial-gradient)
3. Montañas (3 capas con linear-gradient)
4. Fábrica con chimenea blanca/roja + humo animado
5. Estación de tren con cartel "BERAZATEGUI"
6. Árboles pixel art (box-shadow)
7. Tren bala Shinkansen animado (3 vagones)
8. Ruta en perspectiva (clip-path)
9. Speed lines (animación horizontal)
10. Contenido del hero (título, botones)
```

El tren bala tiene un ciclo de 20s:
1. Aparece desde la derecha
2. Frena en la estación (espera ~4s)
3. Arranca y desaparece a la izquierda

---

## 5. Responsive Design

| Breakpoint | Comportamiento |
|:----------:|----------------|
| >1024px | Layout completo, grid de 4 columnas |
| 768-1024px | Grid de 2-3 columnas |
| <640px | 1 columna, header compacto, hero reducido |
| Proyector | Fondos elevados (#14 en vez de #0A) para legibilidad |

---

## 6. Experiencia de Usuario

### Flujo Admin
```
Login → Categorías → [seleccionar] → Autos → [+Nueva Competencia]
  → Generar Fixture → Registrar resultados → Generar Final
  → Ver Podio → (Desempate si empate) → Finalizar → Archivar
```

### Flujo Público
```
Landing → Carreras → [seleccionar categoría]
  → Ver heats en vivo (polling 10s)
  → Ver posiciones y podio
  → Celebrar campeón 🎉
```

### Estados de cada componente
- **Loading**: Skeleton loaders con shimmer animation
- **Empty**: Mensaje informativo + ícono
- **Error**: Banner rojo con mensaje
- **Success**: Feedback visual con transición

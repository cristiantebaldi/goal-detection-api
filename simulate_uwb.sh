#!/bin/bash

# =============================================================================
# Simulador de Âncoras UWB — Bola Inteligente
# =============================================================================
# Simula 3 cenários:
#   1. Gol válido        — bola cruza a linha completamente
#   2. Bola na trave     — bola chega na linha mas volta
#   3. Perda de sinal    — âncoras perdem a bola por alguns frames
#
# A quadra é 40m x 20m. O gol fica no lado direito (x = 40).
# As âncoras estão nos cantos da quadra.
#
# Âncoras:
#   Âncora 1 — (0,  0,  0) — canto superior esquerdo
#   Âncora 2 — (40, 0,  0) — canto superior direito
#   Âncora 3 — (0,  20, 0) — canto inferior esquerdo
#
# Linha do gol: x = 40, y entre 8.5 e 11.5 (largura do gol = 3m)
# =============================================================================

BACKEND_URL="http://localhost:8080/sensor/position"
INTERVAL=0.1  # segundos entre cada leitura (10 leituras/segundo)

# Dimensões do campo visual (colunas x linhas)
FIELD_COLS=42
FIELD_ROWS=11

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

# =============================================================================
# Campo ASCII
# =============================================================================

# Converte coordenadas reais (x: 0-40, y: 0-20) para posição no campo visual
# Uso: to_col x  → coluna (0 a FIELD_COLS-1)
#      to_row y  → linha  (0 a FIELD_ROWS-1)
to_col() {
  echo "scale=0; ($1 * ($FIELD_COLS - 1) / 40) / 1" | bc
}

to_row() {
  echo "scale=0; ($1 * ($FIELD_ROWS - 1) / 20) / 1" | bc
}

# Desenha o campo com a bola na posição (bx, by) real
# Uso: draw_field bx by sinal_perdido
draw_field() {
  local bx=$1
  local by=$2
  local lost=${3:-0}

  local ball_col=$(to_col "$bx")
  local ball_row=$(to_row "$by")

  # Linha do gol começa em y=8.5 (row~4) até y=11.5 (row~6) dos 11 rows
  local goal_top=4
  local goal_bot=6

  # Limpa as linhas do campo antes de redesenhar (FIELD_ROWS + 2 bordas + 1 legenda)
  local total_lines=$((FIELD_ROWS + 3))
  for i in $(seq 1 $total_lines); do
    echo -e "\033[1A\033[2K"
  done

  # Borda superior
  printf "  +"
  printf '%0.s-' $(seq 1 $FIELD_COLS)
  printf "+\n"

  # Linhas do campo
  for row in $(seq 0 $((FIELD_ROWS - 1))); do
    printf "  |"
    for col in $(seq 0 $((FIELD_COLS - 1))); do
      local char=" "

      # Posição da bola
      if [ "$col" -eq "$ball_col" ] && [ "$row" -eq "$ball_row" ]; then
        if [ "$lost" -eq 1 ]; then
          printf "${RED}?${NC}"
        else
          printf "${CYAN}●${NC}"
        fi
        continue
      fi

      # Linha do meio (x=20 → col=21)
      if [ "$col" -eq 21 ]; then
        char="┆"
      fi

      # Linha do gol (coluna final)
      if [ "$col" -eq $((FIELD_COLS - 1)) ]; then
        if [ "$row" -ge "$goal_top" ] && [ "$row" -le "$goal_bot" ]; then
          printf "${GREEN}|${NC}"
          continue
        else
          char="│"
        fi
      fi

      # Âncoras nos cantos (col=0 ou col=41, row=0 ou row=10)
      if { [ "$col" -eq 0 ] || [ "$col" -eq $((FIELD_COLS - 1)) ]; } && \
         { [ "$row" -eq 0 ] || [ "$row" -eq $((FIELD_ROWS - 1)) ]; }; then
        printf "${YELLOW}A${NC}"
        continue
      fi

      printf "%s" "$char"
    done
    printf "|\n"
  done

  # Borda inferior
  printf "  +"
  printf '%0.s-' $(seq 1 $FIELD_COLS)
  printf "+\n"

  # Legenda
  if [ "$lost" -eq 1 ]; then
    echo -e "  ${CYAN}●${NC} bola  ${GREEN}|${NC} gol  ${YELLOW}A${NC} âncora  ${RED}?${NC} sinal perdido   pos: ${RED}(?, ?)${NC}"
  else
    local bx_fmt=$(printf "%.1f" "$bx" 2>/dev/null || echo "$bx")
    local by_fmt=$(printf "%.1f" "$by" 2>/dev/null || echo "$by")
    echo -e "  ${CYAN}●${NC} bola  ${GREEN}|${NC} gol  ${YELLOW}A${NC} âncora                    pos: ${CYAN}(${bx_fmt}, ${by_fmt})${NC}"
  fi
}

# Imprime linhas em branco para reservar o espaço do campo na primeira vez
init_field() {
  local total_lines=$((FIELD_ROWS + 3))
  for i in $(seq 1 $total_lines); do
    echo ""
  done
}

# =============================================================================
# Funções utilitárias
# =============================================================================

distance() {
  echo "scale=4; sqrt(($1-$4)^2 + ($2-$5)^2 + ($3-$6)^2)" | bc
}

send_reading() {
  local x=$1
  local y=$2
  local z=$3
  local label=$4

  local a1x=0  a1y=0  a1z=0
  local a2x=40 a2y=0  a2z=0
  local a3x=0  a3y=20 a3z=0

  local d1=$(distance $a1x $a1y $a1z $x $y $z)
  local d2=$(distance $a2x $a2y $a2z $x $y $z)
  local d3=$(distance $a3x $a3y $a3z $x $y $z)

  local payload=$(cat <<EOF
{
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%S.%3NZ")",
  "positions": [
    { "anchor_id": "anchor-1", "distance_meters": $d1 },
    { "anchor_id": "anchor-2", "distance_meters": $d2 },
    { "anchor_id": "anchor-3", "distance_meters": $d3 }
  ]
}
EOF
)

  local http_status=$(curl -s -o /dev/null -w "%{http_code}" \
    -X POST "$BACKEND_URL" \
    -H "Content-Type: application/json" \
    -d "$payload")

  draw_field "$x" "$y" 0
  if [ -n "$label" ]; then
    echo -e "  $label"
  fi
}

interpolate() {
  echo "scale=4; $1 + ($2 - $1) * $4 / $3" | bc
}

# =============================================================================
# Cenário 1: Gol válido
# =============================================================================
scenario_goal() {
  echo ""
  echo -e "${BOLD}${GREEN}=== CENÁRIO 1: GOL VÁLIDO ===${NC}"
  echo -e "${GREEN}Bola parte do centro e cruza a linha do gol completamente${NC}"
  echo ""

  init_field

  local steps=30
  local start_x=20
  local end_x=41
  local y=10
  local z=0.3

  for i in $(seq 0 $steps); do
    local x=$(interpolate $start_x $end_x $steps $i)

    local label=""
    local crossed=$(echo "$x > 40" | bc)
    if [ "$crossed" = "1" ]; then
      label="${GREEN}⚽ CRUZOU A LINHA — GOL!${NC}"
    fi

    send_reading "$x" "$y" "$z" "$label"
    sleep $INTERVAL
  done

  echo ""
  echo -e "${GREEN}✓ Trajetória de gol enviada${NC}"
}

# =============================================================================
# Cenário 2: Bola na trave
# =============================================================================
scenario_crossbar() {
  echo ""
  echo -e "${BOLD}${YELLOW}=== CENÁRIO 2: BOLA NA TRAVE ===${NC}"
  echo -e "${YELLOW}Bola chega até x=39.8 mas volta — não é gol${NC}"
  echo ""

  init_field

  local steps=20
  local start_x=20
  local peak_x=39.8
  local y=10
  local z=0.3

  for i in $(seq 0 $steps); do
    local x=$(interpolate $start_x $peak_x $steps $i)
    send_reading "$x" "$y" "$z" ""
    sleep $INTERVAL
  done

  for i in $(seq 0 $steps); do
    local x=$(interpolate $peak_x $start_x $steps $i)
    send_reading "$x" "$y" "$z" ""
    sleep $INTERVAL
  done

  echo ""
  echo -e "${YELLOW}✓ Trajetória de trave enviada — backend não deve registrar gol${NC}"
}

# =============================================================================
# Cenário 3: Perda de sinal
# =============================================================================
scenario_signal_loss() {
  echo ""
  echo -e "${BOLD}${RED}=== CENÁRIO 3: PERDA DE SINAL ===${NC}"
  echo -e "${RED}Bola some por ~1 segundo e reaparece — não pode gerar gol fantasma${NC}"
  echo ""

  init_field

  local y=10
  local z=0.3
  local last_x=20

  for x in 5 8 11 14 17 20; do
    send_reading "$x" "$y" "$z" ""
    last_x=$x
    sleep $INTERVAL
  done

  # Perda de sinal — mantém a última posição conhecida com ? no campo
  for i in $(seq 1 10); do
    local payload=$(cat <<EOF
{
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%S.%3NZ")",
  "signal_lost": true,
  "positions": []
}
EOF
)
    local http_status=$(curl -s -o /dev/null -w "%{http_code}" \
      -X POST "$BACKEND_URL" \
      -H "Content-Type: application/json" \
      -d "$payload")

    draw_field "$last_x" "$y" 1
    echo -e "  ${RED}⚠ sinal perdido — frame $i/10 | HTTP ${http_status}${NC}"
    sleep $INTERVAL
  done

  for x in 22 25 28 30; do
    send_reading "$x" "$y" "$z" ""
    sleep $INTERVAL
  done

  echo ""
  echo -e "${RED}✓ Perda de sinal enviada — backend não deve registrar gol${NC}"
}

# =============================================================================
# Main
# =============================================================================
print_header() {
  echo ""
  echo -e "${BOLD}╔══════════════════════════════════════════╗${NC}"
  echo -e "${BOLD}║     Simulador UWB — Bola Inteligente     ║${NC}"
  echo -e "${BOLD}╚══════════════════════════════════════════╝${NC}"
  echo ""
  echo -e "  Backend: ${CYAN}${BACKEND_URL}${NC}"
  echo -e "  Quadra:  ${CYAN}40m x 20m${NC}"
  echo -e "  Gol:     ${CYAN}x=40, y entre 8.5 e 11.5${NC}"
  echo ""
}

main() {
  print_header

  echo -e "  Verificando backend..."
  local health=$(curl -s -o /dev/null -w "%{http_code}" --connect-timeout 2 "$BACKEND_URL" 2>/dev/null || echo "000")
  if [ "$health" = "000" ]; then
    echo -e "  ${RED}⚠ Backend não respondeu em ${BACKEND_URL}${NC}"
    echo -e "  ${YELLOW}Continuando mesmo assim — verifique os logs${NC}"
  else
    echo -e "  ${GREEN}✓ Backend acessível (HTTP ${health})${NC}"
  fi

  scenario_goal

  echo ""
  echo -e "${BLUE}Aguardando 2s antes do próximo cenário...${NC}"
  sleep 2

  scenario_crossbar

  echo ""
  echo -e "${BLUE}Aguardando 2s antes do próximo cenário...${NC}"
  sleep 2

  scenario_signal_loss

  echo ""
  echo -e "${BOLD}${GREEN}════════════════════════════════════════${NC}"
  echo -e "${BOLD}${GREEN}  Simulação completa!${NC}"
  echo -e "${BOLD}${GREEN}════════════════════════════════════════${NC}"
  echo ""
  echo -e "  Resultados esperados no backend:"
  echo -e "  ${GREEN}✓ Cenário 1 → SSE deve disparar${NC}"
  echo -e "  ${YELLOW}✓ Cenário 2 → SSE não deve disparar (trave)${NC}"
  echo -e "  ${RED}✓ Cenário 3 → SSE não deve disparar (sinal perdido)${NC}"
  echo ""
}

main "$@"
# Bola Inteligente
Documento de Contexto — Backend API

---

## 1. Visão Geral

Sistema embarcado de detecção de gol para futsal, utilizando tecnologia UWB (Ultra-Wideband) para rastreamento de alta precisão da bola em tempo real.

O objetivo principal é identificar com precisão quando a bola cruza completamente a linha do gol e notificar imediatamente o árbitro através de um aplicativo.

---

## 2. Arquitetura

### Fluxo principal

- Chip UWB na bola emite sinais continuamente
- Âncoras UWB fixadas na quadra captam o sinal e medem a distância até a bola
- Gateway local (Raspberry Pi) recebe os dados das âncoras via rede local
- Backend Go processa os dados, calcula a posição via trilateração e detecta o gol
- App do árbitro recebe a notificação via SSE (Server-Sent Events)

### Decisões arquiteturais

- Todo o sistema roda localmente — sem dependência de internet
- Gateway e backend podem rodar no mesmo Raspberry Pi
- Comunicação com o app via SSE (unidirecional — servidor notifica cliente)
- WebSocket descartado pois o app do árbitro não precisa enviar dados pelo canal
- gRPC descartado pelo suporte limitado em apps móveis
- Sem gerenciamento de partidas — o sistema só detecta e notifica o gol
- Sem banco de dados — posições são processadas em memória e descartadas
- Posições das âncoras ficam em configuração (config.yaml)

### Ambiente de desenvolvimento

- Simulador shell (`simulate_uwb.sh`) substitui o hardware real durante o desenvolvimento
- Gateway simulado roda na porta `8081`, backend principal na porta `8080`
- No hardware real, `localhost:8081` é substituído pelo IP do Raspberry Pi

---

## 3. Rotas da API

Duas rotas — uma por responsabilidade do sistema.

### POST `/sensor/position`

Endpoint principal. Recebe as leituras brutas das âncoras UWB (ou do simulador) e aciona o pipeline de processamento.

**Payload esperado:**

```json
{
  "timestamp": "2026-09-07T12:00:00.000Z",
  "positions": [
    { "anchor_id": "anchor-1", "distance_meters": 22.36 },
    { "anchor_id": "anchor-2", "distance_meters": 18.50 },
    { "anchor_id": "anchor-3", "distance_meters": 25.10 }
  ]
}
```

---

### GET `/goal/event`

Stream SSE. O app do árbitro conecta neste endpoint e mantém a conexão aberta. Quando um gol é detectado, o servidor envia o evento imediatamente sem necessidade de polling.

**Evento enviado quando gol é detectado:**

```
data: { "event": "goal", "detected_date": "2026-09-07T12:00:00.000Z" }
```

---

## 4. Pipeline de Processamento

Fluxo interno acionado a cada `POST /sensor/position`:

```
POST /sensor/position → Positioning Engine → Goal Detector → SSE
```

- **Positioning Engine** — calcula posição (x, y, z) via trilateração com as distâncias das âncoras
- **Goal Detector** — verifica se `x ≥ 40` e `y` está entre `8.5` e `11.5` (dentro do gol)
- **Se gol detectado** — dispara evento no canal SSE

---

## 5. Configuração (config.yaml)

As posições físicas das âncoras e os limites do gol ficam em arquivo de configuração. Como são fixas na quadra, não há necessidade de CRUD.

```yaml
anchors:
  - id: anchor-1
    x: 0.0
    y: 0.0
    z: 0.0
  - id: anchor-2
    x: 40.0
    y: 0.0
    z: 0.0
  - id: anchor-3
    x: 0.0
    y: 20.0
    z: 0.0

goal:
  x: 40.0
  y_min: 8.5
  y_max: 11.5
```

---

## 6. Cenários de Teste (Simulador)

O script `simulate_uwb.sh` cobre três cenários que o backend deve tratar corretamente:

- **Gol válido** — bola cruza `x = 40` completamente → SSE deve disparar
- **Bola na trave** — bola chega em `x = 39.8` e volta → SSE não deve disparar
- **Perda de sinal** — âncoras perdem a bola por ~1s → nenhum gol fantasma deve ser gerado
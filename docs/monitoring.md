## Monitoring

Для мониторинга состояния API Gateway используется **Prometheus** и **Alertmanager**.

### Поток мониторинга и оповещений

1. **Prometheus** регулярно опрашивает эндпоинт `/metrics` API Gateway и собирает метрики.
2. **Правила алертов** в файле `/etc/prometheus/alert_rules.yml` определяют, при каких условиях должна сработать тревога (например, высокая нагрузка CPU, много 5xx ошибок и т.д.).
3. Когда условие алерта выполняется, **Prometheus отправляет событие в Alertmanager**.
4. **Alertmanager** группирует, подавляет дубликаты и маршрутизирует алерты согласно конфигурации `alertmanager.yml`:
    - критические алерты (`severity: critical`) отправляются основному администратору
    - предупреждения (`severity: warning`) отправляются вторичному администратору
5. **Alertmanager** доставляет уведомления по e-mail через SMTP (например, `smtp.gmail.com`).

Таким образом:
`Prometheus (метрики и правила) → Alertmanager (группировка и маршрутизация) → Email (оповещение).`

### Метрики

| Метрика                                                                                               | Назначение                                                   |
|-------------------------------------------------------------------------------------------------------|--------------------------------------------------------------|
| `up`                                                                                                  | Проверка доступности экземпляра сервиса                      |
| `sum(rate(http_requests_total[1m])) by (instance)`                                                    | Количество HTTP-запросов в минуту по каждому экземпляру      |
| `rate(rate_limited_requests[5m])`                                                                     | Частота запросов, отклонённых Rate Limiter                   |
| `histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le, instance))`      | 95-й перцентиль времени обработки HTTP-запросов               |
| `100 - (avg by(instance) (rate(node_cpu_seconds_total{mode="idle"}[1m])) * 100)`                       | Загрузка CPU по каждому экземпляру                           |
| `sum(rate(http_requests_total{status="429"}[5m])) by (instance)`                                      | Количество ответов с кодом 429 (Rate Limit)                  |
| `sum by(instance) (rate(http_requests_total{status=~"5.."}[5m]))`                                     | Количество 5xx-ошибок по каждому экземпляру                  |

### Алерты

Файл правил: `/etc/prometheus/alert_rules.yml`

| Alert Name              | Условие                                                                                                   | Триггер                | Severity    | Описание                                                   |
|-----------------------|---------------------------------------------------------------------------------------------------------|----------------------|-------------|-----------------------------------------------------------|
| **ServiceDown**        | `up == 0`                                                                                               | 1m                   | critical    | Сервис недоступен                                          |
| **High5xxErrors**      | `sum(rate(http_requests_total{status=~"5.."}[5m])) by (instance) > 0.1`                                  | 5m                   | critical    | Слишком много 5xx ошибок                                   |
| **High429Errors**      | `sum(rate(http_requests_total{status="429"}[5m])) by (instance) > 5`                                     | 5m                   | warning     | Пользователи упираются в лимиты (>5 req/s)                 |
| **HighLatency95**      | `histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le)) > 0.5`            | 5m                   | warning     | 95% запросов обрабатываются дольше 500ms                   |
| **HighRPS**            | `sum(rate(http_requests_total[5m])) by (instance) > 1000`                                               | 5m                   | critical    | Резкий рост RPS, возможный DDoS (>1000 req/s)              |
| **HighRateLimitedRequests** | `sum(rate(rate_limited_requests[5m])) by (instance) > 5`                                                 | 5m                    | warning     | Много отклонённых запросов (>5 req/s)                      |
| **HighCPU**            | `100 - (avg by(instance)(rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 90`                     | 5m                    | warning     | Загрузка CPU > 90%                                         |

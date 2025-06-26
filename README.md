## Scaleway SD

Tiny Go service that turns your Scaleway (Dedibox / Online.net) servers into a Prometheus HTTP Service-Discovery endpoint.*

* Calls the **Online API** (`/api/v1/server`) every *N* minutes
* Looks up each server’s first public IP
* Emits the JSON expected by [`http_sd_configs`](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#http_sd_config)
* Lets Prometheus discover new or rebuilt Dediboxes automatically – no static target lists, no file-reloader

## Quick start

```bash
# 1. build
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o scw_sd scw_sd.go

# 2. copy to server
scp scw_sd user@server:/usr/local/bin/

# 3. set token & run
export ONLINE_API_TOKEN=xxxxxx
/usr/local/bin/scw_sd
```

Visit [http://server:8000/scw-sd](http://server:8000/scw-sd) – you should see JSON targets.

---

## Prometheus config

```yaml
scrape_configs:
  - job_name: 'dedibox-node-exporter'
    http_sd_configs:
      - url: http://scw-sd:8000/scw-sd
        refresh_interval: 5m
    scrape_interval: 15s
    relabel_configs:
      - source_labels: [__address__]
        regex: '(.*):9100'
        target_label: instance
        replacement: '$1'
```

---

## Systemd unit (optional)

```ini
[Unit]
Description=Scaleway SD for Prometheus
After=network-online.target
Wants=network-online.target

[Service]
Environment=ONLINE_API_TOKEN=XXXX
ExecStart=/usr/local/bin/scw_sd
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now scw_sd
```

---

## Configuration

| Env var            | Default | Description                                    |
| ------------------ | ------- | ---------------------------------------------- |
| `ONLINE_API_TOKEN` | –       | **required** bearer token for `api.online.net` |
| `PORT`             | `8000`  | HTTP listen port (optional)                    |

---

## Roadmap

* [ ] Pagination support (when account > 100 servers)
* [ ] In-memory cache to respect Online API rate limits
* [ ] Docker image + Helm chart

PRs welcome!

---

## License

MIT © 2025 — SUDO

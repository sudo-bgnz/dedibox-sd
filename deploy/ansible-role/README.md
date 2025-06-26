# scw_sd_exporter Ansible Role

An Ansible role to install and run `scw-sd`, a lightweight HTTP service exposing Prometheus targets from Scaleway Dedibox API.

## Features

- Downloads `scw-sd` binary from GitHub Releases
- Configures it as a `systemd` service
- Uses an environment variable for `ONLINE_API_TOKEN`

## Requirements

- Ansible 2.9+
- Systemd-enabled Linux host
- Valid `ONLINE_API_TOKEN` (Scaleway Dedibox)

## Role Variables

```yaml
scw_sd_version: "v0.1.0"
scw_sd_os: "linux"
scw_sd_arch: "amd64"
scw_sd_bin_dir: "/usr/local/bin"
scw_sd_token: "{{ lookup('env', 'ONLINE_API_TOKEN') }}"
scw_sd_download_url: "https://github.com/sudo-bngz/dedibox-sd/releases/download/{{ scw_sd_version }}/scw-sd-{{ scw_sd_os }}-{{ scw_sd_arch }}"
````

## Example Playbook

```yaml
- hosts: all
  roles:
    - role: scw_sd_exporter
```

## License

MIT

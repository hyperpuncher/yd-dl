# yd-dl

Download all files from a public Yandex Disk link.

## Install

**Linux / macOS:**

```sh
curl -sSL https://raw.githubusercontent.com/hyperpuncher/yd-dl/main/install.sh | sh
```

**Windows (PowerShell):**

```powershell
irm https://raw.githubusercontent.com/hyperpuncher/yd-dl/main/install.ps1 | iex
```

Or build from source:

```
go install ./...
```

Or cross-compile:

```
just build-all
```

## Usage

```
yd-dl https://disk.yandex.ru/d/IiMVtF9Eo0gDbQ
```

Files land in a directory named after the share, with the same folder structure.

## Build targets

| recipe       | does                     |
|-------------|--------------------------|
| `just fmt`  | format code              |
| `just vet`  | format + vet             |
| `just build-linux` | linux amd64 binary |
| `just build-mac`   | mac arm64 binary   |
| `just build-windows` | windows amd64 exe |
| `just build-all` | vet + all three      |
| `just clean` | remove bin/         |

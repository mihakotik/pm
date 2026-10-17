# pm — Universal Package Manager for Waynux

`pm` is a hybrid package manager designed for the Waynux GNU/Linux distribution. It provides native package management, builds packages from source manifests in an isolated `bwrap` sandbox, and transparently integrates foreign packages (Fedora RPM, Debian DEB, Arch Linux `pkg.tar.zst`, and AUR) safely into the system.

## Key Architecture

### Package Isolation Strategy
- **Native Waynux Packages**: Installed directly into root (`/usr/`, etc.).
- **Local Source Packages (`pm -s`)**: Built from source in `bwrap` sandbox and installed into `/usr/` (tagged `provider=local`).
- **Foreign Packages (Fedora / Debian / Arch / AUR)**: Installed **strictly** into `/opt/pm/<provider>/<pkg>/`. A lightweight wrapper is generated in `/usr/bin/<pkg>` that launches the program inside a `bwrap` container with isolated runtime dependencies.

### Database (`/var/lib/pm/installed.db`)
Metadata is tracked in an embedded CGO-free SQLite database (`modernc.org/sqlite`):
- `installed`: Name, version, provider, dates, license.
- `files`: File paths, package ownership, SHA256 checksums, and modes.
- `depends`: Package dependency relationships and types.
- `available`: Repository sync cache.
- `history`: Transaction audit log for rollbacks.

## CLI Usage

```bash
# Sync repository databases
pm sync

# Install packages
pm install <pkg> [--from waynux|fedora|debian|arch|aur]

# Foreign package installation (.rpm, .deb, .pkg.tar.zst)
pm -i /path/to/package.rpm

# Build and install from source manifest in current directory
pm -s

# Remove installed packages
pm remove <pkg> [--force]

# Querying
pm search <query>
pm info <pkg>
pm list
pm files <pkg>
pm owner /path/to/file

# Maintenance
pm update [<pkg>]
pm clean
pm rollback <pkg>
```

### Global Flags
- `--force`: Bypass non-critical safety checks.
- `--dry-run`: Simulate operations without modifying disk or database.
- `-v, --verbose`: Enable debug logging via `log/slog`.
- `-q, --quiet`: Silence non-error messages.
- `-y, --yes`: Auto-confirm prompts.
- `--no-deps`: Skip dependency resolution.
- `--asdeps`: Mark as dependency.
- `--asexplicit`: Mark as explicitly installed.

## Building from Source

```bash
# Compile binary to bin/pm
make build

# Run unit tests
make test

# Format and vet
make fmt vet
```

#!/usr/bin/env bash
set -euo pipefail

# Simple rsync-based deployment for the static site in ./site
# Required: REMOTE_HOST (IP or DNS)
# Optional: REMOTE_USER (default: current user), REMOTE_PATH (default: /var/www/static-site)
# Optional: SSH_KEY (path to private key), EXTRA_RSYNC_OPTS (e.g. --dry-run), RELOAD_NGINX (0 to skip reload)

REMOTE_HOST=${REMOTE_HOST:-"31.192.110.47"}
REMOTE_USER=${REMOTE_USER:-"root"}
REMOTE_PATH=${REMOTE_PATH:-"/var/www/static-site"}
SSH_KEY=${SSH_KEY:-""}
EXTRA_RSYNC_OPTS=${EXTRA_RSYNC_OPTS:-""}
RELOAD_NGINX=${RELOAD_NGINX:-1}

if [[ -z "$REMOTE_HOST" ]]; then
  echo "REMOTE_HOST is required (IP or domain of the server)." >&2
  exit 1
fi

if [[ ! -d site ]]; then
  echo "site/ directory not found. Run the script from the project root." >&2
  exit 1
fi

ssh_cmd=(ssh)
if [[ -n "$SSH_KEY" ]]; then
  ssh_cmd+=( -i "$SSH_KEY" )
fi

printf "Preparing remote path %s on %s@%s...\n" "$REMOTE_PATH" "$REMOTE_USER" "$REMOTE_HOST"
"${ssh_cmd[@]}" "${REMOTE_USER}@${REMOTE_HOST}" "mkdir -p '$REMOTE_PATH'"

rsync -avz --delete --checksum $EXTRA_RSYNC_OPTS -e "${ssh_cmd[*]}" site/ "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}/"

if [[ "$RELOAD_NGINX" == "1" ]]; then
  echo "Reloading nginx (requires sudo privileges on the server)..."
  "${ssh_cmd[@]}" "${REMOTE_USER}@${REMOTE_HOST}" "sudo systemctl reload nginx" || echo "Could not reload nginx automatically; reload manually if needed."
fi

echo "Deployment finished."

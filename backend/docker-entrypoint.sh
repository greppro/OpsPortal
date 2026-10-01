#!/bin/sh
set -e

# 以 root 启动时（docker compose 默认），先把挂载进来的数据目录和上传目录交给 appuser，
# 再降权运行，避免宿主机上由 root 创建的 ./data、./uploads 导致后端无法写入。
# 以非 root 启动时（例如 Kubernetes 设置了 runAsUser），直接运行。
if [ "$(id -u)" = "0" ]; then
  mkdir -p /app/data /app/uploads/logos
  if ! chown -R appuser:appgroup /app/data /app/uploads; then
    echo "warning: 无法修改 /app/data 或 /app/uploads 的属主，请确认挂载目录对 appuser 可写" >&2
  fi
  exec su-exec appuser:appgroup "$@"
fi

exec "$@"

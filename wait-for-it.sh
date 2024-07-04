#!/bin/bash
# wait-for-it.sh
# Ожидание доступности хоста и порта

set -e

host="$1"
port="$2"
shift 2
cmd="$@"

>&2 echo "Ожидание доступности $host:$port..."

until nc -z "$host" "$port"; do
  >&2 echo "Жду $host:$port..."
  sleep 1
done

>&2 echo "$host:$port доступен"

exec $cmd
#!/usr/bin/env bash

function log() {
  echo "$(date): ${*}"
}

function err() {
  log "${*}"
  exit "${2:-1}"
}

yen=$(date +%Y)
mon=$(date +%m)
dan=$(date +%d)
hn=$(date +%H)
mn=$(date +%M)

dates=("2024-03-02" "2023-02-01" "${yen}-${mon}-${dan}")
times=("15:30" "12:00" "${hn}:${mn}")

[[ -d data ]] || mkdir -p data

cd data || err "Can't cd to 'data'"

for d in "${dates[@]}"; do
  for t in "${times[@]}"; do
    touch -d "${d}${t}:00" "some_file_created_at_${d}_${t/:/-}"
  done
done

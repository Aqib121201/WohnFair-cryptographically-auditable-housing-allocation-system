#!/usr/bin/env bash
set -euo pipefail

echo "[1/3] Running FairRent benchmarks..."
bash scripts/bench_fairrent.sh

echo "[2/3] Running ZK-Lease benchmarks..."
bash scripts/bench_zk.sh

echo "[3/3] Running FairSurvival-GAN benchmarks..."
python src/ml/bench.py \
  --seeds 42 1337 1729 2718 31415 \
  --out experiments/results/fairsurvival.json

echo "All experiments done. Results in experiments/results/"

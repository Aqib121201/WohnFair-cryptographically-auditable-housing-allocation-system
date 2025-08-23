#!/usr/bin/env bash
set -euo pipefail

OUTDIR="experiments/results"
mkdir -p $OUTDIR

# Run ZK proof benchmarks (Halo2/PLONK)
echo "Running ZK-Lease benchmark"
python src/zklease/bench.py \
  --n 1000 \
  --out $OUTDIR/zklease.json

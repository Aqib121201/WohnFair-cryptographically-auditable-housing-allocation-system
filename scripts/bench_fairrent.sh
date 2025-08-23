#!/usr/bin/env bash
set -euo pipefail

OUTDIR="experiments/results"
mkdir -p $OUTDIR

# Seeds
SEEDS=(42 1337 1729 2718 31415)
ALPHAS=(0.5 0.7 0.9)

for SEED in "${SEEDS[@]}"; do
  for A in "${ALPHAS[@]}"; do
    echo "Running FairRent alpha=$A seed=$SEED"
    python src/fairrent/bench.py \
      --alpha $A \
      --n 50000 \
      --m 12000 \
      --seed $SEED \
      --out $OUTDIR/fairrent_${A}_${SEED}.json
  done
done

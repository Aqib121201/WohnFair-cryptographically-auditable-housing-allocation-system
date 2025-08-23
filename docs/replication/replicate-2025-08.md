# Replication Log (2025-08)

- Hardware: Intel i7-12700K, 32GB RAM, RTX 3090
- OS: Ubuntu 22.04, Docker 24.0
- Seeds: 42, 1337, 1729, 2718, 31415

## Results
- FairRent clearance time reduction: 41% ±1.2
- ZK-Lease: proof size 22.5 KB, verify 7.8 ms, prove 0.91 s
- FairSurvival-GAN: C-index 0.852 ±0.02, EOG = 0.048

All reproduced with `make reproduce`.

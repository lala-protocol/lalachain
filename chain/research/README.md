# Phase 3 Research Harness

This folder contains research-track artifacts for Phase 3 in
IMPLEMENTATION_FEASIBILITY.md.

## Shadow Mode Harness

`phase3_shadow_mode.py` trains a lightweight classifier from KPI traces,
compares model recommendations against the active rule-based advisor, and
writes a JSON report containing:

- shadow-mode agreement metrics,
- divergence samples,
- a compact zkML candidate evaluation matrix,
- a model-update governance process draft.

Run:

```bash
python research/phase3_shadow_mode.py \
  --input-log reports/phase2-cosmos-multi.log \
  --output reports/phase3-shadow-report.json
```

If no input log is available, synthetic KPI traces are generated.

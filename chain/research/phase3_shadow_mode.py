"""
Phase 3 research harness for LalaChain.

This script provides a docs-aligned research-track implementation artifact:
- trains a simple model from local testnet KPI traces,
- runs shadow-mode comparison against the rule-based advisor,
- emits a compact zkML evaluation matrix,
- writes a machine-readable JSON report.

It intentionally avoids external dependencies.
"""

from __future__ import annotations

import argparse
import json
import math
import re
from dataclasses import dataclass
from pathlib import Path
from random import Random
from typing import Dict, Iterable, List, Sequence, Tuple

KPI_PATTERN = re.compile(r"\[KPI\]\s+util=([0-9.]+)%\s+fee=([0-9.]+)\s+gwei", re.IGNORECASE)
NODE_PATTERN = re.compile(r"\[node-\d+\].*?util=([0-9.]+)%\s+fee=([0-9.]+)\s+gwei", re.IGNORECASE)

LOW_UTIL = 0.40
HIGH_UTIL = 0.80
MIN_FEE = 0.80
MAX_FEE = 5.00

ACTION_HOLD = "hold"
ACTION_INC_GAS = "increase_block_gas_limit"
ACTION_DEC_GAS = "decrease_block_gas_limit"
ACTION_INC_FEE = "increase_base_fee"
ACTION_DEC_FEE = "decrease_base_fee"

ALL_ACTIONS = [
    ACTION_HOLD,
    ACTION_INC_GAS,
    ACTION_DEC_GAS,
    ACTION_INC_FEE,
    ACTION_DEC_FEE,
]


@dataclass
class KPI:
    util: float
    fee_gwei: float


def parse_kpis_from_log(text: str) -> List[KPI]:
    kpis: List[KPI] = []
    for match in KPI_PATTERN.finditer(text):
        util = float(match.group(1)) / 100.0
        fee = float(match.group(2))
        kpis.append(KPI(util=util, fee_gwei=fee))

    if kpis:
        return kpis

    # Multi-validator reports summarize telemetry per node instead of per epoch.
    for match in NODE_PATTERN.finditer(text):
        util = float(match.group(1)) / 100.0
        fee = float(match.group(2))
        kpis.append(KPI(util=util, fee_gwei=fee))

    return kpis


def synthetic_kpis(epochs: int, seed: int) -> List[KPI]:
    rng = Random(seed)
    rows: List[KPI] = []
    for epoch in range(1, epochs + 1):
        if epoch <= 10:
            base = 0.30
        elif epoch <= 20:
            base = 0.90
        else:
            base = 0.55
        util = max(0.0, min(1.2, base + rng.gauss(0.0, 0.06)))
        # Keep fee range close to prototype traces.
        fee = max(0.10, min(10.0, 1.00 + (util - 0.5) * 0.20 + rng.gauss(0.0, 0.02)))
        rows.append(KPI(util=util, fee_gwei=fee))
    return rows


def rule_actions(kpis: Sequence[KPI]) -> List[str]:
    actions: List[str] = []
    low_streak = 0
    high_streak = 0

    for row in kpis:
        if row.util < LOW_UTIL:
            low_streak += 1
            high_streak = 0
        elif row.util > HIGH_UTIL:
            high_streak += 1
            low_streak = 0
        else:
            low_streak = 0
            high_streak = 0

        action = ACTION_HOLD
        if low_streak >= 3 and row.fee_gwei < MIN_FEE:
            action = ACTION_INC_GAS
        elif high_streak >= 2:
            action = ACTION_DEC_GAS
        elif row.fee_gwei > MAX_FEE:
            action = ACTION_DEC_FEE
        elif row.fee_gwei < MIN_FEE:
            action = ACTION_INC_FEE

        actions.append(action)

    return actions


def split_train_test(rows: Sequence[KPI], labels: Sequence[str], ratio: float = 0.7) -> Tuple[List[KPI], List[str], List[KPI], List[str]]:
    split = max(1, int(len(rows) * ratio))
    split = min(split, len(rows) - 1)
    return list(rows[:split]), list(labels[:split]), list(rows[split:]), list(labels[split:])


def euclidean(a: Tuple[float, float], b: Tuple[float, float]) -> float:
    return math.sqrt((a[0] - b[0]) ** 2 + (a[1] - b[1]) ** 2)


def train_centroid_model(train_rows: Sequence[KPI], train_labels: Sequence[str]) -> Dict[str, Tuple[float, float]]:
    sums: Dict[str, Tuple[float, float, int]] = {}
    for row, label in zip(train_rows, train_labels):
        util_sum, fee_sum, n = sums.get(label, (0.0, 0.0, 0))
        sums[label] = (util_sum + row.util, fee_sum + row.fee_gwei, n + 1)

    centroids: Dict[str, Tuple[float, float]] = {}
    for action in ALL_ACTIONS:
        util_sum, fee_sum, n = sums.get(action, (0.0, 0.0, 0))
        if n == 0:
            continue
        centroids[action] = (util_sum / n, fee_sum / n)

    return centroids


def majority_label(labels: Sequence[str]) -> str:
    counts: Dict[str, int] = {}
    for label in labels:
        counts[label] = counts.get(label, 0) + 1
    return max(counts.items(), key=lambda item: item[1])[0]


def predict_action(centroids: Dict[str, Tuple[float, float]], fallback: str, row: KPI) -> str:
    if not centroids:
        return fallback

    row_vec = (row.util, row.fee_gwei)
    best = fallback
    best_dist = float("inf")
    for action, centroid in centroids.items():
        dist = euclidean(row_vec, centroid)
        if dist < best_dist:
            best_dist = dist
            best = action
    return best


def confusion_matrix(expected: Sequence[str], predicted: Sequence[str]) -> Dict[str, Dict[str, int]]:
    matrix: Dict[str, Dict[str, int]] = {action: {a: 0 for a in ALL_ACTIONS} for action in ALL_ACTIONS}
    for exp, pred in zip(expected, predicted):
        matrix[exp][pred] += 1
    return matrix


def agreement_rate(expected: Sequence[str], predicted: Sequence[str]) -> float:
    if not expected:
        return 0.0
    matches = sum(1 for exp, pred in zip(expected, predicted) if exp == pred)
    return matches / len(expected)


def zkml_matrix() -> List[Dict[str, object]]:
    # Candidate names are explicitly referenced in IMPLEMENTATION_FEASIBILITY.md.
    return [
        {
            "library": "EZKL",
            "focus": "ONNX model proof generation",
            "integration_readiness": "medium",
            "expected_effort": "medium-high",
            "note": "Candidate for small-model inference proofs; benchmark with telemetry-sized inputs.",
        },
        {
            "library": "Modulus Labs",
            "focus": "zkML proving workflows",
            "integration_readiness": "medium",
            "expected_effort": "high",
            "note": "Strong research candidate; assess proving latency versus epoch timing.",
        },
        {
            "library": "Giza",
            "focus": "verifiable ML pipelines",
            "integration_readiness": "medium",
            "expected_effort": "high",
            "note": "Evaluate tooling fit for model governance and reproducibility constraints.",
        },
    ]


def governance_update_process() -> List[str]:
    return [
        "Publish model artifact hash, training metadata, and evaluation report.",
        "Submit model-update proposal through Governance Layer with activation epoch.",
        "Require quorum and approval thresholds equal to parameter-change governance.",
        "Run shadow mode for a fixed observation window before activation.",
        "Activate model only if shadow-mode divergence and stability metrics stay within limits.",
    ]


def read_log_text(path: Path) -> str:
    """Read KPI logs produced by different shells/encodings."""
    for encoding in ("utf-8", "utf-8-sig", "utf-16", "utf-16-le", "utf-16-be"):
        try:
            return path.read_text(encoding=encoding)
        except UnicodeDecodeError:
            continue
    # Last-resort fallback keeps the harness usable with legacy encodings.
    return path.read_text(encoding="cp1252", errors="replace")


def build_report(kpis: Sequence[KPI], source: str) -> Dict[str, object]:
    labels = rule_actions(kpis)
    train_rows, train_labels, test_rows, test_labels = split_train_test(kpis, labels)
    centroids = train_centroid_model(train_rows, train_labels)
    fallback = majority_label(train_labels)

    predicted = [predict_action(centroids, fallback, row) for row in test_rows]
    agreement = agreement_rate(test_labels, predicted)

    top_divergences: List[Dict[str, object]] = []
    for idx, (row, expected, pred) in enumerate(zip(test_rows, test_labels, predicted), start=1):
        if expected == pred:
            continue
        top_divergences.append(
            {
                "test_index": idx,
                "utilization": round(row.util, 4),
                "base_fee_gwei": round(row.fee_gwei, 4),
                "rule_action": expected,
                "model_action": pred,
            }
        )
        if len(top_divergences) >= 10:
            break

    return {
        "source": source,
        "sample_count": len(kpis),
        "train_count": len(train_rows),
        "test_count": len(test_rows),
        "model_type": "nearest-centroid-classifier",
        "shadow_mode": {
            "agreement_rate": round(agreement, 4),
            "confusion_matrix": confusion_matrix(test_labels, predicted),
            "top_divergences": top_divergences,
        },
        "zkml_evaluation_matrix": zkml_matrix(),
        "model_update_governance_process": governance_update_process(),
        "notes": [
            "This is a research-track harness for Phase 3 shadow mode.",
            "Rule-based advisor remains the active decision source.",
            "Replace synthetic/trace-driven training with real multi-validator telemetry datasets for production readiness.",
        ],
    }


def main() -> None:
    parser = argparse.ArgumentParser(description="Phase 3 shadow-mode research harness")
    parser.add_argument("--input-log", type=Path, default=None, help="Optional runtime log file containing KPI lines")
    parser.add_argument("--epochs", type=int, default=30, help="Synthetic epoch count when --input-log is omitted")
    parser.add_argument("--seed", type=int, default=42, help="Random seed for synthetic data")
    parser.add_argument("--output", type=Path, required=True, help="Output JSON report path")
    args = parser.parse_args()

    if args.input_log is not None and args.input_log.exists():
        text = read_log_text(args.input_log)
        rows = parse_kpis_from_log(text)
        source = f"parsed log: {args.input_log}"
        if len(rows) < 8:
            rows = synthetic_kpis(max(args.epochs, 30), args.seed)
            source = f"synthetic fallback (insufficient KPI rows in {args.input_log})"
    else:
        rows = synthetic_kpis(args.epochs, args.seed)
        source = "synthetic"

    report = build_report(rows, source)

    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(json.dumps(report, indent=2), encoding="utf-8")

    print("Phase 3 shadow-mode report generated")
    print(f"  Source        : {report['source']}")
    print(f"  Samples       : {report['sample_count']}")
    print(f"  Agreement     : {report['shadow_mode']['agreement_rate']:.2%}")
    print(f"  Output report : {args.output}")


if __name__ == "__main__":
    main()

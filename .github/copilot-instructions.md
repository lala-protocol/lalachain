# Lala Protocol Copilot Instructions

Use these rules for all work in this repository.

## Source-of-truth hierarchy
1. Latest user request in chat.
2. `IMPLEMENTATION_FEASIBILITY.md` for implementation decisions and roadmap.
3. `lala-protocol-whitepaper.tex` (and generated PDF) for protocol intent and terminology.
4. `simulation/` for currently implemented behavior examples.

If sources conflict, follow the higher-priority item and explicitly note the conflict in your response.

## Anti-hallucination rules
- Do not invent protocol features, modules, equations, timelines, or governance rules that are not present in the sources above.
- Do not claim something is "in the whitepaper" unless it is actually present in repository files.
- If required details are missing, state what is missing and ask for clarification instead of guessing.
- Keep assumptions minimal and label them clearly as assumptions.

## Implementation guardrails
- Treat this repo as **concept + simulation + implementation planning** unless the user explicitly asks for production chain code.
- Default to the Phase 1 path in `IMPLEMENTATION_FEASIBILITY.md`:
  - rule-based advisor first,
  - deterministic telemetry,
  - governance-ratified parameter changes.
- Do not introduce advanced ML/zkML, external oracle dependencies, or extra tokenomics unless explicitly requested.
- Reuse the project's existing terminology exactly (Base Consensus Layer, Network Telemetry Module, AI Advisory Module, Governance Layer).

## Change discipline
- Keep edits tightly scoped to the user request.
- Prefer extending existing files over creating new architecture or process documents unless requested.
- For protocol or simulation behavior changes, include a short traceability note in the response referencing the source file(s) used.

const SOURCE_POLICY =
  "Answer from these graph fields only. For architecture or flow questions, explain the central path; do not exhaust variants, tests, adapters, or generated clients. When the result has file, symbol, relation, and range evidence, answer now. If source text is required, cite the smallest returned sourceSpan/evidence range and stop; do not read files for this answer.";

export function resultGuide(action: string): string {
  return `${action} ${SOURCE_POLICY}`;
}

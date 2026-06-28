import { ITtscGraphDetails } from "./ITtscGraphDetails";
import { ITtscGraphEntrypoints } from "./ITtscGraphEntrypoints";
import { ITtscGraphEscape } from "./ITtscGraphEscape";
import { ITtscGraphLookup } from "./ITtscGraphLookup";
import { ITtscGraphOverview } from "./ITtscGraphOverview";
import { ITtscGraphTrace } from "./ITtscGraphTrace";

/** The typed MCP surface; its single method becomes the single graph tool. */
export interface ITtscGraphApplication {
  /**
   * Inspect TypeScript source flow before shell search or file reads.
   *
   * Use first when an answer depends on TypeScript source structure: symbols,
   * imports, calls, decorators, DI/request lifecycle, refs, types, ranges, or
   * flows such as IPC/process/channel communication, render, validation, load,
   * persist, or propagation. Use shell reads for scripts, configs, docs, exact
   * text, generated output, and non-TypeScript files.
   *
   * It returns index facts and ranges, not file bodies. After calling it,
   * answer from graph fields only or choose `escape` and cite the smallest
   * range.
   *
   * @param props Reasoning plus one graph request
   * @returns Matching `result` union member
   */
  inspect_typescript_source_flow(
    props: ITtscGraphApplication.IProps,
  ): ITtscGraphApplication.IResult;
}

export namespace ITtscGraphApplication {
  /** Plan with graph-specific reasoning, then submit exactly one request. */
  export interface IProps {
    /**
     * User's TypeScript code question.
     *
     * Restate only the code question being answered. Keep it about TypeScript
     * symbols, call flow, type flow, runtime behavior, tests exposed by the
     * graph, or architecture. If the user asks about scripts, config, docs,
     * generated output, exact text, or non-TypeScript files, say so here and
     * choose `escape`.
     */
    question: string;

    /**
     * Why the resident graph is the next evidence source.
     *
     * State the smallest evidence that will settle the answer: central handles,
     * a path, caller, dependency edge, signature, member outline, or range. If
     * current graph evidence is enough, say so and choose `escape`. Do not
     * frame this as an invitation to keep finding every branch.
     *
     * Example: `central public API-to-worker path with range anchors`.
     */
    graphNeed: string;

    /**
     * First request-type decision before arguments are filled.
     *
     * Choose the smallest operation class before filling arguments.
     * `entrypoints` starts natural questions, `lookup` resolves a concrete
     * name, `trace` follows flow, `details` inspects selected handles,
     * `overview` maps broad surfaces, and `escape` stops.
     *
     * Example: reason `find central handles before tracing`, type
     * `entrypoints`.
     */
    draft: IRequestDraft;

    /**
     * Critical gate before the final request.
     *
     * Example inspect gate: reason `no graph evidence yet; default entrypoints
     * is smallest`, decision `inspect`, finish `answer`.
     *
     * Example stop gate: reason `current graph result has ranges; source body
     * would require file text`, decision `escape`, finish `range`.
     */
    review: IRequestReview;

    /** The graph operation chosen from the reasoning above, or a no-op escape. */
    request:
      | ITtscGraphEntrypoints.IRequest
      | ITtscGraphLookup.IRequest
      | ITtscGraphTrace.IRequest
      | ITtscGraphDetails.IRequest
      | ITtscGraphOverview.IRequest
      | ITtscGraphEscape.IRequest;
  }

  /** First-pass operation choice before final request arguments. */
  export interface IRequestDraft {
    /** Why this operation type is the smallest useful next step. */
    reason: string;

    /** Draft discriminator for the intended graph operation. */
    type:
      | ITtscGraphEntrypoints.IRequest["type"]
      | ITtscGraphLookup.IRequest["type"]
      | ITtscGraphTrace.IRequest["type"]
      | ITtscGraphDetails.IRequest["type"]
      | ITtscGraphOverview.IRequest["type"]
      | ITtscGraphEscape.IRequest["type"];
  }

  /** Final gate that prevents graph calls from becoming source-read preludes. */
  export interface IRequestReview {
    /**
     * Why the final request is valid, or why graph use must stop.
     *
     * Reject duplicate evidence, broad limits, neighbor expansion without a
     * named missing edge, test-only hunting, and any plan to fall back to shell
     * search or file reads after this tool returns.
     */
    reason: string;

    /**
     * Whether this MCP call should run a graph request or intentionally stop.
     *
     * Use exactly `inspect` only when the next evidence source is still the
     * resident graph. Use exactly `escape` when current graph evidence is
     * enough, when the user needs non-TypeScript/source-body facts, or when
     * another graph call would be overfetch.
     */
    decision: string;

    /**
     * How the agent must finish after this MCP call returns.
     *
     * Use exactly `answer`, `range`, or `clarify`. `answer` means answer from
     * graph fields only. `range` means cite the smallest returned `sourceSpan`
     * or edge `evidence` range and stop instead of reading the file. `clarify`
     * means ask for a concrete symbol or scope.
     */
    finish: string;
  }

  /** The selected request's output. `result.type` mirrors `request.type`. */
  export interface IResult {
    result:
      | ITtscGraphEntrypoints
      | ITtscGraphLookup
      | ITtscGraphTrace
      | ITtscGraphDetails
      | ITtscGraphOverview
      | ITtscGraphEscape;
  }
}

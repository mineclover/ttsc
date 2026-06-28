import { TtscGraphMemory } from "./model/TtscGraphMemory";
import { resultGuide } from "./server/resultGuide";
import { runDetails } from "./server/runDetails";
import { runEntrypoints } from "./server/runEntrypoints";
import { runLookup } from "./server/runLookup";
import { runOverview } from "./server/runOverview";
import { runTrace } from "./server/runTrace";
import { ITtscGraphApplication } from "./structures/ITtscGraphApplication";
import { ITtscGraphEscape } from "./structures/ITtscGraphEscape";

export type TtscGraphSource = TtscGraphMemory | (() => TtscGraphMemory);

const MAX_GRAPH_CALLS_PER_ANSWER = 4;
const BUDGET_IDLE_RESET_MS = 5 * 60 * 1_000;

/**
 * The MCP tool surface as a plain class over the resident
 * {@link TtscGraphMemory}.
 *
 * Its public method is the MCP tool: `typia.llm.controller` reflects
 * {@link ITtscGraphApplication} to generate the tool's JSON schema and argument
 * validator from the signature and JSDoc, with no hand-written schema. The
 * method delegates to the pure graph functions in `./server`, which are
 * unit-testable without a transport; this class only binds them to the graph.
 *
 * Every method answers from the resident graph; none recompiles. Output is kept
 * compact and bounded so a model can read structure without a file read, which
 * is the token win the redesign exists for.
 */
export class TtscGraphApplication implements ITtscGraphApplication {
  private readonly graph: () => TtscGraphMemory;
  private graphCalls = 0;
  private lastGraphCallAt = 0;

  public constructor(source: TtscGraphSource) {
    this.graph = typeof source === "function" ? source : () => source;
  }

  public inspect_typescript_source_flow(
    props: ITtscGraphApplication.IProps,
  ): ITtscGraphApplication.IResult {
    const decision = props.review.decision.trim().toLowerCase();
    const finish = props.review.finish.trim().toLowerCase();
    if (decision === "escape" && props.request.type !== "escape") {
      return {
        result: this.escape(
          props.review.reason,
          finish === "range"
            ? "cite the smallest returned sourceSpan/evidence range and stop"
            : finish === "clarify"
              ? "ask for a concrete symbol or scope"
              : "answer from prior graph evidence",
        ),
      };
    }
    if (props.request.type === "escape") {
      const result = this.escape(props.request.reason);
      if (props.request.nextStep !== undefined) {
        result.nextStep = props.request.nextStep;
      }
      return {
        result,
      };
    }
    this.refreshBudget();
    if (this.graphCalls >= MAX_GRAPH_CALLS_PER_ANSWER) {
      return {
        result: this.budgetEscape(),
      };
    }
    this.graphCalls++;
    this.lastGraphCallAt = Date.now();
    switch (props.request.type) {
      case "entrypoints":
        return {
          result: this.withBudgetGuide(
            runEntrypoints(this.graph(), props.request),
          ),
        };
      case "lookup":
        return {
          result: this.withBudgetGuide(runLookup(this.graph(), props.request)),
        };
      case "trace":
        return {
          result: this.withBudgetGuide(runTrace(this.graph(), props.request)),
        };
      case "details":
        return {
          result: this.withBudgetGuide(runDetails(this.graph(), props.request)),
        };
      case "overview":
        return {
          result: this.withBudgetGuide(
            runOverview(this.graph(), props.request),
          ),
        };
      default:
        throw new Error("Unknown graph request type");
    }
  }

  private refreshBudget(): void {
    if (
      this.lastGraphCallAt !== 0 &&
      Date.now() - this.lastGraphCallAt > BUDGET_IDLE_RESET_MS
    ) {
      this.graphCalls = 0;
    }
  }

  private withBudgetGuide<T extends { guide: string }>(result: T): T {
    const remaining = Math.max(0, MAX_GRAPH_CALLS_PER_ANSWER - this.graphCalls);
    result.guide = `${result.guide} Graph budget: call ${this.graphCalls}/${MAX_GRAPH_CALLS_PER_ANSWER}; ${
      remaining === 0
        ? "answer now from current graph evidence."
        : `${remaining} graph call(s) remain before you must answer.`
    }`;
    return result;
  }

  private escape(reason: string, nextStep?: string): ITtscGraphEscape {
    return {
      type: "escape",
      skipped: true,
      reason,
      guide: resultGuide(
        "Finish from existing graph evidence, state the graph gap, or ask for clarification.",
      ),
      ...(nextStep !== undefined ? { nextStep } : {}),
    };
  }

  private budgetEscape(): ITtscGraphEscape {
    return {
      type: "escape",
      skipped: true,
      reason: "Graph call budget reached for this answer.",
      guide:
        "Stop using graph and shell tools. Answer now from current graph evidence, or cite the smallest returned range and state the remaining gap.",
      nextStep: "answer from current graph evidence",
    };
  }
}

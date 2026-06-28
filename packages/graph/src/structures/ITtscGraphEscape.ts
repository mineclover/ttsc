/** The no-op result that ends graph use instead of falling back silently. */
export interface ITtscGraphEscape {
  /** Discriminator for the no-op escape route. */
  type: "escape";

  /** Always true so callers can distinguish an intentional no-op. */
  skipped: true;

  /** Why no graph operation should run. */
  reason: string;

  /** How to finish without starting source reads from this no-op result. */
  guide: string;

  /** The final non-graph note, if useful. Not a command to run now. */
  nextStep?: string;
}

export namespace ITtscGraphEscape {
  /** End the graph answer when graph evidence is unnecessary or exhausted. */
  export interface IRequest {
    /** Discriminator for the no-op escape route. */
    type: "escape";

    /**
     * Why no graph operation should run.
     *
     * Use this when the review finds the user is asking about package scripts,
     * config files, generated output, prose documentation, exact text, or an
     * answer that the current graph result already settled. Also use this when
     * source text is required: name the smallest returned sourceSpan and stop
     * instead of reading the file in the same answer.
     */
    reason: string;

    /**
     * The final non-graph note, if useful.
     *
     * Keep this short. Examples: `answer from the prior graph result`, `source
     * body needed at returned sourceSpan`, or `ask the user for a concrete
     * symbol`. Do not put a shell command here.
     */
    nextStep?: string;
  }
}

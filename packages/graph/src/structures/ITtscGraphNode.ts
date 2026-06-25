import { ITtscGraphComponent } from "./ITtscGraphComponent";
import { ITtscGraphDecorator } from "./ITtscGraphDecorator";
import { ITtscGraphEvidence } from "./ITtscGraphEvidence";
import { ITtscGraphRoute } from "./ITtscGraphRoute";
import { TtscGraphNodeKind } from "./TtscGraphNodeKind";
import { TtscGraphNodeModifier } from "./TtscGraphNodeModifier";

/**
 * One node in the graph: a declared symbol, a structural container (file,
 * package), or a synthesized framework node (route, component).
 *
 * The `id` is position-invariant. For a symbol it is `path#qualifiedName:kind`
 * (e.g. `src/order.ts#OrderService.create:method`), so inserting a line above a
 * declaration does not re-key it; for a virtual node it is a stable semantic
 * key (e.g. `route:http:GET:/users/:id`, `component:src/App.tsx#App`). Line and
 * span live in `evidence` and are never part of identity.
 */
export interface ITtscGraphNode {
  /** Position-invariant identity (see the interface doc for the id grammar). */
  id: string;

  /** What this node represents. */
  kind: TtscGraphNodeKind;

  /** The simple, unqualified declared name (`create`, `OrderService`, `App`). */
  name: string;

  /**
   * The owner-qualified name, when the node lives inside another declaration —
   * `OrderService.create`, `Shopping.ISale`. Absent for a top-level or virtual
   * node whose `name` already qualifies it.
   */
  qualifiedName?: string;

  /**
   * Project-relative path of the file that declares this node. For a virtual
   * node it is the file the convention was recognized in.
   */
  file: string;

  /**
   * True when the declaration lives outside the workspace (a dependency). The
   * graph keeps the leaf as a named endpoint but does not expand its
   * internals.
   */
  external: boolean;

  /**
   * True when `file` is git-ignored generated code (a Prisma client, a codegen
   * output). Projections desurface these so generated nodes do not bury the
   * authored graph.
   */
  ignored?: boolean;

  /** True when the symbol is part of its module's export surface. */
  exported?: boolean;

  /** Declaration modifiers, when the declaration pass recorded any. */
  modifiers?: TtscGraphNodeModifier[];

  /**
   * The decorators written on this declaration, in source order, when it has
   * any. A framework pass reads these to synthesize routes (`@Controller`,
   * `@Get`) without re-parsing source.
   */
  decorators?: ITtscGraphDecorator[];

  /** The declaration span, for display and source expansion. */
  evidence?: ITtscGraphEvidence;

  /** Route facts; present iff `kind === "route"`. */
  route?: ITtscGraphRoute;

  /** Component facts; present iff `kind === "component"`. */
  component?: ITtscGraphComponent;
}

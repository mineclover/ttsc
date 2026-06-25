/**
 * The framework-level facts of a `component` virtual {@link IGraphNode}.
 *
 * A component node marks a UI component the framework pass recognized (a React
 * function/class component, a Next page or layout). It is `framework-derived`;
 * the `renders` edge connects it to the components it mounts.
 */
export interface IComponentMetadata {
  /** The originating framework: `react`, `next`, `vue`, … */
  framework: string;

  /**
   * The file-route path for a page or layout component (Next `app/` or `pages/`
   * routing), when this component is route-bound.
   */
  routePath?: string;
}

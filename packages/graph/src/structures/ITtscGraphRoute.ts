/**
 * The framework-level facts of a `route` virtual {@link ITtscGraphNode}.
 *
 * A route node has no declaring symbol; it is synthesized by the framework pass
 * from decorator or file-routing conventions and tagged `framework-derived`.
 * The route's identity is its semantic key (e.g. `route:http:GET:/users/:id`),
 * not a source position, so the resolved `path` here survives controller/module
 * prefix composition.
 */
export interface ITtscGraphRoute {
  /**
   * Transport category: `http`, `graphql`, `websocket`, `microservice`, or
   * `page` for a file-based page/layout route. Left open for frameworks that
   * add their own.
   */
  protocol: string;

  /** The originating framework: `nest`, `next`, `react-router`, `express`, … */
  framework: string;

  /**
   * The fully-resolved route path, with every controller, module, and segment
   * prefix already composed in (e.g. `/shoppings/sellers/sales`).
   */
  path: string;

  /**
   * The HTTP verb, GraphQL operation, or message pattern, when the protocol has
   * one (`GET`, `POST`, `Query`, `Mutation`, `MessagePattern`, …).
   */
  method?: string;

  /**
   * Node id of the handler symbol (controller method, resolver, page component)
   * that serves this route, when it was resolved.
   */
  handler?: string;
}

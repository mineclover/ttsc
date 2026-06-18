import type { SourceFile } from "./SourceFile";

/**
 * A source file that redirects to another source file's content. In the legacy
 * compiler this carries redirect bookkeeping; this package models only the
 * underlying statements, so it emits exactly like its target source file.
 *
 * Built by {@link factory.createRedirectedSourceFile}.
 *
 * @author Jeongho Nam - https://github.com/samchon
 */
export interface RedirectedSourceFile {
  /** Discriminant tag; always `"RedirectedSourceFile"`. */
  kind: "RedirectedSourceFile";

  /** The source file whose content is emitted. */
  redirectTarget: SourceFile;
}

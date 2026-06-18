import nodeResolve from "@rollup/plugin-node-resolve";
import { globSync } from "tinyglobby";

// The package is CommonJS-first: `tsgo -p .` emits the canonical CommonJS
// `lib/**/*.js` plus the `.d.ts` (which resolve under CommonJS rules for any
// consumer). `tsgo -p tsconfig.esm.json` emits a parallel ESM tree into
// `esm-tmp/`. This step re-emits that ESM tree as `lib/**/*.mjs` siblings,
// resolving the extensionless / directory-index imports. Because the input is
// real ESM, the `.mjs` keeps genuine `export { ... }` named exports that ESM
// bundlers (vite, rollup) can statically see.
export default {
  input: globSync("./esm-tmp/**/*.js"),
  output: {
    dir: "./lib",
    format: "esm",
    entryFileNames: "[name].mjs",
    preserveModules: true,
    preserveModulesRoot: "esm-tmp",
    sourcemap: true,
  },
  plugins: [nodeResolve({ extensions: [".js"] })],
};

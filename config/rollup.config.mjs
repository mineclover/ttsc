import typescript from "@rollup/plugin-typescript";

export default {
  input: "src/index.ts",
  output: {
    dir: "lib",
    format: "esm",
    // Emit one .mjs per source module (mirroring tsgo's per-file CJS .js)
    // instead of bundling everything into a single file.
    preserveModules: true,
    preserveModulesRoot: "src",
    entryFileNames: "[name].mjs",
    sourcemap: true,
  },
  plugins: [
    typescript({
      tsconfig: "tsconfig.json",
      module: "esnext",
      moduleResolution: "bundler",
      declaration: false,
      declarationMap: false,
      outDir: "lib",
      sourceMap: true,
    }),
  ],
};

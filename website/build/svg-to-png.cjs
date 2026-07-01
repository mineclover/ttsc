const fs = require("fs");
const os = require("os");
const path = require("path");
const { execFileSync } = require("child_process");

const ROOT = path.resolve(__dirname, "..");
const OUT_DIR = path.join(__dirname, "png-out");
const SCALE = 2;

// The two charts embedded in the dev.to blog post. dev.to only accepts PNG
// uploads, so these occasionally need re-exporting from the published SVGs.
const DEFAULT_SVGS = [
  "public/benchmark/graph-common-codex-gpt-5.5.svg",
  "public/benchmark/graph-zod-common-codex-gpt-5.5.svg",
].map((rel) => path.join(ROOT, rel));

// Locally installed headless browsers, in preference order. Chrome first, Edge
// as a fallback; both render SVG identically via the same --screenshot flag.
const BROWSERS = [
  "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
  "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe",
];

const inputs = process.argv.slice(2);
const svgPaths = (inputs.length > 0 ? inputs : DEFAULT_SVGS).map((p) =>
  path.resolve(p),
);
const browser = findBrowser();

fs.mkdirSync(OUT_DIR, { recursive: true });
for (const svgPath of svgPaths) {
  const out = renderPng(svgPath);
  console.log(`[svg-to-png] wrote ${out.file} (${out.width}x${out.height})`);
}

// Wrap the SVG in a viewport-sized HTML page (white background, no margins) so
// Chrome screenshots exactly the intrinsic dimensions, then scale by SCALE via
// --force-device-scale-factor for a crisp 2x export.
function renderPng(svgPath) {
  if (!fs.existsSync(svgPath)) throw new Error(`SVG not found: ${svgPath}`);

  const svg = fs.readFileSync(svgPath, "utf8");
  const { width, height } = readSvgSize(svg);
  const outFile = path.join(OUT_DIR, `${path.basename(svgPath, ".svg")}.png`);

  const inner = svg.replace(/^\s*<\?xml[^>]*\?>\s*/i, "");
  const html = `<!doctype html>
<html>
 <head>
  <meta charset="utf-8" />
  <style>
   html, body { margin: 0; padding: 0; background: #ffffff; }
   svg { display: block; width: ${width}px; height: ${height}px; }
  </style>
 </head>
 <body>${inner}</body>
</html>`;

  const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), "svg-to-png-"));
  const htmlFile = path.join(tmpDir, "wrapper.html");
  fs.writeFileSync(htmlFile, html);
  try {
    execFileSync(
      browser,
      [
        "--headless=new",
        "--disable-gpu",
        "--hide-scrollbars",
        "--default-background-color=FFFFFFFF",
        `--force-device-scale-factor=${SCALE}`,
        `--screenshot=${outFile}`,
        `--window-size=${width},${height}`,
        toFileUrl(htmlFile),
      ],
      { stdio: "ignore" },
    );
  } finally {
    fs.rmSync(tmpDir, { recursive: true, force: true });
  }

  const rendered = readPngSize(outFile);
  const expected = { width: width * SCALE, height: height * SCALE };
  if (rendered.width !== expected.width || rendered.height !== expected.height)
    throw new Error(
      `${outFile}: expected ${expected.width}x${expected.height}, got ${rendered.width}x${rendered.height}`,
    );
  return { file: outFile, ...rendered };
}

function findBrowser() {
  const found = BROWSERS.find((bin) => fs.existsSync(bin));
  if (!found)
    throw new Error(
      `No headless browser found. Checked:\n  ${BROWSERS.join("\n  ")}`,
    );
  return found;
}

// Intrinsic size from the root <svg> width/height attributes, falling back to
// the viewBox so the PNG never guesses at dimensions.
function readSvgSize(svg) {
  const open = svg.match(/<svg\b[^>]*>/i);
  if (!open) throw new Error("no <svg> root element found");
  const tag = open[0];
  const w = numAttr(tag, "width");
  const h = numAttr(tag, "height");
  if (w > 0 && h > 0) return { width: w, height: h };
  const viewBox = tag.match(/viewBox\s*=\s*"([^"]+)"/i);
  if (viewBox) {
    const parts = viewBox[1]
      .trim()
      .split(/[\s,]+/)
      .map(Number);
    if (parts.length === 4 && parts[2] > 0 && parts[3] > 0)
      return { width: parts[2], height: parts[3] };
  }
  throw new Error("could not determine SVG width/height");
}

function numAttr(tag, name) {
  const match = tag.match(new RegExp(`${name}\\s*=\\s*"([\\d.]+)`, "i"));
  return match ? Math.round(Number(match[1])) : 0;
}

// PNG dimensions live in the IHDR chunk: bytes 16-19 width, 20-23 height, both
// big-endian. Reading them proves the screenshot came out at the exact 2x size.
function readPngSize(file) {
  const buf = fs.readFileSync(file);
  const signature = "89504e470d0a1a0a";
  if (buf.subarray(0, 8).toString("hex") !== signature)
    throw new Error(`${file} is not a PNG`);
  return { width: buf.readUInt32BE(16), height: buf.readUInt32BE(20) };
}

function toFileUrl(file) {
  return `file:///${file.replace(/\\/g, "/")}`;
}

// Positive: near-miss typo on a Next.js data-fetching export.
// expect: nextjs/no-typos error
export function getstaticprops() {
  return { props: {} };
}

// Negative: correctly-cased export name.
export function getStaticProps() {
  return { props: {} };
}

JSON.stringify({ getstaticprops, getStaticProps });

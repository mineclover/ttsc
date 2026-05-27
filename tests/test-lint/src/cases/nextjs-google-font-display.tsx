// Positive: Google Fonts `<link>` without a `font-display` parameter.
// expect: nextjs/google-font-display error
const a = (
  <link
    rel="stylesheet"
    href="https://fonts.googleapis.com/css2?family=Inter"
  />
);

// Negative: Google Fonts `<link>` with a supported `font-display`.
const b = (
  <link
    rel="stylesheet"
    href="https://fonts.googleapis.com/css2?family=Inter&display=swap"
  />
);

JSON.stringify({ a, b });

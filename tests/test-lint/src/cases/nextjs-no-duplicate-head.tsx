import { Head } from "next/document";

// Positive: more than one `<Head>` element in the same render tree.
const a = (
  <>
    <Head />
    {/* expect: nextjs/no-duplicate-head error */}
    <Head />
  </>
);

// Negative: a single `<Head>` element.
const b = <Head />;

JSON.stringify({ a, b });

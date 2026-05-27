declare const html: string;

// Positive: `dangerouslySetInnerHTML` combined with children content.
// expect: react/no-danger-with-children error
const a = (
  <div dangerouslySetInnerHTML={{ __html: html }}>extra</div>
);

// Negative: dangerously-set HTML with no children.
const b = <div dangerouslySetInnerHTML={{ __html: html }} />;

JSON.stringify({ a, b });

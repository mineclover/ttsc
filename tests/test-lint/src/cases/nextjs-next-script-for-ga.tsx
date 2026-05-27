// Positive: hand-written Google Analytics `<script>` tag.
// expect: nextjs/next-script-for-ga error
const a = (
  <script
    async
    src="https://www.googletagmanager.com/gtag/js?id=GA_TRACKING_ID"
  />
);

// Negative: same effect achieved through `next/script` GA helper.
const b = null;

JSON.stringify({ a, b });

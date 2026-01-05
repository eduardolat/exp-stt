## 2024-05-23 - Replaced date-fns with Intl.RelativeTimeFormat

**Learning:** `date-fns` was used in `HistoryItem.svelte` only for `formatDistanceToNow`. Replacing it with the native `Intl.RelativeTimeFormat` API reduces the bundle size and dependency count, aligning with the project's goal of using native browser APIs.
**Action:** Always check if a library is truly needed or if a native API can achieve the same result with minimal effort, especially for simple formatting tasks.

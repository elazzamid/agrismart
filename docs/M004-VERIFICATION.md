# M004 — Farmer Experience Verification

## Slices M004.1–M004.3: Farmer dashboard and crop journal

### Implemented

- Next.js App Router + TypeScript application under `frontend/`.
- Login and registration using the existing Go API.
- Session-scoped bearer token; clear token on sign-out and unauthorized API responses.
- Farmer dashboard lists only farms returned by the authenticated `GET /api/v1/farms` endpoint.
- Farm creation uses `POST /api/v1/farms`; optional location is supported.
- Farm workspace lists and creates plots using the farmer-scoped farm API.
- Plot workspace lists active crops and lists/creates crop cycles with planting dates.
- Crop-cycle journal lists and creates activities, applied inputs, expenses, and harvest records through authenticated APIs.
- Journal forms validate required fields and reflect backend validation; input records remain historical records, not dose recommendations.
- Crop-cycle journal shows activity/input/harvest counts and total recorded expenses, with Indonesian labels for common activity and expense categories.
- Journal panels load only when opened, avoiding unnecessary requests for every crop cycle.
- Explicit loading, empty, success, and error states.
- Responsive desktop and mobile layouts; form labels, keyboard focus styles, and reduced-motion support.
- Same-origin Next.js rewrite proxies `/api/v1/*` to the configured Go API origin.

### Automated checks

CI runs:

- `npm run typecheck`
- `npm run build`
- existing Go tests/build and PostgreSQL migration smoke test

The authoritative result is the GitHub Actions run for the exact frontend branch HEAD; do not infer a pass from code being committed.

### Manual verification checklist

1. Register a new account with a valid email, password, and name.
2. Confirm successful registration leads to login and dashboard.
3. Sign out, then sign in with the same credentials.
4. Verify an account with no farms sees the empty state.
5. Create a farm with a name and optional location; confirm it appears after reload.
6. Open a farm and verify only its plots are shown.
7. Create a plot with a positive area and confirm it appears after reload.
8. Create a crop cycle with an active crop and planting date; confirm it appears under the correct plot.
9. Confirm API failures show a recoverable error rather than fabricated sample data.
10. Check narrow mobile viewport, keyboard navigation, and form labels.
11. Confirm sign-out clears the tab-scoped token.
12. Record one activity, input, expense, and harvest; verify each appears after reload.
13. Check that the default date matches the local calendar date and that input quantity units are explicitly entered.

### Scope and limitations

- Activity/input/expense/harvest forms are included. Browser end-to-end automation and manual browser verification remain outstanding.
- The backend issues bearer tokens; `sessionStorage` is tab-scoped but JavaScript-accessible and is not equivalent to an HttpOnly cookie.
- No production deployment, VPS change, PR merge, or database mutation was performed.
- The feature branch is based on `feature/m003-farm-operations` because the dashboard depends on the farm API; M003 and M004 are not merged into `main`.

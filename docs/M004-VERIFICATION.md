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
- Mobile “Tambah kebun” navigation returns from a selected farm workspace to the dashboard before opening the farm form.
- Responsive desktop and mobile layouts; form labels, keyboard focus styles, and reduced-motion support.
- Same-origin Next.js rewrite proxies `/api/v1/*` to the configured Go API origin.

## Slice M004.4 — Reliability and asynchronous state isolation

- Within the selected farm workspace, the plot empty state appears only after the plot-list API succeeds; failed plot refreshes preserve the last successfully loaded plot list. Failed journal requests remain distinguishable from a genuinely empty journal.
- Plot and crop data remain visible when an individual plot's crop-cycle request fails; the failed cycle list is labelled and has a per-plot retry action.
- A plot shows a cycle-loading state while its first cycle request is pending, not a premature empty-cycle message; successful empty results are tracked separately from never-loaded results.
- A failed initial plot-list request does not display the “Belum ada petak” empty state; after a later refresh failure, existing plots remain visible and their cycle lists are independently retried.
- A crop-catalog API failure does not hide a successfully loaded plot list; crop-cycle creation remains unavailable until the catalog can be loaded.
- If a crop-catalog refresh fails after a successful load, cached crop options remain visible but cycle creation is disabled until a successful refresh.
- Plot data is rendered as soon as the plot request succeeds; a slow crop-catalog request must not block the plot list.
- Journal summary metrics and category counts use placeholders until the initial journal load succeeds.
- If refreshing a plot’s cycle list fails, previously loaded cycles remain visible with a clear stale-data warning.
- Retrying a failed cycle load clears the prior error while pending and restores it only if the retry fails; retrying one plot does not reload unrelated plots.
- A delayed cycle request for one plot does not block successfully loaded cycles or clear the pending state for other plots.
- Switching between farms remounts the workspace so data and load/error states from the prior farm cannot leak into the newly selected farm.
- Overlapping farm-workspace refreshes cannot let an older plot-list or crop-catalog response overwrite a newer refresh or clear its loading state; a per-plot retry started while refresh is pending must not be superseded by that refresh.
- Starting a newer workspace refresh invalidates older per-plot cycle requests; unmounting the workspace invalidates all outstanding cycle requests.
- Dashboard farm-list refreshes ignore stale responses, preserve the latest request's loading state, and invalidate outstanding requests when the dashboard unmounts.
- Crop-cycle and journal save responses are scoped to the form session that initiated them; delayed successes cannot reset newer forms or publish stale success notices, and delayed failures cannot replace newer-session feedback.
- Farm and plot save responses are scoped to the form session that initiated them; delayed successes cannot reset newer forms or publish stale success notices, and delayed failures cannot replace newer-session feedback.
### Automated checks

CI runs:

- `npm test` (Node built-in tests for date-only age, unambiguous stage-range matching, and stale crop-cycle, journal, farm, and plot form-session results)
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
14. On mobile, open a farm workspace, tap “Tambah kebun”, and verify the dashboard form opens.
15. Simulate an initial farm-list API failure and verify the empty-farm message is not shown; after a successful load, fail a refresh and verify the last successful farm list remains visible with a recoverable error.
16. Simulate one crop-cycle API failure and verify the plot/crop list remains visible while that plot's cycle list shows an error.
17. Simulate one journal API failure and verify the journal does not claim there are no records; use “Coba lagi” after restoring the API.
18. Simulate an initial plot-list API failure and verify “Belum ada petak” is not shown; after a successful initial load, fail a later plot-list refresh and verify existing plots remain visible while cycle lists reload independently.
19. Simulate a crop-catalog API failure with successful plot loading and verify existing plots remain visible while crop-cycle creation is unavailable.
20. Simulate an initial journal API failure and verify summary metrics and category counts do not falsely display zero.
21. Delay one crop-cycle API response and verify the pending plot shows a loading state rather than “Belum ada siklus tanam”.
22. Switch between two farms and verify the second farm never displays the first farm’s plots or cycle-loading/error state.
23. After cycles load successfully, simulate a failed refresh and verify the last successful cycle list remains visible with a stale-data warning.
24. Retry a failed cycle request from that plot and verify the old error clears while pending, then returns only if the retry fails; verify other plots are not reloaded.
25. Delay one plot's cycle request and verify other plots render their successfully loaded cycles without waiting for the delayed request.
26. Start overlapping cycle refreshes for one plot and complete them out of order; verify only the newest request updates the list and clears the pending/error state.
27. Load the crop catalog successfully, simulate a failed catalog refresh, and verify cycle creation stays disabled until a refresh succeeds.
28. Start overlapping journal refreshes and complete them out of order; only the newest response may update rows or loading/error state.
29. Refresh the journal after a successful load; verify a pending status explains that the last successful records remain visible until the refresh finishes.
30. Start overlapping workspace refreshes and complete them out of order; verify only the latest refresh controls plots, crop catalog, error, and loading state. While a refresh is pending, retry one plot and verify the refresh preserves that retry result and does not leave its pending indicator stuck.
31. Delay the crop-catalog API while the plot API succeeds; verify plots appear before the catalog request completes.
32. Start a per-plot cycle request, start a newer workspace refresh, then resolve the old request; verify it cannot update cycle data or leave a stale pending indicator. Repeat while unmounting the workspace.
33. Successfully load an empty cycle list, refresh it and fail the request; verify the UI identifies a failed refresh of the last successful empty result rather than an initial-load failure.
34. Start overlapping dashboard farm-list refreshes and complete them out of order; verify only the latest response controls the farm list, error, and loading state, then unmount during a pending request and verify it cannot update state afterward.
35. Submit a farm, plot, or crop-cycle creation request, unmount its view before the POST resolves, and verify completion does not trigger state updates or a follow-up refresh in the unmounted view.
36. Submit a journal record, close the journal before the POST resolves, and verify completion does not trigger state updates or a follow-up journal refresh after unmount.
37. Start registration, unmount the auth screen before the registration response resolves, and verify it does not start login or update auth-screen state afterward.
38. Start login, unmount the auth screen before the login response resolves, and verify it does not save a token, call the authenticated callback, or update auth-screen state afterward.
39. With an existing session token, start the `/api/v1/auth/me` session-restoration request, unmount `HomePage` before it resolves, and verify the response cannot update page state after unmount.
40. Select a crop with active varieties, verify only that crop’s varieties are shown, choose one, save the crop cycle, and verify the request includes the selected `variety_id`.
41. Change crops while a variety request is pending and resolve the requests out of order; verify only the latest crop’s options appear. Simulate a variety-catalog failure and verify retry is available and saving stays disabled until the catalog loads; if the catalog is empty, verify a cycle can be saved without a variety.
42. For an active cycle, verify age is calculated from the local calendar date using date-only arithmetic: planting today is 0 days, yesterday is 1 day, and a future or malformed date does not show a numeric age.
43. Verify planned, harvested, and cancelled cycles do not display an active growth-stage label.
44. For an active cycle with exactly one matching catalog range, verify the indicative sample-stage label appears; gaps, overlapping or open-ended ranges, invalid bounds, and missing ranges must not be guessed.
45. Simulate a growth-stage catalog failure; verify age remains visible, retry is available, and no fallback or fabricated stage is shown.
46. Verify the UI labels stage ranges as sample metadata and does not present them as cultivation advice, plant-condition assessment, or harvest prediction.
47. Start an authenticated request with token A, replace the stored session with token B before the request returns HTTP 401, and verify the stale response neither clears token B nor dispatches `agrismart:unauthorized`. Also verify a 401 from login/registration (a request without a session token) cannot invalidate a stored session.
48. Keep an active crop cycle open across local midnight and verify its age and indicative stage are recalculated for the new local calendar day. Background the tab across midnight, return to it, and verify the age refreshes without requiring a manual reload.
49. Verify the growth-stage catalog is not fetched again merely because the local crop age advances at midnight; changing crop identity or transitioning into an eligible active cycle must still load the corresponding catalog.
50. Open an active crop cycle in a browser whose timezone is ahead of UTC during the local-date/UTC-date crossover (for example, 00:00–07:59 in Asia/Makassar). Verify the initial render does not show a hydration warning or a one-day age discrepancy, and that the displayed age matches the browser’s local calendar date after hydration.
51. Start saving a crop cycle in one plot, then switch to or edit another cycle form before the first POST resolves. Verify the delayed success cannot close or reset the newer form or display a stale success notice; verify a failure remains visible and the user-entered values are preserved.
52. Start a crop-cycle save that will fail, then switch to or edit a newer cycle form before the rejection resolves. Verify the stale failure does not overwrite feedback for the newer form; if the submitted form remains current, verify the failure is still shown and its entered values are preserved.
53. Start a journal save, then switch journal category, close/reopen the form, or edit another journal form before the POST resolves. Verify delayed success cannot close or reset the newer form or display a stale success notice, and delayed failure cannot replace its feedback; when the submitted form remains current, verify normal success reset or recoverable error behavior.
54. For both crop-cycle and journal forms, complete an older save successfully after the user has switched to or edited a newer form. Verify the newer form remains open with its entered values and no success notice from the older request appears; verify a success notice still appears when the submitted form remains current.
55. Start saving a plot, then cancel/reopen or edit the plot form before the POST resolves. Verify delayed success cannot close/reset the newer form or show stale success feedback, and delayed failure cannot replace newer-form feedback; when the submitted form remains current, verify normal reset and success behavior.
56. Start saving a farm, then cancel the form before the POST resolves. Verify cancellation closes the form and clears its draft; while the request remains pending, show a visible status explaining that a new form is temporarily unavailable. Repeated form-entry actions while the save is pending must not create a new session. After the POST settles, reopen the form and verify it is blank. Separately, edit the form before the POST resolves and verify delayed success cannot close/reset the edited form or show stale success feedback, and delayed failure cannot replace newer-form feedback; when the submitted form remains current, verify normal reset and success behavior.
57. Start saving a farm, open an existing farm workspace before the POST resolves, then resolve success while the workspace is active. Verify the submitted form closes/resets without queuing a hidden success notice, the created farm appears after refresh, and returning to the dashboard does not leave a duplicate-ready form. Repeat with a failed POST and verify the entered values and error remain when returning. After that failure, use mobile “Tambah kebun” to start a new form and verify the prior save error is cleared; if the current farm form is edited before either response, its edited values and newer feedback must remain untouched. After the farm POST has settled but its list refresh is still pending, open or edit a newer form and fail the refresh; verify the new form remains usable and the list-refresh error is shown separately without replacing newer-form save feedback. Also verify that canceling while a save is pending clears the canceled draft, and repeated form-entry actions (including mobile “Tambah kebun”) do not create a new form session until the pending save has settled.

### Scope and limitations

- Active crop cycles display elapsed calendar days from the planting date and an indicative label only when exactly one growth-stage catalog range matches. Planned, harvested, cancelled, unknown-status, and future-dated cycles do not receive an active-stage label. The seeded ranges are sample metadata, not agronomic recommendations; the UI explicitly disclaims condition assessment and harvest prediction. Browser end-to-end automation and manual browser verification remain outstanding. The date and range-matching helper is covered by `npm test` using Node’s built-in test runner.
- A 401 response invalidates authentication only when the request used the currently stored non-empty session token; delayed responses from a prior session and unauthenticated login/registration failures cannot clear a newer session.
- Active crop age refreshes at the next local calendar midnight and when a backgrounded tab becomes visible; the timer uses the actual local-day length rather than assuming every day is 24 hours. The initial age is deferred until client mount so server and browser timezone differences cannot produce mismatched hydration markup.
- The backend issues bearer tokens; `sessionStorage` is tab-scoped but JavaScript-accessible and is not equivalent to an HttpOnly cookie.
- No production deployment, VPS change, PR merge, or database mutation was performed.
- The feature branch is based on `feature/m003-farm-operations` because the dashboard depends on the farm API; M003 and M004 are not merged into `main`.

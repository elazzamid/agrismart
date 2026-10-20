# AgriSmart Farmer Web

Mobile-first Next.js interface for the AgriSmart Go API.

## Requirements

- Node.js 20 or newer
- AgriSmart Go API running on `http://127.0.0.1:8080` by default
- PostgreSQL initialized from `database/migrations`

## Run locally

```sh
cd frontend
npm install
npm run dev
```

Open `http://localhost:3000`. The Next.js rewrite proxies `/api/v1/*` to the Go API so browser requests use the same origin and do not need permissive API CORS settings.

Set `AGRISMART_API_ORIGIN` to override the API origin, for example:

```sh
AGRISMART_API_ORIGIN=http://127.0.0.1:8080 npm run dev
```

## Current M004 slices

- Login and account registration with tab-scoped bearer-token storage.
- Farmer dashboard, list/create farm, and farm-scoped plot management.
- List/create crop cycles using the existing active-crop catalog.
- Record and review activities, applied inputs, expenses, and harvests for each crop cycle.
- Loading, empty, success, and recoverable error states.
- Responsive desktop/mobile layouts and accessible form labels.

The interface does not fabricate farm or crop data. Input records document what was actually used; the journal does not calculate or recommend agricultural doses. Browser end-to-end automation and manual browser verification remain outstanding.

## Security note

The current backend returns bearer tokens rather than secure cookies. The frontend keeps the token in tab-scoped `sessionStorage` and clears it on sign-out or a 401 response. This reduces persistence compared with long-lived local storage, but is not equivalent to an HttpOnly cookie session.

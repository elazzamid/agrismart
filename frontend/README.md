# AgriSmart Farmer Web

Mobile-first Next.js interface for the AgriSmart API.

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

## First M004 slice

- Login and account registration against the existing API.
- Session-scoped bearer token storage.
- Farmer dashboard and farmer-owned farm list.
- Create farm with optional location.
- Loading, empty, success, and error states.
- Responsive desktop/mobile layouts and accessible form labels.

The interface intentionally does not fabricate farm, crop, activity, input, expense, or harvest data. Plot, crop-cycle, and operation-record workflows are subsequent M004 slices.

## Security note

The current backend returns bearer tokens rather than secure cookies. The frontend keeps the token in tab-scoped `sessionStorage` and clears it on sign-out or a 401 response. This reduces persistence compared with long-lived local storage, but is not equivalent to an HttpOnly cookie session.

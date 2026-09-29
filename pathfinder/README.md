# Pathfinder frontend

SvelteKit interface for the Pathfinder experiment hosted within Cheap Eats.

```sh
npm ci
npm run dev
```

The Go server must run on port 8080. Open `/pathfinder/` on the Vite development server. Requests to `/pathfinder/api/search` are proxied to Go.

`npm run check` checks types and Svelte components. `npm run build` generates the static site in `build/`; the Cheap Eats server serves it at `/pathfinder/` in production.

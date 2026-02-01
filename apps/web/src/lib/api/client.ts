/**
 * Type-safe API client for MacroGuard backend
 * Generated types come from openapi.yaml via openapi-typescript
 *
 * In development, Vite proxies /api/* to the backend (see vite.config.ts)
 * This allows the app to work when accessed from any device (e.g., mobile testing)
 */
import createClient from "openapi-fetch";
import type { paths } from "./schema";

// Use relative URLs - in dev mode, Vite proxy handles /api/* requests
// In production, the API will be served from the same origin
export const api = createClient<paths>({
  baseUrl: ""
});

// Re-export types for convenience
export type { paths, components, operations } from "./schema";


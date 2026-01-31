/**
 * Type-safe API client for MacroGuard backend
 * Generated types come from openapi.yaml via openapi-typescript
 */
import createClient from "openapi-fetch";
import type { paths } from "./schema";

// Create the API client with base URL
// In development, the API runs on a different port
const API_BASE_URL = import.meta.env.DEV ? "http://localhost:8080" : "";

export const api = createClient<paths>({
    baseUrl: API_BASE_URL
});

// Re-export types for convenience
export type { paths, components, operations } from "./schema";

import type { Handle } from "@sveltejs/kit";

/**
 * Server hook to inject runtime environment variables into the page
 * This allows PUBLIC_API_URL to be configured at runtime, not build time
 */
export const handle: Handle = async ({ event, resolve }) => {
  // Make PUBLIC_API_URL available to the client via $env/dynamic/public
  // This is read from the server's environment at runtime
  return resolve(event);
};

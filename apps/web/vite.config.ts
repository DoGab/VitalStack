import { sveltekit } from "@sveltejs/kit/vite";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";
import { existsSync, readFileSync } from "node:fs";
import { resolve, dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

// Get __dirname equivalent in ESM
const __dirname = dirname(fileURLToPath(import.meta.url));

// Check if local certs exist for HTTPS (needed for mobile camera access)
const certPath = resolve(__dirname, ".cert");
const keyFile = join(certPath, "key.pem");
const certFile = join(certPath, "cert.pem");
const hasLocalCerts = existsSync(keyFile) && existsSync(certFile);

export default defineConfig({
  plugins: [tailwindcss(), sveltekit()],
  server: {
    // Enable access from other devices on the network
    host: true,
    // Use HTTPS if certs are available (required for camera on non-localhost)
    https: hasLocalCerts
      ? {
        key: readFileSync(keyFile),
        cert: readFileSync(certFile)
      }
      : undefined,
    // Proxy API requests to Go backend
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true
      }
    }
  }
});

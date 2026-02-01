#!/bin/bash
# Generate self-signed SSL certificates for local development
# Required for mobile camera access when testing from other devices

set -e

CERT_DIR="$(dirname "$0")/.cert"
mkdir -p "$CERT_DIR"

# Generate certificate
openssl req -x509 \
  -newkey rsa:2048 \
  -keyout "$CERT_DIR/key.pem" \
  -out "$CERT_DIR/cert.pem" \
  -days 365 \
  -nodes \
  -subj "/CN=localhost" \
  -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

echo "✅ SSL certificates generated in $CERT_DIR"
echo ""
echo "📱 To test on mobile, you'll need to trust the certificate:"
echo "   1. Access https://<your-ip>:5173 on your mobile browser"
echo "   2. The browser will warn about the certificate - accept/trust it"
echo "   3. Camera access should now work"

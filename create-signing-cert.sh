#!/bin/bash
# Bash Script to Create Self-Signed Code Signing Certificate
# Works on macOS and Linux

echo "🔐 Creating Self-Signed Code Signing Certificate for Spotify Shuffle"
echo ""

# Set password
PASSWORD="SpotifyShuffleSign2024!"
echo "Using password: $PASSWORD"
echo ""

# Create output directory
mkdir -p cert-output
cd cert-output

echo "Creating private key..."
openssl genrsa -out spotify-shuffle.key 2048

echo "Creating certificate signing request..."
openssl req -new -key spotify-shuffle.key -out spotify-shuffle.csr -subj "/CN=Spotify Shuffle/O=Spotify Shuffle/C=US"

echo "Creating self-signed certificate..."
openssl x509 -req -days 1095 -in spotify-shuffle.csr -signkey spotify-shuffle.key -out spotify-shuffle.crt

echo "Converting to PKCS#12 format..."
openssl pkcs12 -export -out spotify-shuffle-cert.p12 -inkey spotify-shuffle.key -in spotify-shuffle.crt -password pass:$PASSWORD

echo "Converting to Base64 for GitHub..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    base64 -i spotify-shuffle-cert.p12 -o cert-base64.txt
else
    # Linux
    base64 -w 0 spotify-shuffle-cert.p12 > cert-base64.txt
fi

echo ""
echo "🎉 SUCCESS! Files created:"
echo "📁 Certificate: $(pwd)/spotify-shuffle-cert.p12"
echo "📁 Base64: $(pwd)/cert-base64.txt"
echo ""
echo "📋 GitHub Secrets to Add:"
echo "========================="
echo "Secret Name: WINDOWS_CERT_BASE64"
echo "Secret Value: (contents of cert-base64.txt)"
echo ""
echo "Secret Name: WINDOWS_CERT_PASSWORD" 
echo "Secret Value: $PASSWORD"
echo ""
echo "📖 Next Steps:"
echo "1. Go to: https://github.com/petabloc/spotify-shuffle/settings/secrets/actions"
echo "2. Click 'New repository secret'"
echo "3. Add both secrets above"
echo "4. Your next release will be automatically signed!"
echo ""
echo "⚠️  Note: Self-signed certificates will show 'Unknown Publisher'"
echo "   Users can still install by clicking 'More info' -> 'Run anyway'"
echo ""

# Clean up intermediate files
rm spotify-shuffle.key spotify-shuffle.csr spotify-shuffle.crt

echo "🧹 Cleaned up intermediate files (kept .p12 and base64)"
echo "Press Enter to continue..."
read
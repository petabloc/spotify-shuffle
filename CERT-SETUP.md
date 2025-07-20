# Code Signing Certificate Setup

This guide will help you create a self-signed certificate for Windows code signing and add it to GitHub Actions.

## 🚀 Quick Setup

### For Windows Users:
1. **Run PowerShell as Administrator**
2. **Navigate to this directory**
3. **Run the certificate creation script:**
   ```powershell
   .\create-signing-cert.ps1
   ```

### For macOS/Linux Users:
1. **Open Terminal**
2. **Navigate to this directory** 
3. **Run the certificate creation script:**
   ```bash
   ./create-signing-cert.sh
   ```

## 📋 Adding to GitHub

After running the script, you'll have two files in `cert-output/`:
- `spotify-shuffle-cert.p12` - The certificate file
- `cert-base64.txt` - Base64 encoded certificate for GitHub

### Steps to add to GitHub:

1. **Go to Repository Settings:**
   https://github.com/petabloc/spotify-shuffle/settings/secrets/actions

2. **Add Secret #1:**
   - **Name:** `WINDOWS_CERT_BASE64`
   - **Value:** Copy the entire contents of `cert-output/cert-base64.txt`

3. **Add Secret #2:**
   - **Name:** `WINDOWS_CERT_PASSWORD`
   - **Value:** `SpotifyShuffleSign2024!`

## ✅ Verification

Once added, your next release will automatically:
- ✅ Sign the Windows executable
- ✅ Sign the Windows MSI installer  
- ✅ Include signed binaries in the GitHub release

## ⚠️ Important Notes

**Self-Signed Certificate Limitations:**
- Windows will still show "Unknown Publisher" warning
- Users need to click "More info" → "Run anyway" to install
- Provides integrity verification but not publisher trust

**For Production Use:**
- Consider purchasing a commercial certificate from SSL.com (~$199/year)
- Commercial certificates eliminate the "Unknown Publisher" warning
- Provides immediate trust for your users

## 🔒 Security

- The certificate password is stored as a GitHub secret (encrypted)
- The certificate itself is stored as a GitHub secret (encrypted) 
- Signing happens only during GitHub Actions builds
- No certificates are stored in the repository code

## 🆙 Upgrading to Commercial Certificate

When you're ready for a commercial certificate:
1. Purchase from SSL.com or similar provider
2. Follow their verification process
3. Replace the GitHub secrets with your commercial certificate
4. Next release will be fully trusted by Windows

---

**Need help?** Check the GitHub Issues or documentation for more details.
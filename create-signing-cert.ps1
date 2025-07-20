# PowerShell Script to Create Self-Signed Code Signing Certificate
# Run this as Administrator in PowerShell

Write-Host "🔐 Creating Self-Signed Code Signing Certificate for Spotify Shuffle" -ForegroundColor Green
Write-Host ""

# Generate a secure password
$password = "SpotifyShuffleSign2024!"
Write-Host "Using password: $password" -ForegroundColor Yellow
Write-Host ""

try {
    # Create the certificate
    Write-Host "Creating certificate..." -ForegroundColor Blue
    $cert = New-SelfSignedCertificate -Subject "CN=Spotify Shuffle,O=Spotify Shuffle,C=US" `
        -Type CodeSigningCert `
        -KeyUsage DigitalSignature `
        -FriendlyName "Spotify Shuffle Code Signing Certificate" `
        -CertStoreLocation Cert:\CurrentUser\My `
        -HashAlgorithm SHA256 `
        -Provider "Microsoft Enhanced RSA and AES Cryptographic Provider" `
        -KeyLength 2048 `
        -NotAfter (Get-Date).AddYears(3)
    
    Write-Host "✅ Certificate created successfully!" -ForegroundColor Green
    Write-Host "Thumbprint: $($cert.Thumbprint)" -ForegroundColor Gray
    Write-Host ""
    
    # Create output directory
    $outputDir = ".\cert-output"
    if (!(Test-Path $outputDir)) {
        New-Item -ItemType Directory -Path $outputDir | Out-Null
    }
    
    # Export to PFX file
    Write-Host "Exporting certificate to PFX..." -ForegroundColor Blue
    $pwd = ConvertTo-SecureString -String $password -Force -AsPlainText
    $pfxPath = "$outputDir\spotify-shuffle-cert.p12"
    Export-PfxCertificate -Cert $cert -FilePath $pfxPath -Password $pwd | Out-Null
    
    # Convert to Base64 for GitHub
    Write-Host "Converting to Base64 for GitHub..." -ForegroundColor Blue
    $bytes = [System.IO.File]::ReadAllBytes($pfxPath)
    $base64 = [System.Convert]::ToBase64String($bytes)
    
    # Save Base64 to file
    $base64Path = "$outputDir\cert-base64.txt"
    $base64 | Out-File -FilePath $base64Path -Encoding ascii
    
    Write-Host ""
    Write-Host "🎉 SUCCESS! Files created:" -ForegroundColor Green
    Write-Host "📁 Certificate: $pfxPath" -ForegroundColor White
    Write-Host "📁 Base64: $base64Path" -ForegroundColor White
    Write-Host ""
    Write-Host "📋 GitHub Secrets to Add:" -ForegroundColor Yellow
    Write-Host "=========================" -ForegroundColor Yellow
    Write-Host "Secret Name: WINDOWS_CERT_BASE64" -ForegroundColor Cyan
    Write-Host "Secret Value: (contents of $base64Path)" -ForegroundColor Gray
    Write-Host ""
    Write-Host "Secret Name: WINDOWS_CERT_PASSWORD" -ForegroundColor Cyan
    Write-Host "Secret Value: $password" -ForegroundColor Gray
    Write-Host ""
    Write-Host "📖 Next Steps:" -ForegroundColor Yellow
    Write-Host "1. Go to: https://github.com/petabloc/spotify-shuffle/settings/secrets/actions"
    Write-Host "2. Click 'New repository secret'"
    Write-Host "3. Add both secrets above"
    Write-Host "4. Your next release will be automatically signed!"
    Write-Host ""
    Write-Host "⚠️  Note: Self-signed certificates will show 'Unknown Publisher'" -ForegroundColor Red
    Write-Host "   Users can still install by clicking 'More info' -> 'Run anyway'" -ForegroundColor Red
    
} catch {
    Write-Host "❌ Error creating certificate: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Make sure you're running PowerShell as Administrator!" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Press any key to continue..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
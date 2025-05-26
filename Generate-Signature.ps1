# Check if ActiveDirectory module is available
if (-not (Get-Module -ListAvailable -Name ActiveDirectory)) {
    Write-Host "ERROR: Active Directory module for PowerShell is not installed." -ForegroundColor Red
    Write-Host "Please install 'RSAT: Active Directory module for Windows PowerShell'." -ForegroundColor Yellow
    exit 1
}

# Import the module
Import-Module ActiveDirectory

# Get the current user's username (only the username part, no DOMAIN)
$currentUser = [Environment]::UserName

# Query Active Directory
try {
    Write-Host "Querying Active Directory for user '$currentUser'..." -ForegroundColor Gray
    $adUser = Get-ADUser -Identity $currentUser -Properties DisplayName, Mail, telephoneNumber
} catch {
    Write-Host "ERROR: Could not find user '$currentUser' in Active Directory." -ForegroundColor Red
    Write-Host "Details: $_" -ForegroundColor Yellow
    exit 2
}

# Extract fields
$displayName = $adUser.DisplayName
$email = $adUser.Mail
$telephoneNumber = $adUser.telephoneNumber

# Validate fields
$missingFields = @()
if ([string]::IsNullOrEmpty($displayName)) { $missingFields += "DisplayName" }
if ([string]::IsNullOrEmpty($email)) { $missingFields += "Mail" }
if ([string]::IsNullOrEmpty($telephoneNumber)) { 
    Write-Host "WARNING: Phone number is missing. Signature will be generated without it." -ForegroundColor Yellow
}

if ($missingFields.Count -gt 0) {
    Write-Host "ERROR: Required attributes are missing: $($missingFields -join ', ')" -ForegroundColor Red
    exit 3
}

# Output user information
Write-Host "`nUser Information:" -ForegroundColor Cyan
Write-Host "DisplayName: $displayName" -ForegroundColor White
Write-Host "Email: $email" -ForegroundColor White
Write-Host "Phone: $telephoneNumber" -ForegroundColor White

# Check if SignatureInstaller exists
if (-not (Test-Path ".\SignatureInstaller.exe")) {
    Write-Host "ERROR: SignatureInstaller.exe not found in current directory." -ForegroundColor Red
    exit 4
}

# Call SignatureInstaller
Write-Host "`nGenerating signature..." -ForegroundColor Cyan
try {
    & .\SignatureInstaller.exe -name "$displayName" -email "$email" -phone "$telephoneNumber"
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Signature generated successfully!" -ForegroundColor Green
    } else {
        Write-Host "ERROR: SignatureInstaller failed with exit code $LASTEXITCODE" -ForegroundColor Red
        exit 5
    }
} catch {
    Write-Host "ERROR: Failed to run SignatureInstaller." -ForegroundColor Red
    Write-Host "Details: $_" -ForegroundColor Yellow
    exit 6
}

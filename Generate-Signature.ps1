# Check if ActiveDirectory module is available
if (-not (Get-Module -ListAvailable -Name ActiveDirectory)) {
    Write-Host "ERROR: Active Directory module for PowerShell is not installed."
    Write-Host "Please install 'RSAT: Active Directory module for Windows PowerShell'."
    exit 1
}

# Import the module
Import-Module ActiveDirectory

# Get the current user's username (only the username part, no DOMAIN)
$currentUser = [Environment]::UserName

# Query Active Directory
try {
    $adUser = Get-ADUser -Identity $currentUser -Properties *
} catch {
    Write-Host "ERROR: Could not find user '$currentUser' in Active Directory."
    exit 2
}

# Extract fields
$displayName = $adUser.DisplayName
$email = $adUser.Mail
$username = $adUser.SamAccountName
$telephoneNumber = $adUser.telephoneNumber

# Validate fields
if ([string]::IsNullOrEmpty($displayName) -or [string]::IsNullOrEmpty($email) -or [string]::IsNullOrEmpty($username)) {
    Write-Host "ERROR: Some required attributes are missing (DisplayName, Mail, or SamAccountName)."
    exit 3
}

# Output (for debug/log)
Write-Host "DisplayName: $displayName"
Write-Host "Email: $email"
Write-Host "Email: $telephoneNumber"
Write-Host "Username: $username"

# Call Go tool
.\SignatureInstaller.exe -name "$displayName" -email "$email" -phone "$telephoneNumber"

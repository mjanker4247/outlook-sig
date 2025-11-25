# Check if ActiveDirectory module is available
if (-not (Get-Module -ListAvailable -Name ActiveDirectory))
{
    Write-Host "ERROR: Active Directory module for PowerShell is not installed." -ForegroundColor Red
    Write-Host "Please install 'RSAT: Active Directory module for Windows PowerShell'." -ForegroundColor Yellow
    exit 1
}

# Resolve folder where this script is stored
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$installerPath = Join-Path $scriptDir "SignatureInstaller.exe"

# Check if SignatureInstaller.exe exists in the script folder
if (-not (Test-Path $installerPath))
{
    Write-Host "ERROR: SignatureInstaller.exe not found in script directory: $scriptDir" -ForegroundColor Red
    Write-Host "Please place SignatureInstaller.exe in the same folder as this script." -ForegroundColor Yellow
    exit 2
}

# Import the module
Import-Module ActiveDirectory

# Get the current user's username (only the username part, no DOMAIN)
$currentUser = [Environment]::UserName

# Query Active Directory
try
{
    Write-Host "Querying Active Directory for user '$currentUser'..." -ForegroundColor Gray
    $adUser = Get-ADUser -Identity $currentUser -Properties DisplayName, Mail, telephoneNumber
}
catch
{
    Write-Host "ERROR: Could not find user '$currentUser' in Active Directory." -ForegroundColor Red
    Write-Host "Details: $_" -ForegroundColor Yellow
    exit 2
}

# Extract fields
$displayName = $adUser.DisplayName
$email = $adUser.Mail
$telephoneNumber = $adUser.telephoneNumber

# Validate required fields
$missingFields = @()
if ([string]::IsNullOrEmpty($displayName)) { $missingFields += "DisplayName" }
if ([string]::IsNullOrEmpty($email)) { $missingFields += "Mail" }
if ([string]::IsNullOrEmpty($telephoneNumber)) { $missingFields += "telephoneNumber" }

if ($missingFields.Count -gt 0)
{
    Write-Host "ERROR: Required attributes are missing: $($missingFields -join ', ')" -ForegroundColor Red
    Write-Host "All fields (DisplayName, Mail, telephoneNumber) are required for signature generation." -ForegroundColor Yellow
    exit 3
}

# Output user information
Write-Host "`nUser Information:" -ForegroundColor Cyan
Write-Host "DisplayName: $displayName" -ForegroundColor White
Write-Host "Email: $email" -ForegroundColor White
Write-Host "Phone: $telephoneNumber" -ForegroundColor White

# Format phone number for better compatibility with Go validation
if (-not [string]::IsNullOrEmpty($telephoneNumber))
{
    # Remove common separators and ensure proper format
    $telephoneNumber = $telephoneNumber -replace '[\s\-\(\)\.]', ''
    
    # Add country code if missing (assuming German numbers)
    if (-not $telephoneNumber.StartsWith("+"))
    {
        if ($telephoneNumber.StartsWith("0"))
        {
            # Replace leading 0 with +49 for German numbers
            $telephoneNumber = "+49" + $telephoneNumber.Substring(1)
        }
        else
        {
            # Add +49 prefix for numbers without country code
            $telephoneNumber = "+49" + $telephoneNumber
        }
    }
    
    Write-Host "Formatted phone number: $telephoneNumber" -ForegroundColor Gray
}

# Call SignatureInstaller with all required parameters
Write-Host "`nGenerating signature..." -ForegroundColor Cyan
try
{
    # Call with all required parameters (name, email, phone)
    Start-Process $installerPath -ArgumentList @(
        "--name", $displayName,
        "--email", $email,
        "--phone", $telephoneNumber
    )
    -Wait -NoNewWindow
    
    if ($LASTEXITCODE -eq 0)
    {
        Write-Host "Signature generated successfully!" -ForegroundColor Green
    }
    else
    {
        Write-Host "ERROR: SignatureInstaller failed with exit code $LASTEXITCODE" -ForegroundColor Red
        Write-Host "This might be due to validation issues. Please check:" -ForegroundColor Yellow
        Write-Host "- Name format (letters, spaces, dots, hyphens, apostrophes only)" -ForegroundColor Yellow
        Write-Host "- Email format (must be valid email address)" -ForegroundColor Yellow
        Write-Host "- Phone number format (must be valid international format)" -ForegroundColor Yellow
        Write-Host "For phone numbers, ensure they include country code (e.g., +49 for Germany)" -ForegroundColor Yellow
        exit 5
    }
}
catch
{
    Write-Host "ERROR: Failed to run SignatureInstaller." -ForegroundColor Red
    Write-Host "Details: $_" -ForegroundColor Yellow
    exit 6
}

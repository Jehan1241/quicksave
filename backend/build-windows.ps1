# Load environment variables from .env file
if (Test-Path .env) {
    Get-Content .env | ForEach-Object {
        $var = $_ -split '='
        if ($var.Count -eq 2) {
            [System.Environment]::SetEnvironmentVariable($var[0], $var[1])
        }
    }
}

# Set output filename
$OUTPUT = "thismodule.exe"

# Read environment variables
$clientID = [System.Environment]::GetEnvironmentVariable("IGDB_API_KEY")
$clientSecret = [System.Environment]::GetEnvironmentVariable("IGDB_SECRET_KEY")

# Build with embedded API keys
go build -ldflags "-X 'main.clientID=$clientID' -X 'main.clientSecret=$clientSecret'" -o $OUTPUT

Write-Host "✅ Build complete: $OUTPUT"
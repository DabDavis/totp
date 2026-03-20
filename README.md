# totp

Minimal CLI tool for generating 6-digit TOTP 2FA codes ([RFC 6238](https://datatracker.ietf.org/doc/html/rfc6238)). Zero dependencies — Go stdlib only.

Works with GitHub, Google, AWS, and any service that uses standard TOTP.

## Install

```bash
go install github.com/DabDavis/totp@latest
```

Or build from source:

```bash
git clone https://github.com/DabDavis/totp.git
cd totp
go build -o totp .
```

## Usage

```bash
# Generate a code from a secret
totp YOURSECRETHERE
# 482931  (27s remaining)

# Or save the secret and run without args
echo 'YOURSECRETHERE' > ~/.totp && chmod 600 ~/.totp
totp

# Setup instructions
totp setup
```

## GitHub 2FA Setup

1. Go to **github.com → Settings → Password and authentication**
2. Click **Enable two-factor authentication**
3. Click **"setup key"** under the QR code to get the base32 secret
4. Run `totp <secret>` and enter the 6-digit code on GitHub to verify
5. Download and save your recovery codes
6. Save the secret for future use:
   ```bash
   echo '<secret>' > ~/.totp && chmod 600 ~/.totp
   ```

## How It Works

TOTP is HMAC-SHA1 over a time-based counter (30-second intervals), truncated to 6 digits per [RFC 4226](https://datatracker.ietf.org/doc/html/rfc4226). The entire implementation is ~40 lines of Go.

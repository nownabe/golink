# Development

## Proto

```shell
# Update
buf mod update proto

# Lint
buf lint proto

# Build
buf build proto

# Generate
buf generate proto
```

## Test

Run Firestore emulator

```shell
gcloud emulators firestore start --host-port="$FIRESTORE_EMULATOR_HOST"
```

## Extension

### Test

Build extension

```shell
cd extension
npm run build
```

then "Load unpacked" or reload the unpacked extension on chrome://extensions

### Release

Bump up the manifest version and build release

```shell
cd extension
vi manifest.release.json
npm run build:release
```

then upload `extension.zip` to [Chrome Web Store Developer Dashboard](https://chrome.google.com/u/2/webstore/devconsole) and
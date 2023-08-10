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
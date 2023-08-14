#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

region="$1"
root="$(cd "$(dirname "${BASH_SOURCE:-$0}")" && pwd)"

gcloud components install alpha --quiet
gcloud services enable \
	iap.googleapis.com \
	appengine.googleapis.com \
	cloudresourcemanager.googleapis.com \
	firestore.googleapis.com

gcloud app create --region "$region"

cd "$root/redirector"
gcloud app deploy --quiet

cd "$root/api"
gcloud app deploy --quiet

cd "$root"
gcloud app deploy dispatch.yaml --quiet

gcloud alpha firestore databases update --type=firestore-native --quiet

gcloud iap settings set \
	iap-settings.yaml \
	--resource-type=app-engine \
	--project="$(gcloud config get project)"

#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

region="$1"
root="$(cd "$(dirname "${BASH_SOURCE:-$0}")" && pwd)"

deploy_redirector() {
	cd "$root/redirector"
	gcloud app deploy

}

gcloud components install alpha --quiet
gcloud services enable \
	iap.googleapis.com \
	appengine.googleapis.com \
	cloudresourcemanager.googleapis.com

gcloud app create --region "$region"
gcloud alpha firestore databases update --type=firestore-native --quiet

cd "$root/redirector"
gcloud app deploy --quiet

cd "$root/api"
gcloud app deploy --quiet

cd "$root"
gcloud app deploy dispatch.yaml --quiet

gcloud iap settings set \
	iap-settings.yaml \
	--resource-type=app-engine \
	--project="$(gcloud config get project)"

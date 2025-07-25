#!/bin/bash

set -e

# Set variables
# cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/.."
# source vars.sh
az_tenant_id="0b79ac7b-0cc7-4d9f-a549-3b8cc894ac9b"
az_subscription_id="fb5048f4-e435-4bc7-aad8-91382889db7d"
tag_creator_name="Josh Dinndorf"
azure_tagger_function_name="acme-func-azure-tagger-qwdi"
tagger_storage_account="acmestazuretaggerqwdi"
resource_group="acme-rg-azure-tagger-qwdi"
container_name="acme-function-releases"
# Generate datetime and UUID for the filename
datetime=$(date +"%Y%m%d%H%M%S")
uuid=$( uuidgen | tr '[:upper:]' '[:lower:]' )
zip_filename="${datetime}-${uuid}.zip"

# Upload arc.zip to Azure Blob Storage
az storage blob upload \
  --account-name $tagger_storage_account \
  --container-name $container_name \
  --name $zip_filename \
  --file arc-linux-amd64.zip

# Construct the blob URL
blob_url="https://${tagger_storage_account}.blob.core.windows.net/${container_name}/${zip_filename}"

# Update WEBSITE_RUN_FROM_PACKAGE in the function app settings
az functionapp config appsettings set \
  --name $azure_tagger_function_name \
  --resource-group $resource_group \
  --settings WEBSITE_RUN_FROM_PACKAGE=$blob_url

az functionapp restart \
  --name $azure_tagger_function_name \
  --resource-group $resource_group

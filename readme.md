# Automated Resource Tagging with Terraform and Azure Functions

## Overview

This repository provides an automated solution for tagging newly created Azure resources with essential metadata such as the creator's name, email, and date of creation. It leverages Terraform to provision Azure resources and an Azure Function App to handle tagging dynamically. This setup ensures that all new resources are tagged consistently and accurately, enhancing resource management and compliance.

## Benefits

- **Automated Tagging**: Automatically adds tags to newly created resources, ensuring consistency across your Azure environment.
- **Custom Metadata**: Tags include the creator's name, email, and date of creation, which helps in tracking and accountability.
- **Passwordless Authentication**: Utilizes managed identities and app registrations for secure, passwordless authentication.
- **Infrastructure as Code**: Uses Terraform to provision and manage Azure resources, providing reproducibility and ease of management.

## Components

1. **Terraform**: Provisions Azure resources including:

   - Resource Group
   - Function App
   - Storage Account
   - Storage Container
   - Application Insights and Log Analytics
   - Event Grid Topic and Subscription
   - User Assigned Identities and Roles

2. **Azure Function App**: Receives events from Azure Event Grid and applies tags to newly created resources based on the event data.

Calls from Function App to storage account are authenticated using Managed Identity. Storage Blob Data Contributor permission granted to Function App to aquire locks in blob.

Calls from Function App to Azure Management are authenticated using Managed Identity. Minimally required permissions granted (Tag Contributor & Reader) and attached to subscription

Calls from Event Grid to Function App are authenticated using Entra ID. Requests authenticated with corrrect Tenant ID + Application ID (Microsoft.EventGrid) + Audience of Function App are permitted

## Prerequisites

Before using this repository, ensure you have:

- [Terraform](https://www.terraform.io/downloads) installed.
- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) installed and configured.
- An Azure subscription with sufficient permissions to create and manage resources.

## Setup and Configuration

### 1. Clone the Repository

```bash
git clone https://github.com/vlche/azure-tagger.git
cd azure-tagger
```

### 2. Configure Terraform

1. Initialize Terraform: Run the following command to initialize Terraform and download necessary providers.

```bash
cd tf
terraform init
```

2. Configure Variables: Update the `tf/variables.tf` and `vars.sh` files with your specific configuration, including names, regions, and any other required parameters.

3. Plan and Apply: Generate an execution plan and apply it to provision the resources.

```bash
scripts/tf_plan.sh
scripts/tf_apply.sh
```

4. The apply will fail during the initial phase, since no Function is deployed into Function App yet

5. Upload the Function into Function App

6. Plan and Apply: Generate an execution plan and apply it to provision Event Subscription.

Upon next updates regular Plan and Apply sequence should work like a charm.

### 3. Build Function App

Build and pack source code from `src` directory. You can use embedded helper scripts to do it semi-automated.

```bash
scrips/build.sh
```

### 4. Deploy Azure Function App

1. Configure the Function App: Ensure the Azure Function App is set up to use the appropriate environment variables and has the correct permissions to access the Event Grid subscription.

2. Deploy Function Code: Deploy your Azure Function code using technique implemented in `scripts/upload.sh`. ( Upload pre-built code into azure blob storage container, point Function App to start from zip package ).

## Usage

Once the resources are provisioned and the function app is deployed:

1. Create Resources: When new resources are created within the specified resource group, Azure Event Grid will trigger the function app.

2. Tagging: The Azure Function App will receive the event data, process it, and apply tags to the new resources with the creator's name, email, and date of creation.

## Monitoring and Maintenance

- Monitor Function App: Use Azure Monitor and Application Insights to track the performance and logs of your Azure Function App.

- Update Terraform Configuration: Modify the Terraform configuration and redeploy as needed to accommodate changes in your infrastructure.

## Cleanup

Terraform will take care of complete resource removal upon issuing `terraform destroy` command.

## Contributing

Contributions to this repository are welcome. Please open an issue or submit a pull request with your changes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Contact

For any questions or support, please contact me.

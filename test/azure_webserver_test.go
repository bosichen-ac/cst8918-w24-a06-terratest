package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
// var subscriptionID string = "<your-azure-subscription-id"
var subscriptionID string = "890bd13c-ca9a-4500-bfbf-dbfb8609daa3"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			// "labelPrefix": "<your-college-id>",
			"labelPrefix": "chen0854",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variable
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))

	// Confirm NIC exists and is connected to VM
	nic := azure.GetVirtualMachineNics(t, vmName, resourceGroupName, subscriptionID)
	assert.NotEmpty(t, nic)

	// Confirm the VM is running the correct Ubuntu version
	vmImg := azure.GetVirtualMachineImage(t, vmName, resourceGroupName, subscriptionID)
	assert.Equal(t, vmImg.Offer, "0001-com-ubuntu-server-jammy")
	assert.Equal(t, vmImg.Publisher, "Canonical")
	assert.Equal(t, vmImg.SKU, "22_04-lts-gen2")
	// assert.Equal(t, vmImg.Version, "latest"??)
}

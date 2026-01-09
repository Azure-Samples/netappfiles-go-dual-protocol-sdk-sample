// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

// Sample package that is used to obtain an Azure credential
// using Azure Identity's DefaultAzureCredential which supports
// multiple authentication methods including environment variables,
// managed identity, Azure CLI, etc.

package iam

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

// GetAzureCredential gets an Azure credential and subscription ID for use with Azure SDK clients
func GetAzureCredential() (*azidentity.DefaultAzureCredential, string, error) {

	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")
	if subscriptionID == "" {
		return nil, "", fmt.Errorf("AZURE_SUBSCRIPTION_ID environment variable is not set")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to obtain Azure credential: %v", err)
	}

	return cred, subscriptionID, nil
}

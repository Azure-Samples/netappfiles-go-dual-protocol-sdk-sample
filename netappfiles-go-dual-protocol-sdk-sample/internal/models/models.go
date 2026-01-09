// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

// Structs for objects used throughout this sample.

package models

// AzureAuthInfo object definition - retained for backwards compatibility
// Note: With Track 2 SDK, authentication is handled via environment variables:
// - AZURE_SUBSCRIPTION_ID
// - AZURE_TENANT_ID
// - AZURE_CLIENT_ID
// - AZURE_CLIENT_SECRET
// Or via other authentication methods supported by DefaultAzureCredential
type AzureAuthInfo struct {
	ClientID                       *string
	ClientSecret                   *string
	SubscriptionID                 *string
	TenantID                       *string
	ActiveDirectoryEndpointURL     *string
	ResourceManagerEndpointURL     *string
	ActiveDirectoryGraphResourceID *string
	SQLManagementEndpointURL       *string
	GalleryEndpointURL             *string
	ManagementEndpointURL          *string
}

// AzureBasicInfo object definition - retained for backwards compatibility
type AzureBasicInfo struct {
	SubscriptionID             *string
	TenantID                   *string
	ResourceManagerEndpointURL *string
	ManagementEndpointURL      *string
}

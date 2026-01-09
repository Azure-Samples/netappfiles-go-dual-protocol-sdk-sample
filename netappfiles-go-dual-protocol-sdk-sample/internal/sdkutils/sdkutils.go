// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

// This package centralizes any function that directly
// using any of the Azure's (with exception of authentication related ones)
// available SDK packages.

package sdkutils

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure-Samples/netappfiles-go-dual-protocol-sdk-sample/netappfiles-go-dual-protocol-sdk-sample/internal/iam"
	"github.com/Azure-Samples/netappfiles-go-dual-protocol-sdk-sample/netappfiles-go-dual-protocol-sdk-sample/internal/uri"
	"github.com/Azure-Samples/netappfiles-go-dual-protocol-sdk-sample/netappfiles-go-dual-protocol-sdk-sample/internal/utils"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/netapp/armnetapp/v7"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

const (
	nfsv3  = "NFSv3"
	nfsv41 = "NFSv4.1"
	cifs   = "CIFS"
)

var (
	validProtocols = []string{nfsv3, nfsv41, cifs}
)

func validateANFServiceLevel(serviceLevel string) (*armnetapp.ServiceLevel, error) {

	var svcLevel armnetapp.ServiceLevel

	switch strings.ToLower(serviceLevel) {
	case "ultra":
		svcLevel = armnetapp.ServiceLevelUltra
	case "premium":
		svcLevel = armnetapp.ServiceLevelPremium
	case "standard":
		svcLevel = armnetapp.ServiceLevelStandard
	default:
		return nil, fmt.Errorf("invalid service level, supported service levels are: Ultra, Premium, Standard")
	}

	return &svcLevel, nil
}

func getResourcesClient() (*armresources.Client, error) {

	cred, subscriptionID, err := iam.GetAzureCredential()
	if err != nil {
		return nil, err
	}

	client, err := armresources.NewClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getAccountsClient() (*armnetapp.AccountsClient, error) {

	cred, subscriptionID, err := iam.GetAzureCredential()
	if err != nil {
		return nil, err
	}

	client, err := armnetapp.NewAccountsClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getPoolsClient() (*armnetapp.PoolsClient, error) {

	cred, subscriptionID, err := iam.GetAzureCredential()
	if err != nil {
		return nil, err
	}

	client, err := armnetapp.NewPoolsClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getVolumesClient() (*armnetapp.VolumesClient, error) {

	cred, subscriptionID, err := iam.GetAzureCredential()
	if err != nil {
		return nil, err
	}

	client, err := armnetapp.NewVolumesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getSnapshotsClient() (*armnetapp.SnapshotsClient, error) {

	cred, subscriptionID, err := iam.GetAzureCredential()
	if err != nil {
		return nil, err
	}

	client, err := armnetapp.NewSnapshotsClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getSnapshotPoliciesClient() (*armnetapp.SnapshotPoliciesClient, error) {

	cred, subscriptionID, err := iam.GetAzureCredential()
	if err != nil {
		return nil, err
	}

	client, err := armnetapp.NewSnapshotPoliciesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// GetResourceByID gets a generic resource
func GetResourceByID(ctx context.Context, resourceID, APIVersion string) (armresources.ClientGetByIDResponse, error) {

	resourcesClient, err := getResourcesClient()
	if err != nil {
		return armresources.ClientGetByIDResponse{}, err
	}

	return resourcesClient.GetByID(ctx, resourceID, APIVersion, nil)
}

// CreateANFAccount creates an ANF Account resource
func CreateANFAccount(ctx context.Context, location, resourceGroupName, accountName string, activeDirectories []*armnetapp.ActiveDirectory, tags map[string]*string) (armnetapp.Account, error) {

	accountClient, err := getAccountsClient()
	if err != nil {
		return armnetapp.Account{}, err
	}

	accountProperties := armnetapp.AccountProperties{}

	if activeDirectories != nil {
		accountProperties = armnetapp.AccountProperties{
			ActiveDirectories: activeDirectories,
		}
	}

	poller, err := accountClient.BeginCreateOrUpdate(
		ctx,
		resourceGroupName,
		accountName,
		armnetapp.Account{
			Location:   to.Ptr(location),
			Tags:       tags,
			Properties: &accountProperties,
		},
		nil,
	)
	if err != nil {
		return armnetapp.Account{}, fmt.Errorf("cannot create account: %v", err)
	}

	resp, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		return armnetapp.Account{}, fmt.Errorf("cannot get the account create or update future response: %v", err)
	}

	return resp.Account, nil
}

// CreateANFCapacityPool creates an ANF Capacity Pool within ANF Account
func CreateANFCapacityPool(ctx context.Context, location, resourceGroupName, accountName, poolName, serviceLevel string, sizeBytes int64, tags map[string]*string) (armnetapp.CapacityPool, error) {

	poolClient, err := getPoolsClient()
	if err != nil {
		return armnetapp.CapacityPool{}, err
	}

	svcLevel, err := validateANFServiceLevel(serviceLevel)
	if err != nil {
		return armnetapp.CapacityPool{}, err
	}

	poller, err := poolClient.BeginCreateOrUpdate(
		ctx,
		resourceGroupName,
		accountName,
		poolName,
		armnetapp.CapacityPool{
			Location: to.Ptr(location),
			Tags:     tags,
			Properties: &armnetapp.PoolProperties{
				ServiceLevel: svcLevel,
				Size:         to.Ptr(sizeBytes),
			},
		},
		nil,
	)

	if err != nil {
		return armnetapp.CapacityPool{}, fmt.Errorf("cannot create pool: %v", err)
	}

	resp, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		return armnetapp.CapacityPool{}, fmt.Errorf("cannot get the pool create or update future response: %v", err)
	}

	return resp.CapacityPool, nil
}

// CreateANFVolume creates an ANF volume within a Capacity Pool
func CreateANFVolume(ctx context.Context, location, resourceGroupName, accountName, poolName, volumeName, serviceLevel, subnetID, snapshotID string, protocolTypes []string, volumeUsageQuota int64, unixReadOnly, unixReadWrite bool, tags map[string]*string, dataProtectionObject armnetapp.VolumePropertiesDataProtection, securityStyle armnetapp.SecurityStyle) (armnetapp.Volume, error) {

	if len(protocolTypes) > 2 {
		return armnetapp.Volume{}, fmt.Errorf("maximum of two protocol types are supported")
	}

	if len(protocolTypes) > 1 && utils.Contains(protocolTypes, "NFSv4.1") {
		return armnetapp.Volume{}, fmt.Errorf("only cifs/nfsv3 protocol types are supported as dual protocol")
	}

	_, found := utils.FindInSlice(validProtocols, protocolTypes[0])
	if !found {
		return armnetapp.Volume{}, fmt.Errorf("invalid protocol type, valid protocol types are: %v", validProtocols)
	}

	svcLevel, err := validateANFServiceLevel(serviceLevel)
	if err != nil {
		return armnetapp.Volume{}, err
	}

	usingDualProtocol := false
	if len(protocolTypes) == 2 {
		usingDualProtocol = true
	}

	volumeClient, err := getVolumesClient()
	if err != nil {
		return armnetapp.Volume{}, err
	}

	var exportPolicy *armnetapp.VolumePropertiesExportPolicy

	if _, found := utils.FindInSlice(protocolTypes, cifs); !found || usingDualProtocol {
		exportPolicy = &armnetapp.VolumePropertiesExportPolicy{
			Rules: []*armnetapp.ExportPolicyRule{
				{
					AllowedClients: to.Ptr("0.0.0.0/0"),
					Cifs:           to.Ptr(utils.Contains(protocolTypes, cifs)),
					Nfsv3:          to.Ptr(utils.Contains(protocolTypes, nfsv3)),
					Nfsv41:         to.Ptr(utils.Contains(protocolTypes, nfsv41)),
					RuleIndex:      to.Ptr(int32(1)),
					UnixReadOnly:   to.Ptr(unixReadOnly),
					UnixReadWrite:  to.Ptr(unixReadWrite),
				},
			},
		}
	}

	// Convert protocolTypes to []*string
	protocolTypesPtr := make([]*string, len(protocolTypes))
	for i, pt := range protocolTypes {
		protocolTypesPtr[i] = to.Ptr(pt)
	}

	volumeProperties := armnetapp.VolumeProperties{
		ExportPolicy:   exportPolicy,
		ProtocolTypes:  protocolTypesPtr,
		ServiceLevel:   svcLevel,
		SubnetID:       to.Ptr(subnetID),
		UsageThreshold: to.Ptr(volumeUsageQuota),
		CreationToken:  to.Ptr(volumeName),
		SecurityStyle:  to.Ptr(securityStyle),
	}

	if snapshotID != "" {
		volumeProperties.SnapshotID = to.Ptr(snapshotID)
	}

	// Only set DataProtection if it has content
	if dataProtectionObject.Replication != nil || dataProtectionObject.Snapshot != nil {
		volumeProperties.DataProtection = &dataProtectionObject
	}

	poller, err := volumeClient.BeginCreateOrUpdate(
		ctx,
		resourceGroupName,
		accountName,
		poolName,
		volumeName,
		armnetapp.Volume{
			Location:   to.Ptr(location),
			Tags:       tags,
			Properties: &volumeProperties,
		},
		nil,
	)

	if err != nil {
		return armnetapp.Volume{}, fmt.Errorf("cannot create volume: %v", err)
	}

	resp, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		return armnetapp.Volume{}, fmt.Errorf("cannot get the volume create or update future response: %v", err)
	}

	return resp.Volume, nil
}

// UpdateANFVolume update an ANF volume
func UpdateANFVolume(ctx context.Context, location, resourceGroupName, accountName, poolName, volumeName string, volumePropertiesPatch armnetapp.VolumePatchProperties, tags map[string]*string) error {

	volumeClient, err := getVolumesClient()
	if err != nil {
		return err
	}

	poller, err := volumeClient.BeginUpdate(
		ctx,
		resourceGroupName,
		accountName,
		poolName,
		volumeName,
		armnetapp.VolumePatch{
			Location:   to.Ptr(location),
			Tags:       tags,
			Properties: &volumePropertiesPatch,
		},
		nil,
	)

	if err != nil {
		return fmt.Errorf("cannot update volume: %v", err)
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot get the volume update future response: %v", err)
	}

	return nil
}

// MoveANFVolumeToNewPool moves a volume to a new capacity pool within same ANF account
func MoveANFVolumeToNewPool(ctx context.Context, resourceGroupName, accountName, poolName, volumeName string, newPool armnetapp.PoolChangeRequest) error {

	volumeClient, err := getVolumesClient()
	if err != nil {
		return err
	}

	poller, err := volumeClient.BeginPoolChange(
		ctx,
		resourceGroupName,
		accountName,
		poolName,
		volumeName,
		newPool,
		nil,
	)

	if err != nil {
		return fmt.Errorf("cannot move volume: %v", err)
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot get the volume pool change future response: %v", err)
	}

	return nil
}

// AuthorizeReplication - authorizes volume replication
func AuthorizeReplication(ctx context.Context, resourceGroupName, accountName, poolName, volumeName, remoteVolumeResourceID string) error {

	volumeClient, err := getVolumesClient()
	if err != nil {
		return err
	}

	poller, err := volumeClient.BeginAuthorizeReplication(
		ctx,
		resourceGroupName,
		accountName,
		poolName,
		volumeName,
		armnetapp.AuthorizeRequest{
			RemoteVolumeResourceID: to.Ptr(remoteVolumeResourceID),
		},
		nil,
	)

	if err != nil {
		return fmt.Errorf("cannot authorize volume replication: %v", err)
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot get authorize volume replication future response: %v", err)
	}

	return nil
}

// DeleteANFVolumeReplication - deletes volume replication
func DeleteANFVolumeReplication(ctx context.Context, resourceGroupName, accountName, poolName, volumeName string) error {

	volumeClient, err := getVolumesClient()
	if err != nil {
		return err
	}

	poller, err := volumeClient.BeginDeleteReplication(
		ctx,
		resourceGroupName,
		accountName,
		poolName,
		volumeName,
		nil,
	)

	if err != nil {
		return fmt.Errorf("cannot delete volume replication: %v", err)
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot get delete volume replication future response: %v", err)
	}

	return nil
}

// CreateANFSnapshot creates a Snapshot from an ANF volume
func CreateANFSnapshot(ctx context.Context, location, resourceGroupName, accountName, poolName, volumeName, snapshotName string, tags map[string]*string) (armnetapp.Snapshot, error) {

	snapshotClient, err := getSnapshotsClient()
	if err != nil {
		return armnetapp.Snapshot{}, err
	}

	poller, err := snapshotClient.BeginCreate(
		ctx,
		resourceGroupName,
		accountName,
		poolName,
		volumeName,
		snapshotName,
		armnetapp.Snapshot{
			Location: to.Ptr(location),
		},
		nil,
	)

	if err != nil {
		return armnetapp.Snapshot{}, fmt.Errorf("cannot create snapshot: %v", err)
	}

	resp, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		return armnetapp.Snapshot{}, fmt.Errorf("cannot get the snapshot create or update future response: %v", err)
	}

	return resp.Snapshot, nil
}

// DeleteANFSnapshot deletes a Snapshot from an ANF volume
func DeleteANFSnapshot(ctx context.Context, resourceGroupName, accountName, poolName, volumeName, snapshotName string) error {

	snapshotClient, err := getSnapshotsClient()
	if err != nil {
		return err
	}

	poller, err := snapshotClient.BeginDelete(
		ctx,
		resourceGroupName,
		accountName,
		poolName,
		volumeName,
		snapshotName,
		nil,
	)

	if err != nil {
		return fmt.Errorf("cannot delete snapshot: %v", err)
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot get the snapshot delete future response: %v", err)
	}

	return nil
}

// CreateANFSnapshotPolicy creates a Snapshot Policy to be used on volumes
func CreateANFSnapshotPolicy(ctx context.Context, resourceGroupName, accountName, policyName string, policy armnetapp.SnapshotPolicy) (armnetapp.SnapshotPolicy, error) {

	snapshotPolicyClient, err := getSnapshotPoliciesClient()
	if err != nil {
		return armnetapp.SnapshotPolicy{}, err
	}

	snapshotPolicy, err := snapshotPolicyClient.Create(
		ctx,
		resourceGroupName,
		accountName,
		policyName,
		policy,
		nil,
	)

	if err != nil {
		return armnetapp.SnapshotPolicy{}, fmt.Errorf("cannot create snapshot policy: %v", err)
	}

	return snapshotPolicy.SnapshotPolicy, nil
}

// UpdateANFSnapshotPolicy update an ANF snapshot policy
func UpdateANFSnapshotPolicy(ctx context.Context, resourceGroupName, accountName, policyName string, snapshotPolicyPatch armnetapp.SnapshotPolicyPatch) error {

	snapshotPolicyClient, err := getSnapshotPoliciesClient()
	if err != nil {
		return err
	}

	poller, err := snapshotPolicyClient.BeginUpdate(
		ctx,
		resourceGroupName,
		accountName,
		policyName,
		snapshotPolicyPatch,
		nil,
	)

	if err != nil {
		return fmt.Errorf("cannot update snapshot policy: %v", err)
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot get snapshot policy update future response: %v", err)
	}

	return nil
}

// DeleteANFVolume deletes a volume
func DeleteANFVolume(ctx context.Context, resourceGroupName, accountName, poolName, volumeName string) error {

	volumesClient, err := getVolumesClient()
	if err != nil {
		return err
	}

	poller, err := volumesClient.BeginDelete(
		ctx,
		resourceGroupName,
		accountName,
		poolName,
		volumeName,
		nil,
	)

	if err != nil {
		return fmt.Errorf("cannot delete volume: %v", err)
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot get the volume delete future response: %v", err)
	}

	return nil
}

// DeleteANFCapacityPool deletes a capacity pool
func DeleteANFCapacityPool(ctx context.Context, resourceGroupName, accountName, poolName string) error {

	poolsClient, err := getPoolsClient()
	if err != nil {
		return err
	}

	poller, err := poolsClient.BeginDelete(
		ctx,
		resourceGroupName,
		accountName,
		poolName,
		nil,
	)

	if err != nil {
		return fmt.Errorf("cannot delete capacity pool: %v", err)
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot get the capacity pool delete future response: %v", err)
	}

	return nil
}

// DeleteANFSnapshotPolicy deletes a snapshot policy
func DeleteANFSnapshotPolicy(ctx context.Context, resourceGroupName, accountName, policyName string) error {

	snapshotPolicyClient, err := getSnapshotPoliciesClient()
	if err != nil {
		return err
	}

	poller, err := snapshotPolicyClient.BeginDelete(
		ctx,
		resourceGroupName,
		accountName,
		policyName,
		nil,
	)

	if err != nil {
		return fmt.Errorf("cannot delete snapshot policy: %v", err)
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot get the snapshot policy delete future response: %v", err)
	}

	return nil
}

// DeleteANFAccount deletes an account
func DeleteANFAccount(ctx context.Context, resourceGroupName, accountName string) error {

	accountsClient, err := getAccountsClient()
	if err != nil {
		return err
	}

	poller, err := accountsClient.BeginDelete(
		ctx,
		resourceGroupName,
		accountName,
		nil,
	)

	if err != nil {
		return fmt.Errorf("cannot delete account: %v", err)
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot get the account delete future response: %v", err)
	}

	return nil
}

// WaitForNoANFResource waits for a specified resource to don't exist anymore following a deletion.
// This is due to a known issue related to ARM Cache where the state of the resource is still cached within ARM infrastructure
// reporting that it still exists so looping into a get process will return 404 as soon as the cached state expires
func WaitForNoANFResource(ctx context.Context, resourceID string, intervalInSec int, retries int, checkForReplication bool) error {

	var err error

	for i := 0; i < retries; i++ {
		time.Sleep(time.Duration(intervalInSec) * time.Second)
		if uri.IsANFSnapshot(resourceID) {
			client, _ := getSnapshotsClient()
			_, err = client.Get(
				ctx,
				uri.GetResourceGroup(resourceID),
				uri.GetANFAccount(resourceID),
				uri.GetANFCapacityPool(resourceID),
				uri.GetANFVolume(resourceID),
				uri.GetANFSnapshot(resourceID),
				nil,
			)
		} else if uri.IsANFVolume(resourceID) {
			client, _ := getVolumesClient()
			if !checkForReplication {
				_, err = client.Get(
					ctx,
					uri.GetResourceGroup(resourceID),
					uri.GetANFAccount(resourceID),
					uri.GetANFCapacityPool(resourceID),
					uri.GetANFVolume(resourceID),
					nil,
				)
			} else {
				_, err = client.ReplicationStatus(
					ctx,
					uri.GetResourceGroup(resourceID),
					uri.GetANFAccount(resourceID),
					uri.GetANFCapacityPool(resourceID),
					uri.GetANFVolume(resourceID),
					nil,
				)
			}
		} else if uri.IsANFCapacityPool(resourceID) {
			client, _ := getPoolsClient()
			_, err = client.Get(
				ctx,
				uri.GetResourceGroup(resourceID),
				uri.GetANFAccount(resourceID),
				uri.GetANFCapacityPool(resourceID),
				nil,
			)
		} else if uri.IsANFSnapshotPolicy(resourceID) {
			client, _ := getSnapshotPoliciesClient()
			_, err = client.Get(
				ctx,
				uri.GetResourceGroup(resourceID),
				uri.GetANFAccount(resourceID),
				uri.GetANFSnapshotPolicy(resourceID),
				nil,
			)
		} else if uri.IsANFAccount(resourceID) {
			client, _ := getAccountsClient()
			_, err = client.Get(
				ctx,
				uri.GetResourceGroup(resourceID),
				uri.GetANFAccount(resourceID),
				nil,
			)
		}

		// In this case error is expected
		if err != nil {
			return nil
		}
	}

	return fmt.Errorf("exceeded number of retries: %v", retries)
}

// WaitForANFResource waits for a specified resource to be fully ready following a creation operation.
func WaitForANFResource(ctx context.Context, resourceID string, intervalInSec int, retries int, checkForReplication bool) error {

	var err error

	for i := 0; i < retries; i++ {
		time.Sleep(time.Duration(intervalInSec) * time.Second)
		if uri.IsANFSnapshot(resourceID) {
			client, _ := getSnapshotsClient()
			_, err = client.Get(
				ctx,
				uri.GetResourceGroup(resourceID),
				uri.GetANFAccount(resourceID),
				uri.GetANFCapacityPool(resourceID),
				uri.GetANFVolume(resourceID),
				uri.GetANFSnapshot(resourceID),
				nil,
			)
		} else if uri.IsANFVolume(resourceID) {
			client, _ := getVolumesClient()
			if !checkForReplication {
				_, err = client.Get(
					ctx,
					uri.GetResourceGroup(resourceID),
					uri.GetANFAccount(resourceID),
					uri.GetANFCapacityPool(resourceID),
					uri.GetANFVolume(resourceID),
					nil,
				)
			} else {
				_, err = client.ReplicationStatus(
					ctx,
					uri.GetResourceGroup(resourceID),
					uri.GetANFAccount(resourceID),
					uri.GetANFCapacityPool(resourceID),
					uri.GetANFVolume(resourceID),
					nil,
				)
			}
		} else if uri.IsANFCapacityPool(resourceID) {
			client, _ := getPoolsClient()
			_, err = client.Get(
				ctx,
				uri.GetResourceGroup(resourceID),
				uri.GetANFAccount(resourceID),
				uri.GetANFCapacityPool(resourceID),
				nil,
			)
		} else if uri.IsANFSnapshotPolicy(resourceID) {
			client, _ := getSnapshotPoliciesClient()
			_, err = client.Get(
				ctx,
				uri.GetResourceGroup(resourceID),
				uri.GetANFAccount(resourceID),
				uri.GetANFSnapshotPolicy(resourceID),
				nil,
			)
		} else if uri.IsANFAccount(resourceID) {
			client, _ := getAccountsClient()
			_, err = client.Get(
				ctx,
				uri.GetResourceGroup(resourceID),
				uri.GetANFAccount(resourceID),
				nil,
			)
		}

		// In this case, we exit when there is no error
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("resource still not found after number of retries: %v, error: %v", retries, err)
}

## [Azure NetApp Files Dual Protocol SDK Sample for Go] Changelog

<a name="2026-01-09"></a>
<a name="2021-08-03"></a>
<a name="2020-10-21"></a>

# 2.0.0 (2026-01-09)

*Features*
* **Breaking**: Migrated from Azure SDK for Go Track 1 to Track 2
* Updated authentication from file-based (`AZURE_AUTH_LOCATION`) to environment variable-based using `DefaultAzureCredential`
* Updated to use `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/netapp/armnetapp/v7` package
* Updated to use `github.com/Azure/azure-sdk-for-go/sdk/azidentity` for authentication
* Updated to use `github.com/Azure/azure-sdk-for-go/sdk/azcore/to` for pointer helpers
* Updated Go version requirement to 1.21+
* Updated all async operations from `Future.WaitForCompletionRef()` pattern to `Poller.PollUntilDone()` pattern

*Bug Fixes*
* N/A

*Breaking Changes*
* Authentication now requires environment variables (`AZURE_SUBSCRIPTION_ID`, `AZURE_TENANT_ID`, `AZURE_CLIENT_ID`, `AZURE_CLIENT_SECRET`) instead of auth file
* Removed dependency on `github.com/Azure/go-autorest` packages
* All SDK types changed from `netapp.*` to `armnetapp.*`

# 1.0.1 (2021-08-03)

*Features*
* API Version bump up
* sdkutils.go and utils.go updated to be in sync with other ANF Go samples

*Bug Fixes*
* N/A

*Breaking Changes*
* N/A

# 1.0.0 (2020-10-21)

*Features*
* Initial version

*Bug Fixes*
* N/A

*Breaking Changes*
* N/A

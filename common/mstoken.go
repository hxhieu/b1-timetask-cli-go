package common

import (
	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func getCreds() (*azidentity.DeviceCodeCredential, error) {

	// TODO: cache this if possible but highly unlikely because of MFA enabled on our tenant
	cred, err := azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
		TenantID: "c6feb44d-2731-4551-a9e2-dfd519d35042",
		ClientID: "24c73aa6-dc13-46ba-83b2-caf107ad76ef", // is a public client so no need secret
	})
	if err != nil {
		return nil, err
	}

	return cred, nil
}

package common

import (
	"time"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type OutLookCalendarEventTime struct {
	DateTime time.Time `json:"dateTime"`
}

type OutLookCalendarEvent struct {
	Categories []string                 `json:"categories"`
	Subject    string                   `json:"subject"`
	Start      OutLookCalendarEventTime `json:"start"`
	End        OutLookCalendarEventTime `json:"end"`
}

func GetMsGraphCreds(useDeviceCode bool) (any, error) {

	// TODO: cache this if possible but highly unlikely because of MFA enabled on our tenant
	if useDeviceCode {
		cred, err := azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
			TenantID: "c6feb44d-2731-4551-a9e2-dfd519d35042",
			ClientID: "24c73aa6-dc13-46ba-83b2-caf107ad76ef", // is a public client so no need secret
		})
		if err != nil {
			return nil, err
		}

		return cred, nil
	}

	cred, err := azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
		TenantID: "c6feb44d-2731-4551-a9e2-dfd519d35042",
		ClientID: "24c73aa6-dc13-46ba-83b2-caf107ad76ef", // is a public client so no need secret
	})
	if err != nil {
		return nil, err
	}

	return cred, nil
}

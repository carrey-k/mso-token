# mso-token
Fetch microsoft online access token by using microsoft-authentication-library-for-go

## Installation

### Setting Up Go
To install Go, visit [this link](https://golang.org/dl/).

### Installing Module
`go get -u github.com/kangchengkun/mso-token`

## Usage
Before using this Go module, you will need to [register your application with the Microsoft identity platform](https://docs.microsoft.com/azure/active-directory/develop/quickstart-v2-register-an-app).

Get the AAD tenant and application information from your microsoft online Administrator

```
import github.com/kangchengkun/mso-token

msotoken.TenantID = 'your-tenant-id'
msotoken.ClientID = 'your-client-id'
msotoken.ClientSecret = 'your-client-secret'

// Change the default cache time
msotoken.DefaultCacheTime = 10 * time.Minute

// Change the default permission scopes
msotoken.PermissionScopes = []string{"https://outlook.office365.com/.default"}


// Fetch token
accessToken, err := msotoken.GetToken()
if err != nil {
    fmt.Println("Fetch access token failed")
}
```
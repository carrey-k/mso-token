package msotoken

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/patrickmn/go-cache"
)

// https://aad.portal.azure.com/ Azure Active Directory admin center
// https://docs.microsoft.com/en-us/exchange/client-developer/exchange-web-services/how-to-authenticate-an-ews-application-by-using-oauth

var (
	// the tenant id from microsoft Azure Active Directory, ask from administrators
	TenantID string
	// the application client id to fetch token from microsoft online, ask from administrators
	ClientID string
	// the application client secret to fetch token from microsoft online, ask from administrators
	ClientSecret string
	// the permission scope required for microsoft EWS online
	PermissionScopes []string = []string{"https://outlook.office365.com/.default"}
	// the authority base url for microsoft EWS online
	AuthorityBaseUrl string = "https://login.microsoftonline.com/"
	// cache instance with a default expiration time of 5 minutes, and which purges expired items every 10 minutes
	CacheIns *cache.Cache = cache.New(5*time.Minute, 10*time.Minute)
	// the default cache time
	DefaultCacheTime time.Duration = cache.DefaultExpiration
)

func GetToken() (string, error) {
	if accessToken, found := CacheIns.Get("msalToken"); found && accessToken.(string) != "" {
		fmt.Printf("get msalToken from cache: %v\n", accessToken)
		return accessToken.(string), nil
	}

	if TenantID == "" {
		fmt.Println("No valid tenant ID provided")
		return "", errors.New("no valid tenant id provided")
	}

	if ClientID == "" {
		fmt.Println("No valid client ID provided")
		return "", errors.New("no valid client id provided")
	}

	if ClientSecret == "" {
		fmt.Println("No valid client secret provided")
		return "", errors.New("no valid client secret provided")
	}

	cred, err := confidential.NewCredFromSecret(ClientSecret)
	if err != nil {
		fmt.Printf("Create credential from a secret failed with error: %v\n", err)
		return "", err
	}

	authApp, err := confidential.New(ClientID, cred, confidential.WithAuthority(AuthorityBaseUrl+TenantID))
	if err != nil {
		fmt.Printf("Create authentication application failed with error: %v\n", err)
		return "", err
	}

	cacheToken := func(token string) {
		if token != "" {
			go CacheIns.Set("msalToken", token, DefaultCacheTime)
		}
	}

	result, err := authApp.AcquireTokenSilent(context.Background(), PermissionScopes)
	if err != nil {
		fmt.Printf("Fetch silent token failed with error: %v\n", err)
		fmt.Println("Trying to acquire token by credential......")
		result, err = authApp.AcquireTokenByCredential(context.Background(), PermissionScopes)
		if err != nil {
			fmt.Printf("Acquire token by credential failed with error: %v\n", err)
			return "", err
		}

		cacheToken(result.AccessToken) // cache the aquired token by credential
		fmt.Println("Acquired access token by credential: " + result.AccessToken)
		return result.AccessToken, nil
	}

	cacheToken(result.AccessToken) // cache the silently aquired token
	fmt.Println("Silently acquired access token: " + result.AccessToken)
	return result.AccessToken, nil
}

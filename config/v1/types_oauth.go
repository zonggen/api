package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OAuth Server and Identity Provider Config

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OAuth holds cluster-wide information about OAuth.  The canonical name is `cluster`
type OAuth struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              OAuthSpec   `json:"spec"`
	Status            OAuthStatus `json:"status,omitempty"`
}

// OAuthSpec contains desired cluster auth configuration
type OAuthSpec struct {
	//IdentityProviders is an ordered list of ways for a user to identify themselves
	IdentityProviders []OAuthIdentityProvider `json:"identityProviders"`

	// TokenConfig contains options for authorization and access tokens
	TokenConfig TokenConfig `json:"tokenConfig"`

	// Templates allow you to customize pages like the login page.
	// +optional
	Templates OAuthTemplates `json:"templates"`
}

// OAuthStatus shows current known state of OAuth server in the cluster
type OAuthStatus struct {
	// TODO Fill in
}

// TokenConfig holds the necessary configuration options for authorization and access tokens
type TokenConfig struct {
	// AuthorizeTokenMaxAgeSeconds defines the maximum age of authorize tokens
	AuthorizeTokenMaxAgeSeconds int32 `json:"authorizeTokenMaxAgeSeconds"`
	// AccessTokenMaxAgeSeconds defines the maximum age of access tokens
	AccessTokenMaxAgeSeconds int32 `json:"accessTokenMaxAgeSeconds"`
	// AccessTokenInactivityTimeoutSeconds defines the default token
	// inactivity timeout for tokens granted by any client.
	// The value represents the maximum amount of time that can occur between
	// consecutive uses of the token. Tokens become invalid if they are not
	// used within this temporal window. The user will need to acquire a new
	// token to regain access once a token times out.
	// Valid values are integer values:
	//   x < 0  Tokens never timeout (e.g. `-1`)
	//   x = 0  Tokens are disabled (default)
	//   x > 0  Tokens time out if there is no activity for x seconds
	// The current minimum allowed value for X is 300 (5 minutes)
	// +optional
	AccessTokenInactivityTimeoutSeconds int32 `json:"accessTokenInactivityTimeoutSeconds,omitempty"`
}

const (
	// LoginTemplateKey is the default key of the login template
	LoginTemplateKey = "login.html"
	// ProviderSelectionTemplateKey is the default key for the provider selection template
	ProviderSelectionTemplateKey = "providers.html"
	// ErrorsTemplateKey is the default key for the errors template
	ErrorsTemplateKey = "errors.html"
)

// OAuthTemplates allow for customization of pages like the login page
type OAuthTemplates struct {
	// Login is a reference to a secret that specifies a go template to use to render the login page.
	// If a key is not specified, the key `login.html` is used to locate the template data.
	// If unspecified, the default login page is used.
	// +optional
	Login LocalSecretReference `json:"login,omitemtpy"`

	// ProviderSelection is a reference to a secret that specifies a go template to use to render
	// the provider selection page.
	// If a key is not specified, the key `providers.html` is used to locate the template data.
	// If unspecified, the default provider selection page is used.
	// +optional
	ProviderSelection LocalSecretReference `json:"providerSelection,omitempty"`

	// Error is a reference to a secret that specifies a go template to use to render error pages
	// during the authentication or grant flow.
	// If a key is not specified, the key `errrors.html` is used to locate the template data.
	// If unspecified, the default error page is used.
	// +optional
	Error LocalSecretReference `json:"error,omitempty"`
}

// OAuthIdentityProvider provides identities for users authenticating using credentials
type OAuthIdentityProvider struct {
	// Name is used to qualify the identities returned by this provider.
	// - It MUST be unique and not shared by any other identity provider used
	// - It MUST be a vlid path segment: name cannot equal "." or ".." or contain "/" or "%"
	//   Ref: https://godoc.org/k8s.io/apimachinery/pkg/api/validation/path#ValidatePathSegmentName
	Name string `json:"name"`

	// UseAsChallenger indicates whether to issue WWW-Authenticate challenges for this provider
	UseAsChallenger bool `json:"challenge"`
	// UseAsLogin indicates whether to use this identity provider for unauthenticated browsers to login against
	UseAsLogin bool `json:"login"`

	// MappingMethod determines how identities from this provider are mapped to users
	// Defaults to "claim"
	// +optional
	MappingMethod MappingMethodType `json:"mappingMethod"`

	// GrantMethod: allow, deny, prompt
	// This method will be used only if the specific OAuth client doesn't provide a strategy
	// of their own. Valid grant handling methods are:
	//  - auto:   always approves grant requests, useful for trusted clients
	//  - prompt: prompts the end user for approval of grant requests, useful for third-party clients
	//  - deny:   always denies grant requests, useful for black-listed clients
	// Defaults to "prompt" if not set.
	// +optional
	GrantMethod GrantHandlerType `json:"grantMethod"`

	// IdentityProvidersConfig
	ProviderConfig IdentityProviderConfig `json:",inline"`
}

// MappingMethodType specifies how new identities should be mapped to users when they log in
type MappingMethodType string

const (
	// MappingMethodClaim provisions a user with the identity’s preferred user name. Fails if a user
	// with that user name is already mapped to another identity.
	// Default.
	MappingMethodClaim MappingMethodType = "claim"

	// MappingMethodLookup looks up existing users already mapped to an identity but does not
	// automatically provision users or identities. Requires identities and users be set up
	// manually or using an external process.
	MappingMethodLookup MappingMethodType = "lookup"

	// MappingMethodAdd provisions a user with the identity’s preferred user name. If a user with
	// that user name already exists, the identity is mapped to the existing user, adding to any
	// existing identity mappings for the user.
	MappingMethodAdd MappingMethodType = "add"

	// MappingMethodGenerate provisions a user with the identity’s preferred user name. If a user
	// with the preferred user name is already mapped to an existing identity, a unique user name is
	// generated, e.g. myuser2. This method should not be used in combination with external
	// processes that require exact matches between openshift user names and the idp user name
	// such as LDAP group sync.
	MappingMethodGenerate MappingMethodType = "generate"
)

// GrantHandlerType are the valid strategies for handling grant requests
type GrantHandlerType string

const (
	// GrantHandlerAuto auto-approves client authorization grant requests
	GrantHandlerAuto GrantHandlerType = "auto"
	// GrantHandlerPrompt prompts the user to approve new client authorization grant requests
	GrantHandlerPrompt GrantHandlerType = "prompt"
	// GrantHandlerDeny auto-denies client authorization grant requests
	GrantHandlerDeny GrantHandlerType = "deny"
)

type IdentityProviderType string

const (
	// IdentityProviderTypeBasicAuth provides identities for users authenticating with HTTP Basic Auth
	IdentityProviderTypeBasicAuth IdentityProviderType = "BasicAuth"

	// IdentityProviderTypeAllowAll provides identities for all users authenticating using non-empty passwords
	IdentityProviderTypeAllowAll IdentityProviderType = "AllowAll"

	// IdentityProviderTypeDenyAll provides no identities for users
	IdentityProviderTypeDenyAll IdentityProviderType = "DenyAll"

	// IdentityProviderTypeHTPasswd provides identities from an HTPasswd file
	IdentityProviderTypeHTPasswd IdentityProviderType = "HTPasswd"

	// IdentityProviderTypeLDAP provides identities for users authenticating using LDAP credentials
	IdentityProviderTypeLDAP IdentityProviderType = "LDAP"

	// IdentityProviderTypeKeystone provides identitities for users authenticating using keystone password credentials
	IdentityProviderTypeKeystone IdentityProviderType = "Keystone"

	// IdentityProviderTypeRequestHeader provides identities for users authenticating using request header credentials
	IdentityProviderTypeRequestHeader IdentityProviderType = "RequestHeader"

	// IdentityProviderTypeGitHub provides identities for users authenticating using GitHub credentials
	IdentityProviderTypeGitHub IdentityProviderType = "GitHub"

	// IdentityProviderTypeGitLab provides identities for users authenticating using GitLab credentials
	IdentityProviderTypeGitLab IdentityProviderType = "GitLab"

	// IdentityProviderTypeGoogle provides identities for users authenticating using Google credentials
	IdentityProviderTypeGoogle IdentityProviderType = "Google"

	// IdentityProviderTypeOpenID provides identities for users authenticating using OpenID credentials
	IdentityProviderTypeOpenID IdentityProviderType = "OpenID"
)

// IdentityProviderConfig contains configuration for using a specific identity provider
type IdentityProviderConfig struct {
	// IdentityProviderType identifies the identity provider type for this entry.
	Type IdentityProviderType `json:"type"`

	// Provider-specific configuration
	// The json tag MUST match the `Type` specified above, case-insensitively
	// e.g. For `Type: "LDAP"`, the `LDAPPasswordIdentityProvider` configuration should be provided

	// BasicAuth contains configuration options for the BasicAuth IdP
	// +optional
	BasicAuth *BasicAuthPasswordIdentityProvider `json:"basicAuth,omitempty"`

	// AllowAll enables the AllowAllIdentityProvider which provides identities for users
	// authenticating using non-empty passwords.
	// Defaults to `false`, i.e. allowAll set to off
	// +optional
	AllowAll bool `json:"allowAll,omitempty"`

	// DenyAll enables the DenyAllPasswordIdentityProvider which provides no identities for users
	// Defaults to `false`, ie. denyAll set to off
	// +optional
	DenyAll bool `json:"denyAll,omitempty"`

	// HTPasswd enables user authentication using an HTPasswd file to validate credentials
	// +optional
	HTPasswd *HTPasswdPasswordIdentityProvider `json:"htpasswd,omitempty"`

	// LDAP enables user authentication using LDAP credentials
	// +optional
	LDAP *LDAPPasswordIdentityProvider `json:"ldap,omitempty"`

	// Keystone enables user authentication using keystone password credentials
	// +optional
	Keystone *KeystonePasswordIdentityProvider `json:"keystone,omitempty"`

	// RequestHeader enables user authentication using request header credentials
	RequestHeader *RequestHeaderIdentityProvider `json:"requestHeader,omitempty"`

	// GitHub enables  user authentication using GitHub credentials
	// +optional
	GitHub *GitHubIdentityProvider `json:"github,omitempty"`

	// GitLab enables user authentication using GitLab credentials
	// +optional
	GitLab *GitLabIdentityProvider `json:"gitlab,omitempty"`

	// Google enables user authentication using Google credentials
	// +optional
	Google *GoogleIdentityProvider `json:"google,omitempty"`

	// OpenID enables user authentication using OpenID credentials
	// +optional
	OpenID *OpenIDIdentityProvider `json:"openID,omitempty"`
}

// BasicAuthPasswordIdentityProvider provides identities for users authenticating using HTTP basic auth credentials
type BasicAuthPasswordIdentityProvider struct {
	// RemoteConnectionInfo contains information about how to connect to the external basic auth server
	OAuthRemoteConnectionInfo `json:",inline"`
}

// RemoteConnectionInfo holds information necessary for establishing a remote connection
type OAuthRemoteConnectionInfo struct {
	// URL is the remote URL to connect to
	URL string `json:"url"`
	// CA is a reference to a ConfigMap containing the CA for verifying TLS connections
	CA ConfigMapReference `json:"ca"`

	// TLSClientCert references a secret containing the TLS client certificate to present when
	// connecting to the server.
	// Looks under the key "tls.cert" for the data unless a lookup key is specified in the secret ref
	TLSClientCert LocalSecretReference `json:"tlsClientCert"`

	// TLSClientKey references a secret containing the TLS private key for the client certificate
	// Looks under the key "tls.key" for the data unless a lookup key is specified in the secret ref
	TLSClientKey LocalSecretReference `json:"tlsClientKey"`
}

// HTPasswdDataKey is the default key for the htpasswd file data in a secret
const HTPasswdDataKey = "htpasswd"

// HTPasswdPasswordIdentityProvider provides identities for users authenticating using htpasswd credentials
type HTPasswdPasswordIdentityProvider struct {
	// Data is a reference to a secret containing the data to use as the htpasswd file
	// Looks under the key `htpasswd` unless a lookup key is specified in the secret ref
	FileData LocalSecretReference `json:"fileData"`
}

const (
	// BindPasswordKey is default the key for the LDAP bind password in a secret
	BindPasswordKey = "bindPassword"
	// ClientSecretKey is the key for the oauth client secret data in a secret
	ClientSecretKey = "clientSecret"
)

// LDAPPasswordIdentityProvider provides identities for users authenticating using LDAP credentials
type LDAPPasswordIdentityProvider struct {
	// URL is an RFC 2255 URL which specifies the LDAP search parameters to use.
	// The syntax of the URL is:
	//    ldap://host:port/basedn?attribute?scope?filter
	URL string `json:"url"`

	// BindDN is an optional DN to bind with during the search phase.
	// +optional
	BindDN string `json:"bindDN"`

	// BindPasswordSecretRef is a reference to the secret containing an optional password to bind
	// with during the search phase.
	// Looks under the key `bindPassword` unless a lookup key is specified in the secret ref
	// +optional
	BindPasswordSecretRef LocalSecretReference `json:"bindPasswordSecretRef"`

	// Insecure, if true, indicates the connection should not use TLS
	// WARNING: Should not be set to `true` with the URL scheme "ldaps://" as "ldaps://" URLs always
	//          attempt to connect using TLS, even when `insecure` is set to `true`
	// When `true`, "ldap://" URLS connect insecurely. When `false`, "ldap://" URLs are upgraded to
	// a TLS connection using StartTLS as specified in https://tools.ietf.org/html/rfc2830.
	Insecure bool `json:"insecure"`

	// CA is a reference to a ConfigMap containing an optional trusted certificate authority bundle
	// to use when making requests to the server.
	// If empty, the default system roots are used.
	// +optional
	CA ConfigMapReference `json:"ca"`

	// Attributes maps LDAP attributes to identities
	Attributes LDAPAttributeMapping `json:"attributes"`
}

// LDAPAttributeMapping maps LDAP attributes to OpenShift identity fields
type LDAPAttributeMapping struct {
	// ID is the list of attributes whose values should be used as the user ID. Required.
	// First non-empty attribute is used. At least one attribute is required. If none of the listed
	// attribute have a value, authentication fails.
	// LDAP standard identity attribute is "dn"
	ID []string `json:"id"`
	// PreferredUsername is the list of attributes whose values should be used as the preferred username.
	// LDAP standard login attribute is "uid"
	// +optional
	PreferredUsername []string `json:"preferredUsername"`
	// Name is the list of attributes whose values should be used as the display name. Optional.
	// If unspecified, no display name is set for the identity
	// LDAP standard display name attribute is "cn"
	// +optional
	Name []string `json:"name"`
	// Email is the list of attributes whose values should be used as the email address. Optional.
	// If unspecified, no email is set for the identity
	// +optional
	Email []string `json:"email"`
}

// KeystonePasswordIdentityProvider provides identities for users authenticating using keystone password credentials
type KeystonePasswordIdentityProvider struct {
	// RemoteConnectionInfo contains information about how to connect to the keystone server
	OAuthRemoteConnectionInfo `json:",inline"`
	// Domain Name is required for keystone v3
	DomainName string `json:"domainName"`
	// LegacyLookupByUsername flag indicates that user should be authenticated by username, not keystone ID
	// DEPRECATED - only use this option for legacy systems to ensure backwards compatibiity
	// +optional
	LegacyLookupByUsername bool `json:"useKeystoneIdentity"`
}

// RequestHeaderIdentityProvider provides identities for users authenticating using request header credentials
type RequestHeaderIdentityProvider struct {
	// LoginURL is a URL to redirect unauthenticated /authorize requests to
	// Unauthenticated requests from OAuth clients which expect interactive logins will be redirected here
	// ${url} is replaced with the current URL, escaped to be safe in a query parameter
	//   https://www.example.com/sso-login?then=${url}
	// ${query} is replaced with the current query string
	//   https://www.example.com/auth-proxy/oauth/authorize?${query}
	// Required when UseAsLogin is set to true.
	LoginURL string `json:"loginURL"`

	// ChallengeURL is a URL to redirect unauthenticated /authorize requests to
	// Unauthenticated requests from OAuth clients which expect WWW-Authenticate challenges will be
	// redirected here.
	// ${url} is replaced with the current URL, escaped to be safe in a query parameter
	//   https://www.example.com/sso-login?then=${url}
	// ${query} is replaced with the current query string
	//   https://www.example.com/auth-proxy/oauth/authorize?${query}
	// Required when UseAsChallenger is set to true.
	ChallengeURL string `json:"challengeURL"`

	// ClientCA is a reference to a configmap with the trusted signer certs. If empty, no request
	// verification is done, and any direct request to the OAuth server can impersonate any identity
	// from this provider, merely by setting a request header.
	// +optional
	ClientCA ConfigMapReference `json:"ca"`

	// ClientCommonNames is an optional list of common names to require a match from. If empty, any
	// client certificate validated against the clientCA bundle is considered authoritative.
	// +optional
	ClientCommonNames []string `json:"clientCommonNames"`

	// Headers is the set of headers to check for identity information
	Headers []string `json:"headers"`

	// PreferredUsernameHeaders is the set of headers to check for the preferred username
	PreferredUsernameHeaders []string `json:"preferredUsernameHeaders"`

	// NameHeaders is the set of headers to check for the display name
	NameHeaders []string `json:"nameHeaders"`

	// EmailHeaders is the set of headers to check for the email address
	EmailHeaders []string `json:"emailHeaders"`
}

// GitHubIdentityProvider provides identities for users authenticating using GitHub credentials
type GitHubIdentityProvider struct {
	// ClientID is the oauth client ID
	ClientID string `json:"clientID"`

	// ClientSecret is is a reference to the secret containing the oauth client secret
	// The secret referenced must contain a key named `clientSecret` containing the secret data.
	ClientSecret LocalSecretReference `json:"clientSecret"`

	// Organizations optionally restricts which organizations are allowed to log in
	// +optional
	Organizations []string `json:"organizations"`

	// Teams optionally restricts which teams are allowed to log in. Format is <org>/<team>.
	// +optional
	Teams []string `json:"teams"`

	// Hostname is the optional domain (e.g. "mycompany.com") for use with a hosted instance of
	// GitHub Enterprise.
	// It must match the GitHub Enterprise settings value configured at /setup/settings#hostname.
	// +optional
	Hostname string `json:"hostname"`

	// CA is a reference to a ConfigMap containing an optional trusted certificate authority bundle
	// to use when making requests to the server.
	// If empty, the default system roots are used.
	// This can only be configured when hostname is set to a non-empty value.
	// +optional
	CA ConfigMapReference `json:"ca"`
}

// GitLabIdentityProvider provides identities for users authenticating using GitLab credentials
type GitLabIdentityProvider struct {
	// CA is a reference to a ConfigMap containing an optional trusted certificate authority bundle
	// to use when making requests to the server.
	// If empty, the default system roots are used.
	// +optional
	CA ConfigMapReference `json:"ca"`

	// URL is the oauth server base URL
	URL string `json:"url"`

	// ClientID is the oauth client ID
	ClientID string `json:"clientID"`

	// ClientSecret is is a reference to the secret containing the oauth client secret
	// The secret referenced must contain a key named `clientSecret` containing the secret data.
	ClientSecret LocalSecretReference `json:"clientSecret"`

	// LegacyOAuth2 determines that OAuth2 should be used, not OIDC
	// +optional
	LegacyOAuth2 bool `json:"legacy,omitempty"`
}

// GoogleIdentityProvider provides identities for users authenticating using Google credentials
type GoogleIdentityProvider struct {
	// ClientID is the oauth client ID
	ClientID string `json:"clientID"`

	// ClientSecret is is a reference to the secret containing the oauth client secret
	// The secret referenced must contain a key named `clientSecret` containing the secret data.
	ClientSecret LocalSecretReference `json:"clientSecret"`

	// HostedDomain is the optional Google App domain (e.g. "mycompany.com") to restrict logins to
	// +optional
	HostedDomain string `json:"hostedDomain"`
}

// OpenIDIdentityProvider provides identities for users authenticating using OpenID credentials
type OpenIDIdentityProvider struct {
	// CA is a reference to a ConfigMap containing an optional trusted certificate authority bundle
	// to use when making requests to the server.
	// If empty, the default system roots are used.
	// +optional
	CA ConfigMapReference `json:"ca"`

	// ClientID is the oauth client ID
	ClientID string `json:"clientID"`

	// ClientSecret is is a reference to the secret containing the oauth client secret
	// The secret referenced must contain a key named `clientSecret` containing the secret data.
	ClientSecret LocalSecretReference `json:"clientSecret"`

	// ExtraScopes are any scopes to request in addition to the standard "openid" scope.
	// +optional
	ExtraScopes []string `json:"extraScopes"`

	// ExtraAuthorizeParameters are any custom parameters to add to the authorize request.
	// +optional
	ExtraAuthorizeParameters map[string]string `json:"extraAuthorizeParameters"`

	// URLs to use to authenticate
	URLs OpenIDURLs `json:"urls"`

	// Claims mappings
	Claims OpenIDClaims `json:"claims"`
}

// OpenIDURLs are URLs to use when authenticating with an OpenID identity provider
type OpenIDURLs struct {
	// Authorize is the oauth authorization URL
	Authorize string `json:"authorize"`
	// Token is the oauth token granting URL
	Token string `json:"token"`
	// UserInfo is the optional userinfo URL.
	// If present, a granted access_token is used to request claims
	// If empty, a granted id_token is parsed for claims
	// +optional
	UserInfo string `json:"userInfo"`
}

// UserIDClaim is used in the `ID` field for an `OpenIDClaim`
// Per http://openid.net/specs/openid-connect-core-1_0.html#ClaimStability
//  "The sub (subject) and iss (issuer) Claims, used together, are the only Claims that an RP can
//   rely upon as a stable identifier for the End-User, since the sub Claim MUST be locally unique
//   and never reassigned within the Issuer for a particular End-User, as described in Section 2.
//   Therefore, the only guaranteed unique identifier for a given End-User is the combination of the
//   iss Claim and the sub Claim."
const UserIDClaim = "sub"

// OpenIDClaims contains a list of OpenID claims to use when authenticating with an OpenID identity provider
type OpenIDClaims struct {
	// PreferredUsername is the list of claims whose values should be used as the preferred username.
	// If unspecified, the preferred username is determined from the value of the id claim
	// +optional
	PreferredUsername []string `json:"preferredUsername"`
	// Name is the list of claims whose values should be used as the display name. Optional.
	// If unspecified, no display name is set for the identity
	// +optional
	Name []string `json:"name"`
	// Email is the list of claims whose values should be used as the email address. Optional.
	// If unspecified, no email is set for the identity
	// +optional
	Email []string `json:"email"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type OAuthList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OAuth `json:"items"`
}

/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package alicloud

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	credential "github.com/aliyun/credentials-go/credentials"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

// CredentialSource selects how the CPI obtains Alibaba Cloud credentials.
type CredentialSource string

const (
	// CredentialSourceStatic signs requests with the access key (and optional
	// STS token) supplied in the CPI config. This is the default so existing
	// deployments keep working unchanged.
	CredentialSourceStatic CredentialSource = "static"

	// CredentialSourceECSRAMRole signs requests with short-lived STS
	// credentials issued to the RAM role attached to the running ECS instance.
	CredentialSourceECSRAMRole CredentialSource = "ecs_ram_role"
)

// ecsMetadataBaseURL is the link-local instance metadata service. It is a
// variable so tests can point it at a fake metadata server; production code must
// never override it.
var ecsMetadataBaseURL = "http://100.100.100.200"

// ecsMetadataRoleURL lists the RAM roles attached to the running ECS instance.
var ecsMetadataRoleURL = ecsMetadataBaseURL + "/latest/meta-data/ram/security-credentials/"

// ecsMetadataTokenTTL is how long an IMDSv2 token stays valid. Role discovery
// uses the token once and discards it, so the window is deliberately short.
const ecsMetadataTokenTTL = "60"

// ecsMetadataTimeout bounds metadata lookups. The service is link-local, so a
// slow response means it is unreachable rather than busy.
var ecsMetadataTimeout = 5 * time.Second

// ossCredentialGraceWindow is how long the OSS adapter may keep serving the
// last credential it fetched after a refresh failure. It is deliberately much
// shorter than the ~6h STS lifetime so a genuinely expired credential is never
// reused, while a blip in the metadata service does not fail in-flight signing.
const ossCredentialGraceWindow = 5 * time.Minute

// CredentialProvider hands out per-SDK credential objects for the configured
// credential source. The CPI talks to Alibaba Cloud through three SDKs with
// three unrelated credential types, and routing all of them through one
// provider is what guarantees every client authenticates the same way.
//
// Implementations must be safe for concurrent use.
type CredentialProvider interface {
	// LegacyCredential authenticates the ECS, SLB and Location clients built
	// on alibaba-cloud-sdk-go.
	LegacyCredential() (auth.Credential, error)

	// TeaCredential authenticates the ECS and NLB Tea/OpenAPI clients.
	TeaCredential() (credential.Credential, error)

	// OSSCredentialsProvider authenticates the OSS client.
	OSSCredentialsProvider() (oss.CredentialsProvider, error)
}

// NewCredentialProvider builds the provider for the given credential source.
// Callers are expected to have validated the source already; an unknown source
// is still rejected here so a provider can never silently fall back.
func NewCredentialProvider(source CredentialSource, a OpenApi) (CredentialProvider, error) {
	switch source {
	case CredentialSourceStatic:
		return staticCredentialProvider{
			accessKeyId:     a.AccessKeyId,
			accessKeySecret: a.AccessKeySecret,
			securityToken:   a.SecurityToken,
		}, nil
	case CredentialSourceECSRAMRole:
		return &ecsRAMRoleCredentialProvider{roleName: strings.TrimSpace(a.RamRoleName)}, nil
	default:
		return nil, fmt.Errorf("unknown credential_source %q, must be one of %q, %q",
			source, CredentialSourceStatic, CredentialSourceECSRAMRole)
	}
}

// staticCredentialProvider serves the access key from the CPI config.
type staticCredentialProvider struct {
	accessKeyId     string
	accessKeySecret string
	securityToken   string
}

func (p staticCredentialProvider) LegacyCredential() (auth.Credential, error) {
	// An STS token credential with an empty token is what the CPI has always
	// sent for plain access keys; keep it so signing behaviour is unchanged for
	// existing deployments.
	return credentials.NewStsTokenCredential(p.accessKeyId, p.accessKeySecret, p.securityToken), nil
}

func (p staticCredentialProvider) TeaCredential() (credential.Credential, error) {
	config := &credential.Config{
		AccessKeyId:     tea.String(p.accessKeyId),
		AccessKeySecret: tea.String(p.accessKeySecret),
		Type:            tea.String("access_key"),
	}
	if p.securityToken != "" {
		config.SecurityToken = tea.String(p.securityToken)
		config.Type = tea.String("sts")
	}

	c, err := credential.NewCredential(config)
	if err != nil {
		return nil, bosherr.WrapError(scrubCredentialError(err), "Building static credential failed")
	}
	return scrubbingCredential{inner: c}, nil
}

func (p staticCredentialProvider) OSSCredentialsProvider() (oss.CredentialsProvider, error) {
	return staticOSSCredentialsProvider{
		creds: ossCredentials{
			accessKeyId:     p.accessKeyId,
			accessKeySecret: p.accessKeySecret,
			securityToken:   p.securityToken,
		},
	}, nil
}

// ecsRAMRoleCredentialProvider serves short-lived STS credentials issued to the
// RAM role attached to the running ECS instance.
//
// The underlying SDK objects are created once and reused so their refresh state
// and credential caches are shared by every client in the process.
type ecsRAMRoleCredentialProvider struct {
	// roleName is empty when the role is discovered from instance metadata,
	// which lets the same config work on the Concourse worker and on the
	// Director without naming either role.
	roleName string

	mutex           sync.Mutex
	resolvedRole    string
	teaCredential   credential.Credential
	ossCredsAdapter *ossRoleCredentialsProvider
}

func (p *ecsRAMRoleCredentialProvider) LegacyCredential() (auth.Credential, error) {
	// alibaba-cloud-sdk-go's ECS RAM role signer appends the role name to the
	// metadata URL without discovering it, so an empty name has to be resolved
	// here. The signer then owns fetching and refreshing the STS credential.
	roleName, err := p.resolveRoleName()
	if err != nil {
		return nil, err
	}
	return credentials.NewEcsRamRoleCredential(roleName), nil
}

func (p *ecsRAMRoleCredentialProvider) TeaCredential() (credential.Credential, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.teaCredentialLocked()
}

func (p *ecsRAMRoleCredentialProvider) teaCredentialLocked() (credential.Credential, error) {
	if p.teaCredential != nil {
		return p.teaCredential, nil
	}

	config := &credential.Config{Type: tea.String(string(CredentialSourceECSRAMRole))}
	if p.roleName != "" {
		config.RoleName = tea.String(p.roleName)
	}

	c, err := credential.NewCredential(config)
	if err != nil {
		return nil, bosherr.WrapError(scrubCredentialError(err), "Building ecs_ram_role credential failed")
	}

	p.teaCredential = scrubbingCredential{inner: c}
	return p.teaCredential, nil
}

func (p *ecsRAMRoleCredentialProvider) OSSCredentialsProvider() (oss.CredentialsProvider, error) {
	p.mutex.Lock()
	source, err := p.teaCredentialLocked()
	if err != nil {
		p.mutex.Unlock()
		return nil, err
	}
	if p.ossCredsAdapter == nil {
		p.ossCredsAdapter = &ossRoleCredentialsProvider{source: source}
	}
	adapter := p.ossCredsAdapter
	p.mutex.Unlock()

	// Fetch once up front so a missing role or unreachable metadata service
	// fails while building the client instead of surfacing as an opaque 403 on
	// the first stemcell upload.
	if _, err := adapter.fetch(); err != nil {
		return nil, err
	}
	return adapter, nil
}

// resolveRoleName returns the configured role name, or the single role attached
// to this instance when none is configured.
func (p *ecsRAMRoleCredentialProvider) resolveRoleName() (string, error) {
	if p.roleName != "" {
		return p.roleName, nil
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.resolvedRole != "" {
		return p.resolvedRole, nil
	}

	roleName, err := discoverECSRAMRoleName()
	if err != nil {
		return "", err
	}

	p.resolvedRole = roleName
	return roleName, nil
}

// ecsMetadataToken fetches an IMDSv2 token. An instance configured to require
// IMDSv2 rejects unauthenticated reads, so the token is obtained first. An
// instance that still allows IMDSv1 has no token endpoint, in which case this
// returns an empty token and the caller falls back to an unauthenticated read.
func ecsMetadataToken(client *http.Client) string {
	req, err := http.NewRequest(http.MethodPut, ecsMetadataBaseURL+"/latest/api/token", nil)
	if err != nil {
		return ""
	}
	req.Header.Set("X-aliyun-ecs-metadata-token-ttl-seconds", ecsMetadataTokenTTL)

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	token, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(token))
}

// discoverECSRAMRoleName asks the instance metadata service which RAM role is
// attached to this instance.
func discoverECSRAMRoleName() (string, error) {
	client := &http.Client{Timeout: ecsMetadataTimeout}

	req, err := http.NewRequest(http.MethodGet, ecsMetadataRoleURL, nil)
	if err != nil {
		return "", bosherr.WrapError(scrubCredentialError(err),
			"Building the ECS instance metadata request failed")
	}
	if token := ecsMetadataToken(client); token != "" {
		req.Header.Set("X-aliyun-ecs-metadata-token", token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", bosherr.WrapError(scrubCredentialError(err),
			"Discovering the ECS RAM role from instance metadata failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", bosherr.Errorf(
			"Discovering the ECS RAM role from instance metadata failed with HTTP status %d; "+
				"attach a RAM role to this instance or set alicloud.ram_role_name", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		return "", bosherr.WrapError(scrubCredentialError(err),
			"Reading the ECS RAM role from instance metadata failed")
	}

	// The listing endpoint returns role names one per line. More than one means
	// the config has to say which to use rather than us picking arbitrarily.
	var roleNames []string
	for _, line := range strings.Split(string(body), "\n") {
		if name := strings.TrimSpace(line); name != "" {
			roleNames = append(roleNames, name)
		}
	}

	switch len(roleNames) {
	case 0:
		return "", bosherr.Error("No RAM role is attached to this ECS instance; " +
			"attach one or use credential_source 'static'")
	case 1:
		return roleNames[0], nil
	default:
		return "", bosherr.Errorf("%d RAM roles are attached to this ECS instance; "+
			"set alicloud.ram_role_name to select one", len(roleNames))
	}
}

// ossCredentials is the OSS SDK's view of a credential.
type ossCredentials struct {
	accessKeyId     string
	accessKeySecret string
	securityToken   string
}

func (c ossCredentials) GetAccessKeyID() string     { return c.accessKeyId }
func (c ossCredentials) GetAccessKeySecret() string { return c.accessKeySecret }
func (c ossCredentials) GetSecurityToken() string   { return c.securityToken }

type staticOSSCredentialsProvider struct {
	creds ossCredentials
}

func (p staticOSSCredentialsProvider) GetCredentials() oss.Credentials { return p.creds }

// ossRoleCredentialsProvider adapts a refreshing credential to the OSS SDK.
//
// The SDK asks for credentials on every signature through an interface that
// cannot report errors, so this adapter caches the last credential it fetched:
// a transient metadata failure must not break signing mid-upload. The cache is
// bounded by ossCredentialGraceWindow and then dropped, so an expired
// credential is never served and there is no fallback to a static key.
type ossRoleCredentialsProvider struct {
	source credential.Credential

	mutex      sync.Mutex
	cached     ossCredentials
	cacheValid bool
	cacheUntil time.Time
}

func (p *ossRoleCredentialsProvider) GetCredentials() oss.Credentials {
	creds, err := p.fetch()
	if err != nil {
		// Returning empty credentials makes OSS reject the request, which is
		// the correct outcome: we must not sign with a stale credential.
		return ossCredentials{}
	}
	return creds
}

func (p *ossRoleCredentialsProvider) fetch() (ossCredentials, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	model, err := p.source.GetCredential()
	if err == nil && model != nil {
		p.cached = ossCredentials{
			accessKeyId:     tea.StringValue(model.AccessKeyId),
			accessKeySecret: tea.StringValue(model.AccessKeySecret),
			securityToken:   tea.StringValue(model.SecurityToken),
		}
		p.cacheValid = true
		p.cacheUntil = time.Now().Add(ossCredentialGraceWindow)
		return p.cached, nil
	}

	if p.cacheValid && time.Now().Before(p.cacheUntil) {
		return p.cached, nil
	}

	p.cached = ossCredentials{}
	p.cacheValid = false

	if err == nil {
		err = bosherr.Error("credential source returned no credential")
	}
	return ossCredentials{}, bosherr.WrapError(scrubCredentialError(err),
		"Obtaining OSS credentials from the ecs_ram_role credential source failed")
}

// scrubbingCredential wraps a credential so that errors raised while it
// refreshes itself are scrubbed too. Without it, only construction and the OSS
// adapter are covered: the Tea SDK refreshes the credential on its own during a
// call, and credentials-go quotes the raw metadata response in those errors.
type scrubbingCredential struct {
	inner credential.Credential
}

func (c scrubbingCredential) GetAccessKeyId() (*string, error) {
	v, err := c.inner.GetAccessKeyId()
	return v, scrubCredentialError(err)
}

func (c scrubbingCredential) GetAccessKeySecret() (*string, error) {
	v, err := c.inner.GetAccessKeySecret()
	return v, scrubCredentialError(err)
}

func (c scrubbingCredential) GetSecurityToken() (*string, error) {
	v, err := c.inner.GetSecurityToken()
	return v, scrubCredentialError(err)
}

func (c scrubbingCredential) GetBearerToken() *string { return c.inner.GetBearerToken() }
func (c scrubbingCredential) GetType() *string        { return c.inner.GetType() }

func (c scrubbingCredential) GetCredential() (*credential.CredentialModel, error) {
	m, err := c.inner.GetCredential()
	return m, scrubCredentialError(err)
}

// String keeps a wrapped credential from printing its inner state.
func (c scrubbingCredential) String() string {
	return fmt.Sprintf("credential(type=%s)", tea.StringValue(c.inner.GetType()))
}

// String redacts the credential fields so that formatting an OpenApi, a Config
// or anything embedding them with %s or %v cannot put a usable credential into a
// log. Both the access key secret and the STS token are secrets.
func (a OpenApi) String() string {
	redacted := a
	if redacted.AccessKeyId != "" {
		redacted.AccessKeyId = "<redacted>"
	}
	if redacted.AccessKeySecret != "" {
		redacted.AccessKeySecret = "<redacted>"
	}
	if redacted.SecurityToken != "" {
		redacted.SecurityToken = "<redacted>"
	}
	// Format the copy through an alias so this method is not re-entered.
	type openApiFields OpenApi
	return fmt.Sprintf("%+v", openApiFields(redacted))
}

// credentialSecretPattern matches the secret-bearing fields an SDK may echo
// back in an error, plus bare credential values, so they can be redacted.
var credentialSecretPattern = regexp.MustCompile(
	`(?i)("?(?:access[_-]?key[_-]?secret|access[_-]?key[_-]?id|security[_-]?token|sts[_-]?token)"?\s*[:=]\s*"?)([^",}\s]+)` +
		`|\b(?:LTAI|STS\.)[A-Za-z0-9]{6,}`)

// scrubCredentialError redacts credential material from an error before it is
// wrapped and logged. SDK errors can quote the raw metadata response, which
// contains a usable access key, secret and token.
func scrubCredentialError(err error) error {
	if err == nil {
		return nil
	}

	scrubbed := credentialSecretPattern.ReplaceAllStringFunc(err.Error(), func(match string) string {
		groups := credentialSecretPattern.FindStringSubmatch(match)
		if groups[1] != "" {
			return groups[1] + "<redacted>"
		}
		return "<redacted>"
	})

	return bosherr.Error(scrubbed)
}

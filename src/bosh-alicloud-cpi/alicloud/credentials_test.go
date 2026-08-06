/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package alicloud

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	credential "github.com/aliyun/credentials-go/credentials"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// fakeCredential stands in for a refreshing credentials-go credential so the
// OSS adapter can be tested without touching the real metadata service.
type fakeCredential struct {
	mutex sync.Mutex
	calls int32
	next  *credential.CredentialModel
	err   error
}

func (f *fakeCredential) set(model *credential.CredentialModel, err error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.next = model
	f.err = err
}

func (f *fakeCredential) GetCredential() (*credential.CredentialModel, error) {
	atomic.AddInt32(&f.calls, 1)
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.next, f.err
}

// The accessors report the configured error too, so the scrubbing wrapper can be
// exercised on every path the Tea SDK uses to refresh a credential.
func (f *fakeCredential) GetAccessKeyId() (*string, error)     { return tea.String(""), f.current() }
func (f *fakeCredential) GetAccessKeySecret() (*string, error) { return tea.String(""), f.current() }
func (f *fakeCredential) GetSecurityToken() (*string, error)   { return tea.String(""), f.current() }

func (f *fakeCredential) current() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.err
}
func (f *fakeCredential) GetBearerToken() *string              { return tea.String("") }
func (f *fakeCredential) GetType() *string                     { return tea.String("fake") }

func credentialModel(id, secret, token string) *credential.CredentialModel {
	return &credential.CredentialModel{
		AccessKeyId:     tea.String(id),
		AccessKeySecret: tea.String(secret),
		SecurityToken:   tea.String(token),
	}
}

var _ = Describe("CredentialSource resolution", func() {
	It("defaults to static when unset", func() {
		Expect(OpenApi{}.GetCredentialSource()).To(Equal(CredentialSourceStatic))
	})

	It("honours an explicit source", func() {
		Expect(OpenApi{CredentialSource: "ecs_ram_role"}.GetCredentialSource()).
			To(Equal(CredentialSourceECSRAMRole))
		Expect(OpenApi{CredentialSource: " static "}.GetCredentialSource()).
			To(Equal(CredentialSourceStatic))
	})
})

var _ = Describe("Credential configuration validation", func() {
	newConfig := func(a OpenApi) Config {
		a.Region = "cn-beijing"
		return Config{OpenApi: a, Registry: RegistryConfig{Port: "6901"}}
	}

	It("accepts a full static access key", func() {
		err := newConfig(OpenApi{AccessKeyId: "id", AccessKeySecret: "secret"}).Validate()
		Expect(err).NotTo(HaveOccurred())
	})

	It("accepts a static STS triple", func() {
		err := newConfig(OpenApi{
			CredentialSource: "static",
			AccessKeyId:      "id",
			AccessKeySecret:  "secret",
			SecurityToken:    "token",
		}).Validate()
		Expect(err).NotTo(HaveOccurred())
	})

	It("rejects a half-configured static access key", func() {
		err := newConfig(OpenApi{AccessKeyId: "id"}).Validate()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("needs both access_key_id and access_key_secret"))
	})

	It("leaves an entirely empty static access key to the caller", func() {
		// The rendered config looks like this before the integration harness
		// fills the key in, so validation must not reject it.
		err := newConfig(OpenApi{}).Validate()
		Expect(err).NotTo(HaveOccurred())
	})

	It("accepts ecs_ram_role without a role name", func() {
		err := newConfig(OpenApi{CredentialSource: "ecs_ram_role"}).Validate()
		Expect(err).NotTo(HaveOccurred())
	})

	It("accepts ecs_ram_role with an explicit role name", func() {
		err := newConfig(OpenApi{CredentialSource: "ecs_ram_role", RamRoleName: "BoshDirectorRole"}).Validate()
		Expect(err).NotTo(HaveOccurred())
	})

	It("rejects ecs_ram_role combined with static credentials", func() {
		for _, a := range []OpenApi{
			{CredentialSource: "ecs_ram_role", AccessKeyId: "id"},
			{CredentialSource: "ecs_ram_role", AccessKeySecret: "secret"},
			{CredentialSource: "ecs_ram_role", SecurityToken: "token"},
		} {
			err := newConfig(a).Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("can't be combined with"))
		}
	})

	It("rejects an unknown credential source", func() {
		err := newConfig(OpenApi{CredentialSource: "instance_profile"}).Validate()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unknown credential_source"))
	})
})

var _ = Describe("Static credential provider", func() {
	var provider CredentialProvider

	BeforeEach(func() {
		var err error
		provider, err = NewCredentialProvider(CredentialSourceStatic, OpenApi{
			AccessKeyId:     "static-id",
			AccessKeySecret: "static-secret",
			SecurityToken:   "static-token",
		})
		Expect(err).NotTo(HaveOccurred())
	})

	It("builds an STS token credential for the legacy SDK", func() {
		legacy, err := provider.LegacyCredential()
		Expect(err).NotTo(HaveOccurred())

		sts, ok := legacy.(*credentials.StsTokenCredential)
		Expect(ok).To(BeTrue())
		Expect(sts.AccessKeyId).To(Equal("static-id"))
		Expect(sts.AccessKeyStsToken).To(Equal("static-token"))
	})

	It("builds an sts credential for the Tea SDK", func() {
		tea1, err := provider.TeaCredential()
		Expect(err).NotTo(HaveOccurred())
		Expect(tea.StringValue(tea1.GetType())).To(Equal("sts"))
	})

	It("builds an access_key credential for the Tea SDK without a token", func() {
		p, err := NewCredentialProvider(CredentialSourceStatic, OpenApi{
			AccessKeyId:     "static-id",
			AccessKeySecret: "static-secret",
		})
		Expect(err).NotTo(HaveOccurred())

		teaCred, err := p.TeaCredential()
		Expect(err).NotTo(HaveOccurred())
		Expect(tea.StringValue(teaCred.GetType())).To(Equal("access_key"))
	})

	It("serves the configured key to the OSS SDK", func() {
		ossProvider, err := provider.OSSCredentialsProvider()
		Expect(err).NotTo(HaveOccurred())

		creds := ossProvider.GetCredentials()
		Expect(creds.GetAccessKeyID()).To(Equal("static-id"))
		Expect(creds.GetAccessKeySecret()).To(Equal("static-secret"))
		Expect(creds.GetSecurityToken()).To(Equal("static-token"))
	})
})

var _ = Describe("ECS RAM role credential provider", func() {
	var originalURL string

	BeforeEach(func() {
		originalURL = ecsMetadataRoleURL
	})

	AfterEach(func() {
		ecsMetadataRoleURL = originalURL
	})

	// stubMetadata points role discovery at a local server returning body.
	stubMetadata := func(status int, body string) *httptest.Server {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
			fmt.Fprint(w, body)
		}))
		ecsMetadataRoleURL = server.URL + "/"
		return server
	}

	It("uses the configured role name without querying metadata", func() {
		// Any metadata call would fail against this address.
		ecsMetadataRoleURL = "http://127.0.0.1:1/"

		provider, err := NewCredentialProvider(CredentialSourceECSRAMRole,
			OpenApi{RamRoleName: "BoshDirectorRole"})
		Expect(err).NotTo(HaveOccurred())

		legacy, err := provider.LegacyCredential()
		Expect(err).NotTo(HaveOccurred())

		role, ok := legacy.(*credentials.EcsRamRoleCredential)
		Expect(ok).To(BeTrue())
		Expect(role.RoleName).To(Equal("BoshDirectorRole"))
	})

	It("discovers the attached role name from metadata", func() {
		server := stubMetadata(http.StatusOK, "BoshDirectorRole\n")
		defer server.Close()

		provider, err := NewCredentialProvider(CredentialSourceECSRAMRole, OpenApi{})
		Expect(err).NotTo(HaveOccurred())

		legacy, err := provider.LegacyCredential()
		Expect(err).NotTo(HaveOccurred())
		Expect(legacy.(*credentials.EcsRamRoleCredential).RoleName).To(Equal("BoshDirectorRole"))
	})

	It("caches the discovered role name", func() {
		var hits int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&hits, 1)
			fmt.Fprint(w, "BoshDirectorRole")
		}))
		defer server.Close()
		ecsMetadataRoleURL = server.URL + "/"

		provider, err := NewCredentialProvider(CredentialSourceECSRAMRole, OpenApi{})
		Expect(err).NotTo(HaveOccurred())

		for i := 0; i < 3; i++ {
			_, err := provider.LegacyCredential()
			Expect(err).NotTo(HaveOccurred())
		}
		Expect(atomic.LoadInt32(&hits)).To(Equal(int32(1)))
	})

	It("fails when no role is attached", func() {
		server := stubMetadata(http.StatusOK, "\n")
		defer server.Close()

		provider, err := NewCredentialProvider(CredentialSourceECSRAMRole, OpenApi{})
		Expect(err).NotTo(HaveOccurred())

		_, err = provider.LegacyCredential()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("No RAM role is attached"))
	})

	It("fails when several roles are attached", func() {
		server := stubMetadata(http.StatusOK, "RoleA\nRoleB\n")
		defer server.Close()

		provider, err := NewCredentialProvider(CredentialSourceECSRAMRole, OpenApi{})
		Expect(err).NotTo(HaveOccurred())

		_, err = provider.LegacyCredential()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("set alicloud.ram_role_name"))
	})

	It("fails when metadata is unreachable", func() {
		ecsMetadataRoleURL = "http://127.0.0.1:1/"

		provider, err := NewCredentialProvider(CredentialSourceECSRAMRole, OpenApi{})
		Expect(err).NotTo(HaveOccurred())

		_, err = provider.LegacyCredential()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("instance metadata"))
	})

	It("reports a metadata error status", func() {
		server := stubMetadata(http.StatusNotFound, "")
		defer server.Close()

		provider, err := NewCredentialProvider(CredentialSourceECSRAMRole, OpenApi{})
		Expect(err).NotTo(HaveOccurred())

		_, err = provider.LegacyCredential()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("HTTP status 404"))
	})

	It("builds an ecs_ram_role credential for the Tea SDK", func() {
		provider, err := NewCredentialProvider(CredentialSourceECSRAMRole,
			OpenApi{RamRoleName: "BoshDirectorRole"})
		Expect(err).NotTo(HaveOccurred())

		teaCred, err := provider.TeaCredential()
		Expect(err).NotTo(HaveOccurred())
		Expect(tea.StringValue(teaCred.GetType())).To(Equal("ecs_ram_role"))
	})

	It("reuses one Tea credential so refresh state is shared", func() {
		provider, err := NewCredentialProvider(CredentialSourceECSRAMRole,
			OpenApi{RamRoleName: "BoshDirectorRole"})
		Expect(err).NotTo(HaveOccurred())

		first, err := provider.TeaCredential()
		Expect(err).NotTo(HaveOccurred())
		second, err := provider.TeaCredential()
		Expect(err).NotTo(HaveOccurred())
		Expect(first).To(BeIdenticalTo(second))
	})
})

var _ = Describe("OSS role credentials adapter", func() {
	var (
		source  *fakeCredential
		adapter *ossRoleCredentialsProvider
	)

	BeforeEach(func() {
		source = &fakeCredential{}
		source.set(credentialModel("id-A", "secret-A", "token-A"), nil)
		adapter = &ossRoleCredentialsProvider{source: source}
	})

	It("serves the current credential", func() {
		creds := adapter.GetCredentials()
		Expect(creds.GetAccessKeyID()).To(Equal("id-A"))
		Expect(creds.GetAccessKeySecret()).To(Equal("secret-A"))
		Expect(creds.GetSecurityToken()).To(Equal("token-A"))
	})

	It("picks up a rotated credential", func() {
		Expect(adapter.GetCredentials().GetAccessKeyID()).To(Equal("id-A"))

		source.set(credentialModel("id-B", "secret-B", "token-B"), nil)

		creds := adapter.GetCredentials()
		Expect(creds.GetAccessKeyID()).To(Equal("id-B"))
		Expect(creds.GetSecurityToken()).To(Equal("token-B"))
	})

	It("serves the cached credential through a transient refresh failure", func() {
		Expect(adapter.GetCredentials().GetAccessKeyID()).To(Equal("id-A"))

		source.set(nil, errors.New("metadata temporarily unavailable"))

		Expect(adapter.GetCredentials().GetAccessKeyID()).To(Equal("id-A"))
	})

	It("stops serving the cached credential once the grace window closes", func() {
		Expect(adapter.GetCredentials().GetAccessKeyID()).To(Equal("id-A"))

		source.set(nil, errors.New("metadata unavailable"))
		adapter.cacheUntil = time.Now().Add(-time.Second)

		creds := adapter.GetCredentials()
		Expect(creds.GetAccessKeyID()).To(BeEmpty())
		Expect(creds.GetAccessKeySecret()).To(BeEmpty())
		Expect(creds.GetSecurityToken()).To(BeEmpty())
	})

	It("never falls back to a credential it has not fetched", func() {
		source.set(nil, errors.New("no role attached"))

		creds := adapter.GetCredentials()
		Expect(creds.GetAccessKeyID()).To(BeEmpty())
	})

	It("reports the failure to callers that can handle it", func() {
		source.set(nil, errors.New("no role attached"))

		_, err := adapter.fetch()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("ecs_ram_role"))
	})

	It("is safe for concurrent signing", func() {
		var wg sync.WaitGroup
		for i := 0; i < 32; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				if i%4 == 0 {
					source.set(credentialModel(
						fmt.Sprintf("id-%d", i),
						fmt.Sprintf("secret-%d", i),
						fmt.Sprintf("token-%d", i)), nil)
				}
				adapter.GetCredentials()
			}(i)
		}
		wg.Wait()

		Expect(adapter.GetCredentials().GetAccessKeyID()).NotTo(BeEmpty())
	})
})

var _ = Describe("Instance metadata access", func() {
	var originalBase, originalRole string

	BeforeEach(func() {
		originalBase, originalRole = ecsMetadataBaseURL, ecsMetadataRoleURL
	})

	AfterEach(func() {
		ecsMetadataBaseURL, ecsMetadataRoleURL = originalBase, originalRole
	})

	// pointAt makes both the token and the role endpoint resolve to server.
	pointAt := func(server *httptest.Server) {
		ecsMetadataBaseURL = server.URL
		ecsMetadataRoleURL = server.URL + "/latest/meta-data/ram/security-credentials/"
	}

	It("obtains an IMDSv2 token and presents it when reading the role", func() {
		var tokenRequests, roleRequests int32
		var presented string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPut && r.URL.Path == "/latest/api/token" {
				atomic.AddInt32(&tokenRequests, 1)
				Expect(r.Header.Get("X-aliyun-ecs-metadata-token-ttl-seconds")).NotTo(BeEmpty())
				fmt.Fprint(w, "a-metadata-token")
				return
			}
			atomic.AddInt32(&roleRequests, 1)
			presented = r.Header.Get("X-aliyun-ecs-metadata-token")
			fmt.Fprint(w, "BoshDirectorRole")
		}))
		defer server.Close()
		pointAt(server)

		provider, err := NewCredentialProvider(CredentialSourceECSRAMRole, OpenApi{})
		Expect(err).NotTo(HaveOccurred())

		legacy, err := provider.LegacyCredential()
		Expect(err).NotTo(HaveOccurred())
		Expect(legacy.(*credentials.EcsRamRoleCredential).RoleName).To(Equal("BoshDirectorRole"))

		Expect(atomic.LoadInt32(&tokenRequests)).To(Equal(int32(1)))
		Expect(atomic.LoadInt32(&roleRequests)).To(Equal(int32(1)))
		Expect(presented).To(Equal("a-metadata-token"))
	})

	It("still reads the role on an instance that has no token endpoint", func() {
		var presented string
		hadToken := true

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPut && r.URL.Path == "/latest/api/token" {
				// IMDSv1-only instances do not serve this endpoint.
				w.WriteHeader(http.StatusNotFound)
				return
			}
			presented = r.Header.Get("X-aliyun-ecs-metadata-token")
			hadToken = presented != ""
			fmt.Fprint(w, "BoshDirectorRole")
		}))
		defer server.Close()
		pointAt(server)

		provider, err := NewCredentialProvider(CredentialSourceECSRAMRole, OpenApi{})
		Expect(err).NotTo(HaveOccurred())

		legacy, err := provider.LegacyCredential()
		Expect(err).NotTo(HaveOccurred())
		Expect(legacy.(*credentials.EcsRamRoleCredential).RoleName).To(Equal("BoshDirectorRole"))
		Expect(hadToken).To(BeFalse())
	})
})

var _ = Describe("Credential redaction", func() {
	It("redacts the credential fields when an OpenApi is formatted", func() {
		a := OpenApi{
			Region:          "eu-central-1",
			AccessKeyId:     "an-access-key-id",
			AccessKeySecret: "an-access-key-secret",
			SecurityToken:   "an-sts-token",
		}

		for _, rendered := range []string{fmt.Sprintf("%s", a), fmt.Sprintf("%v", a)} {
			Expect(rendered).NotTo(ContainSubstring("an-access-key-id"))
			Expect(rendered).NotTo(ContainSubstring("an-access-key-secret"))
			Expect(rendered).NotTo(ContainSubstring("an-sts-token"))
			// Non-secret fields stay visible so the log is still useful.
			Expect(rendered).To(ContainSubstring("eu-central-1"))
		}
	})

	It("redacts the credential fields when the whole Config is formatted", func() {
		c := Config{OpenApi: OpenApi{Region: "moon", AccessKeySecret: "an-access-key-secret"}}
		Expect(fmt.Sprintf("%v", c)).NotTo(ContainSubstring("an-access-key-secret"))
	})

	It("leaves empty credential fields as they are", func() {
		Expect(fmt.Sprintf("%s", OpenApi{Region: "moon"})).NotTo(ContainSubstring("<redacted>"))
	})
})

var _ = Describe("Runtime refresh error scrubbing", func() {
	// The Tea SDK refreshes the credential during a call, so those errors do not
	// pass through the provider and have to be scrubbed by the wrapper.
	newFailing := func() scrubbingCredential {
		f := &fakeCredential{}
		f.set(nil, errors.New(`refresh failed: {"AccessKeySecret":"9dSaAbcDefGhiJkl","SecurityToken":"CAISleaked"}`))
		return scrubbingCredential{inner: f}
	}

	It("scrubs GetCredential", func() {
		_, err := newFailing().GetCredential()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).NotTo(ContainSubstring("9dSaAbcDefGhiJkl"))
		Expect(err.Error()).NotTo(ContainSubstring("CAISleaked"))
		Expect(err.Error()).To(ContainSubstring("<redacted>"))
	})

	It("scrubs the individual accessors", func() {
		c := newFailing()
		for _, get := range []func() (*string, error){
			c.GetAccessKeyId, c.GetAccessKeySecret, c.GetSecurityToken,
		} {
			_, err := get()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).NotTo(ContainSubstring("9dSaAbcDefGhiJkl"))
		}
	})

	It("does not print the wrapped credential's state", func() {
		Expect(fmt.Sprintf("%s", newFailing())).To(Equal("credential(type=fake)"))
	})
})

var _ = Describe("Credential error scrubbing", func() {
	It("redacts a metadata response echoed back by an SDK", func() {
		raw := errors.New(`refresh Ecs sts token err, json.Unmarshal fail: ` +
			`{"Code":"Success","AccessKeyId":"STS.NUxxxxxxxxxxxxxxxx",` +
			`"AccessKeySecret":"9dSaAbcDefGhiJklMnoPqr","SecurityToken":"CAISxxxxxxxxxxxx"}`)

		scrubbed := scrubCredentialError(raw).Error()

		Expect(scrubbed).NotTo(ContainSubstring("9dSaAbcDefGhiJklMnoPqr"))
		Expect(scrubbed).NotTo(ContainSubstring("STS.NUxxxxxxxxxxxxxxxx"))
		Expect(scrubbed).NotTo(ContainSubstring("CAISxxxxxxxxxxxx"))
		Expect(scrubbed).To(ContainSubstring("<redacted>"))
		// The diagnostic part has to survive so failures stay debuggable.
		Expect(scrubbed).To(ContainSubstring("json.Unmarshal fail"))
	})

	It("redacts a bare access key id", func() {
		scrubbed := scrubCredentialError(errors.New("denied for LTAI5tAbCdEfGhIjKlMnOpQr")).Error()

		Expect(scrubbed).NotTo(ContainSubstring("LTAI5tAbCdEfGhIjKlMnOpQr"))
		Expect(scrubbed).To(ContainSubstring("<redacted>"))
	})

	It("leaves a credential-free error alone", func() {
		Expect(scrubCredentialError(errors.New("connection refused")).Error()).
			To(Equal("connection refused"))
	})

	It("passes nil through", func() {
		Expect(scrubCredentialError(nil)).To(BeNil())
	})
})

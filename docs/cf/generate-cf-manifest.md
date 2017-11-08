## Customizing the Cloud Foundry Deployment Manifest for Alicloud

This topic describes how to customize the Cloud Foundry deployment manifest for Alicloud.

To create a Cloud Foundry manifest, you must perform the following steps:

1. Prepare workspace environment, download cf-release and install [spiff](https://github.com/cloudfoundry-incubator/spiff).

2. Use the BOSH CLI to retrieve your BOSH Director UUID, which you use to customize your manifest stub.

3. Create a manifest stub in YAML format.

4. Use spiff to merge the manifest stub with Alicloud IaaS configuration and other configuration files in the `cf-release` repository to generate your deployment manifest.

5. Deploy Cloud Foundry by the manifest.


### Step 1: Preparation

1. Prepare a work environment

    ```
    $ mkdir workspace
    $ cd workspace
    $ export WORK_DIR=$(pwd)
    ```
2. Clone the `cf-release` GitHub repository. Use `git clone` to clone the latest Cloud Foundry configuration files onto your machine.

    ```
    $ git clone https://github.com/cloudfoundry/cf-release.git
    ```

3. From the `cf-release` directory, run the update script to fetch all the submodules.

    ```
    $ {WORK_DIR}/cf-release/scripts/update
    ```

4. Install [spiff](https://github.com/cloudfoundry-incubator/spiff), a command line tool for generating deployment manifests.


### Step 2: Retrieve Your BOSH Director UUID

1. To perform these procedures, you must have installed the [BOSH CLI2](https://bosh.io/docs/cli-v2.html#install), and [deploy bosh on Alicloud](https://github.com/aliyun/bosh-alicloud-cpi-release#bosh-alibaba-cloud-cpi) successfully.

2. Use the `bosh env` command to target and view information about your BOSH deployment. Record the UUID of the BOSH Director. You use the UUID when customizing the Cloud Foundry deployment manifest stub.

    ```
    $ export BOSH_CLIENT=YOUR_BOSH_CLIENT
    $ export BOSH_CLIENT_SECRET=YOUR_BOSH_CLIENT_SECRET
    $ bosh -e YOUR_BOSH_NAME env
    Using environment '10.0.16.2' as client 'admin'

    Name      my-bosh
    UUID      879e2edf-a9c4-4d2e-be8c-6b23dc530008
    Version   263.2.0 (00000000)
    CPI       alicloud_cpi
    Features  compiled_package_cache: disabled
              config_server: disabled
              dns: disabled
              snapshots: disabled
    User      admin

    Succeeded
    ```


### Step 3: Create Your Manifest Stub

To create your manifest stub, download the example manifest stub for Alicloud from GitHub, then customize it for your deployment.

#### Download the Example Manifest Stub

[Download the manifest stub for Alicloud from GitHub](https://github.com/aliyun/bosh-alicloud-cpi-release/blob/release9-networks/docs/cf/stubs/cf-stub.yml), and you can refer the [example config](https://github.com/aliyun/bosh-alicloud-cpi-release/blob/release9-networks/docs/cf/stubs/cf-stub-example.yml)

#### Customize the Manifest

Follow the instructions in this table to edit the manifest stub for your deployment.


<table border="1" class="nice">
  <tr>
    <th style="width:35%">Deployment Manifest Stub Contents</th>
    <th>Editing Instructions</th>
  </tr>
  <tr>
    <td><pre><code>
director_uuid: DIRECTOR-UUID
    </code></pre></td>
    <td>Replace <code>DIRECTOR-UUID</code> with the BOSH Director UUID. Run the       BOSH CLI command <code>bosh status --uuid</code> to view the BOSH
      Director UUID.
    </td>
  </tr>
  <tr>
    <td><pre><code>
meta:
  environment: ENVIRONMENT
    </code></pre></td>
    <td>Replace <code>ENVIRONMENT</code> with a name of your choice that describes your environment, for example <code>alicloud-prod</code>.
    </td>
  </tr>
  <tr>
    <td><pre><code>
  reserved_static_ips:
  - CLOUD-FOUNDRY-PUBLIC-IP
    </code></pre></td>
    <td>Replace <code>CLOUD-FOUNDRY-PUBLIC-IP</code> with an existing public IP address for your Alicloud network. This is assigned to the <code>ha_proxy</code> job to receive incoming traffic.
    </td>
  </tr>
  <tr>
    <td><pre><code>
networks:
  - name: reserved
    type: vip

  - name: cf1
    type: manual
    subnets:
    - range: 10.0.16.0/20
      gateway: 10.0.16.1
      dns:
      - 168.63.129.16
      reserved:
      - 10.0.16.2 - 10.0.16.3
      static:
      - 10.0.16.4 - 10.0.16.100
      cloud_properties:
        vswitch_id: SUBNET-ID-FOR-CLOUD-FOUNDRY
        security_group_id: NSG-ID-FOR-CLOUD-FOUNDRY
    </code></pre></td>
    <td>Update the values for <code>range</code>, <code>gateway</code>,  <code>reserved</code>, and <code>static</code> to reflect the available networks and IP addresses in your Alicloud network. Replace <code>SUBNET-ID-FOR-CLOUD-FOUNDRY</code> with the Virtual Switch ID,  and <code>NSG-ID-FOR-CLOUD-FOUNDRY</code> with the security group ID of your Alicloud network.
    </td>
  </tr>
  <tr>
    <td><pre><code>
properties:
  system_domain: SYSTEM-DOMAIN
  system_domain_organization: SYSTEM-DOMAIN-ORGANIZATION
    </code></pre></td>
    <td>Replace <code>SYSTEM-DOMAIN</code> with the full domain you want associated with applications pushed to your Cloud Foundry installation. You must have already acquired these domains and configured their DNS records so that these domains resolve to your load balancer. For testing, you can use <code>CLOUD-FOUNDRY-PUBLIC-IP.xip.io</code> if you do not have a domain available. <code>xip.io</code> is not stable, and can sometimes fail to resolve to related IP addresses.
    <br/><br/>
    Replace <code>SYSTEM-DOMAIN-ORGANIZATION</code> with your organization name. This organization will be created and configured to own the <code>SYSTEM-DOMAIN</code>.
    </td>
  </tr>
  <tr>
    <td><pre><code>
  app_ssh:
    host_key_fingerprint: HOST-KEY-FINGERPRINT
    oauth_client_id: ssh-proxy
    </code></pre></td>
    <td>Replace <code>HOST-KEY-FINGERPRINT</code> with the host key fingerprint. Run the following commands to generate the fingerprint:

<pre class=terminal>$ ssh-keygen -N "" -f ssh-proxy-host-key.pem > /dev/null
$ ssh-keygen -lf ssh-proxy-host-key.pem.pub | cut -d ' ' -f2</pre>
    </td>
  </tr>
  <tr>
    <td><pre><code>
  ssl:
    skip_cert_verify: true
    </code></pre></td>
    <td> Set <code>skip_cert_verify</code> to <code>true</code> to skip SSL certificate verification if you do not use a real domain and certificate. If you want to use SSL certificates to secure traffic into your deployment, see the <a href="../../adminguide/securing-traffic.html">Securing Traffic into Cloud Foundry</a> topic.
  </td>
  <tr>
    <td><pre><code>
  cc:
    staging_upload_user: STAGING-UPLOAD-USER
    staging_upload_password: STAGING-UPLOAD-PASSWORD
    bulk_api_password: BULK-API-PASSWORD
    db_encryption_key: DB-ENCRYPTION-KEY
    uaa_skip_ssl_validation: true
    tls_port: CC-MUTUAL-TLS-PORT
    mutual_tls:
      ca_cert: CC-MUTUAL-TLS-CA-CERT
      public_cert: CC-MUTUAL-TLS-PUBLIC-CERT
      private_key: CC-MUTUAL-TLS-PRIVATE-KEY
    </code></pre></td>
    <td>
      The Cloud Controller API endpoint requires basic authentication. Replace <code>STAGING-UPLOAD-USER</code> and <code>STAGING-UPLOAD-PASSWORD</code> with a username and password of your choosing.
      <br /><br />
      Replace <code>BULK-API-PASSWORD</code> with a password of your choosing. The password cannot contain characters that are not allowed for basic authentication unless you encode the characters. Health Manager uses this password to access the Cloud Controller bulk API.
      <br /><br />
      Replace <code>DB-ENCRYPTION-KEY</code> with a secure key that you generate to encrypt sensitive values in the Cloud Controller database.
      You can use any random string. For example, run the following command from a command line to generate a 32-character random string: <code>LC_ALL=C tr -dc 'A-Za-z0-9' < /dev/urandom | head -c 32 ; echo</code>
      <br /><br />
      Replace <code>CC-MUTUAL-TLS-PORT</code> with a tls port, for example: 9023
      <br /><br />

      To generate the certificates and keys for the <code>mutual_tls</code> section, you need the <a href="https://github.com/cloudfoundry/cf-release/blob/master/scripts/generate-cf-diego-certs">generate-cf-diego-certs script</a> from the <code>cf-release</code> repository.
      Run the <code>generate-cf-diego-certs</code> script to generate the certificates and keys for Cloud Controller.
      <br>
      For example, run the following in a terminal window:
      <pre class=terminal>$ ./scripts/generate-cf-diego-certs</pre>
      This script outputs a directory named <code>cf-diego-certs</code> that contains a set of files with the certificates and keys you need.

      <table>
        <tr>
          <th>In the stub, replace&hellip;</th>
          <th>with the contents of this file&hellip;</th>
        </tr>
        <tr>
          <td><code>CC-MUTUAL-TLS-CA-CERT</code></td>
          <td><code>cf-diego-ca.crt</code></td>
        </tr>
        <tr>
          <td><code>CC-MUTUAL-TLS-PUBLIC-CERT</code></td>
          <td><code>cloud_controller.crt</code></td>
        </tr>
        <tr>
          <td><code>CC-MUTUAL-TLS-PRIVATE-KEY</code></td>
          <td><code>cloud_controller.key</code></td>
        </tr>
      </table>
    </td>
  </tr>
  <tr>
    <td><pre><code>
  blobstore:
    admin_users:
      - username: BLOBSTORE-USERNAME
        password: BLOBSTORE-PASSWORD
    secure_link:
        secret: BLOBSTORE-SECRET
    tls:
        port: 443
        cert: BLOBSTORE-TLS-CERT
        private_key: BLOBSTORE-PRIVATE-KEY
        ca_cert: BLOBSTORE-CA-CERT
    </code></pre></td>
    <td>Replace <code>BLOBSTORE-USERNAME</code> and <code>BLOBSTORE-PASSWORD</code> with a username and password of your choosing.
    <br /><br />
    Replace <code>BLOBSTORE-SECRET</code> with a secure secret of your choosing.
    <br /><br />
    Replace <code>BLOBSTORE-TLS-CERT</code>, <code>BLOBSTORE-PRIVATE-KEY</code>, and <code>BLOBSTORE-CA-CERT</code> with the blobstore TLS certificate, private key, and CA certificate. You can use the <code>scripts/generate-blobstore-certs</code> script in the <code>cf-release</code> repository to generate self-signed certificates.
    </td>
  </tr>
  <tr>
    <td><pre><code>
  consul:
    encrypt_keys:
      - CONSUL-ENCRYPT-KEY
    ca_cert: CONSUL-CA-CERT
    server_cert: CONSUL-SERVER-CERT
    server_key: CONSUL-SERVER-KEY
    agent_cert: CONSUL-AGENT-CERT
    agent_key: CONSUL-AGENT-KEY
    </code></pre></td>
    <td>See the <a href="../common/consul-security.html">Security Configuration for Consul</a> topic. You can use the <code>scripts/generate-consul-certs</code> script in the <code>cf-release</code> repository to generate self-signed certificates.
    </td>
  </tr>
  <tr>
    <td><pre><code>
  etcd:
    require_ssl: true
    ca_cert: ETCD-CA-CERT
    client_cert: ETCD-CLIENT-CERT
    client_key: ETCD-CLIENT-KEY
    peer_ca_cert: ETCD-PEER-CA-CERT
    peer_cert: ETCD-PEER-CERT
    peer_key: ETCD-PEER-KEY
    server_cert: ETCD-SERVER-CERT
    server_key: ETCD-SERVER-KEY
    </code></pre></td>
    <td>
      Generate SSL certs for etcd and replace these values.
      You can use the <code>scripts/generate-etcd-certs</code> script
      in the <code>cf-release</code> repository to generate self signed certs.
    </td>
  </tr>

  <tr>
    <td><pre><code>
  loggregator:
    tls:
      ca_cert: LOGGREGATOR-CA-CERT
      doppler:
        cert: DOPPLER-CERT
        key: DOPPLER-KEY
      metron:
        cert: METRON-CERT
        key: METRON-KEY
      trafficcontroller:
        cert: TRAFFICCONTROLLER-CERT
        key: TRAFFICCONTROLLER-KEY
      cc_trafficcontroller:
        cert: CCTRAFFICCONTROLLER-CERT
        key: CCTRAFFICCONTROLLER-KEY
      syslogdrainbinder:
        cert: SYSLOGDRAINBINDER-CERT
        key: SYSLOGDRAINBINDER-KEY
      statsd_injector:
        cert: STATSDINJECTOR-CERT
        key: STATSDINJECTOR-KEY
    </code></pre></td>
    <td>
      To generate the certificates and keys for the Loggregator components, you need the following:
      <ul>
        <li>The original CA certificate and key used to sign the keypairs for TLS between the Cloud Controller and the Diego BBS</li>
        <li>The <a href="https://github.com/cloudfoundry/cf-release/blob/master/scripts/generate-loggregator-certs">generate-loggregator-certs script</a> from the <code>cf-release</code> repository</li>
        <li>The <a href="https://github.com/cloudfoundry/cf-release/blob/master/scripts/generate-statsd-injector-certs">generate-statsd-injector-certs script</a> from the <code>cf-release</code> repository</li>
      </ul>
      Generate certs
      <br>
      For example, run the following in a terminal window:
      <pre class=terminal>$  ./scripts/generate-loggregator-certs cf-diego-certs/cf-diego-ca.crt cf-diego-certs/cf-diego-ca.key
$  ./scripts/generate-statsd-injector-certs loggregator-certs/loggregator-ca.crt loggregator-certs/loggregator-ca.key</pre>
      This script outputs directories named <code>loggregator-certs</code> and <code>statsd-injector-certs</code> that contain a set of files with the certificates and keys you need for Loggregator.

      <table>
        <tr>
          <th>In the stub, replace&hellip;</th>
          <th>with the contents of this file&hellip;</th>
        </tr>
        <tr>
          <td><code>LOGGREGATOR-CA-CERT</code></td>
          <td><code>loggregator-ca.crt</code></td>
        </tr>
        <tr>
          <td><code>DOPPLER-CERT</code></td>
          <td><code>doppler.crt</code></td>
        </tr>
        <tr>
          <td><code>DOPPLER-KEY</code></td>
          <td><code>doppler.key</code></td>
        </tr>
        <tr>
          <td><code>TRAFFICCONTROLLER-CERT</code></td>
          <td><code>trafficcontroller.crt</code></td>
        </tr>
        <tr>
          <td><code>TRAFFICCONTROLLER-KEY</code></td>
          <td><code>trafficontroller.key</code></td>
        </tr>
        <tr>
          <td><code>METRON-CERT</code></td>
          <td><code>metron.crt</code></td>
        </tr>
        <tr>
          <td><code>METRON-KEY</code></td>
          <td><code>metron.key</code></td>
        </tr>
        <tr>
          <td><code>CCTRAFFICCONTROLLER-CERT</code></td>
          <td><code>cc_trafficcontroller.crt</code></td>
        </tr>
        <tr>
          <td><code>CCTRAFFICCONTROLLER-KEY</code></td>
          <td><code>cc_trafficcontroller.key</code></td>
        </tr>
        <tr>
          <td><code>SYSLOGDRAINBINDER-CERT</code></td>
          <td><code>syslogdrainbinder.crt</code></td>
        </tr>
        <tr>
          <td><code>SYSLOGDRAINBINDER-KEY</code></td>
          <td><code>syslogdrainbinder.key</code></td>
        </tr>
        <tr>
          <td><code>STATSDINJECTOR-CERT</code></td>
          <td><code>statsdinjector.crt</code></td>
        </tr>
        <tr>
          <td><code>STATSDINJECTOR-KEY</code></td>
          <td><code>statsdinjector.key</code></td>
        </tr>
      </table>
    </td>
  </tr>
  <tr>
    <td><pre><code>
  loggregator_endpoint:
    shared_secret: LOGGREGATOR-ENDPOINT-SHARED-SECRET
    </code></pre></td>
    <td>Replace <code>LOGGREGATOR-ENDPOINT-SHARED-SECRET</code> with a secure string that you generate.
    </td>
  </tr>
  <tr>
    <td><pre><code>
  nats:
    user: NATS-USER
    password: NATS-PASSWORD
    </code></pre></td>
    <td>Replace <code>NATS-USER</code> and <code>NATS-PASSWORD</code> with a username and secure password of your choosing. Cloud Foundry components use these credentials to communicate with each other over the NATS message bus.
      </td>
  </tr>
  <tr>
    <td><pre><code>
  router:
    status:
      user: ROUTER-USER
      password: ROUTER-PASSWORD
    </code></pre></td>
    <td>
      Replace <code>ROUTER-USER</code> and <code>ROUTER-PASSWORD</code> with a username and secure password of your choosing.
      </td>
  </tr>
  <tr>
    <td><pre><code>
  uaa:
    admin:
      client_secret: ADMIN-SECRET
    cc:
      client_secret: CC-CLIENT-SECRET
    clients:
      cc_service_key_client:
        secret: CC-SERVICE-KEY-CLIENT-SECRET
      cc_routing:
        secret: CC-ROUTING-SECRET
      cloud_controller_username_lookup:
        secret: CLOUD-CONTROLLER-USERNAME-LOOKUP-SECRET
      doppler:
        secret: DOPPLER-SECRET
      gorouter:
        secret: GOROUTER-SECRET
      tcp_emitter:
        secret: TCP-EMITTER-SECRET
      tcp_router:
        secret: TCP-ROUTER-SECRET
      login:
        secret: LOGIN-CLIENT-SECRET
      notifications:
        secret: NOTIFICATIONS-CLIENT-SECRET
      cc-service-dashboards:
        secret: CC-SERVICE-DASHBOARDS-SECRET
    </code></pre></td>
    <td>Replace the values for all <code>secret</code> keys with secure secrets that you generate. <br><br>You can configure container-to-container networking in this section. If you want to deploy container-to-container networking, see the <i>Enable on an IaaS</i> section of the <a href="../../devguide/deploy-apps/cf-networking.html#iaas">Administering Container-to-Container Networking</a> topic, beginning with step 4.
    </td>
  </tr>
  <tr>
    <td><pre><code>
    jwt:
      verification_key: JWT-VERIFICATION-KEY
      signing_key: JWT-SIGNING-KEY
    </code></pre></td>
    <td>Generate a PEM-encoded RSA key pair. Replace <code>JWT-SIGNING-KEY</code> with the private key, and <code>JWT-VERIFICATION-KEY</code> with the corresponding public key. Run <code>openssl genrsa -out jwt-key.pem 2048 && openssl rsa -pubout -in jwt-key.pem -out jwt-key.pem.pub</code> to generate a key pair. This command creates <code>jwt-key.pem.pub</code>, which contains your public key, and <code>jwt-key.pem</code>, which contains your private key.<br/>
Copy in the full keys, including the <code>BEGIN</code> and <code>END</code> delimiter lines.
    </td>
  </tr>
  <tr>
    <td><pre><code>
    scim:
      users:
      - name: admin
        password: ADMIN-PASSWORD
        groups:
        - scim.write
        - scim.read
        - openid
        - cloud_controller.admin
        - doppler.firehose
    </code></pre></td>
    <td>Generate a secure password and replace <code>ADMIN-PASSWORD</code> with that value to set the password for the Admin user of your Cloud Foundry installation.
    </td>
  </tr>
  <tr>
    <td><pre><code>
    ca_cert: UAA-CA-CERT
    sslCertificate: UAA-SERVER-CERT
    sslPrivateKey: UAA-SERVER-KEY
    </code></pre></td>
    <td>Replace <code>UAA-CA-CERT</code>, <code>UAA-SERVER-CERT</code>, and <code>UAA-SERVER-KEY</code> with UAA certificates. You can use the <code>scripts/generate-uaa-certs</code> script in the <code>cf-release</code> repository to generate self-signed certificates.</td>
    </td>
  </tr>
  <tr>
    <td><pre><code>
  ccdb:
    roles:
    - name: ccadmin
      password: CCDB-PASSWORD
  uaadb:
    roles:
      - name: uaaadmin
        password: UAADB-PASSWORD
  databases:
    roles:
    - name: ccadmin
      password: CCDB-PASSWORD
    - name: uaaadmin
      password: UAADB-PASSWORD
    - name: diego
      password: DIEGODB-PASSWORD
    </code></pre></td>
    <td>Replace <code>CCDB-PASSWORD</code>, <code>UAADB-PASSWORD</code>, and <code>DIEGODB-PASSWORD</code> with secure passwords of your choosing.
    </td>
  </tr>
  <tr>
    <td><pre><code>
  login:
    saml:
      serviceProviderKey: SAML-KEY
      serviceProviderKeyPassword: ''
      serviceProviderCertificate: SAML-CERT
    </code></pre></td>
    <td>
        Generate a PEM-encoded RSA key pair. Run <code>openssl req -x509 -nodes -newkey rsa:2048 -days 365 -keyout key.pem -out cert.pem</code> to generate a key pair. This command creates <code>cert.pem</code>, which contains your public key, and <code>key.pem</code>, which contains your private key. Replace <code>SERVICE-PROVIDER-PRIVATE-KEY</code> with the full private key, including the <code>BEGIN</code> and <code>END</code> delimiter lines, under <code>serviceProviderKey</code>.</br>
For RSA keys, you only need to configure the private key.
    </td>
  </tr>
  <tr>
    <td><pre><code>
jobs:
  - name: ha_proxy_z1
    networks:
      - name: cf1
        default:
        - dns
        - gateway
    properties:
      ha_proxy:
        ssl_pem: |
          -----BEGIN RSA PRIVATE KEY-----
          RSA-PRIVATE-KEY
          -----END RSA PRIVATE KEY-----
          -----BEGIN CERTIFICATE-----
          SSL-CERTIFICATE-SIGNED-BY-PRIVATE-KEY
          -----END CERTIFICATE-----
    </code></pre></td>
    <td>Replace <code>RSA-PRIVATE-KEY</code> and <code>SSL-CERTIFICATE-SIGNED-BY-PRIVATE-KEY</code> with the PEM-encoded private key and certificate associated with the system domain and apps domains that you configured to terminate at the floating IP address associated with the <code>ha_proxy</code> job. For self-signed certificates, you can run <code>openssl genrsa -out haproxy_ssl.key 2048</code> to generate a key, and run <code>openssl req -new -x509 -days 365 -key haproxy_ssl.key -out haproxy_ssl.cert</code> to generate the certificate.
    </td>
  </tr>
</table>


### Step 4: Generate Your Manifest

1. Download Alicloud infrastructure stub `cf-infrastructure-alicloud.yml` for Cloud Foundry from this [link](https://github.com/aliyun/bosh-alicloud-cpi-release/blob/release9-networks/docs/cf/stubs/cf-infrastructure-alicloud.yml). You can also set IaaS configuration by changing this file, for example changing instance_type.

2. Run the following command from the `cf-release` directory to create a deployment manifest named `cf-deployment.yml`:

    ```
    $ spiff merge \
      "${WORK_DIR}/cf-release/templates/generic-manifest-mask.yml" \
      "${WORK_DIR}/cf-release/templates/cf.yml" \
      "${WORK_DIR}/cf-infrastructure-alicloud.yml" \
      "${WORK_DIR}/cf-stub.yml" \
      > "${WORK_DIR}/cf-deployment.yml"
    ```

### Step 5: Generate Your Manifest

1. [Download alicloud light stemcell](http://bosh-alicloud.oss-cn-hangzhou.aliyuncs.com/light-bosh-stemcell-1008-alicloud-kvm-ubuntu-trusty-go_agent.tgz), and upload it to the BOSH Director, where STEMCELL-PATH is the location of your downloaded stemcell file:

    ```
    $ bosh -e YOUR_BOSH_NAME us PATH_TO_CF_LIGHT_STEMCELL
    ```

2. [Download cf release](https://bosh.io/d/github.com/cloudfoundry/cf-release?v=275), and upload it to the BOSH Director.

    ```
    $ bosh -e YOUR_BOSH_NAME ur PATH_TO_CF_RELEASE
    ```
3. Deploy. Use bosh deploy to deploy the uploaded Cloud Foundry release.

    ```
    $ bosh -e YOUR_BOSH_NAME -d YOUR_ENVIRONMENT_NAME deploy cf-deployment.yml
    ```
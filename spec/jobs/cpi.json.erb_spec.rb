require 'rspec'
require 'json'
require 'bosh/template/test'

describe 'alicloud_cpi job' do
  let(:release) { Bosh::Template::Test::ReleaseDir.new(File.join(File.dirname(__FILE__), '../..')) }
  let(:job) { release.job('alicloud_cpi') }

  describe 'cpi.json' do
    let(:template) { job.template('config/cpi.json') }

    let(:config) { JSON.parse(template.render(manifest_properties)) }

    let(:manifest_properties) do
      {
        'alicloud' => {
          'region' => 'moon'
        },
        'blobstore' => {
          'address' => 'blobstore-address.example.com',
          'agent' => {
            'user' => 'agent',
            'password' => 'agent-password'
          }
        }
      }
    end

    let(:rendered_alicloud_properties) { config['cloud']['properties']['alicloud'] }

    it 'renders the CPI config properly' do
      expect(rendered_alicloud_properties['region']).to eq('moon')
    end

    describe 'credential_source' do
      context 'when left at its default' do
        it 'renders the static access key keys' do
          expect(rendered_alicloud_properties).not_to have_key('credential_source')
          expect(rendered_alicloud_properties).to have_key('access_key_id')
          expect(rendered_alicloud_properties).to have_key('access_key_secret')
          expect(rendered_alicloud_properties).to have_key('security_token')
        end
      end

      context 'when set to static' do
        before do
          manifest_properties['alicloud']['credential_source'] = 'static'
          manifest_properties['alicloud']['access_key_id'] = 'an-access-key-id'
          manifest_properties['alicloud']['access_key_secret'] = 'an-access-key-secret'
        end

        it 'renders the configured access key' do
          expect(rendered_alicloud_properties['access_key_id']).to eq('an-access-key-id')
          expect(rendered_alicloud_properties['access_key_secret']).to eq('an-access-key-secret')
        end
      end

      context 'when set to ecs_ram_role' do
        before do
          manifest_properties['alicloud']['credential_source'] = 'ecs_ram_role'
        end

        it 'renders the credential source' do
          expect(rendered_alicloud_properties['credential_source']).to eq('ecs_ram_role')
        end

        it 'omits the static credential keys rather than rendering them empty' do
          expect(rendered_alicloud_properties).not_to have_key('access_key_id')
          expect(rendered_alicloud_properties).not_to have_key('access_key_secret')
          expect(rendered_alicloud_properties).not_to have_key('security_token')
        end

        it 'omits ram_role_name so the CPI discovers the attached role' do
          expect(rendered_alicloud_properties).not_to have_key('ram_role_name')
        end

        context 'with an explicit role name' do
          before do
            manifest_properties['alicloud']['ram_role_name'] = 'BoshDirectorRole'
          end

          it 'renders the role name' do
            expect(rendered_alicloud_properties['ram_role_name']).to eq('BoshDirectorRole')
          end
        end
      end

      context 'when set to an unsupported value' do
        before do
          manifest_properties['alicloud']['credential_source'] = 'env_or_profile'
        end

        it 'fails to render' do
          expect { config }.to raise_error(/credential_source must be 'static' or 'ecs_ram_role'/)
        end
      end
    end

    context 'when using a dav blobstore' do
      let(:rendered_blobstore) { config['cloud']['properties']['agent']['blobstore'] }

      it 'renders agent user/password for accessing blobstore' do
          expect(rendered_blobstore['options']['user']).to eq('agent')
          expect(rendered_blobstore['options']['password']).to eq('agent-password')
      end

      context 'when enabling signed URLs' do
        before do
          manifest_properties['blobstore']['agent'].delete('user')
          manifest_properties['blobstore']['agent'].delete('password')
        end

        it 'does not render agent user/password for accessing blobstore' do
          expect(rendered_blobstore['options']['user']).to be_nil
          expect(rendered_blobstore['options']['password']).to be_nil
        end
      end
    end
  end
end

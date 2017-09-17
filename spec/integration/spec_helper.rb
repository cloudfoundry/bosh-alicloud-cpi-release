require 'spec_helper'
require 'integration/helpers/ecs_helper'

RSpec.configure do |rspec_config|
  include IntegrationHelpers

  rspec_config.before(:each) do
    @access_key_id      = ENV['BOSH_AWS_ACCESS_KEY_ID']       || raise('Missing BOSH_AWS_ACCESS_KEY_ID')
    @secret_access_key  = ENV['BOSH_AWS_SECRET_ACCESS_KEY']   || raise('Missing BOSH_AWS_SECRET_ACCESS_KEY')
    @subnet_id          = ENV['BOSH_AWS_SUBNET_ID']           || raise('Missing BOSH_AWS_SUBNET_ID')
    @subnet_zone        = ENV['BOSH_AWS_SUBNET_ZONE']         || raise('Missing BOSH_AWS_SUBNET_ZONE')
    @region             = ENV.fetch('BOSH_AWS_REGION', 'us-west-1')
    @default_key_name   = ENV.fetch('BOSH_AWS_DEFAULT_KEY_NAME', 'bosh')
    @ami                = ENV.fetch('BOSH_AWS_IMAGE_ID', 'ami-866d3ee6')

    logger = Bosh::Cpi::Logger.new(STDERR)
    ec2_client = Aws::EC2::Client.new(
        region:      @region,
        access_key_id: @access_key_id,
        secret_access_key: @secret_access_key,
        logger: logger,
    )
    @ec2 = Aws::EC2::Resource.new(client: ec2_client)

    @registry = instance_double(Bosh::Cpi::RegistryClient).as_null_object
    allow(Bosh::Cpi::RegistryClient).to receive(:new).and_return(@registry)
    allow(@registry).to receive(:read_settings).and_return({})
    allow(Bosh::Clouds::Config).to receive_messages(logger: logger)
    @cpi = Bosh::AwsCloud::Cloud.new(
        'aws' => {
            'region' => @region,
            'default_key_name' => @default_key_name,
            'default_security_groups' => get_security_group_ids,
            'fast_path_delete' => 'yes',
            'access_key_id' => @access_key_id,
            'secret_access_key' => @secret_access_key,
            'max_retries' => 8
        },
        'registry' => {
            'endpoint' => 'fake',
            'user' => 'fake',
            'password' => 'fake'
        }
    )

    @stemcell_id = create_stemcell
    @vpc_id = @ec2.subnet(@subnet_id).vpc_id
  end

  rspec_config.after(:each) do
    delete_stemcell
  end
end

def vm_lifecycle(vm_disks: disks, ami_id: ami, cpi: @cpi)
  stemcell_properties = {
      'encrypted' => false,
      'ami' => {
          @region => ami_id
      }
  }
  stemcell_id = cpi.create_stemcell('/not/a/real/path', stemcell_properties)
  expect(stemcell_id).to end_with(' light')

  instance_id = cpi.create_vm(
      nil,
      stemcell_id,
      vm_type,
      network_spec,
      vm_disks,
      nil,
  )
  expect(instance_id).not_to be_nil

  expect(cpi.has_vm?(instance_id)).to be(true)

  cpi.set_vm_metadata(instance_id, vm_metadata)

  yield(instance_id) if block_given?
ensure
  cpi.delete_vm(instance_id) if instance_id
  cpi.delete_stemcell(stemcell_id) if stemcell_id
  expect(@ec2.image(ami_id)).to exist
end

def get_security_group_names(subnet_id)
  security_groups = @ec2.subnet(subnet_id).vpc.security_groups
  security_groups.map { |sg| sg.group_name }
end

def get_root_block_device(root_device_name, block_devices)
  block_devices.find do |device|
    root_device_name.start_with?(device.device_name)
  end
end


def get_target_group_arn(name)
  elb_v2_client.describe_target_groups(names: [name]).target_groups[0].target_group_arn
end

def route_exists?(route_table, expected_cidr, instance_id)
  4.times do
    route_table.reload
    found_route = route_table.routes.any? { |r| r.destination_cidr_block == expected_cidr && r.instance_id == instance_id }
    return true if found_route
    sleep 0.5
  end

  return false
end

def array_key_value_to_hash(arr_with_keys)
  hash = {}

  arr_with_keys.each do |obj|
    hash[obj.key] = obj.value
  end
  hash
end

class RegisteredInstances < StandardError; end

def ensure_no_instances_registered_with_elb(elb_client, elb_id)
  instances = elb_client.describe_load_balancers({:load_balancer_names => [elb_id]})[:load_balancer_descriptions]
  .first[:instances]

  if !instances.empty?
    raise RegisteredInstances
  end

  true
end
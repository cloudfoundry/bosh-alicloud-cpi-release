module IntegrationHelpers
  def create_stemcell
    @cpi.create_stemcell('/not/a/real/path', { 'ami' => { @region => @ami } })
  end

  def delete_stemcell
    @cpi.delete_stemcell(@stemcell_id) if @stemcell_id
  end

  def create_vm
    network_spec = {
        'default' => {
            'type' => 'dynamic',
            'cloud_properties' => { 'subnet' => @subnet_id }
        }
    }
    vm_type = {
        'instance_type' => 'm3.medium',
        'availability_zone' => @subnet_zone,
    }

    instance_id = @cpi.create_vm(
        nil,
        @stemcell_id,
        vm_type,
        network_spec,
        [],
        nil,
    )
    expect(instance_id).not_to be_nil

    instance_id
  end

  def delete_vm(instance_id)
    @cpi.delete_vm(instance_id) if instance_id
  end

  def get_security_group_ids
    security_groups = @ec2.subnet(@subnet_id).vpc.security_groups
    security_groups.map { |sg| sg.id }
  end
end
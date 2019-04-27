Vagrant.configure(2) do |config|
  config.vm.box_check_update = true
  config.vm.synced_folder File.dirname(__FILE__), "/home/vagrant/expected"
  config.vm.provision "shell", path: "#{File.dirname(__FILE__)}/hack/provision.sh"

  config.vm.box = "ubuntu/xenial64"
  config.vm.network "private_network", type: "dhcp"
  config.vm.network "forwarded_port", guest: 5432, host: 5432, protocol: "tcp"
  config.vm.network "forwarded_port", guest: 5000, host: 5000, protocol: "tcp"
  config.vm.network "forwarded_port", guest: 4222, host: 4222, protocol: "tcp"

  for i in 0..10 do
    config.vm.network "forwarded_port", guest: 3000 + i, host: 3000 + i, protocol: "tcp"
    config.vm.network "forwarded_port", guest: 4000 + i, host: 4000 + i, protocol: "tcp"
  end

  config.vm.provider "virtualbox" do |vb|
    vb.memory = "3072"
  end
end

Date: 2017-09-07
Title: Why to use Terraform and Ansible
cat: ops

####Terraform

Terraform is a great tool for deploying instances on cloud providers. It's better than JSON templates for at least three reasons:

1. You can create a single Terraform configuration that makes use of more than one provider.
2. You can use a variable number of instances, perfect for having to scale your application when automatic scaling isn't available or suitable.
3. Terraform configurations feel way shorter, e.g. an [Azure JSON template](https://raw.githubusercontent.com/azure/azure-quickstart-templates/master/101-vm-sshkey/azuredeploy.json) vs [Azure Terraform configuration](https://www.terraform.io/docs/providers/azure/r/instance.html#example-usage)

####Ansible

Once you have provisioned your instances --however many, whatever size and operating system you need-- you're going to have to configure their state somehow. Ansible is perfect for that because you can choose exactly which instances you're going to affect, and you can ensure you do so most efficiently.

**idempotency** --which most Ansible modules strive to be-- ensures that if you re-run Ansible, you'll achieve the same effect as running it the first time.

Imagine having a cluster of 4 servers, with package *A:0.0.1* installed a long time ago using a startup script `apt-get install A`. You now need to scale up to 4 + X, and X servers install *A:0.0.2* as per the startup script. Your cluster is out of sync.

With Ansible, when you scale up and re-run the *Ansible playbook* (a .yml file containing the recipe for your infrastructure's state) all servers in the cluster will have *A:0.0.2* installed. You can also include a *cache_valid_time* parameter in the Ansible apt module which will ensure that the apt repositories will not have to re-update each time you run Ansible, so you can scale up within a short time many times, efficiently and without problems.

You can assign roles to your instances. Let's say you have an infrastructure with 1-5 web servers and 0-5 workers, you can ensure each instance contains whatever it needs to contain whenever you re-provision your infrastructure. This sort of elasticity is what Cloud Computing was intended for, and *Terraform* and *Ansible* allow you to achieve it easily.

[This](https://yonatankoren.com/post/a3-jenkins-terraform-ansible) post describes how to use *Terraform* and *Ansible* together for any provider.


Date: 2017-10-11
Title: Jenkins pipelines with Terraform and Ansible
cat: ops

You can create powerful pipelines in the cloud using Jenkins, Terraform and Ansible.

![Jenkins ansible terraform pipeline](https://raw.githubusercontent.com/yonkornilov/mmt-autoscale-client/master/pipeline.png)

You can configure Jenkins to modify your cloud infrastructure, then launch Ansible to configure the nodes in your infrastructure, i.e. deploying applications, jobs, etc.

###Configuring Jenkins for Terraform

Terraform doens't like it when another user touches another user's `.tfstate`. You must configure Jenkins to use the same user that manages the terraform state.

In my case, this was the user `ubuntu`. You can do this by modifying what was in my case `/etc/default/jenkins` and changing JENKINS_USER to whichever user manages terraform.

You can refer to [this](http://blog.manula.org/2013/03/running-jenkins-under-different-user-in.html) post.

###Using a dynamic terraform inventory

Your terraform inventory is the list of hosts that are managed by terraform. If you want Ansible to know about them, you must use [terraform-inventory](https://github.com/adammck/terraform-inventory):

```
wget https://github.com/adammck/terraform-inventory/releases/download/v0.7-pre/terraform-inventory_v0.7-pre_linux_amd64.zip -O temp.zip; unzip temp.zip; rm temp.zip
```

If you have a `.tfstate` file in the same directory, you can see your dynamic inventory using:

```
./terraform-inventory --inventory terraform.tfstate
```

###Adding SSH credentials to Jenkins

First, you need to make sure that you have SSH credentials configured in Jenkins. Go to your Jenkins home:

![Jenkins credentials](/images/creds.png)

Click credentials above:

![Jenkins add credentials](/images/globalcreds.png)

Click add credentials above, then add an `SSH username with private key`:

![Jenkins ssh credentials](/images/sshcreds.png)

You can now use this credential in your build environment:

![Jenkins ssh credentials build environment](/images/buildenv.png)

###Delegate roles to instances

Specific instances should have specific or shared roles. We delegate roles to instances by adding tags to them.

Go to the `.tf` file which defines an instance, create a `tags` section:

```
resource "google_compute_instance" "worker" {
  name         = "mmt-worker"
  machine_type = "${var.worker_machine_type}"
  zone         = "${var.zone}"

  boot_disk {
    initialize_params {
      image = "${var.os_image}"
      size = "${var.size}"
    }
  }

  network_interface {
    network = "${var.network}"

    access_config {
    }
  }

  service_account {
    scopes = "${var.scopes}"
  }

  tags = ["role-worker", "env-mmt"]

  count = "${var.worker_count}"
}
```

The tags `role-worker` and `env-mmt` will be filtered by `terraform-inventory`:

```
$ ./terraform-inventory --inventory terraform.tfstate

[role-worker]
xx.xxx.xxx.xxx

[env-mmt]
xx.xxx.xxx.xxx

[all]
xx.xxx.xxx.xxx

[all:vars]

[worker]
xx.xxx.xxx.xxx

[worker.0]
xx.xxx.xxx.xxx

[type_google_compute_instance]
xx.xxx.xxx.xxx
```

Through Ansible, `role-worker` instances will only be affected only by playbooks that look like this:

```
---
- name: worker
  hosts:
    role-worker 
...
```

Specifically, `role-worker` must be declared under the hosts section in the playbook.

###Make Ansible wait for an SSH connection

Servers that were just provisioned may not have an available SSH connection. You must modify your playbook so that it waits for an SSH connection before gathering facts. i.e. you must add:

```
  gather_facts: no
  pre_tasks:
    - name: Wait for target connection to become reachable
      wait_for_connection:

    - name: gather facts
      setup: 
```

At the beginning of your playbook. For example:

```
---
- name: worker
  hosts:
    role-worker 
  gather_facts: no
  pre_tasks:
    - name: Wait for target connection to become reachable
      wait_for_connection:

    - name: gather facts
      setup: 
  vars:
  ...
```

###Adding a count variable to your instances:

In order to provision a specific instance type when it is needed for a particular job, we can add a count variable to an instance in a `.tf` file. For example, `workers.tf` contains:

```
resource "google_compute_instance" "worker" {
  
  ...

  count = "${var.worker_count}"
}
```

Notice the very last attribute, `count`.

A good place to keep the default value for variables is `variables.tf`:

```
...
variable "leader_count" {
  default = 0
}

variable "follower_count" {
  default = 0
}

variable "worker_count" {
  default = 0
}

variable "builder_count" {
  default = 0
}

```

These variables can be manipulated using the `-var` flag in Terraform:

```
$ terraform apply -var worker_count=1
```

As an example, the above command will create a worker instance if none are present, do nothing if one worker instance is present, or scale down horizontally if more than one worker is present.

```
$ terraform apply -var worker_count=0
```

The above command will destroy all worker instances.

###Adding Ansible and Terraform to your pipeline

Now, you can provision a server using Terraform, then immediately manage it using Ansible:

![Jenkins ansible terraform pipeline exec](/images/exec.png)

This pipeline provisions and runs the jobs on your cloud infrastructure, then destroys your instances as soon as the job is done. `-var worker_count=1` ensures a worker exists, and `-var worker_count=0` destroys any worker instances.

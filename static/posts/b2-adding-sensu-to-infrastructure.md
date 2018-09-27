Date: 2017-12-10
Title: Adding Sensu to your infrastructure
cat: ops

Sensu is an application used to monitor infrastructure,
using recurring checks and registering events if something bad happens
or some criteria is met.

Not surprisingly, Ansible has a [role](https://github.com/sensu/sensu-ansible) which deploys a full Sensu stack
to the inventory (the infrastructure).

I had a couple issues with this role. First, `Uchiwa` --the web interface for Sensu-- was not running behind SSL.
Second, I want my infrastructure to look like this:

![Sensu Jenkins Ansible Terraform Infrastructure](/images/sensu-infrastructure.png)

* The master node:
    * uses Terraform and provisions the rest of the infrastructure
    * uses Jenkins to automate the provisioning tasks
    * acts as the Sensu master node to initiate Jenkins projects, i.e.
scaling up the infrastructure horizontally if Sensu checks show high cpu usage / network traffic
    * Manages the state of all nodes using Ansible

The Sensu Ansible role shares variables between the master node
and the clients, and so the master deploys the role onto itself.

In order to fix the issue with `Uchiwa` and SSL, I had to [fork the role](https://github.com/yonkornilov/sensu-ansible)
and edit the deployable templates for the Uchiwa config.

You can install this role using `ansible-galaxy`:

```
$ ansible-galaxy install git+https://github.com/yonkornilov/sensu-ansible.git
```

###The Sensu Master Playbook

Next, I created a playbook for the localhost i.e. master node:

(Assume that yourdomain.com points to your localhost)

```
  - hosts: 127.0.0.1 
    tasks:
      - add_host:
          name: yourdomain.com 
          groups:
            - sensu_masters
            - redis_servers
            - rabbitmq_servers
      - set_fact:
          uchiwa_ssl_certfile: /etc/ssl/uchiwa.pem 
          uchiwa_ssl_keyfile: /etc/ssl/uchiwa.key
      - command: >
                   openssl req \
                    -new \
                    -newkey rsa:4096 \
                    -days 365 \
                    -nodes \
                    -x509 \
                    -subj "/C=CA/ST=ON/L=Toronto/O=None/CN=www.example.com" \
                    -keyout {{ uchiwa_ssl_keyfile }} \
                    -out {{ uchiwa_ssl_certfile }} 
    become: true

  - hosts: yourdomain.com
    roles:
      - { role: geerlingguy.redis }
      - { role: SimpliField.rabbitmq }
      - { role: geerlingguy.ruby }
    connection: local
    tasks:
      - replace:
          path: /etc/redis/redis.conf 
          regexp: '127.0.0.1'
          replace: '0.0.0.0'
    become: true

  - hosts: yourdomain.com 
    roles:
      - { role: sensu-ansible, sensu_master: true, sensu_include_dashboard: true} 
    vars:
      uchiwa_ssl_certfile: /etc/ssl/uchiwa.pem 
      uchiwa_ssl_keyfile: /etc/ssl/uchiwa.key
      uchiwa_dc_name: master
      rabbitmq_server: True
      redis_server: True
      sensu_remote_plugins:
        - sensu-plugins-cpu-checks
      uchiwa_users:
        - username: admin
          password: YOUR_SECRET 
      sensu_api_user_name: admin
      sensu_api_password: YOUR_SECRET
    connection: local
    become: true
```

There a few things in this playbook:

* We add localhost (127.0.0.1) to the `sensu_master`, `redis` and `rabbitmq` groups in the inventory in order for the Sensu config
* We generate the self-signed SSL key and certificate non-interactively, and without having to enter a password
* We point redis to 0.0.0.0, so the port can be reached by other nodes in the infrastructure (but security must be strengthened in production)
* We make the API and portal password-protected

###The Sensu Client Playbook

For the client nodes, I made this playbook:

```
  - hosts: env-mmt
    gather_facts: no
    pre_tasks:
      - name: Wait for target connection to become reachable
      wait_for_connection:
      - name: gather facts
      setup:
    tasks:
      - add_host:
          name: yourdomain.com
          groups:
            - sensu_masters
            - redis_servers
            - rabbitmq_servers

  - hosts: env-mmt 
    roles:
      - role: sensu-ansible
      - role: geerlingguy.ruby
    vars:
      - sensu_remote_plugins:
        - sensu-plugins-cpu-checks
    tasks:
      - command: /etc/init.d/sensu-client restart
    become: True
```

* We wait for the targets to become available, as we run this immediately after provisioning
* All hosts in group `env-mmt` will become Sensu clients of the master node at yourdomain.com
* We restart the `sensu-client` when we run this playbook
* We run this playbook after a new node is added each time we scale

###The CPU check

In order to check the cpu usage of the nodes in the infrastructure, we
create a file `data/static/sensu/definitions/check-cpu.json` in the ansible
deployment directory:

```
{
  "checks": {
    "cron": {
      "command": "check-cpu.rb",
      "subscribers": [
        "env-mmt"
      ],
      "interval": 10
    }
  }
}
```

This file checks the cpu usage of all nodes in group `env-mmt` every 10 seconds.

###Testing the Ansible, Terraform and Sensu components

Deploy the stack on the localhost:

```
$ ansible-playbook deploy/sensu_master.yml -i localhost
```

We provision a server:

```
$ terraform apply -var=worker_count=1
```

Then register this client:

```
$ ansible-playbook --inventory-file=./terraform-inventory deploy/sensu_client.yml
```

We can see Sensu in action on `https://yourdomain.com:3000`

![Uchiwa dashboard with cpu check](/images/uchiwa.png)

Now we can scale up:

```
$ terraform apply -var=worker_count=4
```

Fiannly, we register the new clients:

```
$ ansible-playbook --inventory-file=./terraform-inventory deploy/sensu_client.yml
```

###What's missing?

The only thing we're missing is the Jenkins pipeline to automate scaling when
Sensu checks are critical. We can achieve this by initiating the pipeline
using the RabbitMQ event queue, but this will be done in the following post.

# References

* [Self-signed SSL certificate non-interactive](https://superuser.com/a/226229)

Date: 2018-04-21
Title: Terraform Modules per cluster namespace
cat: ops

I've been managing infrastructure where each Kubernetes cluster (development, staging, production) has multiple namespaces within, each serving a different app.

Within each namespace, the application exists in several pods alongside a redis pod:

```
CLUSTER
  NAMESPACE A
    - APP POD 1 (Serves API A)
      - (CONTAINER)
    - APP POD 2 (Connects to Redis)
      - (CONTAINER)
    - REDIS POD
      - (CONTAINER)
  NAMESPACE B
    - APP POD 1 (Serves API B)
      - (CONTAINER)
    - APP POD 2 (Connects to Redis)
      - (CONTAINER)
    - REDIS POD
      - (CONTAINER)
```

The issue is scaling the Redis instance. The Redis clustering mode named `Redis Cluster` requires an external
program to connect the cluster nodes. This means a potential cluster requires constant administration on both
cluster startup and upon a pod being recreated --defeating the purpose of a kubernetes Deployment of a redis cluster
as a solution for single instance failure.

A solution is to use Elasticache per-namespace:

```
CLUSTER
  NAMESPACE A
    - APP POD 1 (Serves API A)
      - (CONTAINER)
    - APP POD 2 (Connects to Redis)
      - (CONTAINER)
  NAMESPACE B
    - APP POD 1 (Serves API B)
      - (CONTAINER)
    - APP POD 2 (Connects to Redis)
      - (CONTAINER)
ELASTICACHE CLUSTER A
ELASTICACHE CLUSTER B
```

We will need to make use of Terraform interpolation in order to create a cluster for each namespace.
First, let's look at the `namespaces` variable in the `main.tf` file of our cluster:

```
module "kubernetes" {
  ...
  namespaces                = ["namespace_a,namespace_b"]
  ...
}

...

module "elasticache" {
  source = "../../modules/elasticache"

  zone                      = "${module.vpc.route53_zone}"
  cluster_name              = "${var.cluster_name}"
  cluster_short_code        = "${var.cluster_short_code}"
  namespaces                = [${module.Kubernetes.namespaces}]
  allowed_security_group    = "${var.allowed_security_group}"
  private_subnets           = ["${module.vpc.private_subnets}"]
  node_groups               = "3"
  vpc                       = "${module.vpc.vpc_id}"
}
```

We supply our Elasticache module with the vpc, name of our cluster, and the kubernetes namespaces. Now our Elasticache module
has to iterate over each namespace and create the resources appropriately:

```
resource "aws_elasticache_replication_group" "ec_replication_group" {
  replication_group_id          = "${lower(var.cluster_short_code)}-${lower(element(var.namespaces,count.index))}"
  replication_group_description = "ec replication group ${lower(var.cluster_short_code)}-${lower(element(var.namespaces,count.index))}"
  automatic_failover_enabled    = true
  engine                        = "redis"
  engine_version                = "3.2.6"
  node_type                     = "cache.t2.medium"
  parameter_group_name          = "default.redis3.2.cluster.on"
  port                          = 6379
  at_rest_encryption_enabled    = true
  transit_encryption_enabled    = true
  subnet_group_name             = "${aws_elasticache_subnet_group.ec_subnet_group.name}"
  security_group_ids            = ["${aws_security_group.ec_cluster.id}"]
  auth_token                    = "${chomp(file(format("%s%s","/secrets_folder/",lower(element(var.namespaces,count.index)))))}"
  cluster_mode {
    replicas_per_node_group     = 0
    num_node_groups             = "${var.node_groups}"
  }
  count                         = "${length(var.namespaces)}"
  depends_on                    = ["aws_elasticache_subnet_group.ec_subnet_group"]
}
```

In particular, the `count` field iteration ensures we create each Elasticache cluster, so we supply it
with the number of elements in `namespaces`. We then use the `element`interpolation to reference the current
namespace being iterated on, i.e. calling `element` on `namespaces` at the current `count.index`.

Using the high-level `namespaces` variable allows us to achieve a level of abstraction akin to deployments, in that
we can have a hands-off assurance that each kubernetes namespace has a Redis connection.

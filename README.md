
# Terraform Hiera Provider
This provider implements data sources that can be used to perform hierachical data lookups with Hiera.

This is useful for providing configuration values in an environment with a high level of dimensionality or for making values from an existing Puppet deployment available in Terraform.


## Requirements
* Terraform 0.10.x
* Go 1.9 (to build the provider plugin)

## Usage
Hiera must be installed locally and configured correctly.

### Provider

```
provider "hiera" {}
```
Hiera isn't very useful without scope variables to inform its lookups:
```
provider "hiera" {
  scope {
    environment = "live"
    service     = "api"
  }
}
```
By default, the provider expects `hiera` to be available in your `$PATH`. It will also look for the configuration file at `/etc/puppetlabs/puppet/hiera.yaml`. You can override both those values if you so wish:
```
provider "hiera" {
  config = "~/hiera.yaml"
  bin    = "/usr/local/bin/hiera"
  scope {
    environment = "live"
    service     = "api"
  }
}
```
### Data Sources
This provider only implements data sources.
* `hiera` - retrieves any flat value (string, int, bool) and returns it as a string
* `hiera_array`- returns an array as a list
* `hiera_hash` - returns a hash as a map

All of the data sources accept only one parameter: the key to lookup.
```
data "hiera_hash" "tags" {
  key = "aws_tags"
}
```

All flat values are returned as strings because Terraform doesn't support typing. Instead, it performs implicit conversions depending on the usage of a particular variable.

Terraform doesn't support complex data structures either (lists of maps, maps of lists, nested maps) so neither do these data sources. If you feed `hiera_array` or `hiera_hash` a complicated value their behaviour could be unpredictable but in my testing I've found they will drop the keys/items they don't understand and keep the ones they do.

## Example
Here's an example of setting different values and data types at multiple levels in Hiera and then retrieving those values as data sources for use in outputs.

**/etc/puppetlabs/puppet/hiera.yaml**
```yaml
---
:backends:
  - yaml

:hierarchy:
  - service/%{service}
  - environment/%{environment}
  - common

:yaml:
  :datadir: /data/
```
**/data/common.yaml**
```yaml
---
aws_cloudwatch_enable: false
aws_instance_size: t2.micro
aws_tags: {}
java_opts: []
```
**/data/environment/live.yaml**
```yaml
---
aws_cloudwatch_enable: true
aws_tags:
  tier: 1
java_opts:
  - '-Dspring.profiles.active=live'
```
**/data/service/api.yaml**
```yaml
---
aws_instance_size: t2.large
aws_tags:
  team: A
java_opts:
  - '-Xms512m'
  - '-Xmx2g'
```
**main.tf**
```
provider "hiera" {
  scope {
    environment = "live"
    service = "api"
  }
}

data "hiera" "aws_instance_size" {
  key = "aws_instance_size"
}

data "hiera" "aws_cloudwatch_enable" {
  key = "aws_cloudwatch_enable"
}

data "hiera_hash" "aws_tags" {
  key = "aws_tags"
}

data "hiera_array" "java_opts" {
  key = "java_opts"
}

output "aws_instance_size" {
  value = "${data.hiera.aws_instance_size.value}"
}

output "aws_cloudwatch_enable" {
  value = "${data.hiera.aws_cloudwatch_enable.value}"
}

output "aws_tags" {
  value = "${data.hiera_hash.aws_tags.value}"
}

output "java_opts" {
  value = "${data.hiera_array.java_opts.value}"
}

```
Then, plan and apply:
```
$ terraform plan
Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.

data.hiera_hash.aws_tags: Refreshing state...
data.hiera_array.java_opts: Refreshing state...
data.hiera.aws_cloudwatch_enable: Refreshing state...
data.hiera.aws_instance_size: Refreshing state...

------------------------------------------------------------------------

No changes. Infrastructure is up-to-date.

This means that Terraform did not detect any differences between your
configuration and real physical resources that exist. As a result, no
actions need to be performed.


$ terraform apply
data.hiera.aws_cloudwatch_enable: Refreshing state...
data.hiera_array.java_opts: Refreshing state...
data.hiera_hash.aws_tags: Refreshing state...
data.hiera.aws_instance_size: Refreshing state...

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

aws_cloudwatch_enable = true
aws_instance_size = t2.large
aws_tags = {
  team = A
  tier = 1
}
java_opts = [
    -Xms512m, 
    -Xmx2g,
    -Dspring.profiles.active=live
]
```

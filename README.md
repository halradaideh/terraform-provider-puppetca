Puppet CA Terraform Provider
=============================

[![Terraform Registry Version](https://img.shields.io/badge/dynamic/json?color=blue&label=registry&query=%24.version&url=https%3A%2F%2Fregistry.terraform.io%2Fv1%2Fproviders%2Fcamptocamp%2Fpuppetca)](https://registry.terraform.io/providers/camptocamp/puppetca)
[![Go Report Card](https://goreportcard.com/badge/github.com/camptocamp/terraform-provider-puppetca)](https://goreportcard.com/report/github.com/camptocamp/terraform-provider-puppetca)
[![Build Status](https://travis-ci.org/camptocamp/terraform-provider-puppetca.svg?branch=master)](https://travis-ci.org/camptocamp/terraform-provider-puppetca)
[![By Camptocamp](https://img.shields.io/badge/by-camptocamp-fb7047.svg)](http://www.camptocamp.com)

This Terraform provider allows to connect to a Puppet Certificate Authority to verify that node certificates were signed, and clean them upon decommissioning the node.


Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/camptocamp/terraform-provider-puppetca`

```sh
$ mkdir -p $GOPATH/src/github.com/camptocamp; cd $GOPATH/src/github.com/camptocamp
$ git clone git@github.com:camptocamp/terraform-provider-puppetca
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/camptocamp/terraform-provider-puppetca
$ make build
```

Using the provider
----------------------

```hcl
provider puppetca {
  url = "https://puppetca.example.com:8140"
  cert = "certs/puppet.crt"
  key = "certs/puppet.key"
  ca = "certs/ca.pem"
}

resource "puppetca_certificate" "test" {
  name = "0a7842c26ad0.foo.com"
}

resource "puppetca_certificate" "ec2instance" {
  name   = "0a7842c26ad1.foo.com"
  usedby = aws_instance.ec2instance.id
}

# Example: Passing a CSR to the Puppet CA
resource "tls_private_key" "example" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_cert_request" "example" {
  private_key_pem = tls_private_key.example.private_key_pem

  subject {
    common_name  = "foo.example.com"
    organization = "Example Org"
  }
}

resource "puppetca_certificate" "csr_example" {
  name = "foo.example.com"
  csr  = tls_cert_request.example.cert_request_pem
  sign = true
}

resource "vm_provider" "vm" {
  name = "example-instance"
  tags = {
    Name = "example-instance"
  }
  data = <<EOF
  BLA BLA
  key = "${puppetca_certificate.csr_example.cert}"
}
EOF
```

The first `puppetca_certificate` resource, `test`, will remove the certificate if a destroy plan is run.
The second `puppetca_certificate` resource, `ec2instance`, will remove the certificate if Terraform destroys the EC2 instance.
The third `puppetca_certificate` resource, `csr_example`, will submit a CSR to the Puppet CA and sign it automatically.

The `usedby` parameter can be populated as a resource parameter to drive the removal of the certificate from the Puppet CA at the desired time.  In the example above, if a Terraform plan has to recreate the EC2 instance, the certificate will be removed when the EC2 instance is destroyed since each EC2 instance is assigned a new instance id.

The `csr` parameter allows you to pass a Certificate Signing Request (CSR) to the Puppet CA. In the example above:
- A private key is generated using the `tls_private_key` resource.
- A CSR is created using the `tls_cert_request` resource.
- The CSR is passed to the `puppetca_certificate` resource using the `csr` attribute.
- The `sign` parameter ensures the certificate is signed automatically after submission.

The provider can also be configured using environment variables:

```sh
export PUPPETCA_URL="https://puppetca.example.com:8140"
export PUPPETCA_CA=$(cat certs/ca.pem)
export PUPPETCA_CERT=$(cat certs/puppet.crt)
export PUPPETCA_KEY=$(cat certs/puppet.key)
```

The provider needs to be configured with a certificate. This certificate
should be signed by the CA, and have specific rights to list and delete
certificates. See [the Puppet docs](https://puppet.com/docs/puppetserver/5.3/config_file_auth.html)
for how to configure your Puppet Master to give these rights to your
certificate. For example, if your certificate uses the `pp_employee` extension,
you could add a rule like the following:

```ruby
{                                                                         
    match-request: {
        path: "^/puppet-ca/v1/certificate(_status|_request)?/([^/]+)$"
        type: regex
        method: [delete]
    }
    allow: [
      {extensions:{pp_employee: "true"}},
      ]
    sort-order: 500
    name: "let employees delete certs"
},
```


Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-puppetca
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

In order to test the code locally, you can use `make local-install` to deploy the generated provider.

```hcl
terraform {
  required_providers {
    puppetca = {
      source  = "local/puppetca/puppetca"
      version = "0.1.0"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.1"
    }
  }
}

locals {
  FQDN = "foo.example.com"
}

provider puppetca {
  url = "https://puppetca.example.com:8140"
  cert = "certs/puppet.crt"
  key = "certs/puppet.key"
  ca = "certs/ca.pem"
}

# Example: Passing a CSR to the Puppet CA
resource "tls_private_key" "example" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_cert_request" "example" {
  private_key_pem = tls_private_key.example.private_key_pem

  subject {
    common_name  = local.FQDN
    organization = "Example Org"
  }
}

resource "puppetca_certificate" "csr_example" {
  name = local.FQDN
  csr  = tls_cert_request.example.cert_request_pem
}
```
Puppet CA Terraform Provider
=============================

[![CI/CD Pipeline](https://github.com/halradaideh/terraform-provider-puppetca/actions/workflows/ci-cd.yml/badge.svg)](https://github.com/halradaideh/terraform-provider-puppetca/actions/workflows/ci-cd.yml)
[![Terraform Registry Version](https://img.shields.io/badge/dynamic/json?color=blue&label=registry&query=%24.version&url=https%3A%2F%2Fregistry.terraform.io%2Fv1%2Fproviders%2Fhalradaideh%2Fpuppetca)](https://registry.terraform.io/providers/halradaideh/puppetca)
[![Go Report Card](https://goreportcard.com/badge/github.com/halradaideh/terraform-provider-puppetca)](https://goreportcard.com/report/github.com/halradaideh/terraform-provider-puppetca)
[![GitHub Release](https://img.shields.io/github/v/release/halradaideh/terraform-provider-puppetca)](https://github.com/halradaideh/terraform-provider-puppetca/releases)
[![License](https://img.shields.io/github/license/halradaideh/terraform-provider-puppetca)](https://github.com/halradaideh/terraform-provider-puppetca/blob/master/LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/halradaideh/terraform-provider-puppetca)](https://github.com/halradaideh/terraform-provider-puppetca/blob/master/go.mod)

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

  timeouts {
    create = "60m"
    update = "30m"
    delete = "10m"
  }
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

## Timeouts

**New in v1.0.0:** The provider now supports configurable timeouts for certificate operations to prevent hanging operations.

```hcl
resource "puppetca_certificate" "example" {
  name = "example-node"
  env  = "production"
  sign = true

  timeouts {
    create = "60m"  # Default: 20m
    update = "30m"  # Default: 20m
    delete = "10m"  # Default: 20m
  }
}
```

Timeout values can be specified in:
- `s` for seconds
- `m` for minutes  
- `h` for hours

If no timeout is specified, operations default to 20 minutes. This prevents certificate operations from hanging indefinitely when the Puppet CA is slow or unresponsive.

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

## Local Development and Installation

For local development and testing, you can use the following make targets:

```sh
# Install to development path (recommended for local testing)
$ make local-install-dev

# Install to legacy path
$ make local-install

# Install to modern registry path
$ make local-install-modern

# Install to all paths
$ make install-local
```

After installation, use the provider in your Terraform configuration:

```hcl
terraform {
  required_providers {
    puppetca = {
      source  = "local/puppetca/puppetca"
      version = "1.0.0"
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

  timeouts {
    create = "45m"
    update = "20m"
    delete = "5m"
  }
}
```

## CI/CD and Development Workflows

This project uses GitHub Actions for continuous integration and deployment:

### 🔄 **Unified CI/CD Pipeline**
- **CI Triggers**: Pull requests with `ready-to-test` label, pushes to master
- **CD Triggers**: Pull requests with `ready-to-deploy` label, manual workflow dispatch
- **Quality Checks**: Go formatting, linting (golangci-lint), static analysis
- **Testing**: Unit tests on Go 1.20 & 1.21, race detection
- **Multi-Architecture Builds**: Linux, macOS, Windows, FreeBSD (amd64, arm64, etc.)
- **Release Process**: GoReleaser with GPG signing, automatic tagging
- **Registry**: Automatic publication to `halradaideh/puppetca` provider
- **Artifacts**: Signed binaries, SHA256 checksums, release notes

### 📋 **Development Commands**
```bash
# Run all CI checks locally
make ci-test

# Run quality checks (format, lint, vet)
make quality

# Build for multiple architectures
make ci-build

# Prepare for release
make pre-release

# Clean build artifacts
make clean
```

### 🏷️ **Creating Releases**
```bash
# Method 1: Using PR labels (recommended)
# 1. Create a PR with title containing version (e.g., "Release v1.0.1")
# 2. Add the "ready-to-deploy" label to trigger release

# Method 2: Manual workflow dispatch
# Go to Actions tab → CI/CD Pipeline → Run workflow → Enable "Force deployment"
```

For detailed CI/CD documentation, see [docs/CICD.md](docs/CICD.md).

### 📦 **Required Secrets for Releases**
- `GPG_PRIVATE_KEY`: GPG private key for signing releases
- `GPG_PASSPHRASE`: Passphrase for the GPG private key
# Resource: puppetca_certificate

Manages Puppet CA certificates.  
**New in v2.2.0:** You can now submit a Certificate Signing Request (CSR) directly from Terraform.

## Example Usage

### Submit and Sign a CSR

```hcl
resource "puppetca_certificate" "mycert" {
  name = "myhost.example.com"
  csr  = file("myhost.csr")
}
```

### Retrieve a Signed Certificate

```hcl
output "cert" {
  value = puppetca_certificate.mycert.cert
}
```

### Using Timeouts

```hcl
resource "puppetca_certificate" "mycert" {
  name = "myhost.example.com"
  csr  = file("myhost.csr")

  timeouts {
    create = "60m"
    update = "30m"
    delete = "10m"
  }
}
```

## Argument Reference

- `name` (String, Required): The node name for the certificate.
- `csr` (String, Optional): The PEM-encoded CSR to submit to the Puppet CA. If omitted, the provider will only attempt to retrieve or sign an existing certificate request. Cannot be used together with `sign`.
- `sign` (Bool, Optional): Whether to sign an existing certificate request. Defaults to `false`. Cannot be used together with `csr`.
- `env` (String, Optional): Puppet environment name.
- `usedby` (String, Optional): An optional string to indicate who or what uses this certificate.

### Timeouts

- `create` (String, Optional): Timeout for certificate creation operations. Defaults to `20m`.
- `update` (String, Optional): Timeout for certificate update operations. Defaults to `20m`.
- `delete` (String, Optional): Timeout for certificate deletion operations. Defaults to `20m`.

## Attribute Reference

- `cert` (String): The signed certificate in PEM format.

## Import

Import is supported using:

```
terraform import puppetca_certificate.example "nodename,environment"
```

## Notes

- The `csr` field must contain a valid PEM-encoded CSR. Use `file("path/to/file.csr")` to load from disk.
- The `csr` and `sign` attributes cannot be used together. Use `csr` to submit a new certificate request, or use `sign` to sign an existing certificate request.
- When using `csr`, the certificate signing must be handled outside of Terraform (e.g., through Puppet CA admin tools).
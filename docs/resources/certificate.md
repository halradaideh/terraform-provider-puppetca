# Resource: puppetca_certificate

Manages Puppet CA certificates.  
**New in v1.0.0:** You can now submit a Certificate Signing Request (CSR) directly from Terraform.

## Example Usage

### Submit and Sign a CSR

```hcl
resource "puppetca_certificate" "mycert" {
  name = "myhost.example.com"
  csr  = file("myhost.csr")
  sign = true
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
  sign = true

  timeouts {
    create = "60m"
    update = "30m"
    delete = "10m"
  }
}
```

## Argument Reference

- `name` (String, Required): The node name for the certificate.
- `csr` (String, Optional): The PEM-encoded CSR to submit to the Puppet CA. If omitted, the provider will only attempt to retrieve or sign an existing CSR.
- `sign` (Bool, Optional): Whether to sign the certificate after CSR submission. Defaults to `false`.
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
- If `sign` is `true`, the provider will attempt to sign the submitted CSR after submission.
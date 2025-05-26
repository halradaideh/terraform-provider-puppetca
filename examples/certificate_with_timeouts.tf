resource "puppetca_certificate" "example" {
  name = "example-node"
  env  = "production"
  sign = true

  timeouts {
    create = "60m"
    update = "30m"
    delete = "10m"
  }
}

# Example with CSR
resource "puppetca_certificate" "example_with_csr" {
  name = "example-node-csr"
  env  = "production"
  csr  = file("path/to/certificate.csr")

  timeouts {
    create = "45m"
    update = "20m"
    delete = "5m"
  }
} 
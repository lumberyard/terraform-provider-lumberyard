# Terraform Provider Lumberyard

This repository is the home of the `lumberyard` Terraform provider. This provider allows you to interact with a variety of data sources and resources, with the first being the `dirmap` data source.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23

## Using the Provider

To use the `lumberyard` provider, you will need to add it to your Terraform configuration.

```hcl
terraform {
  required_providers {
    lumberyard = {
      source  = "lumberyard/lumberyard"
      version = "~> 1.0"
    }
  }
}

provider "lumberyard" {}
```

## `dirmap` Data Source

The `dirmap` data source allows you to traverse a nested directory of YAML and JSON files and construct a nested object from the directory structure and file contents.

### Usage

```hcl
data "lumberyard_dirmap" "example" {
  path = "./environments"
}

output "dirmap" {
  value = data.lumberyard_dirmap.example.result
}
```

### Arguments

- `path` - (Required) The base directory to traverse.
- `filter` - (Optional) A glob pattern to filter files.

### Attributes

- `result` - The constructed object from the directory structure.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request if you have a feature request or bug report.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`. You can also use `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

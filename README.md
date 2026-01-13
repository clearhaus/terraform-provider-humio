# Terraform provider for Humio

A (maintained) fork of https://github.com/humio/terraform-provider-humio which was deleted.

## Currently tested with

- [Terraform](https://opentofu.org) v1.10.7
- [Go](https://golang.org/doc/install) 1.25 (to build the provider plugin)

## Installing the provider

The provider is published in [the official Hashicorp Terraform registry](https://registry.terraform.io/providers/clearhaus/humio), so it can used like any other published provider:

```hcl
terraform {
  required_providers {
    humio = {
      source  = "clearhaus/humio"
      version = "1.0.0"
    }
  }
}
```

Alternatively, the provider be built locally for development.
To do so:


1. Clone the git repository of the provider

    ```bash
    git clone https://github.com/clearhaus/terraform-provider-humio
    cd terraform-provider-humio
    ```

2. Build the provider plugin

    ```bash
    go build -o terraform-provider-humio
    ```

3. Create `~/.terraformrc` with the following content:

    ```hcl
    provider_installation {

    dev_overrides {
        "clearhaus/humio" = "<path-to-repo>"
    }

    # For all other providers, install them directly from their origin provider
    # registries as normal. If you omit this, Terraform will _only_ use
    # the dev_overrides block, and so no other providers will be available.
    direct {}
    }
    ```

    Make sure to update the path in the above to the recently cloned `terraform-provider-humio` repository.

When running a Terraform command with the publisher defined in `~/.terraformrc`, the following warning will be shown:

```hcl
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - clearhaus/humio in <path-to-repo>
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
```

## Using the provider

### Authentication

The provider is configured through the environment variables `HUMIO_ADDR` and `HUMIO_API_TOKEN`.

Alternatively, the configurations can be hardcoded in the Terraform provider definition:

```hcl
provider "humio" {
  addr      = "https://humio.example.com/"
  api_token = "XXXXXXXXXXXXXXXXXXXXXXXXX"
}
```

If no `addr` is defined, `https://cloud.humio.com/` will be used.

It's recommended to configure the address directly in the Terraform provider and the API key using the environment variable.

### Supported resources and examples

See [examples directory](examples/).

terraform {
  required_providers {
    humio = {
      source  = "clearhaus/humio"
      version = "0.5.0"
    }
  }
}

provider "humio" {
  addr = "https://cloud.humio.com/"
}

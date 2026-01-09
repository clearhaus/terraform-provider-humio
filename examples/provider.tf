terraform {
  required_providers {
    humio = {
      source = "clearhaus/humio"
    }
  }
}

provider "humio" {
  addr = "https://cloud.humio.com/"
}

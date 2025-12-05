terraform {
  required_providers {
    humio = {
      source  = "clearhaus/humio"
      version = "1.0.0"
    }
  }
}

provider "humio" {
  addr = "https://cloud.humio.com/"
}

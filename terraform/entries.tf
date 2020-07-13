locals {
  entries = [
    {
      zone      = ".example.com."
      name      = "dns"
      addresses = ["10.10.1.1"]
      tags = {
        team = "ops"
      }
      description = "internal DNS"
    },
    {
      zone      = ".example.com."
      name      = "nexus"
      addresses = ["10.10.1.2"]
      tags = {
        team = "devs"
      }
      description = "Nexus Server"
    },
    {
      zone      = ".example.com."
      name      = "bookmarks"
      addresses = ["10.10.1.5"]
      tags = {
        team = "devs"
      }
      description = "Bookmarks Service"
    },
    {
      zone      = ".example.com."
      name      = "crm"
      addresses = ["10.10.1.10"]
      tags = {
        team = "marketing"
      }
      description = "Marketing CRM"
    },
  ]
}

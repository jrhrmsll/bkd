# A Bookmarks Service synced from Terraform

This project is a proof of concept for a Bookmarks service synced from Terraform upon changes on DNS entries.

## Install

This project is conceived to be installed in a local environment, with Docker.

### Bookmarks Service
After cloning the repository the first step is build the "Bookmarks Service" (bkd) container with a multi-stage build.

```
docker build -t bkd .
```

### Docker Compose
A Docker Compose file is provided, then `docker-compose up` bring **bkd**, **MongoDB** and **CoreDNS**.

### Making changes
To see all the pieces in action the way to go is modifying the [entries.tf](https://github.com/jrhrmsll/bkd/blob/master/terraform/entries.tf) file and execute a `plan` followed by an `apply` inside the `terraform` directory. For example, adding:

```
{
      zone      = ".example.com."
      name      = "jenkins"
      addresses = ["10.10.1.15"]
      tags = {
        team = "devs"
      }
      description = "Jenkins Server"
},
```

To execute it run:

```
terraform init
terraform plan -out output.plan
terraform apply output.plan
```

To check the changes do a request to the [bookmarks list](http://localhost:8080/bookmarks).

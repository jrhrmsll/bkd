locals {
  bookmarks_service_url = "http://localhost:8080/bookmarks"
  content_type          = "Content-Type: application/json"

  version = formatdate("YYMMhhmmss", timestamp())

  records = flatten([
    for entry in local.entries : [
      for address in entry.addresses :
      join(" ", [join("", [entry.name, entry.zone]), "IN", "A", address])
    ]
  ])
}

data "template_file" "zone_db_tpl" {
  template = "${file("${path.module}/zone.db.tpl")}"

  vars = {
    version = local.version,
    records = join("\n", local.records)
  }
}

resource "local_file" "zone_db" {
  content  = data.template_file.zone_db_tpl.rendered
  filename = "${path.module}/../coredns/example.db"
}

locals {
  update_payloads = [
    for entry in local.entries :
    jsonencode({
      version     = local.version
      url         = trimsuffix(join("", [entry.name, entry.zone]), ".")
      description = entry.description
      tags        = entry.tags
      mode        = "auto"
    })
  ]
}

resource "null_resource" "bookmarks_update" {
  count = length(local.update_payloads)

  depends_on = [local_file.zone_db]

  triggers = {
    db = local.version
  }

  provisioner "local-exec" {
    command = "curl -d '${local.update_payloads[count.index]}' -H '${local.content_type}' -X POST ${local.bookmarks_service_url}"
  }
}

locals {
  cleanup_payload = jsonencode({ version = local.version })
}

resource "null_resource" "bookmarks_cleanup" {
  depends_on = [null_resource.bookmarks_update]

  triggers = {
    db = local.version
  }

  provisioner "local-exec" {
    command = "curl -d '${local.cleanup_payload}' -H '${local.content_type}' -X DELETE ${local.bookmarks_service_url}"
  }
}

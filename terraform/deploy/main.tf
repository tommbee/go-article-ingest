resource "helm_release" "article-ingest-k8s" {
    name      = "article-ingest-k8s"
    chart     = "../article-ingest-k8s"
    namespace = "${var.namespace}"

    set {
        name  = "image.repository"
        value = "${var.image_repository}"
    }
    set {
        name  = "image.tag"
        value = "${var.image_tag}"
    }
    set {
        name  = "sources"
        value = "${var.sources}"
    }
    set {
        name  = "server"
        value = "${var.server}"
    }
    set {
        name  = "db"
        value = "${var.db}"
    }
    set {
        name  = "configFileLocation"
        value = "${var.config_file_location}"
    }
    set {
        name  = "articleCollection"
        value = "${var.article_collection}"
    }
    set {
        name  = "dbUser"
        value = "${var.db_user}"
    }
    set {
        name  = "dbPassword"
        value = "${var.db_password}"
    }
    set {
        name  = "authDb"
        value = "${var.auth_db}"
    }
    set {
        name  = "dbSsl"
        value = "${var.db_ssl}"
    }
}

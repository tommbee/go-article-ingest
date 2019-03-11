terraform {
  backend "gcs" {
    bucket = "article-app-storage"
    prefix = "terraform/state"
  }
}

module "article-app" {
  source = "git@github.com:tommbee/article-app-terraform.git"

  config_file = "${var.config_file}"
}

module "deploy" {
    source = "./deploy"

    client_certificate = "${base64decode(module.article-app.client_certificate)}"
    client_key = "${base64decode(module.article-app.client_key)}"
    cluster_ca_certificate = "${base64decode(module.article-app.cluster_ca_certificate)}"
    host = "${module.article-app.host}"
    token = "${module.article-app.token}"
    service_account_name = "${module.article-app.helm_service_account}"

    ## app specific
    image_repository = "${var.image_repository}"
    image_tag = "${var.image_tag}"
    sources = "${var.sources}"
    server = "${var.server}"
    db = "${var.db}"
    config_file_location = "${var.config_file_location}"
    article_collection = "${var.article_collection}"
    db_user = "${var.db_user}"
    db_password = "${var.db_password}"
    auth_db = "${var.auth_db}"
    db_ssl = "${var.db_ssl}"
    namespace = "${var.namespace}"
}

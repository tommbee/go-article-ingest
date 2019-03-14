terraform {
  backend "gcs" {
    bucket = "article-app-storage"
    prefix = "terraform/state"
  }
}

module "article-app-cluster" {
  source = "git@github.com:tommbee/k8s-prometheus-terraform-module.git"

  config_file = "${var.config_file}"
  region = "europe-west1-c"
  projet_name = "temporal-parser-229715"
}

module "deploy" {
    source = "./deploy"

    client_certificate = "${module.article-app-cluster.client_certificate}"
    client_key = "${module.article-app-cluster.client_key}"
    cluster_ca_certificate = "${module.article-app-cluster.cluster_ca_certificate}"
    host = "${module.article-app-cluster.host}"
    helm_service_account = "${module.article-app-cluster.helm_service_account}"
    helm_namespace = "${module.article-app-cluster.helm_namespace}"
    token = "${module.article-app-cluster.token}"
    helm_init_id = "${module.article-app-cluster.helm_init_id}"

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

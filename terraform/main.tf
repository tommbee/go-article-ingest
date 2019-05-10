# terraform {
#   backend "gcs" {
#     bucket = "article-app-storage"
#     prefix = "terraform/state"
#   }
# }

provider "google" {
  credentials = "${file("${var.config_file}")}"
}

data "google_storage_bucket_object" "kubeconfig" {
  name   = "kubeconfig"
  bucket = "${var.gcs_bucket}"
}

resource "local_file" "kubeconfig" {
    content     = "${data.google_storage_bucket_object.kubeconfig.md5hash}"
    filename    = "${path.module}/kubeconfig"
}

module "deploy" {
    source = "./deploy"

    helm_service_account = "${var.helm_service_account}"
    helm_namespace = "${var.helm_namespace}"
    kubeconfig = "${local_file.kubeconfig.filename}"

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

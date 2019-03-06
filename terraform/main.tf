module "article-app" {
  source = "git@github.com:tommbee/article-app-terraform.git"

  config_file = "${var.config_file}"
}

module "deploy" {
    source = "./deploy"

    client_certificate = "${module.article-app.client_certificate}"
    client_key = "${module.article-app.client_key}"
    cluster_ca_certificate = "${base64decode(module.article-app.cluster_ca_certificate)}"
}

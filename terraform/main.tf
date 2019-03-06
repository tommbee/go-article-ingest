module "article-app" {
  source = "git@github.com:tommbee/article-app-terraform.git"

  config_file = "${var.config_file}"
}

module "deploy" {
    source = "./deploy"
    cluster_client_certificate ="${module.article-app.client_certificate}"
    cluster_client_key = "${module.article-app.client_key}"
    cluster_ca_certificate = "${module.article-app.cluster_ca_certificate}"
    host = "${module.article-app.client_certificate}"
}

provider "helm" {
  version = "~> 0.6"

  kubernetes {
    host = "https://${data.article-app.article-app-initial-primary.endpoint}"
    token = "${data.article-app.default.access_token}"
    cluster_ca_certificate = "${base64decode(data.article-app.article-app-initial-primary.master_auth.0.cluster_ca_certificate)}"
  }
}

provider "kubernetes" {
  load_config_file = false

  host = "https://${data.article-app.article-app-initial-primary.endpoint}"
  token = "${data.article-app.default.access_token}"
  cluster_ca_certificate = "${base64decode(data.article-app.article-app-initial-primary.master_auth.0.cluster_ca_certificate)}"
}

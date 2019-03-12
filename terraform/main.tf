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

# module "deploy" {
#     source = "./deploy"

#     client_certificate = "${module.article-app-cluster.client_certificate}"
#     client_key = "${module.article-app-cluster.client_key}"
#     cluster_ca_certificate = "${module.article-app-cluster.cluster_ca_certificate}"
#     host = "${module.article-app-cluster.host}"
#     #service_account_name = "${module.article-app-cluster.service_account}"
#     token = "${module.article-app-cluster.token}"

#     ## app specific
#     image_repository = "${var.image_repository}"
#     image_tag = "${var.image_tag}"
#     sources = "${var.sources}"
#     server = "${var.server}"
#     db = "${var.db}"
#     config_file_location = "${var.config_file_location}"
#     article_collection = "${var.article_collection}"
#     db_user = "${var.db_user}"
#     db_password = "${var.db_password}"
#     auth_db = "${var.auth_db}"
#     db_ssl = "${var.db_ssl}"
#     namespace = "${var.namespace}"
# }

#### new here


# variable "projet_name" {
#     default = "temporal-parser-229715"
# }
# variable "region" {
#   default = "europe-west1-c"
# }

# provider "google" {
#   credentials = "${file("${var.config_file}")}"
#   project = "${var.projet_name}"
# }

# data "google_client_config" "current" {}

# resource "google_container_cluster" "default" {
#   name = "${var.projet_name}-initial-primary"

#   zone = "${var.region}"
#   initial_node_count = 2

#   min_master_version = 1.11
#   node_version = 1.11

#   node_config {
#     machine_type = "n1-standard-2"

#     oauth_scopes = [
#       "https://www.googleapis.com/auth/compute",
#       "https://www.googleapis.com/auth/devstorage.read_only",
#       "https://www.googleapis.com/auth/logging.write",
#       "https://www.googleapis.com/auth/monitoring",
#     ]
#   }
# }

provider "kubernetes" {
  #client_certificate = "${module.article-app-cluster.client_certificate}"
  #client_key = "${module.article-app-cluster.client_key}"
  cluster_ca_certificate = "${module.article-app-cluster.cluster_ca_certificate}"
  host = "${module.article-app-cluster.host}"
  token = "${module.article-app-cluster.token}"

  load_config_file = false
}

resource "kubernetes_service_account" "tiller" {
  metadata {
    name      = "tiller"
    namespace = "kube-system"
  }
  
  automount_service_account_token = true
}

resource "kubernetes_cluster_role_binding" "tiller" {
  metadata {
    name = "tiller"
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = "cluster-admin"
  }

  subject {
    api_group = ""
    kind      = "ServiceAccount"
    name      = "tiller"
    namespace = "kube-system"
  }
}

variable "helm_version" {
  default = "v2.9.1"
}

provider "helm" {
  service_account = "${kubernetes_service_account.tiller.metadata.0.name}"
  namespace       = "${kubernetes_service_account.tiller.metadata.0.namespace}"
  #install_tiller  = false

  kubernetes {
    #client_certificate     = "${var.client_certificate}"
    #client_key             = "${var.client_key}"
    cluster_ca_certificate = "${module.article-app-cluster.cluster_ca_certificate}"
    host = "${module.article-app-cluster.host}"
    token = "${module.article-app-cluster.token}"
  }
}

resource "helm_release" "article-ingest-k8s" {
    #depends_on = ["kubernetes_service_account.tiller"]
    name      = "article-ingest-k8s"
    chart     = "../article-ingest-k8s"
    namespace = "${var.namespace}"
    depends_on = [
      "kubernetes_service_account.tiller",
      "kubernetes_cluster_role_binding.tiller",
    ]

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

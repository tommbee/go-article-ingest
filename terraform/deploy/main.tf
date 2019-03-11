provider "helm" {
  version = "~> 0.6"
  service_account = "${kubernetes_service_account.tiller.metadata.0.name}"

  kubernetes {
    host                   = "${var.host}"
    token                  = "${var.token}"
    #client_certificate     = "${var.client_certificate}"
    #client_key             = "${var.client_key}"
    cluster_ca_certificate = "${var.cluster_ca_certificate}"
  }
}

provider "kubernetes" {
  host                   = "${var.host}"
  token                  = "${var.token}"
  #client_certificate     = "${var.client_certificate}"
  #client_key             = "${var.client_key}"
  cluster_ca_certificate = "${var.cluster_ca_certificate}"

  load_config_file = false
}

resource "kubernetes_service_account" "tiller" {
  metadata {
    name      = "tiller"
    namespace = "kube-system"
  }
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

resource "helm_release" "prometheus_operator" {
  name  = "monitoring"
  chart = "stable/prometheus-operator"

  values = [
    "${file("${path.module}/monitoring/prometheus.yml")}",
  ]
}

resource "helm_release" "article-ingest-k8s" {
    depends_on = ["kubernetes_service_account.tiller"]
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

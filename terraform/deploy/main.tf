provider "helm" {
  version = "~> 0.6"

  kubernetes {
    client_certificate     = "${var.cluster_client_certificate}"
    client_key             = "${var.cluster_client_key}"
    cluster_ca_certificate = "${var.cluster_ca_certificate}"
    host                   = "${var.host}"
  }
}

provider "kubernetes" {
  client_certificate     = "${var.cluster_client_certificate}"
  client_key             = "${var.cluster_client_key}"
  cluster_ca_certificate = "${var.cluster_ca_certificate}"
  host                   = "${var.host}"
}

resource "helm_release" "article-ingest-k8s" {
 
}

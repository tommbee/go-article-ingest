variable "config_file" {
  default="article-app-google-account-creds.json"
}
variable "helm_namespace" {}
variable "helm_service_account" {}
variable "gcs_bucket" {}

## app specific
variable "image_repository" {}
variable "image_tag" {}
variable "sources" {}
variable "server" {}
variable "db" {}
variable "config_file_location" {}
variable "article_collection" {}
variable "db_user" {}
variable "db_password" {}
variable "auth_db" {}
variable "db_ssl" {}
variable "namespace" {}

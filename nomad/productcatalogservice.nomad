job "productcatalogservice-job" {
  type        = "service"
  datacenters = ["eu-west-1a"]

  update {
    canary            = 1
    max_parallel      = 1
    min_healthy_time  = "5s"
    healthy_deadline  = "1m"
    progress_deadline = 0 # fail immediately
    auto_revert       = true
    auto_promote      = true
  }

  group "productcatalogservice-group" {
    count = 2

    network {
      mode = "bridge"
      port "grpc" {
        to = 8502
      }
    }

    service {
      name = "productcatalogservice"
      port = 8502

      connect {
        sidecar_service {}
      }

      check {
        type     = "grpc"
        interval = "5s"
        timeout  = "2s"
        port     = "grpc"
      }
    }

    task "server" {
      driver = "docker"

      config {
        image      = "ghcr.io/simonnordberg/productcatalogservice:main"
        force_pull = true
      }

      resources {
        cpu    = 200
        memory = 100
      }
    }
  }
}
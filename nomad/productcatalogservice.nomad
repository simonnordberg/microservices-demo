job "productcatalogservice-job" {
  type        = "service"
  datacenters = ["eu-west-1a"]

  update {
    max_parallel     = 1
    min_healthy_time = "5s"
    healthy_deadline = "1m"
  }

  group "productcatalogservice-group" {
    count = 1

    network {
      mode = "bridge"

      port "grpc" {
        to = 8502
      }
    }

    service {
      name = "productcatalogservice"
      port = "grpc"

      connect {
        sidecar_service {}
      }

      check {
        type     = "grpc"
        interval = "5s"
        timeout  = "5s"
      }
    }

    task "server" {
      driver = "docker"

      config {
        image      = "ghcr.io/simonnordberg/productcatalogservice:main"
        ports      = ["grpc"]
        force_pull = true
      }

      resources {
        cpu    = 100
        memory = 32
      }
    }
  }
}
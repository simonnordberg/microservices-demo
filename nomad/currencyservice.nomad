job "currencyservice-job" {
  type = "service"
  datacenters = ["dc1"]

  update {
    max_parallel     = 1
    min_healthy_time = "5s"
    healthy_deadline = "1m"
  }

  group "currencyservice-group" {
    count = 3

    network {
      port "grpc" {
        to = 3550
      }
    }

    service {
      name = "currencyservice"
      port = "grpc"
      check {
        type     = "grpc"
        interval = "15s"
        timeout  = "5s"
      }
    }

    task "server" {
      driver = "docker"

      config {
        image = "ghcr.io/simonnordberg/currencyservice:main"
        ports = ["grpc"]
      }

      resources {
        cpu    = 140
        memory = 64
      }
    }
  }
}
job "frontend-job" {
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

  group "frontend-group" {
    count = 1

    network {
      mode = "bridge"

      port "http" {
        to     = 8080
        static = 8080
      }
    }

    service {
      name = "frontend"
      port = "http"

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "currencyservice"
              local_bind_port  = 5000
            }
            upstreams {
              destination_name = "productcatalogservice"
              local_bind_port  = 5001
            }
          }
        }
      }

      check {
        type     = "http"
        path     = "/"
        interval = "5s"
        timeout  = "5s"
      }
    }

    task "server" {
      driver = "docker"

      env {
        CURRENCY_SERVICE_ADDR        = "${NOMAD_UPSTREAM_ADDR_currencyservice}"
        PRODUCT_CATALOG_SERVICE_ADDR = "${NOMAD_UPSTREAM_ADDR_productcatalogservice}"
        GRPC_GO_LOG_SEVERITY_LEVEL   = "info"
        GRPC_GO_LOG_VERBOSITY_LEVEL  = 2
        LOG_LEVEL                    = "trace"
      }

      config {
        image      = "ghcr.io/simonnordberg/frontend:main"
        #        ports      = ["http"]
        force_pull = true
      }

      resources {
        cpu    = 200
        memory = 100
      }
    }
  }
}
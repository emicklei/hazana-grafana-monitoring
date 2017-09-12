# hazana-grafana-monitoring

# install

    go get github.com/emicklei/hazana-grafana-monitoring

# example

    package main

    import (
        "github.com/emicklei/hazana"
        monitoring "github.com/emicklei/hazana-grafana-monitoring"
    )

    func main() {
        attack := NewZombieAttack("zombies.com:50051")
        hazana.Run(monitoring.NewMonitor(attack), hazana.ConfigFromFlags())
    }

# grafana

    docker run -d -p 8181:80 -p 8125:8125/udp -p 8126:8126 --publish=2003:2003 --name kamon-grafana-dashboard kamon/grafana_graphite

- Add new Datasource type Graphite, http://localhost:81
- Add new Dashboard, call it Hazana
- Add Graph Panel, select datasource 
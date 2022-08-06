# Easycoding chart

helm template easycoding deploy/chart/ --values deploy/chart/values.yaml
helm install -n easycoding -f deploy/chart/values.yaml easycoding deploy/chart

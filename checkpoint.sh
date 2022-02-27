# export WORKER=worker-m-rp7q7-6zkxv
export IP=${WORKER}.subdomain:5747/checkpoint 

envsubst < checkpoint.yaml | kubectl apply -f - # | at now + 6 hours



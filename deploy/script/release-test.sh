#!/bin/bash

# Define an array of shell scripts to run
need_start_server_shell=(
  # rpc
  "user-rpc-test.sh"

  # api
  "user-api-test.sh"
)

# Iterate over each script in the array
for i in "${need_start_server_shell[@]}" ; do
    chmod +x "$i"
    ./"$i"
done

# Check running Docker containers
docker ps

# Execute a command inside the etcd container
docker exec -it etcd etcdctl get --prefix ""


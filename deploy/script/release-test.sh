need_start_server_shell=(
  # rpc
  user-rpc-test.sh
  social-rpc-test.sh

  # api
  user-api-test.sh
  social-api-test.sh
)

for i in ${need_start_server_shell[*]} ; do
    chmod +x $i
    ./$i
done


docker ps

docker exec -it etcd etcdctl get --prefix ""
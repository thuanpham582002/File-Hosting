docker run -d \
  -it \
  --name nginx_file_hosting \
  --mount type=bind,source="$(pwd)",target=/.file,readonly \
  hehe:lastest

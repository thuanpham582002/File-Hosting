docker run \
-itd \
-e PATH_TO_UPLOAD_DIR=/app/uploads \
-v $(pwd)/uploads:/app/uploads \
-p 8080:8080 \
file_host_upload:v2
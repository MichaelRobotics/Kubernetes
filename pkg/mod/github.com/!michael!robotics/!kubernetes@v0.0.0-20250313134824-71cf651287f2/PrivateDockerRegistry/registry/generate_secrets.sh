kubectl create secret tls certs-secret \
  --cert=${BASE_PATH}/certs/tls.crt \
  --key=${BASE_PATH}/certs/tls.key

kubectl create secret generic auth-secret \
  --from-file=${BASE_PATH}/auth/htpasswd

kubectl create secret docker-registry nginx-secret \
  --docker-server=my-registry:30000 \
  --docker-username=myuser \
  --docker-password=mypasswd \
  --docker-email=mail@gmail.com


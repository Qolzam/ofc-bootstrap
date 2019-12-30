#!/bin/bash

cp ./tmp/config/generated-gateway_config.yml ./tmp/telar-social/config/gateway_config.yml
cp ./tmp/config/generated-server_web_config.yml ./tmp/telar-social/config/server_web_config.yml


# Update builder for any ECR secrets needed

cd ./tmp/telar-social

echo "Creating payload-secret in openfaas-fn"

export PAYLOAD_SECRET=$(kubectl get secret -n openfaas payload-secret -o jsonpath='{.data.payload-secret}'| base64 --decode)

kubectl create secret generic payload-secret -n openfaas-fn --from-literal payload-secret="$PAYLOAD_SECRET"

export ADMIN_PASSWORD=$(kubectl get secret -n openfaas basic-auth -o jsonpath='{.data.basic-auth-password}'| base64 --decode)

faas-cli template pull 

kubectl port-forward svc/gateway -n openfaas 31111:8080 &
sleep 2

for i in {1..60};
do
    echo "Checking if OpenFaaS GW is up."

    curl -if 127.0.0.1:31111
    if [ $? == 0 ];
    then
        break
    fi

    sleep 1
done


export OPENFAAS_URL=http://127.0.0.1:31111
echo -n $ADMIN_PASSWORD | faas-cli login --username admin --password-stdin

cp ../generated-stack.yml ./stack.yml

faas-cli deploy

if [ "$GITLAB" = "true" ] ; then
    cp ../generated-gitlab.yml ./gitlab.yml
    echo "Deploying gitlab functions..."
    faas deploy -f ./gitlab.yml
fi

if [ "$ENABLE_AWS_ECR" = "true" ] ; then
    echo "Deploying AWS ECR functions (register-image)..."
    faas deploy -f ./aws.yml
fi

cd ./dashboard
faas-cli template store pull node10-express
faas-cli deploy

sleep 2

# This `ServiceAccount` needs to be patched in place so that the function can perform create / get and update on the SealedSecret CRD:
#kubectl patch -n openfaas-fn deploy import-secrets -p '{"spec":{"template":{"spec":{"serviceAccountName":"sealedsecrets-importer-rw"}}}}'
# This is now applied through an annotation in stack.yml

# Close the kubectl port-forward
kill %1

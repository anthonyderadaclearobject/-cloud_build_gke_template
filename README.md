# Cloud Build Pipeline for GKE

Google Cloud Build executes a series of steps defined in a YAML file in the repository, that way the pipeline is versioned alongside the code. Each step runs a docker container with the content of the repository mounted as a volume in the working directory.

## CI setup:
We will use a Cloud build trigger to build a docker image from our Dockerfile, run tests, and push this image to gcr using the following template:
```yaml
steps:

  # build only the first stage, so we can run tests with it
  - id: build-test-image
    name: gcr.io/cloud-builders/docker
    entrypoint: bash
    args:
      - -c
      - |
        docker image build --target build --tag go-todo:test .
        
  # use initial image to run tests
  - id: run-tests
    name: gcr.io/cloud-builders/docker
    entrypoint: bash
    args:
      - -c
      - |
        docker container run go-todo:test go test
        
  # build app image after tests pass
  - id: build-app
    name: gcr.io/cloud-builders/docker
    entrypoint: bash
    args:
      - -c
      - |
        docker image build --tag gcr.io/${PROJECT_ID}/go-todo:${COMMIT_SHA} .
        
# listing images to push to gcr
images:
  - gcr.io/${PROJECT_ID}/go-todo:${COMMIT_SHA}
```

## CD setup:
We will then use another trigger to upgrade our helm deployment with our new image.
```yaml
steps:

  # Get kubernetes configuration credentials
  - id: get-kube-config
    name: gcr.io/cloud-builders/kubectl
    env:
    - CLOUDSDK_CORE_PROJECT=${_CLOUDSDK_CORE_PROJECT}
    - CLOUDSDK_COMPUTE_ZONE=${_CLOUDSDK_COMPUTE_ZONE}
    - CLOUDSDK_CONTAINER_CLUSTER=${_CLOUDSDK_CONTAINER_CLUSTER}
    - KUBECONFIG=/workspace/.kube/config
    args:
       - cluster-info

  # Update image tags with commit sha and environment
  - id: update-deploy-tag
    name: gcr.io/cloud-builders/gcloud
    args:
      - container
      - images
      - add-tag
      - gcr.io/${PROJECT_ID}/go-todo:${COMMIT_SHA}
      - gcr.io/${PROJECT_ID}/go-todo:${TAG_NAME}

  # Deploy app image using a custon helm image
  - id: deploy
    name: 'gcr.io/${PROJECT_ID}/helm'
    env:
      - KUBECONFIG=/workspace/.kube/config
    args:
      - upgrade
      - --install
      - ${TAG_NAME}-go-todo
      - --namespace=${TAG_NAME}-go-todo
      - --values
      - k8s/go-todo/${TAG_NAME}-values.yaml
      - --set
      - container.image=gcr.io/${PROJECT_ID}/go-todo
      - --set
      - container.tag=${COMMIT_SHA}
      - ./k8s/go-todo
```

> As you can see we use base images available on GCR for docker and kubectl, but not for Helm. We have to build and push the helm image from this repository. To do so run the following commands:

``` bash
git clone https://github.com/GoogleCloudPlatform/cloud-builders-community.git
cd cloud-builders-community/helm
docker build -t gcr.io/<project_name>/helm .
docker push gcr.io/<project_name>/helm
```

# Cloud Build Pipeline for GKE

Google Cloud Build executes a series of steps defined in a YAML file in the repository, that way the pipeline is versioned alongside the code. Each step runs a docker container with the content of the repository mounted as a volume in the working directory. In our case, we will want to build a docker image from our Dockerfile, push this image to a container registry, and upgrade our helm deployment with this new image. Create the following cloudbuild.yml file at the root of your repository:

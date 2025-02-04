# Commands for building the docker image and start a docker container by default.

TAG=$1
DOCKER_FILE_PATH=$2
CURRENT_PATH = $(pwd)

cd $DOCKER_FILE_PATH
# Build the docker image
# t: Assign a tag to the image, but its like name:tag.
docker build -t $TAG .
cd $CURRENT_PATH
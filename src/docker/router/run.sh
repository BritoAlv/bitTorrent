# Commands for building the docker image and start a docker container by default.

TAG=router:1.0

# Build the docker image
# t: Assign a tag to the image.
docker build -t $TAG .

# Start a docker container and open a shell
# rm: Automatically remove the container once it stops (useful for cleanup).
# it: Run the container in interactive mode.
# privileged: Give the container full access to the host system.
# name: Assign a name to the container.
# open sh.
docker run --rm --privileged --name router-container $TAG

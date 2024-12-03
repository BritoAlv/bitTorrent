To run the project, using Docker:

```shell
docker compose up -d
```

That will set up everything and starts the docker containers needed server, client and router.

All of them are using alpine so to access them run :
    
```shell
docker exec -it <containername> sh
# container name : bitclient, bitserver or bitrouter.
```

To stop the containers run:
```shell
docker compose down
```

If done changes in the Dockerfiles, then 

```shell
docker compose down --rmi all
```

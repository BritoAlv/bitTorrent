This project has three types of network communication:

- client with client: after two clients agreed on sharing a common .torrent they establish a connection to share pieces of the file.
- client with tracker: when a client wants to find other clients related to a specific .torrent it should ask a tracker for that specific information.
- tracker with tracker: trackers will communicate between each other, to find the information asked by a client related to a specific .torrent.

Given those constraints there could be as many networks in the project as long as they are intercommunicated between each other. 

For this purpose a router will be used. It will be responsible for forward packages from a network to the other.

Based on that is needed the following on the project:

- have a router ready to intercommunicate new networks.
- create a network and tell the router about its existence.
- create a client in a specific network.
- create a tracker in a specific network.





To run the project via Docker:

To set up everything and starts the docker containers needed: server, client and router.
```shell
docker compose up -d
```

All of them use alpine, so to access them run:
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
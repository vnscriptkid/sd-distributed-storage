up:
	docker-compose up -d

down:
	docker-compose down --remove-orphans --volumes

cli:
	docker exec -it cassandra1 cqlsh

status1:
	docker exec -it cassandra1 nodetool status

status2:
	docker exec -it cassandra2 nodetool status

status3:
	docker exec -it cassandra3 nodetool status

ring1:
	docker exec -it cassandra1 nodetool ring

ring2:
	docker exec -it cassandra2 nodetool ring

ring3:
	docker exec -it cassandra3 nodetool ring

log1:
	docker logs -f cassandra1

log2:
	docker logs -f cassandra2

log3:
	docker logs -f cassandra3
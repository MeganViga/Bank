postgres:
	    docker run --name postgres16  -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres:16-alpine
createdb:
	docker exec postgres16 createdb --username=root --owner=root bank
dropdb:
	docker exec postgres16 dropdb bank
imu:
	 migrate -database 'postgres://root:secret@localhost:5432/bank?sslmode=disable' -path db/migrations -verbose up
imd:
	 migrate -database 'postgres://root:secret@localhost:5432/bank?sslmode=disable' -path db/migrations -verbose down
.PHONY: createdb dropdb postgres imu imd
## ABOUT

The project is message queue implemented over tcp.

### EXTRA POINTS

- The project is message queue where allow you selected the storage driver between sqlite and in memory.

## WARNING

- The project is only to study and learning purposes. PS: no use in production.


# TECH

- Golang
- SQLite
- In Memory


# CONCEPTS APPLIED

- TCP
- Linked List
- Go(Goroutines and Mutex)
- SQLite(to avoid lost messages)


## HOW TO RUN THE SERVER

- To run the server you need to have the go installed in your machine.
```bash
go run cmd/server/main.go -storage=sqlite
```
OR
```bash
go run cmd/server/main.go -storage=inmemory
```





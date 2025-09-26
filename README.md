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

- Before start the server you need to create the .env file where contains:
```txt
USER_ADMIN=user_here
PASSWORD=password_here
```
- The folder **examples** has examples how to use publish and consume the messages from the **tqueue**

- To run the server you need to have the go installed in your machine.
```bash
go run cmd/server/main.go -storage=sqlite // The queue is storage on sqlite, so if you restart the server you don't lose the messages
```
OR
```bash
go run cmd/server/main.go -storage=inmemory // The queue is storage in memory, so if you restart the server you will lose the messages.
```





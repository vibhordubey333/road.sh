Link: https://roadmap.sh/projects/broadcast-server

#### You are required to create a simple broadcast server that will allow clients to connect to it, send messages that will be broadcasted to all connected clients.

**Goal**

The goal of this project is to help you understand how to work with websockets and implement real-time communication between clients and servers. This will help you understand how the real-time features of applications like chat applications, live scoreboards, etc., work.

**Requirements**

You are required to build a CLI based application that can be used to either start the server or connect to the server as a client. Here are the sample commands that you can use:

- broadcast-server start - This command will start the server. <br/>
- broadcast-server connect - This command will connect the client to the server. <br/>
- When the server is started using the broadcast-server start command, it should listen for client connections on a specified port (you can configure that using command options or hardcode for simplicity). When a client connects and sends a message, the server should broadcast this message to all connected clients. <br/>

The server should be able to handle multiple clients connecting and disconnecting gracefully.

**Implementation** <br>
You can use any programming language to implement this project. Here are some of the steps that you can follow to implement this project:

1. Create a server that listens for incoming connections.
2. When a client connects, store the connection in a list of connected clients.
3. When a client sends a message, broadcast this message to all connected clients.
4. Handle client disconnections and remove the client from the list of connected clients.
5. Implement a client that can connect to the server and send messages.
6. Test the server by connecting multiple clients and sending messages.
7. Implement error handling and graceful shutdown of the server.
8. This project will help you understand how to work with websockets and implement real-time communication between clients and servers. You can extend this project by adding features like authentication, message history, etc.


# Run the server:

1. Start the server: `go run cmd/main.go start --port=9090`
2. Connect to server with client 1: `go run cmd/main.go connect --server=localhost:9090 --name vibhor`
3. Connect to server with client 2: `go run cmd/main.go connect --server=localhost:9090 --name vibhor-1`

# Load Testing Server:

1. Start the server: `go run cmd/main.go start --port=9090`
2. Navigate to `load-testing` directory.
3. Execute `k6 run broadcast-load-testing.js`
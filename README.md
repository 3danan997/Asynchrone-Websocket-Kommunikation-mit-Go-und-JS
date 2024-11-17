# **Asynchronous WebSocket Communication**

This project is part of the "Konzepte moderner Programmiersprachen" course at FH Aachen, led by Prof. Dr. Andreas Claßen. The goal is to implement real-time communication using WebSockets for a simulated trouble ticket system.

### Overview
**Language:** Go (Backend Server), JavaScript (Clients)

**Functionality:** Real-time communication for ticket assignment using WebSockets.

**Use Case:** Simulated trouble ticket system where clients can self-assign tickets.

### Key Features
**Server:** Written in Go, handles ticket creation and assignment, communicates with clients via WebSockets.

**Clients:** Written in JavaScript, can run as command-line clients (Node.js) or web clients (browser).

**Concurrency:** Asynchronous handling of WebSocket connections to manage real-time updates.

### Requirements
**Server:** Go command-line program, capable of adding new tickets via keyboard input.

**Clients:** JavaScript code, either as Node.jscommand-line clients or simple web clients with console.logoutputs.

### Usage
**Server:** Run the Go server to start accepting WebSocket connections.

**Clients:** Connect multiple clients to the server, each client can self-assign tickets and receive real-time updates.

### Notes
**No Authentication:** Simplified for educational purposes, no complex login required.

**Error Handling:** Focus on positive flow, basic error handling for client inputs.

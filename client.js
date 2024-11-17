const WebSocket = require("ws");
const readline = require("readline");
const { log } = require("console");

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
});

const ws = new WebSocket("ws://localhost:8080/ws");
let clientIdInput;

ws.on("open", async function () {
  log("\nClient gestartet");
  rl.question("Bitte Client-ID eingeben: ", (clientIdInput) => {
    log(`Client-ID: ${clientIdInput}`);
    ws.send(
      JSON.stringify({ message_type: "client_id", client_name: clientIdInput })
    );
    // ws.send(JSON.stringify({ message_type: "client_ready" }));
  });
});

ws.on("message", function (data) {
  //log("hii");
  const message = JSON.parse(data);
  switch (message.message_type) {
    case "ticket_info":
      if (message.message_type === "ticket_info" && message.tickets !== null) {
        const tickets = message.tickets;
        log("Tickets:");
        tickets.forEach(ticket => {
          log(`id: ${ticket.id}, assigned_to: '${ticket.assigned_to}'`);
        });
      } else {
        log("keine Tickets.");
      }
      log("n: Selbstzuweisung q: quit\n");
      break;

    case "assign_error":
      log(message.error);
      rl.question(
        "Bitte geben Sie andere Ticketsnummer ein!: ",
        (ticketId) => {
          ws.send(
            JSON.stringify({
              message_type: "assign_ticket",
              ticket_id: parseInt(ticketId, 10),
              client_id: clientIdInput,
            })
          );
        }
      );
      rl.close
      break;
  }
});

ws.on("close", function () {
  log("Server wurde geschlossen.");
  rl.close();
});

ws.on("error", function (error) {
  console.error("Server Fehler:", error);
  rl.close();
});

function writeTicketNumber() {
  rl.question("Bitte geben Sie die gewünschte Ticketsnummer ein! : ", (ticketId) => {
    if (isNaN(ticketId)) {
      log("Ungültige Eingabe.");
      writeTicketNummer();
    } else {
      log(ticketId);
      ws.send(
        JSON.stringify({message_type: "assign_ticket", ticket_id: parseInt(ticketId, 10) })
      );
    }
  });
}

rl.on("line", (input) => {
  switch (input) {
    case "n":
      writeTicketNumber();
      break;
    case "q":
      rl.close();
      process.exit(0);
  }
});

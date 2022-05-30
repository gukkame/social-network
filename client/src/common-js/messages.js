export { connectToWS, closeWsConnection, ws }

let ws = null
const wsAddress = "ws://localhost:8080/WSconnect"


function connectToWS() {
    ws = new WebSocket(wsAddress);
}

function closeWsConnection() {
    ws.close()
}


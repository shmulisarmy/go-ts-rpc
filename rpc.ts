const ws = new WebSocket("ws://localhost:8080/ws");

    ws.onopen = () => {
      console.log("Connected to WebSocket server");
      ws.send("Hello, server!");
    };

    ws.onmessage = (event) => {
      console.log("Message from server:", event.data);
    };

    ws.onclose = () => {
      console.log("Disconnected from WebSocket server");
    };
    
let id_counter = 0

function rpc_call(functionName, ...args) {
    return new Promise((resolve, reject) => {
        const current_call_id = id_counter++;
        ws.send(JSON.stringify({ "type": "rpc-call", function: functionName, args: args, "id": current_call_id }));
        ws.onmessage = (event) => {
            if (JSON.parse(event.data).id == current_call_id) {
                resolve(JSON.parse(event.data));
            }
        };
    });
}
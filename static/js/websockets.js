// Пошла работа с websockets
function websocket_connect(){
    ws = new WebSocket('ws://127.0.0.1:8081/connect');

    ws.onopen(websocket_open);
    ws.onclose(websocket_close);
    ws.onmessage(websocket_message);
}

function websocket_open(event){
    console.log('websocket is open!');
}

function websocket_close(event){
    console.log('websocket is close!');
}

function websocket_message(event){
    var data = JSON.parse(event.data)
    handlers[data['handler_name']](data["data"])
}

// ------------- incoming ------------------
function set_grid(data){
    grid_coordinats = data['grid']
}

// ------------ outcoming ------------------
function select_chanck(){
    
}
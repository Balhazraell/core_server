'use strict';

var chunck = require('./chunck');

// Набор функций получаемых от сервера
var handlers = {
    'set_grid': set_grid
};

// Пошла работа с websockets
function connect(){
    // Здесь происходит падение если не можем подключится, надо красиво обработать...
    var ws = new WebSocket('ws://127.0.0.1:8081/appgame');
    ws.onopen = open;
    ws.onclose = close;
    ws.onmessage = message;
}

// websocket стартанул.
function open(event){
    console.log('websocket is open!');
    console.log(event);
}

// websocket закрылся.
function close(event){
    console.log('websocket is close!');
    console.log(event);
}

// пришло сообщение по websocket.
function message(event){
    var data = JSON.parse(event.data);
    handlers[data['handler_name']](data['data']);
}

// ------------- incoming ------------------
// Пришла сетка.
function set_grid(data){
    new_grid = data['grid'];
    // По сути очищаем список.
    grid_coordinats = [];
    for (let i = 0; i < new_grid.length; i++){
        let chunck = new chunck.Chunck(
            // первым аргументом потом будет id.
            new_grid[i]
        );
        grid_coordinats.append(chunck);
    } 
}

exports.connect = connect;
var canvas;
var ctx;

// Заглушки.
var grid_coordinats = [
    [[0,0],     [0, 100],   [100, 100], [100, 0]],
    [[0,100],   [0, 200],   [100, 200], [100, 100]],
    [[0,200],   [0, 300],   [100, 300], [100, 200]],

    [[100,0],   [100, 100], [200, 100], [200, 0]],
    [[100,100], [100, 200], [200, 200], [200, 100]],
    [[100,200], [100, 300], [200, 300], [200, 200]],

    [[200,0],   [200, 100], [300, 100], [300, 0]],
    [[200,100], [200, 200], [300, 200], [300, 100]],
    [[200,200], [200, 300], [300, 300], [300, 200]],
]

// Набор функций получаемых от сервера
var handlers = {
    'set_grid': set_grid
};

function app_start(){
    canvas = document.getElementById('myCanvas');
    ctx = canvas.getContext('2d');
    ctx.strokeStyle = '#000000';

    websocket_connect();
    draw_grid();
}

function draw_grid(){
    
    for (var i = 0; i < grid_coordinats.length; i++){
        ctx.beginPath();
        ctx.moveTo(grid_coordinats[i][0][0], grid_coordinats[i][0][1]);
        for (var pos_index = 1; pos_index < grid_coordinats[i].length; pos_index++){
            ctx.lineTo(grid_coordinats[i][pos_index][0], grid_coordinats[i][pos_index][1])
        }
        ctx.closePath();
        ctx.stroke();
    }

    requestAnimationFrame(draw_grid);
}

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

function set_grid(data){
    grid_coordinats = data['grid']
}
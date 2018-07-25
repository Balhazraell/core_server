'use strict';
var chunck = require('./chunck');
var websocket = require('./websockets');

var canvas;
var ctx;

// Заглушки.
var grid_coordinats = [
    new chunck.Chunck([[0,0],     [0, 100],   [100, 100], [100, 0]]),
    new chunck.Chunck([[0,100],   [0, 200],   [100, 200], [100, 100]]),
    new chunck.Chunck([[0,200],   [0, 300],   [100, 300], [100, 200]]),

    new chunck.Chunck([[100,0],   [100, 100], [200, 100], [200, 0]]),
    new chunck.Chunck([[100,100], [100, 200], [200, 200], [200, 100]]),
    new chunck.Chunck([[100,200], [100, 300], [200, 300], [200, 200]]),

    new chunck.Chunck([[200,0],   [200, 100], [300, 100], [300, 0]]),
    new chunck.Chunck([[200,100], [200, 200], [300, 200], [300, 100]]),
    new chunck.Chunck([[200,200], [200, 300], [300, 300], [300, 200]]),
]

// Запуск приложения после загрузки html страницы.
if (document.readyState != 'loading'){
    app_start();
} else {
    document.addEventListener('DOMContentLoaded', app_start);
}
//-------------------------------------------------------------//

function app_start(){
    canvas = document.getElementById('myCanvas');
    ctx = canvas.getContext('2d');
    websocket.connect();
    draw_grid();
}

function draw_grid(){
    for (let i = 0; i < grid_coordinats.length; i++){
        ctx.strokeStyle = grid_coordinats.color;
        let chunck_points = grid_coordinats[i].draw_poins
        ctx.beginPath();
        ctx.moveTo(chunck_points[0][0], chunck_points[0][1]);
        for (let pos_index = 1; pos_index < chunck_points.length; pos_index++){
            ctx.lineTo(chunck_points[pos_index][0], chunck_points[pos_index][1])
        }
        ctx.closePath();
        ctx.stroke();
    }

    requestAnimationFrame(draw_grid);
}

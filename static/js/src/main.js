'use strict';
var chunck = require('./chunck');
var websocket = require('./websockets');
var mouse_manager = require('./mouse_manager');

var canvas;
var ctx;
var MousManager;
var selectChunck;

// Заглушки.
var grid_coordinats = [
    // new chunck.Chunck(0, [[0, 0],     [100, 0],   [100, 100], [0, 100]]),
    // new chunck.Chunck(0, [[100, 0],   [200, 0],   [200, 100], [100, 100]]),
    // new chunck.Chunck(0, [[200, 0],   [300, 0],   [300, 100], [200, 100]]),

    // new chunck.Chunck(0, [[0,100],   [100, 100], [100, 200], [0, 200]]),
    // new chunck.Chunck(0, [[100,100], [200, 100], [200, 200], [200, 100]]),
    // new chunck.Chunck(0, [[200,100], [300, 100], [300, 200], [200, 200]]),

    // new chunck.Chunck(0, [[0,200],   [100, 200], [100, 300], [0, 300]]),
    // new chunck.Chunck(0, [[100,200], [200, 200], [200, 300], [100, 300]]),
    // new chunck.Chunck(0, [[200,200], [300, 200], [300, 300], [200, 300]]),
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
    MousManager = new mouse_manager.MouseManager(canvas);
    game_loop();
}

function game_loop(){
    set_chuncks_color();
    draw_grid();

    requestAnimationFrame(game_loop);
}

function draw_grid(){
    for (let i = 0; i < grid_coordinats.length; i++){
        ctx.strokeStyle = grid_coordinats[i].color;
        let chunck_points = grid_coordinats[i].draw_poins
        ctx.beginPath();
        ctx.moveTo(chunck_points[0][0], chunck_points[0][1]);
        for (let pos_index = 1; pos_index < chunck_points.length; pos_index++){
            ctx.lineTo(chunck_points[pos_index][0], chunck_points[pos_index][1])
        }
        ctx.closePath();
        ctx.stroke();
    }

    
}

function set_grid(new_map) {
    grid_coordinats = [];
    for (let i = 0; i < new_map.length; i++){
        let newChunck = new chunck.Chunck(
            new_map[i].state,
            new_map[i].coordinates
        );

        grid_coordinats.push(newChunck);
    } 
}

function set_chuncks_color(){
    // TODO: интересно, что правильнее, собрать информацию о том, с каким чанком мы пересечены и после
    // пробегатся - задавая цвет, или все делать в одном цикле???
    for (let i = 0; i < grid_coordinats.length; i++){
        let is_collision = check_collision(grid_coordinats[i])
        selectChunck = grid_coordinats[i]

        if (is_collision){
            if (grid_coordinats[i].state == 0){
                grid_coordinats[i].color = grid_coordinats[i].positive_color
            } else{
                grid_coordinats[i].color = grid_coordinats[i].negative_color
            }
        } else {
            grid_coordinats[i].color = grid_coordinats[i].normal_color
        }
    }
} 

// Данный метод - некоторого рода защита на клиенте -
// мы не должны отправлять запрос, если клетка занята.
function check_collision(cunck_for_check){
    var result = false;
    var resultList = [];

    for (let i = 0; i < cunck_for_check.draw_poins.length - 1; i++){ 
        let pos1X = cunck_for_check.draw_poins[i][0];
        let pos1Y = cunck_for_check.draw_poins[i][1];

        let pos2X = cunck_for_check.draw_poins[i + 1][0];
        let pos2Y = cunck_for_check.draw_poins[i + 1][1];

        resultList.push(
            (pos2X - pos1X) * (MousManager.Y - canvas.offsetTop - pos1Y) - (pos2Y - pos1Y) * (MousManager.X - canvas.offsetLeft - pos1X)
        )
    }

    if (resultList.length > 0){
        result = resultList.every(function (element){
            return element >= 0;
        })
    }    

    return result
}

function set_chunck_state(){
    if(selectChunck != undefined){
        websocket.set_chunck_state(selectChunck.id)
    }
}

exports.set_grid = set_grid;
exports.set_chunck_state = set_chunck_state;
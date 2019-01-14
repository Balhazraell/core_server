'use strict';

var main = require('./main')

class MouseManager {
    // не понимаю почему импорт канваса не дает результата...
    constructor(canvas){
        var X = 0;
        var Y = 0;

        canvas.addEventListener('mousemove', this.mouseMove.bind(this));
        canvas.addEventListener('click', this.mouseClick.bind(this));
    }

    mouseMove(event){
        this.X = event.pageX;
        this.Y = event.pageY;
    }

    mouseClick(event){
        main.set_chunck_state();
    }
}



exports.MouseManager = MouseManager;
const Emitter = require('events').EventEmitter;

const emitter = new Emitter();

emitter.on('mul', (data) => {
    let res = data.x * data.y;
    console.log('on 1: ', res);
});

emitter.on('mul', (data) => {
    let res = data.x * data.y;
    console.log('on 2: ', res);
});

emitter.emit('mul', {x: 10, y: 5});
emitter.emit('mul', {x: 1000, y: 5});
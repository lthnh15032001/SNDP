var net = require('net');
console.log('Listening')
var server = net.createServer((socket) => {
    socket.write('Echo server\r\n');

    // socket is defined here and in scope
    socket.on('data', function (data) {

        console.log('data', Buffer.from(data).toString('hex'))

        socket.pipe(socket);
    }) 
}) 

server.listen(8081, 'localhost');

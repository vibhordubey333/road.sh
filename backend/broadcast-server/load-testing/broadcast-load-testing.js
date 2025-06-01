import { check } from 'k6';
import ws from 'k6/ws';

export default function() {
    const url = 'ws://localhost:9090/ws?username=loadtest-' + __VU;

    const res = ws.connect(url, function(socket) {
        socket.on('open', () => {
            console.log('Connected');

            // Send a message every second for 10 seconds
            const interval = setInterval(() => {
                socket.send('Message from VU ' + __VU);
                socket.send(randomString(1000)); // 1KB message
                socket.send(randomString(10000)); // 10KB message
                sleep(5);
                socket.close();
            }, 1000);

            // Close after 10 seconds
            setTimeout(() => {
                clearInterval(interval);
                socket.close();
            }, 10000);
        });

        socket.on('message', (data) => {
            // Parse the message to check format
            try {
                const parsed = JSON.parse(data);
                check(parsed, {
                    'has sender field': (obj) => obj.sender !== undefined,
                    'has content field': (obj) => obj.content !== undefined,
                    'has timestamp field': (obj) => obj.timestamp !== undefined
                });
            } catch (e) {
                console.error('Failed to parse message:', e);
            }
        });

        socket.on('close', () => console.log('disconnected'));

        socket.on('error', (e) => {
            console.error('Error:', e);
        });
    });

    check(res, { 'status is 101': (r) => r && r.status === 101 });
}

// Test configuration
export const options = {
    // Start with a small number of VUs and gradually increase
    stages: [
        { duration: '30s', target: 1000000 },   // Ramp up to 1000000 users over 30 seconds
        { duration: '1m', target: 50 },    // Ramp up to 50 users over 1 minute
        { duration: '30s', target: 100 },  // Ramp up to 100 users over 30 seconds
        { duration: '1m', target: 100 },   // Stay at 100 users for 1 minute
        { duration: '30s', target: 0 },    // Ramp down to 0 users over 30 seconds
    ],
    thresholds: {
        // Define thresholds for acceptable performance
        'checks': ['rate>0.95'], // 95% of checks should pass
    }
};
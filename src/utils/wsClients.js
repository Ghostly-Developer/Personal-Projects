const WS_URL = 'ws://localhost:8080/games';

class WSClient {
    constructor(url = WS_URL) {
        this.url = url;
        this.ws = null;
        this.listeners = new Map();
        this.requestId = 0;
    }

    connect() {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            console.log('WebSocket already connected');
            return;
        }

        this.ws = new WebSocket(this.url);

        this.ws.onopen = () => {
            console.log('WebSocket connected');
        };

        this.ws.onmessage = (event) => {
            try {
                const response = JSON.parse(event.data);
                console.log('Received:', response);
                
                // Notify all registered listeners
                this.listeners.forEach((callback) => {
                    callback(response);
                });
            } catch (error) {
                console.error('Failed to parse message:', error);
            }
        };

        this.ws.onclose = () => {
            console.log('WebSocket disconnected');
        };

        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
    }

    // Send a request to the server
    // Flow types: 'join', 'game', 'chat'
    // Action types: 'createRoom', 'joinRoom', 'move', 'message'
    send(flow, action, payload, playerId) {
        if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
            console.error('WebSocket is not connected');
            return;
        }

        const request = {
            PlayerId: playerId,
            Flow: flow,
            Action: action,
            Payload: payload
        };

        try {
            this.ws.send(JSON.stringify(request));
        } catch (error) {
            console.error('Failed to send message:', error);
        }
    }

    // Join or create a game room
    joinRoom(roomId, playerName, game, playerId) {
        const payload = {
            RoomId: roomId,
            PlayerName: playerName,
            Game: game
        };
        this.send('join', 'joinRoom', payload, playerId);
    }

    // Create a new game room
    createRoom(roomId, playerName, game, playerId) {
        const payload = {
            RoomId: roomId,
            PlayerName: playerName,
            Game: game
        };
        this.send('join', 'createRoom', payload, playerId);
    }

    // Send a chat message
    sendMessage(roomId, playerName, message, playerId) {
        const payload = {
            RoomId: roomId,
            PlayerName: playerName,
            Message: message
        };
        this.send('chat', 'message', payload, playerId);
    }

    // Add a listener for incoming messages
    addListener(id, callback) {
        this.listeners.set(id, callback);
    }

    // Remove a listener
    removeListener(id) {
        this.listeners.delete(id);
    }

    // Close the WebSocket connection
    disconnect() {
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
        this.listeners.clear();
    }
}

export default WSClient;

/* Example usage:
const ws = new WSClient();
ws.connect();

// Add a listener for all messages
ws.addListener('main', (response) => {
    console.log('Response:', response);
    // response: { Status: string, Message: string, Payload: any }
});

// Create a room
ws.createRoom('room1', 'Player1', 'tictactoe', 'player1_id');

// Join a room
ws.joinRoom('room1', 'Player2', 'tictactoe', 'player2_id');

// Send a chat message
ws.sendMessage('room1', 'Player1', 'Hello!', 'player1_id');

// Remove listener when done
ws.removeListener('main');

// Disconnect when done
ws.disconnect();
*/

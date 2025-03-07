class ChatRoom {
    constructor() {
        this.ws = new WebSocket('ws://localhost:8080/chats');
        this.messageList = document.getElementById('messageList');
        this.messageInput = document.getElementById('messageInput');
        this.sendButton = document.getElementById('sendButton');

        this.setupWebSocket();
        this.setupEventListeners();

        // Auto focus the input
        this.messageInput.focus();
        this.setInputsEnabled(false); // Start disabled until connection is established
    }

    setupWebSocket() {
        this.ws.onopen = () => {
            console.log('Connected to WebSocket server');
            this.setInputsEnabled(true);

            this.ws.send(JSON.stringify({
                type: 'join',
                userId: 'User'
            }));
        };

        this.ws.onmessage = (event) => {
            console.log("onmessage", event);
            const message = JSON.parse(event.data);
            this.displayMessage(message, 'received');
        };

        this.ws.onclose = () => {
            console.warn('Disconnected from WebSocket server');
            this.setInputsEnabled(false);
        };
    }

    setupEventListeners() {
        this.sendButton.addEventListener('click', () => this.sendMessage());

        this.messageInput.addEventListener('keypress', (event) => {
            if (event.key === 'Enter') {
                this.sendMessage();
            }
        });
    }

    sendMessage() {
        const messageText = this.messageInput.value.trim();
        if (messageText) {
            const message = {
                type: 'message',
                userId: 'User', // This will be replaced with actual username later
                content: messageText,
                timestamp: new Date().toISOString()
            };

            this.ws.send(JSON.stringify(message));
            this.displayMessage(message, 'sent');
            this.messageInput.value = '';
        }
    }

    displayMessage(message, type) {
        const container = document.createElement('div');
        container.classList.add('message-container');

        // Create and format timestamp
        const timestamp = new Date(message.timestamp);
        const timeString = timestamp.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });

        const timestampElement = document.createElement('div');
        timestampElement.classList.add('message-timestamp');
        timestampElement.textContent = timeString;

        const messageElement = document.createElement('div');
        messageElement.classList.add('message', type);

        const messageContent = `
            <span class="username">${message.userId}</span>
            <span class="text">${message.content}</span>
        `;

        messageElement.innerHTML = messageContent;

        container.appendChild(timestampElement);
        container.appendChild(messageElement);
        this.messageList.appendChild(container);
        this.messageList.scrollTop = this.messageList.scrollHeight;
    }

    setInputsEnabled(enabled) {
        this.messageInput.disabled = !enabled;
        this.sendButton.disabled = !enabled;

        if (enabled) {
            this.messageInput.placeholder = "Type your message...";
            this.messageInput.focus();
        } else {
            this.messageInput.placeholder = "Disconnected from chat server...";
        }
    }
}

// Initialize the chat room when the page loads
document.addEventListener('DOMContentLoaded', () => {
    new ChatRoom();
}); 
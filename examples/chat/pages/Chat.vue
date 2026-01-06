<template>
  <div class="chat-container">
    <h1>Real-time Chat</h1>
    
    <div class="username-input" v-if="!username">
      <input
        v-model="tempUsername"
        placeholder="Enter your name"
        @keyup.enter="setUsername"
      />
      <button @click="setUsername">Join Chat</button>
    </div>

    <div v-else class="chat-room">
      <div class="messages">
        <div
          v-for="msg in allMessages"
          :key="msg.id"
          :class="['message', msg.user === username ? 'own-message' : '']"
        >
          <strong>{{ msg.user }}:</strong>
          <span>{{ msg.text }}</span>
          <small>{{ formatTime(msg.timestamp) }}</small>
        </div>
      </div>

      <form @submit.prevent="sendMessage" class="message-form">
        <input
          v-model="newMessage"
          placeholder="Type a message..."
          :disabled="!connected"
        />
        <button type="submit" :disabled="!connected || !newMessage.trim()">
          Send
        </button>
      </form>

      <div class="connection-status">
        <span :class="connected ? 'connected' : 'disconnected'">
          {{ connected ? '🟢 Connected' : '🔴 Disconnected' }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useLiveUpdate } from '@toutaio/inertia-vue';

interface Message {
  id: number;
  user: string;
  text: string;
  timestamp: string;
}

interface Props {
  messages: Message[];
}

const props = defineProps<Props>();

const username = ref('');
const tempUsername = ref('');
const newMessage = ref('');
const liveMessages = ref<Message[]>([]);

const wsUrl = `ws://${window.location.host}/ws`;
const { connected, on } = useLiveUpdate(wsUrl);

// Listen for new messages
on('chat', (message: Message) => {
  liveMessages.value.push(message);
});

const allMessages = computed(() => {
  return [...props.messages, ...liveMessages.value];
});

function setUsername() {
  if (tempUsername.value.trim()) {
    username.value = tempUsername.value.trim();
  }
}

async function sendMessage() {
  if (!newMessage.value.trim()) return;

  try {
    await fetch('/messages', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        user: username.value,
        text: newMessage.value,
      }),
    });

    newMessage.value = '';
  } catch (error) {
    console.error('Failed to send message:', error);
  }
}

function formatTime(timestamp: string) {
  const date = new Date(timestamp);
  return date.toLocaleTimeString();
}
</script>

<style scoped>
.chat-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.username-input {
  display: flex;
  gap: 10px;
  justify-content: center;
  margin-top: 100px;
}

.username-input input {
  padding: 10px;
  font-size: 16px;
  border: 1px solid #ddd;
  border-radius: 4px;
  width: 300px;
}

.username-input button {
  padding: 10px 20px;
  font-size: 16px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.chat-room {
  display: flex;
  flex-direction: column;
  height: 600px;
}

.messages {
  flex: 1;
  overflow-y: auto;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 20px;
  margin-bottom: 10px;
  background: #f9f9f9;
}

.message {
  margin-bottom: 15px;
  padding: 10px;
  background: white;
  border-radius: 4px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.message.own-message {
  background: #e3f2fd;
  margin-left: 20%;
}

.message strong {
  color: #333;
  margin-right: 8px;
}

.message small {
  display: block;
  color: #999;
  font-size: 12px;
  margin-top: 5px;
}

.message-form {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
}

.message-form input {
  flex: 1;
  padding: 12px;
  font-size: 16px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.message-form button {
  padding: 12px 24px;
  font-size: 16px;
  background: #28a745;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.message-form button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.connection-status {
  text-align: center;
  padding: 10px;
}

.connected {
  color: #28a745;
}

.disconnected {
  color: #dc3545;
}
</style>

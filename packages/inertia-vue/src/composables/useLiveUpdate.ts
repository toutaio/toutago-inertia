import { ref, onUnmounted, Ref, getCurrentInstance } from 'vue';

export interface LiveUpdateOptions {
  reconnect?: boolean;
  reconnectDelay?: number;
  maxReconnectAttempts?: number;
}

export interface Message {
  channel: string;
  type: string;
  data: any;
}

type MessageHandler = (data: any) => void;

export function useLiveUpdate(url: string, options: LiveUpdateOptions = {}) {
  const {
    reconnect = true,
    reconnectDelay = 1000,
    maxReconnectAttempts = 10,
  } = options;

  const ws: Ref<WebSocket | null> = ref(null);
  const connected = ref(false);
  const error: Ref<Error | null> = ref(null);
  const reconnectAttempts = ref(0);

  const handlers = new Map<string, Set<MessageHandler>>();
  let reconnectTimeout: ReturnType<typeof setTimeout> | null = null;

  function connect() {
    try {
      ws.value = new WebSocket(url);

      ws.value.onopen = () => {
        connected.value = true;
        error.value = null;
        reconnectAttempts.value = 0;
      };

      ws.value.onmessage = (event: MessageEvent) => {
        try {
          const message: Message = JSON.parse(event.data);
          const channelHandlers = handlers.get(message.channel);
          
          if (channelHandlers) {
            channelHandlers.forEach((handler) => {
              handler(message.data);
            });
          }

          // Also trigger wildcard handlers
          const wildcardHandlers = handlers.get('*');
          if (wildcardHandlers && message.channel !== '*') {
            wildcardHandlers.forEach((handler) => {
              handler(message.data);
            });
          }
        } catch (err) {
          error.value = err as Error;
        }
      };

      ws.value.onerror = (event: Event) => {
        error.value = new Error('WebSocket error occurred');
        connected.value = false;
      };

      ws.value.onclose = () => {
        connected.value = false;

        if (reconnect && reconnectAttempts.value < maxReconnectAttempts) {
          reconnectTimeout = setTimeout(() => {
            reconnectAttempts.value++;
            connect();
          }, reconnectDelay);
        }
      };
    } catch (err) {
      error.value = err as Error;
      connected.value = false;
    }
  }

  function subscribe(channel: string) {
    if (ws.value && ws.value.readyState === WebSocket.OPEN) {
      ws.value.send(
        JSON.stringify({
          type: 'subscribe',
          channel,
        })
      );
    }
  }

  function unsubscribe(channel: string) {
    if (ws.value && ws.value.readyState === WebSocket.OPEN) {
      ws.value.send(
        JSON.stringify({
          type: 'unsubscribe',
          channel,
        })
      );
    }
  }

  function on(channel: string, handler: MessageHandler) {
    if (!handlers.has(channel)) {
      handlers.set(channel, new Set());
    }
    handlers.get(channel)!.add(handler);

    // Auto-subscribe when adding first handler for a channel
    if (handlers.get(channel)!.size === 1 && channel !== '*') {
      subscribe(channel);
    }
  }

  function off(channel: string, handler: MessageHandler) {
    const channelHandlers = handlers.get(channel);
    if (channelHandlers) {
      channelHandlers.delete(handler);
      
      // Auto-unsubscribe when removing last handler for a channel
      if (channelHandlers.size === 0 && channel !== '*') {
        unsubscribe(channel);
        handlers.delete(channel);
      }
    }
  }

  function disconnect() {
    if (reconnectTimeout) {
      clearTimeout(reconnectTimeout);
      reconnectTimeout = null;
    }

    if (ws.value) {
      ws.value.close();
      ws.value = null;
    }

    connected.value = false;
  }

  // Initialize connection
  connect();

  // Cleanup on unmount (only in component context)
  if (getCurrentInstance()) {
    onUnmounted(() => {
      disconnect();
    });
  }

  return {
    connected,
    error,
    subscribe,
    unsubscribe,
    on,
    off,
    disconnect,
  };
}

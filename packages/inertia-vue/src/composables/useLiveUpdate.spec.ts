import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { ref, nextTick } from 'vue';
import { useLiveUpdate } from './useLiveUpdate';

// Mock WebSocket
class MockWebSocket {
  public onopen: ((event: Event) => void) | null = null;
  public onmessage: ((event: MessageEvent) => void) | null = null;
  public onerror: ((event: Event) => void) | null = null;
  public onclose: ((event: CloseEvent) => void) | null = null;
  public readyState = WebSocket.CONNECTING;
  public sent: string[] = [];

  constructor(public url: string) {
    // Simulate connection opening
    setTimeout(() => {
      this.readyState = WebSocket.OPEN;
      if (this.onopen) {
        this.onopen(new Event('open'));
      }
    }, 0);
  }

  send(data: string) {
    this.sent.push(data);
  }

  close() {
    this.readyState = WebSocket.CLOSED;
    if (this.onclose) {
      this.onclose(new CloseEvent('close'));
    }
  }

  // Helper to simulate receiving a message
  simulateMessage(data: any) {
    if (this.onmessage) {
      this.onmessage(new MessageEvent('message', { data: JSON.stringify(data) }));
    }
  }
}

describe('useLiveUpdate', () => {
  let mockWebSocket: MockWebSocket;
  let websocketInstances: MockWebSocket[] = [];

  beforeEach(() => {
    websocketInstances = [];
    // Mock global WebSocket
    global.WebSocket = vi.fn((url: string) => {
      mockWebSocket = new MockWebSocket(url);
      websocketInstances.push(mockWebSocket);
      return mockWebSocket;
    }) as any;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('should establish WebSocket connection', async () => {
    const { connected } = useLiveUpdate('ws://localhost:8080/ws');

    await nextTick();
    await new Promise((resolve) => setTimeout(resolve, 10));

    expect(connected.value).toBe(true);
  });

  it('should subscribe to channels', async () => {
    const { subscribe } = useLiveUpdate('ws://localhost:8080/ws');

    await nextTick();
    await new Promise((resolve) => setTimeout(resolve, 10));

    subscribe('test-channel');

    await nextTick();

    expect(mockWebSocket.sent).toContainEqual(
      JSON.stringify({
        type: 'subscribe',
        channel: 'test-channel',
      })
    );
  });

  it('should unsubscribe from channels', async () => {
    const { subscribe, unsubscribe } = useLiveUpdate('ws://localhost:8080/ws');

    await nextTick();
    await new Promise((resolve) => setTimeout(resolve, 10));

    subscribe('test-channel');
    unsubscribe('test-channel');

    await nextTick();

    const messages = mockWebSocket.sent;
    expect(messages).toContainEqual(
      JSON.stringify({
        type: 'unsubscribe',
        channel: 'test-channel',
      })
    );
  });

  it('should receive and handle messages', async () => {
    const { on } = useLiveUpdate('ws://localhost:8080/ws');

    await nextTick();
    await new Promise((resolve) => setTimeout(resolve, 10));

    const handler = vi.fn();
    on('test-channel', handler);

    // Simulate receiving a message
    mockWebSocket.simulateMessage({
      channel: 'test-channel',
      type: 'update',
      data: { id: 123, name: 'Test' },
    });

    await nextTick();

    expect(handler).toHaveBeenCalledWith({
      id: 123,
      name: 'Test',
    });
  });

  it('should handle multiple handlers for same channel', async () => {
    const { on } = useLiveUpdate('ws://localhost:8080/ws');

    await nextTick();
    await new Promise((resolve) => setTimeout(resolve, 10));

    const handler1 = vi.fn();
    const handler2 = vi.fn();

    on('test-channel', handler1);
    on('test-channel', handler2);

    mockWebSocket.simulateMessage({
      channel: 'test-channel',
      type: 'update',
      data: { value: 42 },
    });

    await nextTick();

    expect(handler1).toHaveBeenCalledWith({ value: 42 });
    expect(handler2).toHaveBeenCalledWith({ value: 42 });
  });

  it('should remove handlers with off', async () => {
    const { on, off } = useLiveUpdate('ws://localhost:8080/ws');

    await nextTick();
    await new Promise((resolve) => setTimeout(resolve, 10));

    const handler = vi.fn();
    on('test-channel', handler);
    off('test-channel', handler);

    mockWebSocket.simulateMessage({
      channel: 'test-channel',
      type: 'update',
      data: { value: 42 },
    });

    await nextTick();

    expect(handler).not.toHaveBeenCalled();
  });

  it('should auto-reconnect on disconnect', async () => {
    const { connected } = useLiveUpdate('ws://localhost:8080/ws', {
      reconnect: true,
      reconnectDelay: 10,
    });

    await nextTick();
    await new Promise((resolve) => setTimeout(resolve, 10));

    expect(connected.value).toBe(true);
    const firstWs = websocketInstances[0];

    // Simulate disconnect
    firstWs.close();

    await nextTick();
    expect(connected.value).toBe(false);

    // Wait for reconnection
    await new Promise((resolve) => setTimeout(resolve, 100));

    // A new WebSocket should have been created
    expect(websocketInstances.length).toBe(2);
    expect(connected.value).toBe(true);
  });

  it('should not reconnect when disabled', async () => {
    const { connected } = useLiveUpdate('ws://localhost:8080/ws', {
      reconnect: false,
    });

    await nextTick();
    await new Promise((resolve) => setTimeout(resolve, 10));

    mockWebSocket.close();

    await nextTick();
    await new Promise((resolve) => setTimeout(resolve, 50));

    expect(connected.value).toBe(false);
  });

  it('should cleanup on disconnect', async () => {
    const { disconnect } = useLiveUpdate('ws://localhost:8080/ws');

    await nextTick();
    await new Promise((resolve) => setTimeout(resolve, 10));

    const closeSpy = vi.spyOn(mockWebSocket, 'close');

    disconnect();

    expect(closeSpy).toHaveBeenCalled();
  });

  it('should handle connection errors', async () => {
    const { error } = useLiveUpdate('ws://localhost:8080/ws');

    await nextTick();
    await new Promise((resolve) => setTimeout(resolve, 10));

    // Simulate error
    if (mockWebSocket.onerror) {
      mockWebSocket.onerror(new Event('error'));
    }

    await nextTick();

    expect(error.value).not.toBeNull();
  });

  it('should merge updates with reactive state', async () => {
    const { on } = useLiveUpdate('ws://localhost:8080/ws');

    await nextTick();
    await new Promise((resolve) => setTimeout(resolve, 10));

    const state = ref({ id: 1, name: 'Original' });

    on('items', (data: any) => {
      Object.assign(state.value, data);
    });

    mockWebSocket.simulateMessage({
      channel: 'items',
      type: 'update',
      data: { name: 'Updated' },
    });

    await nextTick();

    expect(state.value).toEqual({ id: 1, name: 'Updated' });
  });

  it('should handle broadcast messages', async () => {
    const { on } = useLiveUpdate('ws://localhost:8080/ws');

    await nextTick();
    await new Promise((resolve) => setTimeout(resolve, 10));

    const handler = vi.fn();
    on('*', handler);

    mockWebSocket.simulateMessage({
      channel: '*',
      type: 'announcement',
      data: 'System message',
    });

    await nextTick();

    expect(handler).toHaveBeenCalledWith('System message');
  });
});

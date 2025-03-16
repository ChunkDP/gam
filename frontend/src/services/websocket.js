class WebSocketService {
  constructor() {
    this.socket = null;
    this._isConnected = false;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectInterval = 3000;
    this.listeners = new Map();
    this.messageQueue = [];
  }

  isConnected() {
    return this._isConnected && this.socket && this.socket.readyState === WebSocket.OPEN;
  }

  connect(token) {
    if (this.isConnected()) {
      return Promise.resolve();
    }

    return new Promise((resolve, reject) => {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const host = window.location.host;
      const url = `${protocol}//${host}/ws/notifications?token=${encodeURIComponent(token)}`;
     // const url = `ws://127.0.0.1:8080/ws/notifications?token=${encodeURIComponent(token)}`;
      this.socket = new WebSocket(url);
      // 创建WebSocket连接，并添加token到请求头
     
     
      this.socket.onopen = () => {
        console.log('WebSocket connected');
        this._isConnected = true;
        this.reconnectAttempts = 0;
        
        // 发送队列中的消息
        while (this.messageQueue.length > 0) {
          const message = this.messageQueue.shift();
          this.socket.send(JSON.stringify(message));
        }
        
        resolve();
      };

      this.socket.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          this.handleMessage(data);
        } catch (error) {
          console.error('Error parsing WebSocket message:', error);
        }
      };

      this.socket.onclose = (event) => {
        console.log('WebSocket disconnected:', event.code, event.reason);
        this._isConnected = false;
        
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
          this.reconnectAttempts++;
          setTimeout(() => this.connect(token), this.reconnectInterval);
        }
      };

      this.socket.onerror = (error) => {
        console.error('WebSocket error:', error);
        reject(error);
      };
    });
  }

  disconnect() {
    if (this.socket) {
      this.socket.close();
      this.socket = null;
      this._isConnected = false;
    }
  }

  send(message) {
    if (this.isConnected()) {
      this.socket.send(JSON.stringify(message));
    } else {
      this.messageQueue.push(message);
    }
  }

  handleMessage(data) {
    const { type } = data;
    
    if (this.listeners.has(type)) {
      const callbacks = this.listeners.get(type);
      callbacks.forEach(callback => callback(data));
    }
    
    // 触发全局监听器
    if (this.listeners.has('*')) {
      const callbacks = this.listeners.get('*');
      callbacks.forEach(callback => callback(data));
    }

    if (type === 'notification') {
      switch (data.action) {
        case 'recall':
          // 触发撤回事件
          this.emit('notification-recall', {
            id: data.id,
            message: data.message
          });
          break;
        // ... 其他情况处理
      }
    }
  }

  on(type, callback) {
    if (!this.listeners.has(type)) {
      this.listeners.set(type, []);
    }
    
    this.listeners.get(type).push(callback);
    
    return () => {
      const callbacks = this.listeners.get(type);
      const index = callbacks.indexOf(callback);
      if (index !== -1) {
        callbacks.splice(index, 1);
      }
    };
  }

  off(type, callback) {
    if (!this.listeners.has(type)) {
      return;
    }
    
    const callbacks = this.listeners.get(type);
    const index = callbacks.indexOf(callback);
    if (index !== -1) {
      callbacks.splice(index, 1);
    }
  }

  reconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      const token = sessionStorage.getItem('token');
      if (token) {
        console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
        return this.connect(token);
      }
    }
    return Promise.reject(new Error('Max reconnection attempts reached'));
  }

  emit(type, data) {
    if (this.listeners.has(type)) {
      const callbacks = this.listeners.get(type);
      callbacks.forEach(callback => callback(data));
    }
  }
}

export default new WebSocketService(); 
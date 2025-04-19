/**
 * PingService - Handles connection status monitoring between frontend and backend
 * 
 * This service implements a ping mechanism with WebSocket as primary method and
 * HTTP keepalive as fallback. If both fail, it will clear the user session.
 */

class PingService {
  constructor(options = {}) {
    this.baseUrl = options.baseUrl || window.location.origin;
    this.wsUrl = options.wsUrl || `${this.baseUrl.replace(/^http/, 'ws')}/ws`;
    this.keepaliveUrl = options.keepaliveUrl || `${this.baseUrl}/api/v1/keepalive`;
    this.pingInterval = options.pingInterval || 30000; // 30 seconds
    this.reconnectAttempts = options.reconnectAttempts || 3;
    this.reconnectDelay = options.reconnectDelay || 2000; // 2 seconds
    
    this.ws = null;
    this.pingTimer = null;
    this.attemptCount = 0;
    this.onDisconnect = options.onDisconnect || this.defaultDisconnectHandler;
    this.isConnected = false;
    
    // Bind methods to maintain 'this' context
    this.connect = this.connect.bind(this);
    this.disconnect = this.disconnect.bind(this);
    this.reconnect = this.reconnect.bind(this);
    this.httpKeepAlive = this.httpKeepAlive.bind(this);
    this.clearSession = this.clearSession.bind(this);
  }

  /**
   * Start the ping service
   * @param {string} token - JWT token for authentication
   * @returns {Promise} - Resolves when connection is established
   */
  connect(token) {
    this.token = token;
    return new Promise((resolve, reject) => {
      try {
        // Create WebSocket connection
        this.ws = new WebSocket(this.wsUrl);
        
        this.ws.onopen = () => {
          console.log('WebSocket connection established');
          this.isConnected = true;
          this.attemptCount = 0;
          this.startPingTimer();
          resolve(true);
        };
        
        this.ws.onclose = () => {
          console.log('WebSocket connection closed');
          this.isConnected = false;
          this.reconnect();
        };
        
        this.ws.onerror = (error) => {
          console.error('WebSocket error:', error);
          this.isConnected = false;
          this.reconnect();
        };
        
        // Handle incoming messages (not expected in ping service, but just in case)
        this.ws.onmessage = (event) => {
          console.log('Received message:', event.data);
        };
      } catch (error) {
        console.error('Failed to create WebSocket connection:', error);
        this.reconnect();
        reject(error);
      }
    });
  }

  /**
   * Start the ping timer to keep the connection alive
   */
  startPingTimer() {
    // Clear any existing timer
    if (this.pingTimer) {
      clearInterval(this.pingTimer);
    }
    
    // Set up new timer
    this.pingTimer = setInterval(() => {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        // Send a ping message
        this.ws.send('ping');
      } else {
        this.reconnect();
      }
    }, this.pingInterval);
  }

  /**
   * Attempt to reconnect when WebSocket connection fails
   */
  reconnect() {
    // Clear existing timer
    if (this.pingTimer) {
      clearInterval(this.pingTimer);
      this.pingTimer = null;
    }
    
    // If we've exceeded max attempts, try HTTP keepalive
    if (this.attemptCount >= this.reconnectAttempts) {
      console.log('Max WebSocket reconnect attempts reached, falling back to HTTP keepalive');
      this.httpKeepAlive();
      return;
    }
    
    // Increment attempt counter
    this.attemptCount++;
    
    // Try to reconnect after delay
    setTimeout(() => {
      console.log(`Attempting to reconnect (${this.attemptCount}/${this.reconnectAttempts})`);
      this.connect(this.token).catch(() => {
        // Connect method will handle reconnection if it fails
      });
    }, this.reconnectDelay);
  }

  /**
   * Fallback to HTTP keepalive when WebSocket fails
   */
  httpKeepAlive() {
    console.log('Attempting HTTP keepalive');
    
    fetch(this.keepaliveUrl, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${this.token}`,
        'Content-Type': 'application/json'
      }
    })
    .then(response => {
      if (response.ok) {
        console.log('HTTP keepalive successful');
        // Schedule next HTTP keepalive
        setTimeout(this.httpKeepAlive, this.pingInterval);
      } else {
        console.error('HTTP keepalive failed with status:', response.status);
        this.clearSession();
      }
    })
    .catch(error => {
      console.error('HTTP keepalive request failed:', error);
      this.clearSession();
    });
  }

  /**
   * Clear user session when all connection attempts fail
   */
  clearSession() {
    console.log('Clearing session due to connection failure');
    
    // Call the disconnect handler
    if (typeof this.onDisconnect === 'function') {
      this.onDisconnect();
    }
  }

  /**
   * Default disconnect handler - can be overridden in constructor options
   */
  defaultDisconnectHandler() {
    console.log('Connection to server lost. Logging out...');
    
    // Clear local storage
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    
    // Redirect to login page
    window.location.href = '/login';
  }

  /**
   * Manually disconnect the ping service
   */
  disconnect() {
    if (this.pingTimer) {
      clearInterval(this.pingTimer);
      this.pingTimer = null;
    }
    
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    
    this.isConnected = false;
    this.attemptCount = 0;
    console.log('Ping service disconnected');
  }
}

export default PingService;

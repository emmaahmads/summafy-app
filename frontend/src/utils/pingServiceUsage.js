/**
 * Example usage of PingService in a Vue.js application
 * 
 * This example shows how to integrate the ping service into your frontend application.
 * You can include this in your main.js or App.vue file.
 */

import PingService from './pingService';

// Example in Vue.js application
export function setupPingService(store, router) {
  // Create ping service instance
  const pingService = new PingService({
    // Optional custom configuration
    baseUrl: process.env.VUE_APP_API_URL || window.location.origin,
    pingInterval: 30000, // 30 seconds
    reconnectAttempts: 3,
    
    // Custom disconnect handler
    onDisconnect: () => {
      // Clear user state in Vuex store
      store.dispatch('auth/logout');
      
      // Show notification to user
      store.dispatch('notification/show', {
        type: 'error',
        message: 'Connection to server lost. You have been logged out.'
      });
      
      // Redirect to login page
      router.push('/login');
    }
  });

  // Start ping service when user logs in
  store.subscribe((mutation, state) => {
    if (mutation.type === 'auth/SET_TOKEN') {
      const token = state.auth.token;
      
      if (token) {
        // Start ping service with the token
        pingService.connect(token)
          .then(() => {
            console.log('Ping service started successfully');
          })
          .catch(error => {
            console.error('Failed to start ping service:', error);
          });
      } else {
        // Disconnect ping service when token is cleared
        pingService.disconnect();
      }
    }
  });

  // Return the ping service instance for direct access if needed
  return pingService;
}

/**
 * Example usage in React application
 */
export function setupReactPingService(token, logoutCallback) {
  const pingService = new PingService({
    onDisconnect: () => {
      // Clear local storage
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      
      // Call the logout callback
      if (typeof logoutCallback === 'function') {
        logoutCallback();
      }
    }
  });

  // Start ping service with the token
  if (token) {
    pingService.connect(token)
      .catch(error => {
        console.error('Failed to start ping service:', error);
      });
  }

  // Return the ping service instance
  return pingService;
}

/**
 * Example usage in vanilla JavaScript application
 */
export function initPingService() {
  // Get token from localStorage or wherever it's stored
  const token = localStorage.getItem('token');
  
  if (!token) {
    console.warn('No token available, ping service not started');
    return null;
  }
  
  const pingService = new PingService();
  
  // Start ping service
  pingService.connect(token)
    .then(() => {
      console.log('Ping service connected');
    })
    .catch(error => {
      console.error('Failed to connect ping service:', error);
    });
    
  return pingService;
}

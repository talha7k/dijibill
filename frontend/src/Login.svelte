<script>
  import { createEventDispatcher } from 'svelte';
  import { Login } from '../wailsjs/go/main/App';

  const dispatch = createEventDispatcher();

  let username = '';
  let password = '';
  let rememberMe = false;
  let loading = false;
  let error = '';

  // Load saved credentials on component mount
  function loadSavedCredentials() {
    try {
      const savedCredentials = localStorage.getItem('dijibill_saved_credentials');
      if (savedCredentials) {
        const credentials = JSON.parse(savedCredentials);
        username = credentials.username || '';
        rememberMe = true;
      }
    } catch (error) {
      console.error('Error loading saved credentials:', error);
    }
  }

  // Save credentials to localStorage
  function saveCredentials() {
    try {
      if (rememberMe && username) {
        const credentials = {
          username: username,
          savedAt: new Date().toISOString()
        };
        localStorage.setItem('dijibill_saved_credentials', JSON.stringify(credentials));
      } else {
        localStorage.removeItem('dijibill_saved_credentials');
      }
    } catch (error) {
      console.error('Error saving credentials:', error);
    }
  }

  async function handleLogin() {
    if (!username || !password) {
      error = 'Please enter both username and password';
      return;
    }

    loading = true;
    error = '';

    try {
      const authContext = await Login(username, password);
      
      // Save credentials if remember me is checked
      saveCredentials();
      
      dispatch('login-success', authContext);
    } catch (err) {
      console.error('Login error:', err);
      error = err.message || 'Login failed';
    } finally {
      loading = false;
    }
  }

  function handleKeyPress(event) {
    if (event.key === 'Enter') {
      handleLogin();
    }
  }

  function handleRememberMeChange() {
    if (!rememberMe) {
      // If unchecking remember me, remove saved credentials
      localStorage.removeItem('dijibill_saved_credentials');
    }
  }

  // Load saved credentials when component mounts
  loadSavedCredentials();
</script>

<div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
  <div class="max-w-md w-full space-y-8 p-8 bg-white rounded-xl shadow-lg">
    <div class="text-center">
      <h2 class="text-3xl font-bold text-gray-900 mb-2">DijiBill</h2>
      <p class="text-gray-600">Sign in to your account</p>
    </div>

    <form on:submit|preventDefault={handleLogin} class="space-y-6">
      {#if error}
        <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
          {error}
        </div>
      {/if}

      <div>
        <label for="username" class="block text-sm font-medium text-gray-700 mb-2">
          Username
        </label>
        <div class="relative">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <i class="fas fa-user text-gray-400"></i>
          </div>
          <input
            id="username"
            type="text"
            bind:value={username}
            on:keypress={handleKeyPress}
            disabled={loading}
            class="w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100 transition-colors"
            placeholder="Enter your username"
            required
          />
        </div>
      </div>

      <div>
        <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
          Password
        </label>
        <div class="relative">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <i class="fas fa-lock text-gray-400"></i>
          </div>
          <input
            id="password"
            type="password"
            bind:value={password}
            on:keypress={handleKeyPress}
            disabled={loading}
            class="w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100 transition-colors"
            placeholder="Enter your password"
            required
          />
        </div>
      </div>

      <!-- Remember Me Checkbox -->
      <div class="flex items-center">
        <input
          id="remember-me"
          type="checkbox"
          bind:checked={rememberMe}
          on:change={handleRememberMeChange}
          class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
        />
        <label for="remember-me" class="ml-2 block text-sm text-gray-700">
          Remember my username
        </label>
      </div>

      <button
        type="submit"
        disabled={loading}
        class="w-full flex justify-center py-2 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {#if loading}
          <i class="fas fa-spinner fa-spin -ml-1 mr-3 h-5 w-5 text-white"></i>
          Signing in...
        {:else}
          Sign in
        {/if}
      </button>
    </form>

    <div class="text-center text-sm text-gray-600">
      <p>Default credentials:</p>
      <p class="font-mono text-xs bg-gray-100 p-2 rounded mt-1">
        Username: admin<br>
        Password: password
      </p>
      <div class="mt-4 pt-4 border-t border-gray-200">
        <p class="text-gray-600">Don't have an account?</p>
        <button
          type="button"
          on:click={() => dispatch('switch-to-signup')}
          class="text-blue-600 hover:text-blue-700 font-medium"
        >
          Create a new account
        </button>
      </div>
    </div>
  </div>
</div>
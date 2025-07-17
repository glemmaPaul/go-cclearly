<script>
  import Header from './components/Header.svelte';
  import HistoryPane from './components/HistoryPane.svelte';
  import RequestPane from './components/RequestPane.svelte';
  import ResponsePane from './components/ResponsePane.svelte';
  import { ExecuteRequest, GetHistory, GetRequestByID, ParseCurlCommand } from '../wailsjs/go/main/App.js';

  // State management
  let selectedHistoryItem = null;
  let requestData = {
    method: 'GET',
    url: '',
    headers: [],
    body: ''
  };
  let responseData = {
    status: null,
    body: '',
    headers: []
  };

  // History data
  let history = [];
  let loading = false;

  // Load history on startup
  async function loadHistory() {
    try {
      console.log('Loading history...');
      history = await GetHistory(50); // Get last 50 requests
      console.log('History loaded:', history);
    } catch (error) {
      console.error('Failed to load history:', error);
      history = [];
    }
  }

  // Load history when component mounts
  loadHistory();

  async function handleSendRequest() {
    console.log('Send Request clicked!');
    console.log('Current request data:', requestData);
    
    if (!requestData.url.trim()) {
      alert('Please enter a URL');
      return;
    }

    loading = true;
    try {
      // Convert frontend format to backend format
      const backendRequest = {
        method: requestData.method,
        url: requestData.url,
        headers: {},
        body: requestData.body
      };

      // Convert headers array to map
      requestData.headers.forEach(header => {
        if (header.key && header.value) {
          backendRequest.headers[header.key] = header.value;
        }
      });

      console.log('Sending backend request:', backendRequest);
      const response = await ExecuteRequest(backendRequest);
      console.log('Received response:', response);
      
      // Convert backend response to frontend format
      responseData = {
        status: response.statusCode,
        body: response.body,
        headers: Object.entries(response.headers || {}).map(([key, value]) => ({ key, value }))
      };
      
      console.log('Converted response data:', responseData);
      console.log('Status code:', responseData.status);

      // Reload history to show the new request
      await loadHistory();
    } catch (error) {
      console.error('Request failed:', error);
      responseData = {
        status: 0,
        body: `Error: ${error.message}`,
        headers: []
      };
    } finally {
      loading = false;
    }
  }

  function handleRequestDataChange(newData) {
    requestData = newData;
  }

  async function handleSelectHistoryItem(id) {
    selectedHistoryItem = id;
    try {
      console.log('Loading history item with ID:', id);
      const historyItem = await GetRequestByID(id);
      console.log('History item loaded:', historyItem);
      
      if (historyItem) {
        // Parse the full command to populate the request form
        console.log('Parsing cURL command:', historyItem.fullCommand);
        const parsedRequest = await ParseCurlCommand(historyItem.fullCommand);
        console.log('Parsed request:', parsedRequest);
        
        if (parsedRequest) {
          // Convert headers from map to array format
          const headerArray = [];
          if (parsedRequest.headers) {
            Object.entries(parsedRequest.headers).forEach(([key, value]) => {
              headerArray.push({ key, value });
            });
          }
          
          requestData = {
            method: parsedRequest.method,
            url: parsedRequest.url,
            headers: headerArray,
            body: parsedRequest.body
          };
          console.log('Updated request data from history:', requestData);
        }
      }
    } catch (error) {
      console.error('Failed to load history item:', error);
    }
  }

  // Handle paste events for cURL commands
  async function handlePaste(event) {
    console.log('Paste event detected');
    const text = event.clipboardData?.getData('text') || '';
    console.log('Pasted text:', text);
    
    if (text.trim().toLowerCase().startsWith('curl ')) {
      console.log('cURL command detected, parsing...');
      try {
        const parsedRequest = await ParseCurlCommand(text);
        console.log('Parsed request:', parsedRequest);
        if (parsedRequest) {
          // Convert headers from map to array format
          const headerArray = [];
          if (parsedRequest.headers) {
            Object.entries(parsedRequest.headers).forEach(([key, value]) => {
              headerArray.push({ key, value });
            });
          }
          
          requestData = {
            method: parsedRequest.method,
            url: parsedRequest.url,
            headers: headerArray,
            body: parsedRequest.body
          };
          console.log('Updated request data from paste:', requestData);
        }
      } catch (error) {
        console.error('Failed to parse cURL command:', error);
      }
    }
  }

  // Add paste event listener to the document when component mounts
  import { onMount } from 'svelte';
  
  onMount(() => {
    console.log('Component mounted, adding paste listener');
    document.addEventListener('paste', handlePaste);
    
    return () => {
      document.removeEventListener('paste', handlePaste);
    };
  });
</script>

<main class="h-screen flex flex-col bg-gray-50">
  <Header onSendRequest={handleSendRequest} />
  
  <div class="flex-1 flex overflow-hidden">
    <div class="w-1/4 bg-white border-r border-gray-200 flex flex-col">
      <HistoryPane 
        {history}
        {selectedHistoryItem}
        onSelectHistoryItem={handleSelectHistoryItem}
      />
    </div>
    
    <div class="flex-1 flex overflow-hidden">
      <div class="w-1/2 bg-white border-r border-gray-200">
        <RequestPane 
          {requestData}
          onRequestDataChange={handleRequestDataChange}
        />
      </div>
      
      <div class="w-1/2 bg-white">
        <ResponsePane {responseData} />
      </div>
    </div>
  </div>

  {#if loading}
    <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white p-6 rounded-lg shadow-lg">
        <div class="flex items-center space-x-3">
          <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
          <span class="text-gray-700">Sending request...</span>
        </div>
      </div>
    </div>
  {/if}
</main>

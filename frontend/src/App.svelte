<script>
  import Header from './components/Header.svelte';
  import HistoryPane from './components/HistoryPane.svelte';
  import RequestPane from './components/RequestPane.svelte';
  import ResponsePane from './components/ResponsePane.svelte';
  import { ExecuteRequest, GetHistory, GetRequestByID, ParseCurlCommand } from '../wailsjs/go/main/App.js';
  import { responseStore } from './stores/responseStore.js';

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
  
  // Force reactivity trigger
  let responseKey = 0;

  // History data
  let history = [];
  let loading = false;
  let databaseStatus = { connected: true, message: 'Initializing...' };

  // Load history on startup
  async function loadHistory() {
    try {
      console.log('Loading history...');
      const historyResult = await GetHistory(50); // Get last 50 requests
      
      // Handle case where GetHistory returns null
      if (historyResult === null) {
        console.log('GetHistory returned null, retrying in 1 second...');
        // Retry once after a short delay in case database is still initializing
        setTimeout(async () => {
          try {
            const retryResult = await GetHistory(50);
            if (retryResult === null) {
              history = [];
              databaseStatus = { connected: false, message: 'Database not available' };
            } else {
              history = retryResult;
              databaseStatus = { connected: true, message: `Database ready (${history.length} items)` };
            }
          } catch (retryError) {
            console.error('Retry failed:', retryError);
            history = [];
            databaseStatus = { connected: false, message: `Database error: ${retryError.message}` };
          }
        }, 1000);
        return;
      }
      
      history = historyResult;
      console.log('History loaded:', history);
      
      // Check if database is working by seeing if we get a proper response
      if (history.length === 0) {
        databaseStatus = { connected: true, message: 'Database ready (no history yet)' };
      } else {
        databaseStatus = { connected: true, message: `Database ready (${history.length} items)` };
      }
    } catch (error) {
      console.error('Failed to load history:', error);
      history = [];
      databaseStatus = { connected: false, message: `Database error: ${error.message}` };
    }
  }

  // Load history when component mounts
  import { onMount } from 'svelte';
  
  onMount(() => {
    console.log('Component mounted, adding paste listener');
    document.addEventListener('paste', handlePaste);
    
    // Add a small delay to ensure database is initialized
    setTimeout(() => {
      loadHistory();
    }, 500);
    
    return () => {
      document.removeEventListener('paste', handlePaste);
    };
  });

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
      
      // Convert backend response to frontend format and force reactivity
      const newResponseData = {
        status: response.statusCode,
        body: response.formattedBody || response.body || '',
        formattedBody: response.formattedBody || response.body || '',
        responseType: response.responseType || 'raw',
        headers: Object.entries(response.headers || {}).map(([key, value]) => ({ key, value }))
      };
      
      // Update the store to trigger reactivity
      responseStore.set(newResponseData);
      
      console.log('Converted response data:', newResponseData);
      console.log('Status code:', newResponseData.status);
      console.log('Response body length:', newResponseData.body?.length || 0);
      console.log('Response headers count:', newResponseData.headers?.length || 0);

      // Reload history to show the new request
      await loadHistory();
      
      // Auto-select the newest item (first in the list)
      if (history.length > 0) {
        selectedHistoryItem = history[0].id;
        console.log('Auto-selected new history item:', selectedHistoryItem);
        // Update database status to show it's working
        databaseStatus = { connected: true, message: `Database ready (${history.length} items)` };
      }
    } catch (error) {
      console.error('Request failed:', error);
      responseStore.set({
        status: 0,
        body: `Error: ${error.message}`,
        formattedBody: `Error: ${error.message}`,
        responseType: 'raw',
        headers: []
      });
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

        // Load and display the response data
        if (historyItem.responseBody) {
          console.log('Loading response data for history item');
          const responseData = {
            status: historyItem.statusCode,
            body: historyItem.responseBody,
            formattedBody: historyItem.responseBody,
            responseType: historyItem.responseType || 'raw',
            headers: []
          };

          // Parse headers if available
          if (historyItem.responseHeaders) {
            try {
              const headersObj = JSON.parse(historyItem.responseHeaders);
              responseData.headers = Object.entries(headersObj).map(([key, value]) => ({ key, value }));
            } catch (error) {
              console.error('Failed to parse response headers:', error);
            }
          }

          // Update the response store
          responseStore.set(responseData);
          console.log('Response data loaded for history item:', responseData);
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
</script>

<main class="h-screen flex flex-col bg-gray-900">
  <Header onSendRequest={handleSendRequest} />
  
  <div class="flex-1 flex overflow-hidden">
    <div class="w-1/4 bg-gray-800 border-r border-gray-700 flex flex-col">
      <HistoryPane 
        {history}
        {selectedHistoryItem}
        {databaseStatus}
        onSelectHistoryItem={handleSelectHistoryItem}
      />
    </div>
    
    <div class="flex-1 flex overflow-hidden">
      <div class="w-1/2 bg-gray-800 border-r border-gray-700">
        <RequestPane 
          {requestData}
          onRequestDataChange={handleRequestDataChange}
          onSendRequest={handleSendRequest}
        />
      </div>
      
      <div class="w-1/2 bg-gray-800">
        <ResponsePane {loading} />
      </div>
    </div>
  </div>
</main>

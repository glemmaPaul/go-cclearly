<script>
  import { getStatusColor, getStatusText } from '../utils/statusColors.js';
  import { responseStore } from '../stores/responseStore.js';
  
  // Props
  export let loading = false;
  
  // Subscribe to the store
  $: responseData = $responseStore;

  // Debug logging
  $: {
    console.log('ResponsePane received data:', responseData);
    console.log('Status:', responseData.status);
    console.log('Body length:', responseData.body?.length || 0);
    console.log('Response type:', responseData.responseType);
  }

  // State for collapsible headers
  let headersExpanded = false;

  function getBodyClass() {
    if (responseData.responseType === 'json') {
      return 'text-sm font-mono text-gray-100 whitespace-pre-wrap';
    }
    return 'text-sm font-mono text-gray-100 whitespace-pre-wrap break-all';
  }

  function toggleHeaders() {
    headersExpanded = !headersExpanded;
  }
</script>

<div class="flex flex-col h-full relative">
  <div class="p-4 bg-gray-800">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-lg font-semibold text-white font-mono">Response</h2>
      </div>
      {#if responseData.status}
        <div class="flex items-center space-x-2">
          <span class="px-3 py-1 text-sm font-semibold rounded-full {getStatusColor(responseData.status)}">
            {responseData.status}
          </span>
          <span class="text-sm text-gray-300">{getStatusText(responseData.status)}</span>
          {#if responseData.responseType === 'json'}
            <span class="px-2 py-1 text-xs bg-blue-900 text-blue-300 rounded-full">JSON</span>
          {/if}
        </div>
      {:else}
        <div class="text-xs text-gray-400">No status</div>
      {/if}
    </div>
  </div>
  
  <div class="flex-1 px-4 pb-4 flex flex-col overflow-hidden">
    {#if responseData.status}
      <div class="flex flex-col h-full space-y-4">
        <!-- Response Headers (Collapsible) -->
        {#if responseData.headers && responseData.headers.length > 0}
          <div class="flex-shrink-0">
            <button 
              on:click={toggleHeaders}
              class="flex items-center justify-between w-full p-3 bg-gray-700 border border-gray-600 rounded-lg hover:bg-gray-600 transition-colors"
            >
              <div class="flex items-center space-x-2">
                <span class="text-sm font-medium text-gray-300">Response Headers ({responseData.headers.length})</span>
              </div>
              <div class="flex items-center space-x-2">
                <span class="text-xs text-gray-400">
                  {headersExpanded ? 'Hide' : 'Show'}
                </span>
                <svg 
                  class="w-4 h-4 text-gray-400 transition-transform {headersExpanded ? 'rotate-180' : ''}" 
                  fill="none" 
                  stroke="currentColor" 
                  viewBox="0 0 24 24"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
                </svg>
              </div>
            </button>
            
            {#if headersExpanded}
              <div class="mt-2 bg-gray-700 border border-gray-600 rounded-lg p-4 max-h-48 overflow-auto">
                <div class="space-y-2">
                  {#each responseData.headers as header}
                    <div class="flex text-sm">
                      <span class="font-medium text-gray-300 min-w-0 flex-1">{header.key}:</span>
                      <span class="text-gray-400 ml-2 break-all">{header.value}</span>
                    </div>
                  {/each}
                </div>
              </div>
            {/if}
          </div>
        {/if}

        <!-- Response Body (Fills remaining height) -->
        <div class="flex-1 flex flex-col min-h-0">
          {#if responseData.responseType === 'json'}
            <div class="flex-shrink-0 mb-2">
              <span class="text-xs text-blue-400 bg-blue-900 px-2 py-1 rounded">Formatted JSON</span>
            </div>
          {/if}
          <div class="flex-1 bg-gray-700 border border-gray-600 rounded-lg p-4 overflow-auto">
            <pre class="{getBodyClass()} h-full">{responseData.formattedBody || responseData.body || 'No response body'}</pre>
          </div>
        </div>
      </div>
    {:else}
      <div class="flex-1 flex items-center justify-center text-gray-400">
        <div class="text-center">
          <div class="text-4xl mb-4">📡</div>
          <p class="text-lg font-medium mb-2">No response yet</p>
          <p class="text-sm">Send a request to see the response here</p>
          <div class="mt-6 p-4 bg-gray-700 border border-gray-600 rounded-lg">
            <p class="text-xs text-gray-300">
              💡 <strong>Tip:</strong> Try a simple GET request to https://httpbin.org/json to test JSON formatting
            </p>
          </div>
        </div>
      </div>
    {/if}
  </div>

  {#if loading}
    <div class="absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-gray-800 p-6 rounded-lg border border-gray-600">
        <div class="flex items-center space-x-3">
          <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-400"></div>
          <span class="text-gray-200">Sending request...</span>
        </div>
      </div>
    </div>
  {/if}
</div> 
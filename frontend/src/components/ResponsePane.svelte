<script>
  export let responseData = { status: null, body: '', headers: [] };

  function getStatusColor(status) {
    console.log('Getting status color for:', status, 'Type:', typeof status);
    if (status >= 200 && status < 300) return 'bg-green-50 text-green-700';
    if (status >= 400 && status < 500) return 'bg-yellow-50 text-yellow-700';
    if (status >= 500) return 'bg-red-50 text-red-700';
    return 'bg-gray-50 text-gray-700';
  }

  function getStatusText(status) {
    if (status >= 200 && status < 300) return 'Success';
    if (status >= 400 && status < 500) return 'Client Error';
    if (status >= 500) return 'Server Error';
    return 'Unknown';
  }
</script>

<div class="flex flex-col h-full">
  <div class="p-4 bg-white">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-lg font-semibold text-gray-900">Response</h2>
        <p class="text-sm text-gray-500 mt-1">View response details</p>
      </div>
      {#if responseData.status}
        <div class="flex items-center space-x-2">
          <span class="px-3 py-1 text-sm font-semibold rounded-full {getStatusColor(responseData.status)}">
            {responseData.status}
          </span>
          <span class="text-sm text-gray-600">{getStatusText(responseData.status)}</span>
        </div>
      {:else}
        <div class="text-xs text-gray-400">No status</div>
      {/if}
    </div>
  </div>
  
  <div class="flex-1 px-4 pb-4 overflow-y-auto">
    {#if responseData.status}
      <div class="space-y-6">
        <!-- Response Body -->
        <div class="space-y-2">
          <h3 class="text-sm font-medium text-gray-700">Response Body</h3>
          <div class="bg-gray-50 border border-gray-200 rounded-lg p-4 max-h-64 overflow-auto">
            <pre class="text-sm font-mono text-gray-800 whitespace-pre-wrap">{responseData.body || 'No response body'}</pre>
          </div>
        </div>

        <!-- Response Headers -->
        {#if responseData.headers && responseData.headers.length > 0}
          <div class="space-y-2">
            <h3 class="text-sm font-medium text-gray-700">Response Headers</h3>
            <div class="bg-gray-50 border border-gray-200 rounded-lg p-4 max-h-48 overflow-auto">
              <div class="space-y-2">
                {#each responseData.headers as header}
                  <div class="flex text-sm">
                    <span class="font-medium text-gray-700 min-w-0 flex-1">{header.key}:</span>
                    <span class="text-gray-600 ml-2 break-all">{header.value}</span>
                  </div>
                {/each}
              </div>
            </div>
          </div>
        {/if}
      </div>
    {:else}
      <div class="text-center py-12 text-gray-500">
        <div class="text-4xl mb-4">📡</div>
        <p class="text-lg font-medium mb-2">No response yet</p>
        <p class="text-sm">Send a request to see the response here</p>
        <div class="mt-6 p-4 bg-blue-50 border border-blue-200 rounded-lg">
          <p class="text-xs text-blue-700">
            💡 <strong>Tip:</strong> Try a simple GET request to https://httpbin.org/get to test the app
          </p>
        </div>
      </div>
    {/if}
  </div>
</div> 
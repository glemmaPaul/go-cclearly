<script>
  export let requestData;
  export let onRequestDataChange;
  export let onSendRequest;

  function addHeader() {
    const newHeaders = [...requestData.headers, { key: '', value: '' }];
    onRequestDataChange({ ...requestData, headers: newHeaders });
  }

  function removeHeader(index) {
    const newHeaders = requestData.headers.filter((_, i) => i !== index);
    onRequestDataChange({ ...requestData, headers: newHeaders });
  }

  function updateHeader(index, field, value) {
    const newHeaders = [...requestData.headers];
    newHeaders[index] = { ...newHeaders[index], [field]: value };
    onRequestDataChange({ ...requestData, headers: newHeaders });
  }

  function handleSendClick() {
    if (onSendRequest) {
      onSendRequest();
    }
  }

  function handleKeyPress(event) {
    if (event.key === 'Enter' && (event.ctrlKey || event.metaKey)) {
      handleSendClick();
    }
  }
</script>

<div class="flex-1 bg-gray-800 flex flex-col h-full">
  <div class="p-4 bg-gray-800">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-lg font-semibold text-white font-mono">Request</h2>
      </div>
      <button 
        on:click={handleSendClick}
        class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium"
      >
        Send Request
      </button>
    </div>
  </div>
  
  <div class="flex-1 px-4 pb-4 space-y-6 overflow-y-auto">
    <!-- Method and URL -->
    <div class="space-y-2">
      <label class="text-sm font-medium text-gray-300">Request URL</label>
      <div class="flex space-x-3">
        <select 
          bind:value={requestData.method}
          class="px-4 py-3 border border-gray-600 rounded-lg bg-gray-700 text-gray-100 text-sm font-medium"
        >
          <option value="GET">GET</option>
          <option value="POST">POST</option>
          <option value="PUT">PUT</option>
          <option value="DELETE">DELETE</option>
          <option value="PATCH">PATCH</option>
        </select>
        <input 
          type="text" 
          bind:value={requestData.url}
          on:keydown={handleKeyPress}
          placeholder="https://api.example.com/endpoint (Ctrl+Enter to send)"
          class="flex-1 px-4 py-3 border border-gray-600 rounded-lg bg-gray-700 text-gray-100 text-sm"
        />
      </div>
    </div>

    <!-- Headers -->
    <div class="space-y-3">
      <div class="flex items-center justify-between">
        <label class="text-sm font-medium text-gray-300">Headers</label>
        <button 
          on:click={addHeader}
          class="text-sm text-blue-400 hover:text-blue-300 font-medium flex items-center space-x-1"
        >
          <span>+</span>
          <span>Add Header</span>
        </button>
      </div>
      
      {#if requestData.headers.length === 0}
        <div class="text-center py-8 text-gray-400 border-2 border-dashed border-gray-600 rounded-lg bg-gray-700">
          <div class="text-2xl mb-2">🔧</div>
          <p class="text-sm">No headers added</p>
          <p class="text-xs mt-1">Click "Add Header" to get started</p>
        </div>
      {:else}
        <div class="space-y-3">
          {#each requestData.headers as header, index}
            <div class="flex space-x-3 items-center">
              <input 
                type="text" 
                bind:value={header.key}
                on:input={() => updateHeader(index, 'key', header.key)}
                placeholder="Content-Type"
                class="flex-1 px-3 py-2 border border-gray-600 rounded-md bg-gray-700 text-gray-100 text-sm"
              />
              <input 
                type="text" 
                bind:value={header.value}
                on:input={() => updateHeader(index, 'value', header.value)}
                placeholder="application/json"
                class="flex-1 px-3 py-2 border border-gray-600 rounded-md bg-gray-700 text-gray-100 text-sm"
              />
              <button 
                on:click={() => removeHeader(index)}
                class="p-2 text-red-400 hover:text-red-300 hover:bg-red-900 rounded-md"
                title="Remove header"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                </svg>
              </button>
            </div>
          {/each}
        </div>
      {/if}
    </div>

    <!-- Body -->
    {#if requestData.method !== 'GET'}
      <div class="space-y-2">
        <label class="text-sm font-medium text-gray-300">Request Body</label>
        <textarea 
          bind:value={requestData.body}
          placeholder="Enter JSON, form data, or raw text"
          rows="8"
          class="w-full px-4 py-3 border border-gray-700 rounded-lg bg-gray-600 text-gray-100 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 resize-none placeholder-gray-400"
        ></textarea>
        <p class="text-xs text-gray-400">Enter JSON, form data, or raw text</p>
      </div>
    {/if}
  </div>
</div> 
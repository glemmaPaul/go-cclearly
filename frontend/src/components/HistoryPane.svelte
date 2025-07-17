<script>
  export let history = [];
  export let selectedHistoryItem = null;
  export let onSelectHistoryItem;

  function getMethodColor(method) {
    switch (method) {
      case 'GET': return 'bg-green-100 text-green-800';
      case 'POST': return 'bg-blue-100 text-blue-800';
      case 'PUT': return 'bg-yellow-100 text-yellow-800';
      case 'DELETE': return 'bg-red-100 text-red-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  }

  function truncateUrl(url, maxLength = 50) {
    if (url.length <= maxLength) return url;
    
    const parts = url.split('/');
    if (parts.length <= 3) {
      // For simple URLs, truncate from the left to show the endpoint
      return '...' + url.substring(url.length - maxLength + 3);
    }
    
    // For complex URLs, keep domain and endpoint, truncate middle
    const domain = parts.slice(0, 3).join('/');
    const lastPart = parts[parts.length - 1];
    const availableSpace = maxLength - domain.length - lastPart.length - 6; // 6 for "..."
    
    if (availableSpace <= 0) {
      return domain + '/.../' + lastPart.substring(0, maxLength - domain.length - 6);
    }
    
    const middle = parts.slice(3, -1).join('/');
    if (middle.length <= availableSpace) {
      return domain + '/' + middle + '/' + lastPart;
    }
    
    return domain + '/.../' + lastPart;
  }

  function getStatusColor(status) {
    if (status >= 200 && status < 300) return 'text-green-600';
    if (status >= 400 && status < 500) return 'text-yellow-600';
    if (status >= 500) return 'text-red-600';
    return 'text-gray-600';
  }

  function getStatusBackground(status) {
    if (status >= 200 && status < 300) return 'bg-green-50 text-green-700';
    if (status >= 400 && status < 500) return 'bg-yellow-50 text-yellow-700';
    if (status >= 500) return 'bg-red-50 text-red-700';
    return 'bg-gray-50 text-gray-700';
  }

  function calculateMetadata(item) {
    let headerCount = 0;
    let payloadSize = 0;
    
    try {
      // Parse the full command to get headers and body
      const command = item.fullCommand;
      
      // Count headers (simple regex approach)
      const headerMatches = command.match(/-H\s+["']([^"']+)["']/g);
      headerCount = headerMatches ? headerMatches.length : 0;
      
      // Calculate payload size from body
      const dataMatches = command.match(/-d\s+["']([^"']+)["']/g);
      if (dataMatches) {
        dataMatches.forEach(match => {
          const body = match.replace(/-d\s+["']/, '').replace(/["']$/, '');
          payloadSize += new Blob([body]).size;
        });
      }
    } catch (error) {
      console.error('Error calculating metadata:', error);
    }
    
    return { headerCount, payloadSize };
  }

  function formatBytes(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
  }
</script>

<div class="flex flex-col h-full">
  <div class="p-4 bg-white">
    <h2 class="text-lg font-semibold text-gray-900">History</h2>
    <p class="text-sm text-gray-500 mt-1">Recent requests</p>
  </div>
  
  <div class="flex-1 overflow-y-auto">
    {#if history.length === 0}
      <div class="p-6 text-center text-gray-500">
        <div class="text-3xl mb-3">📝</div>
        <p class="text-sm">No requests yet</p>
        <p class="text-xs mt-1">Your request history will appear here</p>
      </div>
    {:else}
      {#each history as item}
        {@const metadata = calculateMetadata(item)}
        {@const statusColor = getStatusBackground(item.statusCode)}
        {@const debug = console.log('History item status:', item.statusCode, 'Color:', statusColor)}
        <button 
          class="w-full text-left p-4 border-b border-gray-100 hover:bg-gray-50 transition-colors duration-150 {selectedHistoryItem === item.id ? 'bg-blue-50 border-blue-200' : ''}"
          on:click={() => onSelectHistoryItem(item.id)}
        >
          <div class="flex items-center justify-between mb-2">
            <span class="px-2 py-1 text-xs font-semibold rounded {getMethodColor(item.method)}">
              {item.method}
            </span>
            <span class="px-2 py-1 text-xs font-mono rounded {statusColor}">
              {item.statusCode}
            </span>
          </div>
          
          <div class="mb-2">
            <p class="text-sm text-gray-900 font-medium truncate" title={item.url}>
              {truncateUrl(item.url)}
            </p>
          </div>
          
          <div class="flex items-center justify-between text-xs text-gray-500">
            <span class="truncate">{new Date(item.createdAt).toLocaleTimeString()}</span>
            <div class="flex items-center space-x-2">
              {#if metadata.headerCount > 0}
                <span class="text-xs" title="{metadata.headerCount} headers">📄 {metadata.headerCount}</span>
              {/if}
              {#if metadata.payloadSize > 0}
                <span class="text-xs" title="{formatBytes(metadata.payloadSize)}">📀 {formatBytes(metadata.payloadSize)}</span>
              {/if}
            </div>
          </div>
        </button>
      {/each}
    {/if}
  </div>
</div> 
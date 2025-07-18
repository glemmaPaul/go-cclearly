<script>
  export let history = [];
  export let selectedHistoryItem = null;
  export let onSelectHistoryItem;
  export let databaseStatus = { connected: true, message: '' };

  function getMethodColor(method) {
    return 'bg-gray-600 text-gray-200 text-xs font-bold uppercase px-2 py-0.5 rounded';
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

  function formatUrlDisplay(url) {
    try {
      const urlObj = new URL(url);
      const domain = urlObj.protocol + '//' + urlObj.host;
      const path = urlObj.pathname + urlObj.search + urlObj.hash;
      
      return {
        domain: domain,
        path: path || '/',
        fullUrl: url
      };
    } catch (error) {
      // If URL parsing fails, fall back to simple truncation
      return {
        domain: url.substring(0, 30) + (url.length > 30 ? '...' : ''),
        path: '',
        fullUrl: url
      };
    }
  }

  function getStatusColor(status) {
    if (status >= 200 && status < 300) return 'text-green-600';
    if (status >= 400 && status < 500) return 'text-yellow-600';
    if (status >= 500) return 'text-red-600';
    return 'text-gray-600';
  }

  function getStatusBackground(status) {
    if (status >= 200 && status < 300) return 'bg-emerald-900 text-emerald-300';
    if (status >= 400 && status < 500) return 'bg-amber-900 text-amber-300';
    if (status >= 500) return 'bg-red-900 text-red-300';
    return 'bg-gray-600 text-gray-200';
  }

  function calculateMetadata(item) {
    let headerCount = 0;
    
    try {
      // Parse the full command to get headers
      const command = item.fullCommand;
      
      // Count headers (simple regex approach)
      const headerMatches = command.match(/-H\s+["']([^"']+)["']/g);
      headerCount = headerMatches ? headerMatches.length : 0;
    } catch (error) {
      console.error('Error calculating metadata:', error);
    }
    
    return { headerCount };
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
  <div class="p-4 bg-gray-800">
    <div class="flex items-center justify-between mb-2">
      <h2 class="text-lg font-semibold text-white">History</h2>
      {#if !databaseStatus.connected}
        <span class="text-xs text-red-400 bg-red-900 px-2 py-1 rounded" title={databaseStatus.message}>
          ⚠️ DB
        </span>
      {/if}
    </div>
    {#if !databaseStatus.connected}
      <div class="text-xs text-red-300 bg-red-900/20 p-2 rounded border border-red-800">
        {databaseStatus.message}
      </div>
    {/if}
  </div>
  
  <div class="flex-1 overflow-y-auto">
    {#if history.length === 0}
      <div class="p-6 text-center text-gray-400">
        <div class="text-3xl mb-3">📝</div>
        <p class="text-sm">No requests yet</p>
        <p class="text-xs mt-1">Your request history will appear here</p>
      </div>
    {:else}
      {#each history as item}
        {@const metadata = calculateMetadata(item)}
        {@const statusColor = getStatusBackground(item.statusCode)}
        {@const urlParts = formatUrlDisplay(item.url)}
        {@const debug = console.log('History item status:', item.statusCode, 'Color:', statusColor)}
        <button 
          class="w-full text-left p-4 border-b border-gray-700 hover:bg-gray-700 transition-colors duration-150 {selectedHistoryItem === item.id ? 'bg-blue-600 border-blue-500' : ''}"
          on:click={() => onSelectHistoryItem(item.id)}
        >
          <div class="flex items-center justify-between mb-2">
            <div class="flex items-center space-x-3 flex-1 min-w-0">
              <span class="px-2 py-1 text-xs font-semibold rounded {getMethodColor(item.method)} flex-shrink-0">
                {item.method}
              </span>
              <div class="flex-1 min-w-0" title={item.url}>
                <div class="text-sm text-gray-200 font-medium truncate">
                  {(urlParts.path !== '/' && urlParts.path !== '') ? urlParts.path : urlParts.domain}
                </div>
                <div class="text-xs text-gray-400 truncate">
                  {(urlParts.path !== '/' && urlParts.path !== '') ? urlParts.domain : ''}
                </div>
              </div>
            </div>
            <span class="px-2 py-1 text-xs font-mono rounded {statusColor} flex-shrink-0">
              {item.statusCode}
            </span>
          </div>
          
          <div class="flex items-center justify-between text-xs text-gray-400">
            <span class="truncate">{new Date(item.createdAt).toLocaleTimeString()}</span>
            <div class="flex items-center space-x-2">
              {#if metadata.headerCount > 0}
                <span class="text-xs" title="{metadata.headerCount} headers">📄 {metadata.headerCount}</span>
              {/if}
              {#if item.responseBody}
                <span class="text-xs" title="Has response data">📋</span>
              {/if}
            </div>
          </div>
        </button>
      {/each}
    {/if}
  </div>
</div> 
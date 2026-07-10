document.addEventListener('DOMContentLoaded', async () => {
    const meta = document.querySelector('meta[name="api-key"]');
    if (!meta) { console.error('Missing API Key meta tag'); return; }
    const apiKey = meta.content;

    const res = await fetch('/api/jobs', {
        headers: {
            'x-api-key': apiKey
        }
    })
    if(!res.ok) throw new Error(`HTTP ${res.status}`);
    const jobs = await res.json();

    function formatDuration(ms) {
        if (ms<1000) return ms + 'ms';
        return (ms/1000).toFixed(1) + 's';
    }

    function truncate(str, max) {
        if(!str || str.length <= max) return str || '';
        return str.slice(0, max) + '...';
    }

    if (jobs.length === 0) {
        document.getElementById('app').innerHTML = '<p>No jobs found</p>';
        return;
    }

    let html = `<div class="table-container"><table class="table">
            <thead><tr>
                <th>Name</th><th>Status</th><th>Duration</th><th>Created At</th><th>Error</th>
            </tr></thead><tbody>`;
    
    for (const job of jobs) {
        const badgeClass = `status-badge status-badge--${job.status}`;
        const duration = formatDuration(job.duration);
        const createdAt = new Date(job.created_at).toLocaleString();
        const error = truncate(job.error, 50);

        html += `<tr>
            <td>${job.name}</td>
            <td><span class="${badgeClass}">${job.status}</span></td>
            <td>${duration}</td>
            <td>${createdAt}</td>
            <td title="${job.error || ''}">${error}</td>
            </tr>`
    }

    html += '</tbody></table></div>';
    document.getElementById('app').innerHTML = html;

});

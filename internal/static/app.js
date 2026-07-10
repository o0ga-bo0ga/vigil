document.addEventListener('DOMContentLoaded', async () => {
    const apiKey = readMetaKey();
    if (!apiKey) {
        document.getElementById('app').innerHTML = '<p>Missing API key</p>';
        return;
    }

    renderFilters();
    await renderStats(apiKey);
    await renderJobs(apiKey);

    const statusSelect = document.getElementById('status-filter');
    const tenantInput = document.getElementById('tenant-filter');

    statusSelect.addEventListener('change', () => {
        renderJobs(apiKey, statusSelect.value, tenantInput.value);
    });

    tenantInput.addEventListener('keydown', (e) => {
        if (e.key === "Enter") {
            renderJobs(apiKey, statusSelect.value, tenantInput.value);
        }
    });

});

function formatDuration(ms) {
    if (ms<1000) return ms + 'ms';
    return (ms/1000).toFixed(1) + 's';
}

function truncate(str, max) {
    if(!str || str.length <= max) return str || '';
    return str.slice(0, max) + '...';
}

function readMetaKey() {
    const meta = document.querySelector('meta[name="api-key"]');
    if (!meta) { console.error('Missing API Key meta tag'); return; }
    return meta.content;
}

async function renderStats(apiKey) {
    try {
        const res = await fetch('/api/stats', {
            headers: {
                'x-api-key': apiKey
            }
        });
        if(!res.ok) throw new Error(`HTTP ${res.status}`);
        const stats = await res.json();
        const html = `
            <div class="stats-bar">
                <div class="stat-card">
                    <div class="stat-card__value">${stats.total}</div>
                    <div class="stat-card__label">Total</div>
                </div>
                <div class="stat-card">
                    <div class="stat-card__value" style="color: #22c55e">${stats.succeeded}</div>
                    <div class="stat-card__label">Succeeded</div>
                </div>
                <div class="stat-card">
                    <div class="stat-card__value" style="color: #ef4444">${stats.failed}</div>
                    <div class="stat-card__label">Failed</div>
                </div>
                <div class="stat-card">
                    <div class="stat-card__value">${formatDuration(Math.round(stats.avgDuration))}</div>
                    <div class="stat-card__label">Avg Duration</div>
                </div>
            </div>`;
        document.getElementById('stats').innerHTML = html
    } catch (e) {
        document.getElementById('stats').innerHTML = `<p>Failed to load stats: ${e.message}</p>`;
    }

}

async function renderJobs(apiKey, status, tenant) {
    const params = new URLSearchParams();
    if (status) params.set('status', status);
    if (tenant) params.set('tenant', tenant);

    try {
        const res = await fetch(`/api/jobs?${params}`, {
                headers: {
                    'x-api-key': apiKey
                }
            });
    
        if(!res.ok) throw new Error(`HTTP ${res.status}`);
        const jobs = await res.json();
    
        if (jobs.length === 0) {
            document.getElementById('jobs').innerHTML = '<p>No jobs found</p>';
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
        document.getElementById('jobs').innerHTML = html;
    } catch (e) {
        document.getElementById('jobs').innerHTML = `<p>Failed to load jobs: ${e.message}</p>`;
    }
}

function renderFilters() {
    let html = `
        <div class="filter-bar">
            <select id="status-filter">
                <option value="">All</option>
                <option value="started">Started</option>
                <option value="succeeded">Succeeded</option>
                <option value="failed">Failed</option>
                <option value="retried">Retried</option>
            </select>
            <input id="tenant-filter" type="text" placeholder="Filter by tenant...">
        </div>`;
    document.getElementById('filters').innerHTML = html;
}


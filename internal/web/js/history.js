let searchHistory = [];
const MAX_HISTORY_ITEMS = 20;

function loadHistory() {
    const savedHistory = localStorage.getItem('searchHistory');
    if (savedHistory) {
        try {
            searchHistory = JSON.parse(savedHistory) || [];
        } catch (e) {
            console.error("Ошибка загрузки истории:", e);
            searchHistory = [];
        }
    }
}

function saveHistory() {
    localStorage.setItem('searchHistory', JSON.stringify(searchHistory));
}

function clearHistory() {
    if (confirm('Очистить всю историю поиска?')) {
        searchHistory = [];
        saveHistory();
        renderHistory();
    }
}

function renderHistory() {
    const historyList = document.getElementById('historyList');
    if (!historyList) return;

    historyList.innerHTML = '';

    if (searchHistory.length === 0) {
        historyList.innerHTML = '<div>История запросов пуста</div>';
        return;
    }

    searchHistory.forEach((item, index) => {
        const time = new Date(item.timestamp).toLocaleTimeString();
        const div = document.createElement('div');
        div.className = 'history-item';
        div.innerHTML = `
            <span>${item.query}</span>
            <span class="time">${time}</span>
        `;
        div.onclick = () => loadFromHistory(index);
        historyList.appendChild(div);
    });
}

function loadFromHistory(index) {
    const item = searchHistory[index];
    document.getElementById('expression').value = item.query;
    document.getElementById('topN').value = item.topN || 20;
    performSearch(true);
}

function addToHistory(query, topN) {
    searchHistory = searchHistory.filter(item => item.query !== query);

    searchHistory.unshift({
        query: query,
        topN: topN,
        timestamp: Date.now()
    });

    if (searchHistory.length > MAX_HISTORY_ITEMS) {
        searchHistory.pop();
    }

    saveHistory();
    renderHistory();
}
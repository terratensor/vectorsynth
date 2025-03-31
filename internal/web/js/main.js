document.addEventListener('DOMContentLoaded', function () {
    initTheme();
    loadHistory();
    renderHistory();

    document.getElementById('search').addEventListener('click', function () {
        performSearch();
    });

    document.getElementById('expression').addEventListener('keypress', function (e) {
        if (e.key === 'Enter') {
            performSearch();
        }
    });

    window.onpopstate = function (event) {
        if (event.state) {
            document.getElementById('expression').value = event.state.query;
            document.getElementById('topN').value = event.state.topN || 20;
            performSearch(false);
        }
    };

    const params = new URLSearchParams(window.location.search);
    const query = params.get('q');
    const topN = params.get('n') || 20;

    if (query) {
        document.getElementById('expression').value = decodeURIComponent(query);
        document.getElementById('topN').value = topN;
        performSearch(false);
    }
});
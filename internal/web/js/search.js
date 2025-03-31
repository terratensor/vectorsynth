function performSearch(pushToHistory = true) {
    let expression = document.getElementById('expression').value.trim();
    expression = expression.toLowerCase();
    const topN = document.getElementById('topN').value;

    if (!expression) {
        document.getElementById('results').innerHTML =
            '<div class="error">Пожалуйста, введите выражение</div>';
        return;
    }

    if (pushToHistory) {
        addToHistory(expression, topN);

        const state = { query: expression, topN: topN };
        const url = `?q=${encodeURIComponent(expression)}&n=${topN}`;
        history.pushState(state, '', url);
    }

    showLoading();

    fetch('/api/similar', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            expression: expression,
            topN: parseInt(topN)
        })
    })
        .then(handleResponse)
        .then(displayResults)
        .catch(handleError);
}

function showLoading() {
    const resultsEl = document.getElementById('results');
    if (resultsEl) {
        resultsEl.innerHTML = '<div>Идёт поиск...</div>';
    }
}

function handleResponse(response) {
    if (!response.ok) {
        return response.text().then(text => {
            throw new Error(text || 'Неизвестная ошибка сервера')
        });
    }
    return response.json();
}

function handleError(error) {
    console.error('Ошибка поиска:', error);
    const resultsEl = document.getElementById('results');
    if (resultsEl) {
        resultsEl.innerHTML =
            `<div class="error">Ошибка: ${error.message || 'неизвестная ошибка'}</div>`;
    }
}

function displayResults(results) {
    const resultsDiv = document.getElementById('results');
    if (!resultsDiv) return;

    if (results.length === 0) {
        resultsDiv.innerHTML = '<div>Ничего не найдено</div>';
        return;
    }

    let html = '<h3>Результаты:</h3>';
    results.forEach((item, index) => {
        html += `
            <div class="result-item">
                ${index + 1}. 
                <span class="clickable-word" onclick="useWordInSearch('${item.word.replace(/'/g, "\\'")}')">
                    ${item.word}
                </span>
                <span class="similarity">(сходство: ${item.similarity.toFixed(4)})</span>
                <span class="word-actions">
                    <span class="word-action" onclick="addToSearch('+', '${item.word.replace(/'/g, "\\'")}')">+</span>
                    <span class="word-action" onclick="addToSearch('-', '${item.word.replace(/'/g, "\\'")}')">−</span>
                    <span class="word-action copy" title="Копировать" onclick="copyToClipboard('${item.word.replace(/'/g, "\\'")}')">⎘</span>
                </span>
            </div>
        `;
    });

    resultsDiv.innerHTML = html;
}

function useWordInSearch(word) {
    const input = document.getElementById('expression');
    if (input) {
        input.value = word;
        performSearch();
    }
}

function addToSearch(operator, word) {
    const input = document.getElementById('expression');
    if (!input) return;

    const currentValue = input.value.trim();

    if (currentValue === '') {
        input.value = word;
    } else {
        const lastChar = currentValue[currentValue.length - 1];
        if (lastChar === '+' || lastChar === '-') {
            input.value = `${currentValue} ${word}`;
        } else {
            input.value = `${currentValue} ${operator} ${word}`;
        }
    }

    input.focus();
}

function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(() => {
        const notification = document.createElement('div');
        notification.textContent = 'Скопировано!';
        notification.style.position = 'fixed';
        notification.style.bottom = '20px';
        notification.style.right = '20px';
        notification.style.backgroundColor = '#28a745';
        notification.style.color = 'white';
        notification.style.padding = '10px';
        notification.style.borderRadius = '5px';
        notification.style.zIndex = '1000';
        document.body.appendChild(notification);

        setTimeout(() => {
            document.body.removeChild(notification);
        }, 2000);
    }).catch(err => {
        console.error('Ошибка копирования: ', err);
    });
}
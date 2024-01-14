document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('reset-errors').addEventListener('click', function() {
        fetch('http://http://95.142.45.133:23873/avaryreset', {method: 'POST'})
            .then(response => response.text())
            .then(data => {
                alert('Результат сброса: ' + data);
            })
            .catch(error => console.error('Ошибка сброса ошибок: ', error));
    });
});

document.addEventListener('DOMContentLoaded', function() {
    const resetButton = document.getElementById('reset-errors');
    resetButton.addEventListener('click', function() {
        fetch('http://95.142.45.133:23873/avaryreset', {method: 'POST'})
            .then(response => response.text())
            .then(data => {
                alert('Результат сброса: ' + data);
            })
            .catch(error => console.error('Ошибка сброса ошибок: ', error));
    });

    // Инициализация данных при загрузке страницы
    initializeData();
    const boilerNames = [
        "Котельная «Склады Мищенко»",                   //0   кот№1 Склады Мищенко
        "Котельная «Выставка Ендальцева»",              //1   кот№2 Ендальцев         (датчик на базе)
        "Котельная «ЧукотОптТорг»",                     //2   кот№3 ЧукотОптТорг      (датчик на базе)
        "Котельная «ЧСБК новая»",                       //3   кот№4 "ЧСБК Новая"
        "Котельная «Офис СВТ»",                         //4   кот№5 офис "СВТ"
        "Котельная «Общежитие на Южной»",               //5   кот№6 общежитие на Южной
        "Котельная «Офис ЧСБК»",                        //6   кот№7 офис ЧСБК
        "Котельная «Рынок»",                            //7   кот№8 "Рынок"
        "Котельная «Макатровых»",                       //8   кот№9 Макатровых
        "Котельная ДС «Сказка»",                        //9   кот№10  "Д/С Сказка"
        "Котельная «Полярный»",                         //10  кот№11 Полярный
        "Котельная «Департамент»",                      //11  кот№12 Департамент
        "Котельная «Офис ЧСБК квартиры»",               //12  кот№13 квартиры в офисе
        "Котельная Шишкина"                             //13  кот№14 ТО Шишкина
    ];
    // Функция для инициализации данных
    function initializeData() {
        fetch('http://95.142.45.133:23873/getparams')
            .then(response => response.json())
            .then(boilers => {
                const container = document.querySelector('.boiler-room-container');
                container.innerHTML = ''; // Очищаем текущий контент

                boilers.forEach(boiler => {
                    const boilerName = boilerNames[boiler.id] || `Котельная ${boiler.id + 1}`;
                    const boilerDiv = document.createElement('div');
                    boilerDiv.className = 'boiler';
                    boilerDiv.setAttribute('data-boiler-id', boiler.id);

                    const boilerHeaderDiv = document.createElement('div');
                    boilerHeaderDiv.className = 'boiler-header';

                    const boilerImage = document.createElement('img');
                    boilerImage.src = `/static/images/boiler_icon_${boiler.id}.png`;
                    boilerImage.alt = `Изображение котельной ${boiler.id + 1}`;
                    boilerImage.className = 'boiler-image';

                    const boilerTitle = document.createElement('h2');
                    boilerTitle.textContent = boilerName;

                    boilerHeaderDiv.appendChild(boilerImage);
                    boilerHeaderDiv.appendChild(boilerTitle);

                    const parametersDiv = document.createElement('div');
                    parametersDiv.className = 'parameters';

                    const tPodCanvas = document.createElement('canvas');
                    tPodCanvas.id = `tPodChart${boiler.id}`;
                    tPodCanvas.width = 400;
                    tPodCanvas.height = 200;

                    const pPodCanvas = document.createElement('canvas');
                    pPodCanvas.id = `pPodChart${boiler.id}`;
                    pPodCanvas.width = 400;
                    pPodCanvas.height = 200;

                    boilerDiv.appendChild(boilerHeaderDiv);
                    boilerDiv.appendChild(parametersDiv);
                    boilerDiv.appendChild(tPodCanvas);
                    boilerDiv.appendChild(pPodCanvas);

                    container.appendChild(boilerDiv);

                    // Запрос данных и создание графиков для каждого котла
                    loadAndCreateChart(boiler.id, `tPodChart${boiler.id}`, 'Temperature (tPod)', 'tpod');
                    loadAndCreateChart(boiler.id, `pPodChart${boiler.id}`, 'Pressure (pPod)', 'ppod');
                });
            })
            .then(() => {
                updateData(); // Обновляем данные после инициализации
            })
            .catch(error => console.error('Ошибка инициализации данных: ', error));
    }
    function createChart(ctx, rawData, label) {
        const now = new Date();
        const data = rawData.filter((_, index) => index % 2 === 0);

        const totalDuration = 16 * 60 * 60 * 1000; // Общая продолжительность данных в миллисекундах (16 часов)
        const startTime = new Date(now.getTime() - totalDuration); // Начальное время для первой точки данных

        const timeLabels = data.map((_, index) => {
            // Время для каждой точки данных
            const date = new Date(startTime.getTime() + (index * (totalDuration / data.length)));
            return moment(date).format('HH:mm');
        });

        return new Chart(ctx, {
            type: 'line',
            data: {
                labels: timeLabels,
                datasets: [{
                    label: label,
                    data: data,
                    fill: false,
                    borderColor: 'rgb(75, 192, 192)',
                    tension: 0.1,
                    pointRadius: 0
                }]
            },
            options: {
                scales: {
                    xAxes: [{
                        type: 'time',
                        time: {
                            parser: 'HH:mm',
                            unit: 'hour',
                            displayFormats: {
                                hour: 'HH:mm'
                            },
                        }
                    }]
                }
            }
        });
    }

// Функция для запроса данных и создания графика
    function loadAndCreateChart(boilerId, canvasId, label, dataKey) {
        fetch(`http://95.142.45.133:23874/boilers/get${dataKey}/${boilerId}`)
            .then(response => response.json())
            .then(data => {
                const ctx = document.getElementById(canvasId).getContext('2d');
                createChart(ctx, data, label);
            })
            .catch(error => console.error('Ошибка загрузки данных: ', error));
    }
    // Функция для обновления данных
    function updateData() {
        fetch('http://95.142.45.133:23873/getparams')
            .then(response => response.json())
            .then(boilers => {
                boilers.forEach(boiler => {
                    const boilerDiv = document.querySelector(`.boiler[data-boiler-id="${boiler.id}"]`);
                    if (boilerDiv) {
                        const parametersDiv = boilerDiv.querySelector('.parameters');
                        parametersDiv.innerHTML = `
                            <span>Температура подачи: ${boiler.tPod} °C</span>
                            <span>Давление подачи: ${boiler.pPod} МПа</span>
                            <span>Температура улицы: ${boiler.tUlica} °C</span>
                            <!-- Добавьте другие параметры по мере необходимости -->
                        `;
                    }
                });
            })
            .catch(error => console.error('Ошибка обновления данных: ', error));
    }

    // Установка интервала для обновления данных
    setInterval(updateData, 3000); // Обновление каждые 3 секунды
});
const chart2Ctx = document.getElementById("chart2").getContext("2d");
const chart2 = new Chart(chart2Ctx, {
    type: 'bar',
    data: {
        labels: ['Time 1', 'Time 2', 'Time 3'],
        datasets: [{
            label: 'Pre-Filtered Data',
            data: [30, 50, 70],
            backgroundColor: 'rgba(255, 99, 132, 0.5)',
            borderColor: 'rgba(255, 99, 132, 1)',
            borderWidth: 1
        }]
    },
    options: {
        responsive: true,
        plugins: {
            legend: { display: true },
        }
    }
});

const chart3Ctx = document.getElementById("chart3").getContext("2d");
const chart3 = new Chart(chart3Ctx, {
    type: 'bar',
    data: {
        labels: ['Time 1', 'Time 2', 'Time 3'],
        datasets: [{
            label: 'Post-Filtered Data',
            data: [28, 48, 65],
            backgroundColor: 'rgba(75, 192, 192, 0.5)',
            borderColor: 'rgba(75, 192, 192, 1)',
            borderWidth: 1
        }]
    },
    options: {
        responsive: true,
        plugins: {
            legend: { display: true },
        }
    }
});
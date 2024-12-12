const chart1Ctx = document.getElementById("chart1").getContext("2d");
const chart1 = new Chart(chart1Ctx, {
    type: 'bar',
    data: {
        labels: ['10 ml', '20 ml', '30 ml', '40 ml', '50 ml'],
        datasets: [{
            label: 'Milivolt',
            data: [100, 200, 300, 400, 500],
            backgroundColor: 'rgba(54, 162, 235, 0.5)',
            borderColor: 'rgba(54, 162, 235, 1)',
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

// Add table data dynamically
const table = document.querySelector("table");
for (let i = 1; i <= 26; i++) {
    const row = table.insertRow();
    for (let j = 0; j < 7; j++) {
        const cell = row.insertCell();
        cell.textContent = `Row ${i} Col ${j + 1}`;
    }
}
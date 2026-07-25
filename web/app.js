let allocationChart = null;
let gainsLossesChart = null;

// Initialize on page load
document.addEventListener('DOMContentLoaded', () => {
	loadPortfolio();
});

// Tabs
function showTab(tabName) {
	document.querySelectorAll('.tab-content').forEach(tab => {
		tab.classList.remove('active');
	});
	document.querySelectorAll('.tab-btn').forEach(btn => {
		btn.classList.remove('active');
	});

	document.getElementById(`${tabName}-tab`).classList.add('active');
	event.target.closest('.tab-btn').classList.add('active');
}

// Load portfolio
async function loadPortfolio() {
	try {
		const portfolio = await getPortfolio();
		displayPortfolio(portfolio);
		displayAllocationChart(portfolio.holdings);
		displayGainsLossesChart(portfolio.holdings);
	} catch (error) {
		showToast('Failed to load portfolio: ' + error.message, 'error');
	}
}

// Display portfolio
function displayPortfolio(portfolio) {
	const { holdings, total_cost, total_value, total_gain_loss, gain_loss_percent } = portfolio;

	// Update summary
	document.getElementById('totalCost').textContent = formatCurrency(total_cost);
	document.getElementById('totalValue').textContent = formatCurrency(total_value);

	const gainLossEl = document.getElementById('totalGainLoss');
	gainLossEl.textContent = formatCurrency(total_gain_loss);
	gainLossEl.classList.toggle('positive', total_gain_loss >= 0);
	gainLossEl.classList.toggle('negative', total_gain_loss < 0);

	const gainLossPercentEl = document.getElementById('gainLossPercent');
	gainLossPercentEl.textContent = formatPercent(gain_loss_percent);
	gainLossPercentEl.classList.toggle('positive', gain_loss_percent >= 0);
	gainLossPercentEl.classList.toggle('negative', gain_loss_percent < 0);

	// Update table
	const tbody = document.getElementById('holdingsBody');
	tbody.innerHTML = '';

	holdings.forEach(holding => {
		const row = document.createElement('tr');
		row.innerHTML = `
			<td><strong>${holding.symbol}</strong></td>
			<td>${holding.quantity.toFixed(2)}</td>
			<td>${formatCurrency(holding.cost_per_share)}</td>
			<td>${formatCurrency(holding.current_price)}</td>
			<td>${formatCurrency(holding.total_cost)}</td>
			<td>${formatCurrency(holding.current_value)}</td>
			<td class="${holding.gain_loss >= 0 ? 'positive' : 'negative'}">
				${formatCurrency(holding.gain_loss)}
			</td>
			<td class="${holding.gain_loss_percent >= 0 ? 'positive' : 'negative'}">
				${holding.gain_loss_percent.toFixed(2)}%
			</td>
			<td>
				<button class="btn btn-danger" onclick="removeHolding('${holding.symbol}')">Remove</button>
			</td>
		`;
		tbody.appendChild(row);
	});

	document.getElementById('loading').style.display = 'none';
}

// Allocation Pie Chart
function displayAllocationChart(holdings) {
	const ctx = document.getElementById('allocationChart').getContext('2d');
	
	const labels = holdings.map(h => h.symbol);
	const data = holdings.map(h => h.current_value);
	const colors = [
		'#3b82f6', '#10b981', '#f59e0b', '#ef4444',
		'#8b5cf6', '#ec4899', '#06b6d4', '#f97316'
	];

	if (allocationChart) allocationChart.destroy();

	allocationChart = new Chart(ctx, {
		type: 'doughnut',
		data: {
			labels: labels,
			datasets: [{
				data: data,
				backgroundColor: colors.slice(0, labels.length),
				borderColor: '#ffffff',
				borderWidth: 2
			}]
		},
		options: {
			responsive: true,
			maintainAspectRatio: false,
			plugins: {
				legend: {
					position: 'bottom',
					labels: {
						font: { size: 12, weight: '500' },
						padding: 15,
						usePointStyle: true
					}
				},
				tooltip: {
					callbacks: {
						label: function(context) {
							const total = context.dataset.data.reduce((a, b) => a + b, 0);
							const percentage = ((context.parsed / total) * 100).toFixed(1);
							return context.label + ': ' + formatCurrency(context.parsed) + ' (' + percentage + '%)';
						}
					}
				}
			}
		}
	});
}

// Gains & Losses Chart
function displayGainsLossesChart(holdings) {
	const ctx = document.getElementById('gainsLossesChart').getContext('2d');
	
	const labels = holdings.map(h => h.symbol);
	const gainLossData = holdings.map(h => h.gain_loss);
	const colors = holdings.map(h => h.gain_loss >= 0 ? '#10b981' : '#ef4444');

	if (gainsLossesChart) gainsLossesChart.destroy();

	gainsLossesChart = new Chart(ctx, {
		type: 'bar',
		data: {
			labels: labels,
			datasets: [{
				label: 'Gain/Loss ($)',
				data: gainLossData,
				backgroundColor: colors,
				borderRadius: 6,
				borderSkipped: false
			}]
		},
		options: {
			indexAxis: 'y',
			responsive: true,
			maintainAspectRatio: false,
			plugins: {
				legend: {
					display: true,
					labels: {
						font: { size: 12, weight: '500' }
					}
				},
				tooltip: {
					callbacks: {
						label: function(context) {
							return formatCurrency(context.parsed.x);
						}
					}
				}
			},
			scales: {
				x: {
					beginAtZero: true,
					ticks: {
						callback: function(value) {
							return formatCurrency(value);
						}
					}
				}
			}
		}
	});
}

// Add holding
async function addHolding(event) {
	event.preventDefault();

	const symbol = document.getElementById('addSymbol').value;
	const quantity = document.getElementById('addQuantity').value;
	const costBasis = document.getElementById('addCostBasis').value;

	try {
		await addHoldingAPI(symbol, quantity, costBasis);
		showToast('Holding added successfully!', 'success');
		document.getElementById('addSymbol').value = '';
		document.getElementById('addQuantity').value = '';
		document.getElementById('addCostBasis').value = '';
		loadPortfolio();
	} catch (error) {
		showToast('Failed to add holding: ' + error.message, 'error');
	}
}

// Remove holding
async function removeHolding(symbol) {
	if (!confirm(`Remove ${symbol} from portfolio?`)) {
		return;
	}

	try {
		await removeHoldingAPI(symbol);
		showToast('Holding removed successfully!', 'success');
		loadPortfolio();
	} catch (error) {
		showToast('Failed to remove holding: ' + error.message, 'error');
	}
}

// Search quote
async function searchQuote() {
	const symbol = document.getElementById('quoteSymbol').value;

	if (!symbol) {
		showToast('Please enter a symbol', 'error');
		return;
	}

	try {
		const quote = await getQuote(symbol);
		displayQuote(quote);
	} catch (error) {
		showToast('Failed to fetch quote: ' + error.message, 'error');
	}
}

// Display quote
function displayQuote(quote) {
	const result = document.getElementById('quoteResult');
	const changeClass = quote.change >= 0 ? 'positive' : 'negative';

	result.innerHTML = `
		<h3>${quote.symbol}</h3>
		<div class="quote-info">
			<div class="quote-item">
				<label>Price</label>
				<div class="value">${formatCurrency(quote.price)}</div>
			</div>
			<div class="quote-item">
				<label>Change</label>
				<div class="value ${changeClass}">${formatCurrency(quote.change)}</div>
			</div>
			<div class="quote-item">
				<label>Change %</label>
				<div class="value ${changeClass}">${quote.change_percent.toFixed(2)}%</div>
			</div>
			<div class="quote-item">
				<label>Volume</label>
				<div class="value">${quote.volume.toLocaleString()}</div>
			</div>
		</div>
	`;
	result.classList.add('show');
}

// Refresh prices
async function refreshPrices() {
	try {
		await refreshPricesAPI();
		showToast('Prices refreshed!', 'success');
		loadPortfolio();
	} catch (error) {
		showToast('Failed to refresh prices: ' + error.message, 'error');
	}
}

// Utilities
function formatCurrency(value) {
	return new Intl.NumberFormat('en-US', {
		style: 'currency',
		currency: 'USD',
	}).format(value);
}

function formatPercent(value) {
	return value.toFixed(2) + '%';
}

function showToast(message, type = 'info') {
	const toast = document.getElementById('toast');
	toast.textContent = message;
	toast.className = `toast ${type} show`;
	
	setTimeout(() => {
		toast.classList.remove('show');
	}, 3000);
}

const API_BASE = '/api';

async function apiCall(endpoint, options = {}) {
	const url = `${API_BASE}${endpoint}`;
	const defaultOptions = {
		headers: {
			'Content-Type': 'application/json',
		},
	};

	const config = { ...defaultOptions, ...options };

	try {
		const response = await fetch(url, config);
		const data = await response.json();

		if (!response.ok) {
			throw new Error(data.error || 'API error');
		}

		return data.data;
	} catch (error) {
		console.error('API Error:', error);
		throw error;
	}
}

// Portfolio API calls
async function getPortfolio() {
	return await apiCall('/portfolio');
}

async function addHoldingAPI(symbol, quantity, costBasis) {
	return await apiCall('/holdings', {
		method: 'POST',
		body: JSON.stringify({
			symbol,
			quantity: parseFloat(quantity),
			cost_basis: parseFloat(costBasis),
		}),
	});
}

async function removeHoldingAPI(symbol) {
	return await apiCall('/holdings', {
		method: 'DELETE',
		body: JSON.stringify({ symbol }),
	});
}

async function getQuote(symbol) {
	return await apiCall(`/quote?symbol=${symbol}`);
}

async function refreshPricesAPI() {
	return await apiCall('/prices/refresh', { method: 'POST' });
}
var results = {
	'headers': [
		{'key': 'id',           'display_header': 'Name'},
		{'key': 'content',      'display_header': 'Content'},
		{'key': 'signal',       'display_header': 'Signal'},
		{'key': 'total_btc', 'display_header': 'Bitcoin'},
	],
	'signals': [
	{
		'id': 'monolith',
		'content': 'If bitcoin could speak, what would it say?',
		'signal': '.002 btc/byte',
		'total_btc': '0.0011',
	},
	],
}

function genHeaders(headers) {
	var headerHTML = '';
	for (header of headers) {
		headerHTML += '<th>' + header.display_header + '</th>\n';
	}
	document.getElementById('tbl-header-tr').innerHTML = headerHTML;
	console.log(headerHTML);
}

function genRow(row, headers) {
	// follow same order specified for headers
	var rowHTML = '';
	for (header of headers) {
		rowHTML += '<td>' + row[header.key] + '</td>';
	}
	return rowHTML
}

function genRows(headers, rows) {
	var modalTemplateHtml = '';
	for (row of rows) {
		document.getElementById('tbl-body-content').innerHTML += '<tr class="trigger-modal id="' + row.id + '>' + genRow(row, headers) + '</tr>';
	}
}

function genTable() {
	genHeaders(results.headers);
	genRows(results.headers, results.signals);
}

genTable()

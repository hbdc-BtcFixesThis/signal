var recordTableID = 'record-tbl-body-content';
var recordTableHeaderID = 'record-tbl-header-tr';
var results = {
	'headers': [
		{'key': 'name',           'display_header': 'Key'},
		// {'key': 'value',         'display_header': 'Value'},
	],
};

function genHeaders(headers) {
	var headerHTML = '';
	for (header of headers) {
		headerHTML += '<th>' + header.display_header + '</th>\n';
	}
	document.getElementById(recordTableHeaderID).innerHTML = headerHTML;
	console.log(headerHTML);
}

function genRow(row, headers) {
	// follow same order specified for headers
	var rowHTML = ''
	for (header of headers) {
		rowHTML += '<td>' + row[header.key] + '</td>';
	}
	return rowHTML
}

function genRows(headers, rows) {
	// var modalTemplateHtml = '';
	var tbl = document.getElementById(recordTableID);
	for (row of rows) {
		tbl.innerHTML += '<div class="top-left">' +
					  		'sats: ' + row.sats +
				  		 '</div>' +
						 '<div class="top-left">' +
					  		'vBytes: ' + row.vbytes +
				  		 '</div>' +
						 '<tr class="trigger-modal" id="' + row[headers[0].key]+ '">' +
							genRow(row, headers) +
						 '</tr>';
	}
}

function  clearTableRows() {
	var tbl = document.getElementById(recordTableID);
	var tableRows = tbl.getElementsByTagName('tr');
	var rowCount = tableRows.length;

	for (var i = rowCount - 1; i >= 0; i--) {
		tbl.deleteRow(i);
	}
}

function successfullyRetrievedPage(xhr) {
    // removeLoadingClass(submitNewRecord);
    // toggleCreateRecordModal();
	resp = JSON.parse(xhr.response);
	results.signals = resp;
	clearTableRows()
	genRows(results.headers, results.signals);
    // showSuccessBanner(xhr.responseText)
}

function failedToRetrievePage(xhr) {
    // removeLoadingClass(submitNewRecord);
    showErrorBanner(xhr.responseText);
} 

function genTable() {
	genHeaders(results.headers);
	results.signals = sendJsonPost(
		routes.getPage, "GET", null,
		successfullyRetrievedPage, failedToRetrievePage,
	);
	// genRows(results.headers, results.signals);
}

genTable()

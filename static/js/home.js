var recordTableID = 'record-tbl-body-content';
var recordDetailsID = 'record-details';

// used for signatures/sigals and preview formatting
var detailsDialog = document.getElementById(recordDetailsID);
document.addEventListener('click', ({target}) => target === detailsDialog && detailsDialog.close());

function failedToRetrieve(xhr) {
    showErrorBanner(xhr.responseText);
} 

function satsOrBtcRounding(amount) {
	var retAmt = amount;
	var unit = 'BTC';
	if (amount > 10000000) {
		retAmt = (amount/100000000);
	} else if (amount > 1000000) {
		retAmt = (amount/100000000);
	} else if (amount > 100000) {
		retAmt = (amount/100000000);
	} else {
		unit = 'SATS';
	}
	return [parseFloat(retAmt.toFixed(3)), unit];
}

function formatBytes(bytes, decimals = 2) {
    if (!+bytes) return '0 Bytes';

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];

    const i = Math.floor(Math.log(bytes) / Math.log(k));	

    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
}

function genRows(rows) {
	// var modalTemplateHtml = '';
	var tbl = document.getElementById(recordTableID);
	for (row of rows) {
		var [fmattedBtcAmnt, unitsBtc] = satsOrBtcRounding(row.sats);
		fmattedBtcPerByte = satsOrBtcRounding(row.sats/row.vbytes);
		fmattedBtc = satsOrBtcRounding(row.sats);
		tbl.innerHTML += '<tr id="' + row.rid + '">' +
							'<td class="w40 overflow-x-scroll">' +
								row.name +
								'<div class="w100">' +
								'<button id="' + row.rid + '" ' +
									    'class="previewIcon orange-button "' +
									    'onclick="showRecordValue(this)"></button>' +
								'<button id="' + row.rid + '" ' +
									    'class="bitcoinIcon orange-button "' +
									    'onclick="showRecordSignals(this)"></button>' +
								'<button id="' + row.rid + '" ' +
									    'class="signatureIcon orange-button "' +
									    'onclick="showRecordDetails(this)"></button>' +
								"</div>" +
							'</td>' +
							'<td class="record-stats-text">' +
							 	'<div>Value</div><div>Size</div><div>signal</div><div>Signal count</div>' +
							'</td>' +
							'<td class="record-stats-text">' +
							 	'<div>' + fmattedBtc[0] + ' ' + fmattedBtc[1] +  '</div>' +
							 	'<div>' + formatBytes(row.vbytes) + '</div>' +
							 	'<div><b>' + fmattedBtcPerByte[0] + '</b> ' +
											 fmattedBtcPerByte[1] + '/B</div>' +
							 	'<div>' + row.sids.length + '</div>' +
							'</td>' +
						 '</tr>';
	}
}

function successfullyRetrievedRecordValue(xhr) {
	var response = JSON.parse(xhr.response);
	detailsDialog = document.getElementById(recordDetailsID);
	detailsDialog.innerHTML ='<pre class="grey-text" style="text-wrap: balance;">' + response.value + "</pre>";
	detailsDialog.showModal();
}

function successfullyRetrievedRecordSignals(xhr) {
	var response = JSON.parse(xhr.response);
	detailsDialog = document.getElementById(recordDetailsID);
	
	var rows = '';
	for (signal of response) {
		var [fmattedBtcAmnt, unitsBtc] = satsOrBtcRounding(signal.sats);
		rows += '<tr><td>' + signal.btc_address + '</td><td>' +
							 fmattedBtcAmnt + unitsBtc + '</td><td>' +
							 signal.signature +
				'</td></tr>'
	}
	detailsDialog.innerHTML = `<div style="min-width: 25rem;">
        <table cellpadding="0" cellspacing="0" border="0" class="w100">
          <tbody id="record-details-tbl-body-content">
		  <tr><th>BTC Address</th><th>Value</th><th>Signature</th></tr>`
		  	+ rows +
          `</tbody>
        </table>
      </div>`;
	detailsDialog.showModal();
}

function showRecordValue(e) {
	console.log(e.id)
	sendJsonPost(
		routes.recordValue+'?rid='+e.id, "GET", null,
		successfullyRetrievedRecordValue, failedToRetrieve,
	);
}

function showRecordSignals(e) {
	console.log(e.id)
	sendJsonPost(
		routes.getRecordSignals+'?rid='+e.id, "GET", null,
		successfullyRetrievedRecordSignals, failedToRetrieve,
	);
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
	clearTableRows()
	genRows(JSON.parse(xhr.response));
	initIcons();
    // showSuccessBanner(xhr.responseText)
}

function genTable() {
	// results.signals = 
	sendJsonPost(
		routes.getPage, "GET", null,
		successfullyRetrievedPage, failedToRetrieve,
	);
	// genRows(results.headers, results.signals);
}

genTable();

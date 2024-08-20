var recordTableID = 'record-tbl-body-content';
var recordDetailsID = 'record-details';

function satsOrBtcRounding(amount) {
	var retAmt = amount.toFixed(5);
	var unit = 'btc';
	if (amount > 10000000) {
		retAmt = (amount/100000000).toFixed(2);
	} else if (amount > 1000000) {
		retAmt = (amount/100000000).toFixed(3);
	} else if (amount > 100000) {
		retAmt = (amount/100000000).toFixed(4);
	} else {
		unit = 'sats';
	}
	return [parseFloat(retAmt), unit];
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
							'<td class="w40 overflow-x-scroll">' + row.name + '</td>' +
							'<td class="record-stats-text">' +
							 	'<div>' + fmattedBtcPerByte[1] + '/byte</div>' +
								'<div>' + fmattedBtc[1] + ' </div><div>size </div><div>signals </div>' +
								'<button id="' + row.rid + '" ' +
									    'class="previewIcon orange-button w100 tbl-btn-spacing"' +
									    'onclick="showRecordDetails(this)"></button>' +
							'</td>' +
							'<td class="record-stats-text">' +
							 	'<div><h2>' + fmattedBtcPerByte[0] + '</h2></div>' +
							 	'<div>' + fmattedBtc[0] + '</div>' +
							 	'<div>' + formatBytes(row.vbytes) + '</div>' +
							 	'<div>' + row.sids.length + '</div>' +
								'<button id="' + row.rid + '" ' +
									    'class="bitcoinIcon orange-button w100 tbl-btn-spacing"' +
									    'onclick="showRecordDetails(this)"></button>' +
							'</td>' +
						 '</tr>';
	}
}

// used for signatures/sigals and preview formatting
var detailsDialog = document.getElementById(recordDetailsID);
document.addEventListener('click', ({target}) => target === detailsDialog && detailsDialog.close());

function showRecordDetails(e) {
	console.log(e.id)
	detailsDialog = document.getElementById(recordDetailsID);
	detailsDialog.innerHTML = `<div>
        <table cellpadding="0" cellspacing="0" border="0" class="w100">
          <tbody id="record-details-tbl-body-content">
		  	<tr><td>address</td><td>21</td></tr>
          </tbody>
        </table>
      </div>`;
	detailsDialog.showModal();
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

function failedToRetrievePage(xhr) {
    // removeLoadingClass(submitNewRecord);
    showErrorBanner(xhr.responseText);
} 

function genTable() {
	// results.signals = 
	sendJsonPost(
		routes.getPage, "GET", null,
		successfullyRetrievedPage, failedToRetrievePage,
	);
	// genRows(results.headers, results.signals);
}

genTable();

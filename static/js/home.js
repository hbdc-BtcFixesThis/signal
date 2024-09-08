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
	var satsPerBtc = 100000000;
	var unit = '₿';
	if ((amount / satsPerBtc) > 1) {
		retAmt = (amount/satsPerBtc);
	} else if ((amount / satsPerBtc) > 0.001) {
		unit = 'm₿';
		retAmt = (amount/satsPerBtc)/.001;
	} else if ((amount / satsPerBtc) > 0.000001) {
		unit = 'μ₿';
		retAmt = (amount/satsPerBtc)/0.000001;
	} else if ((amount / satsPerBtc) > 0.00000001) {
		unit = 'SATS';
	} else {
		unit = 'mSATS';
		retAmt = (amount/satsPerBtc)/0.00000000001;
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
		tbl.innerHTML += `<tr id="${row.rid}">
							<td class="w40 overflow-x-scroll">${row.name}
								<div class="w100">
								<button id="${row.rid}"
										title="Preview record value"
								       	class="previewIcon orange-button"
								       	onclick="showRecordValue(this)"></button>
								<button id="${row.rid}"
										title="Sats signed for"
								       	class="bitcoinIcon orange-button"
								       	onclick="showRecordSignals(this)"></button>
								<button id="${row.rid}"
										title="Sign for record"
								       	class="signatureIcon orange-button"
								       	onclick="toggleRecordSignalModal(this)"></button>
								<button id="${row.rid}"
										title="Record ID"
									    class="fingerPrintIcon orange-button"
									    onclick="showRecordId(this)"></button>
								</div>
							</td>
							<td class="record-stats-text">
								<div>Value</div><div>Size</div><div>signal</div><div>Signal count</div>
							</td>
							<td class="record-stats-text">
							 	<div title="${row.sats.toLocaleString('en', {useGrouping:true})} Sats">
									${fmattedBtc[0]} ${fmattedBtc[1]}
								</div>
							 	<div>${formatBytes(row.vbytes)}</div>
							 	<div><b>${fmattedBtcPerByte[0]}</b> ${fmattedBtcPerByte[1]}/B</div>
							 	<div>${row.sids.length}</div>
							</td>
						 </tr>`;
	}
}

function showRecordId(e) {
	document.getElementById(recordDetailsID).innerHTML = '<pre class="grey-text" style="text-wrap: balance;">' + e.id + "</pre>";
	detailsDialog.showModal();
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
		rows += `<tr>
					<td title="₿ address">${signal.btc_address}</td>
					<td title="${signal.sats.toLocaleString('en', {useGrouping:true})} SATS">
						${fmattedBtcAmnt + unitsBtc}
					</td>
					<td title="₿ signature">${signal.signature}</td>
				</tr>`
	}
	detailsDialog.innerHTML = `<div style="min-width: 25rem;">
        <table cellpadding="0" cellspacing="0" border="0" class="w100">
          <tbody id="record-details-tbl-body-content">
		  <tr><th>BTC Address</th><th>Value</th><th>Signature</th></tr>
		  	${rows}
          </tbody>
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

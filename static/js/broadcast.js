let broadcast = document.forms.broadcast;

let createRecordModal = document.getElementById("record-modal");
let createRecordTrigger = document.getElementById("create-record-trigger");
let cancelCreateTrigger = document.getElementById("cancel-create-trigger");
let submitNewRecord = document.getElementById('submit-create-record');

let signalSignaturerMessage = document.getElementById('signal-signature-message')

let newSignalSignatureMessage = ''

let recordState = {
	rid: '',
	name: document.getElementById("new-record-name"),
	content: document.getElementById("new-record-content"),
	address: document.getElementById("signal-wallet-address"),
	numSats: document.getElementById('new-record-signal-sats'),
	signature: document.getElementById("signal-signature"),
	signatureMessage: document.getElementById('signal-signature-message'),
};

function clearRecordState() {
	recordState.rid = '';
	recordState.name.value = '';
	recordState.content.value = '';
	recordState.address.value = '';
	recordState.numSats.value = null;
	recordState.signature.value = '';
	recordState.signatureMessage.value = '';
}

function toggleCreateRecordModal() {
	clearRecordState();
	recordState.name.style.display = 'block';
	recordState.content.style.display = 'block';
	toggleShowModal(createRecordModal);
}

function toggleRecordSignalModal(e) {
	recordState.rid = e.id;
	recordState.name.style.display = 'none';
	recordState.content.style.display = 'none';
	toggleShowModal(createRecordModal);
}

function successfullyAddedRecord(xhr) {
    removeLoadingClass(submitNewRecord);
	toggleCreateRecordModal();
	showSuccessBanner(xhr.responseText)
	clearRecordState();
}

function failedToAddRecord(xhr) {
	removeLoadingClass(submitNewRecord);
	showErrorBanner(xhr.responseText);
}

function updateNewSignalSignatureMessage(recordId, numSats, btcAddr) {
	var newMessageToSign = newSignalSignatureMessage.slice().replace(
		"{record id}", recordId).replace(
		"{numSats}", numSats).replace(
		"{bitcion address}", btcAddr,
	);
	console.log(newMessageToSign)
	recordState.signatureMessage.value = newMessageToSign;
}

function updateSignalSignaturerMessage() {
	var rid = recordState.rid;
	if (rid.length == 0) {
		var name = recordState.name.value;
		var content = recordState.content.value;
		rid = sha256(sha256(name) + '::' + sha256(content));
	}
	updateNewSignalSignatureMessage(
		rid,
		recordState.numSats.value,
		recordState.address.value,
	);
}

recordState.name.addEventListener('input', (e)=> {updateSignalSignaturerMessage();}, false);
recordState.content.addEventListener('input', (e)=> {updateSignalSignaturerMessage();}, false);
recordState.address.addEventListener('input', (e)=> {updateSignalSignaturerMessage();}, false);
recordState.numSats.addEventListener('input', (e)=> {updateSignalSignaturerMessage();}, false);
createRecordTrigger.addEventListener("click", toggleCreateRecordModal);
cancelCreateTrigger.addEventListener("click", toggleCreateRecordModal);
submitNewRecord.addEventListener('click', (e)=> {
	// dont reload page
	e.preventDefault();

	// fire up spinner
	addLoadingClass(submitNewRecord);

	// make request
	if (recordState.rid.length == 0) {
		sendJsonPost(routes.newRecord, "POST", {
			name: recordState.name.value,
			value: recordState.content.value,
			signals: [{
				btc_address: recordState.address.value,
				sats: new Number(recordState.numSats.value),
				signature: recordState.signature.value,
			},],
		}, successfullyAddedRecord, failedToAddRecord);
	} else {
		sendJsonPost(routes.newSignal, "POST", {
			rid: recordState.rid,
			signals: [{
				rid: recordState.rid,
				btc_address: recordState.address.value,
				sats: new Number(recordState.numSats.value),
				signature: recordState.signature.value,
			}],
		}, successfullyAddedRecord, failedToAddRecord);
		recordState.content.rid = '';
	}
	// TODO when paginating, check last rank to see
	// if table needs to be regenerated
	genTable()
}, false);

function successfullyRetrievedTemplate(xhr) {
	// resp = JSON.parse(xhr.response);
	// console.log(resp)
	console.log(xhr.response)
	newSignalSignatureMessage = xhr.response;
}

function failedToRetrieveTemplate(xhr) {
    showErrorBanner(xhr.responseText);
} 

function getMessageTemplate() {
	sendJsonPost(
		routes.getMessageTemplate, "GET", null,
		successfullyRetrievedTemplate, failedToRetrieveTemplate,
	);
}

getMessageTemplate();

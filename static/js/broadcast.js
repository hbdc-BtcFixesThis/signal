let broadcast = document.forms.broadcast;

let createRecordModal = document.getElementById("create-record-modal");
let createRecordTrigger = document.getElementById("create-record-trigger");
let cancelCreateTrigger = document.getElementById("cancel-create-trigger");
let submitNewRecord = document.getElementById('submit-create-record');

let signalSignaturerMessage = document.getElementById('signal-signature-message')

const newSignalSignatureMessage = `This is not a bitcoin transaction!

For as long as the there are at least

SATS: {numSats}

unspent in Bitcoin

ADDRESS: {bitcion address}

may they be used to spread

RECORD ID: {record id}


Peace and love freaks`

let recordState = {
	name: document.getElementById("new-record-name"),
	content: document.getElementById("new-record-content"),
	address: document.getElementById("signal-wallet-address"),
	numSats: document.getElementById('new-record-signal-sats'),
	signature: document.getElementById("signal-signature"),
	signatureMessage: document.getElementById('signal-signature-message'),
};

function toggleCreateRecordModal() { toggleShowModal(createRecordModal); }

function successfullyAddedRecord(xhr) {
    removeLoadingClass(submitNewRecord);
	toggleCreateRecordModal();
	showSuccessBanner(xhr.responseText)
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
	recordState.signatureMessage.value = newMessageToSign;
}

function updateSignalSignaturerMessage() {
	var name = recordState.name.value;
	var content = recordState.content.value;
	updateNewSignalSignatureMessage(
		sha256(sha256(name) + '::' + sha256(content)),
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
	sendJsonPost(routes.newRecord, "POST", {
		key: recordState.name.value,
		value: recordState.content.value,
		signals: [{
			btc_address: recordState.address.value,
			sats: new Number(recordState.numSats.value),
			signature: recordState.signature.value,
		},],
	}, successfullyAddedRecord, failedToAddRecord);
	genTable()
}, false);

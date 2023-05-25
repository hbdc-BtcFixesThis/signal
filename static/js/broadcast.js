let broadcast = document.forms.broadcast;

let createRecordModal = document.getElementById("create-record-modal");
let createRecordTrigger = document.getElementById("create-record-trigger");
let cancelCreateTrigger = document.getElementById("cancel-create-trigger");
let submitNewRecord = document.getElementById('submit-create-record');

let signalSignaturerMessage = document.getElementById('signal-signature-message')

const newSignalSignatureMessage = `This is not a bitcoin transaction!

For as long as the there are funds left
unspent in bitcoin wallet address

{bitcion address}

may {percent}% of the balance be used
to spread the following record

{name label}:
{name}
{content label}:
{content}

Peace and love freaks`

let recordState = {
	name: document.getElementById("new-record-name"),
	content: document.getElementById("new-record-content"),
	address: document.getElementById("signal-wallet-address"),
	percent: document.getElementById('percent-slider-input'),
	signature: document.getElementById("signal-signature"),
	signatureMessage: document.getElementById('signal-signature-message'),
};

function toggleCreateRecordModal() {
    createRecordModal.classList.toggle("show-modal");
}

function successfullyAddedRecord(xhr) {
    removeLoadingClass(submitNewRecord);
	toggleCreateRecordModal();
	showSuccessBanner(xhr.responseText)
}

function failedToAddRecord(xhr) {
	removeLoadingClass(submitNewRecord);
	showErrorBanner(xhr.responseText);
}

function updateNewSignalSignatureMessage(name, content, percent, nLabel, cLabel, btcAddr) {
	var newMessageToSign = newSignalSignatureMessage.slice().replace(
		"{name}", name).replace(
		"{content}", content).replace(
		"{percent}", percent).replace(
		"{name label}", nLabel).replace(
		"{content label}", cLabel).replace(
		"{bitcion address}", btcAddr,
	);
	recordState.signatureMessage.value = newMessageToSign;
}

function updateSignalSignaturerMessage() {
	const hashLength = 64;

	var name = recordState.name.value;
	var nameLabel = 'RECORD NAME';

	var content = recordState.content.value;
	var contentLabel = 'RECORD CONTENT';

	var hashLabel = ' HASH/FINGERPRINT'

	if (name.length > hashLength) {
		name = sha256(recordState.name.value);
		nameLabel += hashLabel;
	}
	if (content.length > hashLength) {
		content = sha256(recordState.content.value)
		contentLabel += hashLabel
	}
	updateNewSignalSignatureMessage(
		name, content, recordState.percent.value,
		nameLabel, contentLabel, recordState.address.value,
	);
}

recordState.name.addEventListener('input', (e)=> {updateSignalSignaturerMessage();}, false);
recordState.content.addEventListener('input', (e)=> {updateSignalSignaturerMessage();}, false);
recordState.address.addEventListener('input', (e)=> {updateSignalSignaturerMessage();}, false);
createRecordTrigger.addEventListener("click", toggleCreateRecordModal);
cancelCreateTrigger.addEventListener("click", toggleCreateRecordModal);
broadcast.addEventListener('submit', (e)=> {
	// dont reload page
	e.preventDefault();

	// fire up spinner
	addLoadingClass(submitNewRecord);

	// make request
	sendJsonPost(routes.createRecord, {
		name: recordState.name.value,
		content: recordState.content.value,
		address: recordState.address.value,
		percent: recordState.percent.value,
		signature: recordState.signature.value,
	}, successfullyAddedRecord, failedToAddRecord);
}, false);

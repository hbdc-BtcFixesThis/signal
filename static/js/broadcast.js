let broadcast = document.forms.broadcast;

let createRecordModal = document.getElementById("create-record-modal");
let createRecordTrigger = document.getElementById("create-record-trigger");
let cancelCreateTrigger = document.getElementById("cancel-create-trigger");
let submitNewRecord = document.getElementById('submit-create-record');

let signalSignaturerMessage = document.getElementById('signal-signature-message')

const newSignalSignatureMessage = `This is not a bitcoin transaction!

For as long as the there are funds left
unspent, may {percent}% of the balance in
{bitcion address}, 
be used to spread this record
whose name and content hash is 
{hash of name}
{hash of content}

Peace and love freaks`

let recordState = {
	name: document.getElementById("new-record-name"),
	content: document.getElementById("new-record-content"),
	address: document.getElementById("signal-wallet-address"),
	percent: document.getElementById('percent-slider-input'),
	signature: document.getElementById("signal-signature"),
	signatureMessage: document.getElementById('signal-signature-message'),
};


function updateSignalSignaturerMessage() {
	var newMessageToSign = newSignalSignatureMessage.slice().replace(
		"{hash of name}", recordState.name.value).replace(
		"{hash of content}", recordState.content.value).replace(
		"{bitcion address}", recordState.address.value).replace(
		"{percent}", recordState.percent.value,
	);
	recordState.signatureMessage.value = newMessageToSign;
}

recordState.name.addEventListener('input', (e)=> {updateSignalSignaturerMessage();}, false);
recordState.content.addEventListener('input', (e)=> {updateSignalSignaturerMessage();}, false);
recordState.address.addEventListener('input', (e)=> {updateSignalSignaturerMessage();}, false);

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

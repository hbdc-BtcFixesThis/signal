let broadcast = document.forms.broadcast;

let createRecordModal = document.getElementById("create-record-modal");
let createRecordTrigger = document.getElementById("create-record-trigger");
let cancelCreateTrigger = document.getElementById("cancel-create-trigger");
let submitNewRecord = document.getElementById('submit-create-record');

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

	sendJsonPost(routes.createRecord, {
		name: document.getElementById("new-record-name").value,
		content: document.getElementById("new-record-content").value,
	}, successfullyAddedRecord, failedToAddRecord);
}, false);

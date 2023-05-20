let broadcast = document.forms.broadcast;

let createRecordModal = document.getElementById("create-record-modal");
let createRecordTrigger = document.getElementById("create-record-trigger");
let cancelCreateButton = document.getElementById("cancel-create-trigger");

function toggleCreateRecordModal() {
    createRecordModal.classList.toggle("show-modal");
}

createRecordTrigger.addEventListener("click", toggleCreateRecordModal);
cancelCreateButton.addEventListener("click", toggleCreateRecordModal);

broadcast.addEventListener('submit', (e)=> {
	// dont reload page
	e.preventDefault();
	console.log(JSON.stringify({
		name: document.getElementById("new-record-name").value,
		content: document.getElementById("new-record-content").value,
	}));
    createRecordModal.classList.toggle("show-modal");
}, false);

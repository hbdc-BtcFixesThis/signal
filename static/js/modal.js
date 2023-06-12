function windowOnClick(event) {
    if (event.target === createRecordModal) {
        toggleCreateRecordModal();
    }
}

function toggleShowModal(modal) { modal.classList.toggle("show-modal"); }

window.addEventListener("click", windowOnClick);

function windowOnClick(event) {
    if (event.target === createRecordModal) {
        toggleCreateRecordModal();
    }
}

window.addEventListener("click", windowOnClick);

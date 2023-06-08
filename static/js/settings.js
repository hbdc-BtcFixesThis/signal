let serverKey = '';

let authContainer = document.getElementById('admin-key-container');
let expandSettings = document.getElementById('expand-settings');
let settingsContainer = document.getElementById('settings');
let unlockButton = document.getElementById('unlock');

// to access server settings
let adminKey = document.getElementById('admin-key');

// change abo
let oldAdminKey = document.getElementById('password');
let newAdminKey = document.getElementById('new-password');

function hide(e) {e.style.display = 'none';}
function show(e) {e.style.display = '';}

function successfullyUnlocked(xhr) {
	// todo: fill in vals from server
	// todo: set browser state
	if (JSON.parse(xhr.response).success) {
		addLoadingClass(unlockButton);
		hide(authContainer);
		show(expandSettings);
	} else {
		failedToUnlock(xhr);
	}
}

function failedToUnlock(xhr) {
	removeLoadingClass(unlockButton);
	showErrorBanner(xhr.responseText);
}

// only show login at first
hide(expandSettings);
hide(settingsContainer);

adminKey.addEventListener("keypress", (e)=> {
	if (e.key === "Enter") {
		e.preventDefault();
		unlock.click();
	}
});

unlock.addEventListener("click", (e)=> {
	// e.preventDefault();
	path = routes.verifyAuthToken + "?" + objectToEncodedQueryString({
		key: genToken(adminKey.value),
	});
	sendJsonPost(path, "GET", null, successfullyUnlocked, failedToUnlock)
});

expandSettings.addEventListener("click", (e)=> {
	e.preventDefault();
	if (settingsContainer.style.display === 'none') {
		show(settingsContainer);
	} else {
		hide(settingsContainer);
	}
});

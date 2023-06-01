let serverKey = '';

let authContainer = document.getElementById('admin-key');
let expandSettings = document.getElementById('expand-settings');
let settingsContainer = document.getElementById('settings');
let unlockButton = document.getElementById('unlock');

let password = document.getElementById('admin-key');
let newPassword = document.getElementById('new-password');

function hide(e) {e.style.display = 'none';}
function show(e) {e.style.display = '';}

function successfullyUnlocked(data) {
	// todo: fill in vals from server
	// todo: set browser state

	hide(authContainer);
	show(expandSettings);
}

// only show login at first
hide(expandSettings);
hide(settingsContainer);

password.addEventListener("keypress", (e)=> {
	if (e.key === "Enter") {
		e.preventDefault();
		unlock.click();
	}
});

unlock.addEventListener("click", (e)=> {
	// e.preventDefault();
	successfullyUnlocked(e);
});

expandSettings.addEventListener("click", (e)=> {
	e.preventDefault();
	if (settingsContainer.style.display === 'none') {
		show(settingsContainer);
	} else {
		hide(settingsContainer);
	}
});

let expandSettings = document.getElementById('expand-settings');
let settingsContainer = document.getElementById('settings');

let oldAdminKey = document.getElementById('password');
let newAdminKey = document.getElementById('new-password');

function hide(e) {e.style.display = 'none';}
function show(e) {e.style.display = '';}

expandSettings.addEventListener("click", (e)=> {
	e.preventDefault();
	if (isLoggedIn()) {
		if (settingsContainer.style.display === 'none') {
			show(settingsContainer);
		} else {
			hide(settingsContainer);
		}
	} else {
		toggleLoginModal();
	}
});

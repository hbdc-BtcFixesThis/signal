let settingsContainer = document.getElementById('settings');

let expandDBSelect = document.getElementById('expand-db-select-settings');
let dbSelectSettings = document.getElementById('db-select-settings');

let expandDBSettings = document.getElementById('expand-database-settings');
let dbSettings = document.getElementById('database-settings');

let expandPeerSettings = document.getElementById('expand-peer-settings');
let peerSettings = document.getElementById('peer-settings');

let expandServerSettings = document.getElementById('expand-server-settings');
let serverSettings = document.getElementById('server-settings');

let expandBitcoinSettings = document.getElementById('expand-bitcoin-settings');
let bitcoinSettings = document.getElementById('bitcoin-settings');

///////////////////
let oldAdminKey = document.getElementById('password');
let newAdminKey = document.getElementById('new-password');
//////////////////

var lastOpened = null;

function expandSettings(opening, e) {
	e.preventDefault();
	if (lastOpened !== null) {
		settingsContainer.classList.toggle("visible");
		lastOpened.classList.toggle("visible");
		lastOpened.classList.toggle("hide");
		if (opening === lastOpened) {
			lastOpened = null;
			return;
		}
	}
	settingsContainer.classList.toggle("visible");
	opening.classList.toggle("hide");
	opening.classList.toggle("visible");
	if (!opening.classList.contains("visible")) {
		lastOpened = null;
	} else {
		lastOpened = opening;
	}
}

expandDBSelect.addEventListener("click", (e)=> {expandSettings(dbSelectSettings, e);});
expandDBSettings.addEventListener("click", (e)=> {expandSettings(dbSettings, e);});
expandPeerSettings.addEventListener("click", (e)=> {expandSettings(peerSettings, e);});
expandServerSettings.addEventListener("click", (e)=> {expandSettings(serverSettings, e);});
expandBitcoinSettings.addEventListener("click", (e)=> {expandSettings(bitcoinSettings, e);});

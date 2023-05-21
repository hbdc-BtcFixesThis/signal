var closeBanners = {
	error: document.getElementById("ack-error"),
	warninig: document.getElementById("ack-warning"),
	info: document.getElementById("ack-info"),
	success: document.getElementById("ack-success"),
};

var messageBanners = {
	error: document.getElementById("banner-error-message"),
	warninig: document.getElementById("banner-warning-message"),
	info: document.getElementById("banner-info-message"),
	success: document.getElementById("banner-success-message"),
};

function setBanner(close, banner, msg, hidden) {
	close.parentNode.hidden = hidden;
	close.parentNode.parentNode.hidden = hidden;
	close.parentNode.parentNode.parentNode.hidden = hidden;
	banner.innerHTML = msg
}

function showBanner(key, msg) {
	// ex: msg = wtf; key = error;
	// <strong>Error</strong>wft
	let msgHtml = '<strong>' + key.charAt(0).toUpperCase() + key.slice(1) + '</strong> ' + msg;
	setBanner(closeBanners[key], messageBanners[key], msgHtml, false);
}

function hideBanner(key) {
	setBanner(closeBanners[key], messageBanners[key], '', true);
}

function showErrorBanner(msg) {showBanner('error',  msg);}
function showWarningBanner(msg) {showBanner('warning', msg);}
function showInfoBanner(msg) {showBanner('info', msg);}
function showSuccessBanner(msg) {showBanner('success', msg);}

for (const [key, value] of Object.entries(closeBanners)) {
	// console.log(`${key}: ${value}`);
	value.addEventListener('click', function() {hideBanner(key)});
}

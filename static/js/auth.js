const AUTH_TOKEN = "settingsToken";
const AUTH_USER = "settingsUser"

let authKey = document.getElementById("auth-key");
let authUser = document.getElementById("auth-user");
let loginModal = document.getElementById ("login-modal")
let cancelLogin = document.getElementById ("cancel-login")
let unlockButton = document.getElementById("unlock");

function padTo2Digits(num) {
	return num.toString().padStart(2, '0');
}

function formatDate(date) {
	return (
		[
			date.getFullYear(),
			padTo2Digits(date.getUTCMonth() + 1),
			padTo2Digits(date.getUTCDate()),
		].join('-')
	);
}

function genToken(pw) {
	return sha256(sha256(pw)+" "+formatDate(new Date()));
}

function toggleLoginModal() { toggleShowModal(loginModal); }

function isLoggedIn() {
	// If a token is stored here it will only be
	// valid for a day. The token is a hash of
	// your password and the date of the token
	// creation (check genToken below for impl)
	sk = localStorage.getItem(AUTH_TOKEN);
	su = localStorage.getItem(AUTH_USER);

	return Boolean(sk) && Boolean(su);
}	

function successfullyUnlocked(xhr) {
	// todo: fill in vals from server
	// todo: set browser state
	resp = JSON.parse(xhr.response);
	if (resp.status_code === 200) {
		showSuccessBanner(resp.message);
		removeLoadingClass(unlockButton);
		// show(expandSettings);
		expandSettings.click();

		// set in browser to avoid loosing creds on refresh
		localStorage.setItem(AUTH_TOKEN, genToken(authKey.value));
		localStorage.setItem(AUTH_USER, authUser.value);

		// clear auth data from inputs
		authKey.value = '';
		authUser.value = '';
	} else {
		failedToUnlock(xhr);
	}
}

function failedToUnlock(xhr) {
	removeLoadingClass(unlockButton);
	showErrorBanner(JSON.parse(xhr.response).error);
	
	// clear any auth creds stored in browser
	localStorage.removeItem(AUTH_TOKEN);
	localStorage.removeItem(AUTH_USER);
}

authKey.addEventListener("keypress", (e)=> {
	if (e.key === "Enter") {
		e.preventDefault();
		unlock.click();
	}
});

unlock.addEventListener("click", (e)=> {
	e.preventDefault();
	addLoadingClass(unlockButton);
	var key = authKey.value;
	var user = authUser.value;
	if (Boolean(key)) {
		key = genToken(key);
	} else { 
		key = localStorage.getItem(AUTH_TOKEN);
	}
	if (!Boolean(user)) {
		user = localStorage.getItem(AUTH_USER);
	}

	sendJsonPost(routes.verifyAuthToken + "?" + objectToEncodedQueryString({
		user: user, key: key,
	}), "GET", null, successfullyUnlocked, failedToUnlock)
});

cancelLogin.addEventListener("click", toggleLoginModal);

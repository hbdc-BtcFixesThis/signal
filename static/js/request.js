var loadingClassName = 'loading';
var regexLoadingClassStr = '(?:^|\\s)'+ loadingClassName + '(?!\\S)';

function removeLoadingClass(elem) {
	elem.className = elem.className.replace(new RegExp(regexLoadingClassStr), '');
}

function addLoadingClass(elem) {
	elem.className += ' ' + loadingClassName;
	// setTimeout(removeLoadingClass, 2000, elem);
}

function sendJsonPost(path, data, success, fail) {
	const xhr = new XMLHttpRequest();
	xhr.open("POST", routes.createRecord);

	// Send the proper header information along with the request
	xhr.setRequestHeader("Content-Type", "application/json");

	xhr.onreadystatechange = () => {
		// Call a function when the state changes.
		if (xhr.readyState === XMLHttpRequest.DONE) {
			// Request finished. Do processing here.
			if (xhr.status === 200) {
				success(xhr)
			} else {
				fail(xhr)
			}
		}
	};
	xhr.send(JSON.stringify(data));
}

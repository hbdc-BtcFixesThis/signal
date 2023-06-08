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

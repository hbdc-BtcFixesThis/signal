/*$(window).on("load resize ", function() {
  var scrollWidth = $('.tbl-content').width() - $('.tbl-content table').width();
  $('.tbl-header').css({'padding-right':scrollWidth});
}).resize();
*/

var results = {
	'headers': [
		{'key': 'id',           'display_header': 'Name'},
		{'key': 'content',      'display_header': 'Content'},
		//{'key': 'content_size', 'display_header': 'Content Size'},
		//{'key': 'record_size',  'display_header': 'Record Size'},
		//{'key': 'sum_positive', 'display_header': '+ sats'},
		//{'key': 'sum_negative', 'display_header': '- sats'},
		//{'key': 'total_sats',   'display_header': 'Total Sats'},
		//{'key': 'sum_sats',     'display_header': 'Sum'},
		{'key': 'signal',       'display_header': 'Signal (total/size)'},
		//{'key': 'signal_count', 'display_header': 'Number of Signals'},
		//{'key': 'addr',         'display_header': 'Funding Address'},
		{'key': 'created_at', 'display_header': 'Created'},
	],
	'signals': [
	{
		'id': 'unique user specified name',
		'type': 'text',
		'content': 'test',
		'content_size': '32bytes',
		'record_size': '74bytes',
		'sum_positive': '100 sats',
		'sum_negative': '100 sats',
		'total_sats': '200 sats',
		'sum_sats': '0 sats',
		'signal': '200/74',
		'signal_count': '2',
		'addr': 'N/A',
		'created_at': '111111',
	},
	{
		'id': 'unique user specified name',
		'type': 'text',
		'content': 'test',
		'content_size': '32bytes',
		'record_size': '74bytes',
		'sum_positive': '100 sats',
		'sum_negative': '100 sats',
		'total_sats': '200 sats',
		'sum_sats': '0 sats',
		'signal': '200/74',
		'signal_count': '2',
		'addr': 'N/A',
		'created_at': '111111',
	},
	],
}

function genHeaders(headers) {
	var headerHTML = '';
	for (header of headers) {
		headerHTML += '<th>' + header.display_header + '</th>\n';
	}
	document.getElementById('tbl-header-tr').innerHTML = headerHTML;
	console.log(headerHTML);
}

function genRow(row, headers) {
	// follow same order specified for headers
	var rowHTML = '';
	for (header of headers) {
		rowHTML += '<td>' + row[header.key] + '</td>';
	}
	return rowHTML
}

function genRows(headers, rows) {
	var modalTemplateHtml = '';
	for (row of rows) {
		document.getElementById('tbl-body-content').innerHTML += '<tr class="trigger-modal">' + genRow(row, headers) + '</tr>';
	}
}

function genTable() {
	genHeaders(results.headers);
	genRows(results.headers, results.signals);
}

genTable()


var modal = document.querySelector(".modal");
var trigger = document.querySelector(".trigger-modal");
var closeButton = document.querySelector(".close-button");

function toggleModal() {
    modal.classList.toggle("show-modal");
}

function windowOnClick(event) {
    if (event.target === modal) {
        toggleModal();
    }
}

trigger.addEventListener("click", toggleModal);
closeButton.addEventListener("click", toggleModal);
window.addEventListener("click", windowOnClick);




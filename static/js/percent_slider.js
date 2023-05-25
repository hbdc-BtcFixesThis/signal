var percentNewRecordInput = document.querySelector("#percent-slider-input");
var newRecordSliderBar1 = document.getElementById("percent-slider-new-rec-1");
var newRecordSliderBar2 = document.getElementById("percent-slider-new-rec-2");

var percentSliderStorage = document.querySelector("#percent-slider-storage");
var serverStorageSliderBar1 = document.getElementById("percent-slider-server-storage-1");
var serverStorageSliderBar2 = document.getElementById("percent-slider-server-storage-2");


function initSliderPercentInput(sliderBar, elem) {
	elem.style.height = getComputedStyle(sliderBar).height;
	elem.style.bottom = getComputedStyle(sliderBar).height
}

function initStorageCapacityPercent() {
	var sliderBars = document.querySelector(".percent-slider-bars");
	percentSliderStorage.style.height = getComputedStyle(sliderBars).height;
	percentSliderStorage.style.bottom = getComputedStyle(sliderBars).height
}

function percentChange(bar1, bar2, e) {
	var applied = e.target.value + "%";
	var remaining = Number(100 - e.target.value) + "%";

	// left side of pct slider
	bar1.style.width = bar1.innerHTML = applied;
	// right side of pct slider
	bar2.style.width = bar2.innerHTML = remaining;
}

percentNewRecordInput.addEventListener("input", (e) => {
	percentChange(newRecordSliderBar1, newRecordSliderBar2, e);
	updateSignalSignaturerMessage();
});

percentSliderStorage.addEventListener("input", (e) => {
	percentChange(serverStorageSliderBar1, serverStorageSliderBar2, e);
	updateSignalSignaturerMessage();
});

initSliderPercentInput(newRecordSliderBar1, percentNewRecordInput);
initSliderPercentInput(serverStorageSliderBar1, percentSliderStorage);

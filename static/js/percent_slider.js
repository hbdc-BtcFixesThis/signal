var percentInput = document.querySelector("#percent-slider-input");

function initPercentInput() {
	var sliderBars = document.querySelector(".percent-slider-bars");
	percentInput.style.height = getComputedStyle(sliderBars).height;
	percentInput.style.bottom = getComputedStyle(sliderBars).height
}

percentInput.addEventListener("input", (e) => {
	var bars = document.querySelectorAll(".percent-slider-bar");

	var applied = e.target.value + "%";
	var remaining = Number(100 - e.target.value) + "%";

	// left side of pct slider
	bars[0].style.width = bars[0].innerHTML = applied;

	// right side of pct slider
	bars[1].style.width = bars[1].innerHTML = remaining;

	updateSignalSignaturerMessage();
});

initPercentInput()

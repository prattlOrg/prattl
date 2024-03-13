window.addEventListener("DOMContentLoaded", () => {
	addElement();
});

function addElement() {
	transcriptionWrapper = document.getElementById("transcription-wrapper");
	const transcriptionText = document.createElement("marquee");
	transcriptionText.setAttribute("id", "transcription-text");
	const textContent = document.createTextNode("Hello world");
	transcriptionText.appendChild(textContent);
	transcriptionWrapper.appendChild(transcriptionText);
}

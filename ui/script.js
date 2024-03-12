function addElement() {
	// create a new div element
	const transcriptionWrapper = document.createElement("div", {
		class: "transcription-wrapper",
	});

	// and give it some content
	const newContent = document.createTextNode("Hello world");

	// add the text node to the newly created div
	transcriptionWrapper.appendChild(newContent);

	// add the newly created element and its content into the DOM
	const currentDiv = document.getElementById("div1");
	document.body.insertBefore(transcriptionWrapper, currentDiv);
}

addElement();

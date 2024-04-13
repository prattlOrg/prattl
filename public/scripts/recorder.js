// Set up basic variables for app
const record = document.querySelector(".record");
const rec_stop = document.querySelector(".stop");
const soundClips = document.querySelector(".sound-clips");
const canvas = document.querySelector(".visualizer");
const mainSection = document.querySelector(".main-controls");

// Disable stop button while not recording
rec_stop.disabled = true;

// Visualiser setup - create web audio api context and canvas
let audioCtx;
const canvasCtx = canvas.getContext("2d");

// Main block for doing the audio recording
if (navigator.mediaDevices.getUserMedia) {
	console.log("The mediaDevices.getUserMedia() method is supported.");

	const constraints = { audio: true };
	let chunks = [];

	let onSuccess = function (stream) {
		const mediaRecorder = new MediaRecorder(stream);

		visualize(stream);

		record.onclick = function () {
			mediaRecorder.start();
			console.log(mediaRecorder.state);
			console.log("Recorder started.");
			record.style.background = "red";

			rec_stop.disabled = false;
			record.disabled = true;
		};

		rec_stop.onclick = function () {
			mediaRecorder.stop();
			console.log(mediaRecorder.state);
			console.log("Recorder stopped.");
			record.style.background = "";
			record.style.color = "";

			rec_stop.disabled = true;
			record.disabled = false;
		};

		mediaRecorder.onstop = function (e) {
			console.log("Last data to read (after MediaRecorder.stop() called).");

			const clipName = prompt(
				"Enter a name for your sound clip?",
				"My unnamed clip"
			);

			const clipContainer = document.createElement("article");
			const clipLabel = document.createElement("p");
			const audio = document.createElement("audio");
			const deleteButton = document.createElement("button");
			const transcribeButton = document.createElement("button");

			clipContainer.classList.add("clip");
			audio.setAttribute("controls", "");
			deleteButton.textContent = "Delete";
			deleteButton.className = "delete";

			transcribeButton.textContent = "Transcribe";

			if (clipName === null) {
				clipLabel.textContent = "My unnamed clip";
			} else {
				clipLabel.textContent = clipName;
			}

			clipContainer.appendChild(audio);
			clipContainer.appendChild(clipLabel);
			clipContainer.appendChild(deleteButton);
			clipContainer.appendChild(transcribeButton);
			soundClips.appendChild(clipContainer);

			audio.controls = true;
			const blob = new Blob(chunks, { type: mediaRecorder.mimeType });
			chunks = [];
			const audioURL = window.URL.createObjectURL(blob);
			audio.src = audioURL;
			console.log("recorder stopped");

			deleteButton.onclick = function (e) {
				e.target.closest(".clip").remove();
			};

			transcribeButton.onclick = async function (e) {
				const file = new File(
					[audio.src],
					`${clipLabel.textContent.replaceAll(" ", "_")}.mp3`,
					{
						type: "audio/mpeg",
					}
				);
				console.log(file);

				const formData = new FormData();
				formData.append("file", file);

				fetch("http://localhost:8080/transcribe/", {
					method: "POST",
					body: formData,
				})
					.then((response) => {
						if (!response.ok) {
							throw new Error("Network response was not ok");
						}
						return response.json(); // Assuming the server returns JSON
					})
					.then((data) => {
						console.log("File uploaded successfully:", data);
					})
					.catch((error) => {
						console.error("Error uploading file:", error);
					});
			};

			clipLabel.onclick = function () {
				const existingName = clipLabel.textContent;
				const newClipName = prompt("Enter a new name for your sound clip?");
				if (newClipName === null) {
					clipLabel.textContent = existingName;
				} else {
					clipLabel.textContent = newClipName;
				}
			};
		};

		mediaRecorder.ondataavailable = function (e) {
			chunks.push(e.data);
		};
	};

	let onError = function (err) {
		console.log("The following error occured: " + err);
	};

	navigator.mediaDevices.getUserMedia(constraints).then(onSuccess, onError);
} else {
	console.log("MediaDevices.getUserMedia() not supported on your browser!");
}

function visualize(stream) {
	if (!audioCtx) {
		audioCtx = new AudioContext();
	}

	const source = audioCtx.createMediaStreamSource(stream);

	const analyser = audioCtx.createAnalyser();
	analyser.fftSize = 2048;
	const bufferLength = analyser.frequencyBinCount;
	const dataArray = new Uint8Array(bufferLength);

	source.connect(analyser);

	draw();

	function draw() {
		const WIDTH = canvas.width;
		const HEIGHT = canvas.height;

		requestAnimationFrame(draw);

		analyser.getByteTimeDomainData(dataArray);

		canvasCtx.fillStyle = "rgb(200, 200, 200)";
		canvasCtx.fillRect(0, 0, WIDTH, HEIGHT);

		canvasCtx.lineWidth = 2;
		canvasCtx.strokeStyle = "rgb(0, 0, 0)";

		canvasCtx.beginPath();

		let sliceWidth = (WIDTH * 1.0) / bufferLength;
		let x = 0;

		for (let i = 0; i < bufferLength; i++) {
			let v = dataArray[i] / 128.0;
			let y = (v * HEIGHT) / 2;

			if (i === 0) {
				canvasCtx.moveTo(x, y);
			} else {
				canvasCtx.lineTo(x, y);
			}

			x += sliceWidth;
		}

		canvasCtx.lineTo(canvas.width, canvas.height / 2);
		canvasCtx.stroke();
	}
}

window.onresize = function () {
	canvas.width = mainSection.offsetWidth;
};

window.onresize();

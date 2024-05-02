// Set up basic variables for app
const record = document.querySelector(".record");
const rec_stop = document.querySelector(".stop");

// Disable stop button while not recording
rec_stop.disabled = true;

// Visualiser setup - create web audio api context and canvas
let audioCtx;

const blobToBase64 = (blob) => {
	const reader = new FileReader();
	reader.readAsDataURL(blob);
	return new Promise((resolve) => {
		reader.onloadend = () => {
			resolve(reader.result);
		};
	});
};

if (
	navigator.mediaDevices.getUserMedia({
		audio: {
			channelCount: 1,
			sampleRate: 44100,
		},
	})
) {
	console.log("The mediaDevices.getUserMedia() method is supported.");

	const constraints = { audio: true };
	let chunks = [];

	let onSuccess = function (stream) {
		const options = { mimeType: "audio/webm" };
		const mediaRecorder = new MediaRecorder(stream, options);
		var socket = new WebSocket("ws://localhost:8080/transcribe/");

		record.onclick = function () {
			mediaRecorder.start();
			record.style.background = "red";
			rec_stop.disabled = false;
			record.disabled = true;
		};

		rec_stop.onclick = function () {
			mediaRecorder.stop();
			record.style.background = "";
			record.style.color = "";
			rec_stop.disabled = true;
			record.disabled = false;
		};

		mediaRecorder.onstop = function (e) {
			console.log("recorder stopped");
		};

		mediaRecorder.ondataavailable = async function (e) {
			chunks.push(e.data);
			let reader = e.data.stream().getReader();
			reader.read().then(async function processText({ done, value }) {
				const blob = new Blob(chunks, { type: "audio/wav" });
				const base64Blob = await blobToBase64(blob);
				socket.send(base64Blob);
				if (done) {
					return;
				}
			});
		};
	};

	let onError = function (err) {
		console.log("The following error occured: " + err);
	};

	navigator.mediaDevices.getUserMedia(constraints).then(onSuccess, onError);
} else {
	console.log("MediaDevices.getUserMedia() not supported on your browser!");
}

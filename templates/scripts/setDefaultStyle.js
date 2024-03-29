import { getTextOptions, storeTextOptions } from "./localStorage.js";
import changeTextStyle from "./changeTextStyle.js";

const setDefaultStyle = (textId) => {
	const textOptions = getTextOptions();
	console.log(textOptions);
	if (!textOptions) {
		storeTextOptions({
			font: "system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif",
			textColor: "greenyellow",
			textSize: "4rem",
		});
	}
	changeTextStyle(textId);
};

export default setDefaultStyle;

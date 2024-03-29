import { getTextOptions } from "./localStorage.js";

export const changeTextStyle = (textId) => {
	const transcriptionText = document.getElementById(textId);
	const styles = getTextOptions();
	transcriptionText.style.color = styles.textColor;
	transcriptionText.style.fontSize = styles.textSize;
	transcriptionText.style.fontFamily = styles.font;
};

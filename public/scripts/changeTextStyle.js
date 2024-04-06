import { getTextOptions } from "./localStorage.js";

const changeTextStyle = (textId) => {
	const transcriptionText = document.getElementById(textId);
	const styles = getTextOptions();
	transcriptionText.style.fontFamily = styles.font;
	transcriptionText.style.color = styles.textColor;
	transcriptionText.style.fontSize = `${styles.textSize}px`;
};

export default changeTextStyle;

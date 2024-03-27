export const storeTextOptions = (options) => {
	window.localStorage.setItem("textOptions", JSON.stringify(options));
};

export const getTextOptions = () => {
	const options = window.localStorage.getItem("textOptions");
	return JSON.parse(options);
};

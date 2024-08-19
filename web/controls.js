function useControl(id, onChange) {
    const element = document.getElementById(id);
    if (element.value !== undefined) {
        // Capture initial state
        onChange(getControlElementValue(element));
    }
    element.addEventListener("change", function (evt) {
        onChange(getControlElementValue(evt.target));
    });
}

function controlCompressionOption(name) {
    return function (value) {
        appState.compressionOptions[name] = value;
    };
}

function getControlElementValue(element) {
    let value = element.value;
    if (element.getAttribute("type") == "number") {
        value = Number(value);
    }
    return value;
}

function useControl(id, onChange) {
    const element = document.getElementById(id);
    element.addEventListener("change", function (evt) {
        let value = evt.target.value;
        if (evt.target.getAttribute("type") == "number") {
            value = Number(value);
        }
        onChange(value);
    });
}

function controlCompressionOption(name) {
    return function (value) {
        appState.compressionOptions[name] = value;
    };
}

function useControl(id, onChange) {
    const element = document.getElementById(id);
    element.addEventListener("change", function (evt) {
        onChange(evt.target.value);
    });
}

function controlCompressionOption(name) {
    return function (value) {
        appState.compressionOptions[name] = value;
    };
}

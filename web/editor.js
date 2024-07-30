function setupEditor() {
    editorCompressedImage.style.backgroundImage = `url(${userState.sourceImage})`;
    editorCompressedImage.style.backgroundSize = "100%";
}

function handleEditorMouseWheel(evt) {
    evt.preventDefault();

    userState.editorZoom += evt.deltaY * -0.01;
    userState.editorZoom = Math.min(Math.max(0.3, userState.editorZoom), 4);

    EventEditorZoom.fire({ zoom: userState.editorZoom });
}

function applyImageZoom({ zoom }) {
    zoom = (zoom * 100).toFixed(0);
    editorCompressedImage.style.backgroundSize = `${zoom}%`;
}
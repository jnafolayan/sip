function startImagePanning(evt) {
    editorState.panning = true;
    editorState.pan.oldX = evt.pageX;
    editorState.pan.oldY = evt.pageY;
}

function panImage(evt) {
    const { pan, panning } = editorState;
    if (!panning) return;

    evt.preventDefault();
    const dx = evt.pageX - pan.oldX;
    const dy = evt.pageY - pan.oldY;
    pan.x += dx;
    pan.y += dy;

    pan.oldX = evt.pageX;
    pan.oldY = evt.pageY;
}

function endImagePanning(_evt) {
    if (!editorState.panning) return;
    editorState.panning = false;
}

function handleEditorMouseWheel(evt) {
    evt.preventDefault();

    const oldScale = editorState.scale;

    editorState.scale += evt.deltaY * -0.005;
    editorState.scale = Math.min(Math.max(0.15, editorState.scale), 3);

    const delta = editorState.scale - oldScale;

    EventEditorZoom.fire({
        pageX: evt.pageX,
        pageY: evt.pageY,
        scale: editorState.scale,
        delta,
    });
}

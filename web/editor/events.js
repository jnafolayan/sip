function startImagePanning(evt) {
    const { editor } = userState;
    editor.panning = true;
    editor.pan.oldX = evt.pageX;
    editor.pan.oldY = evt.pageY;
}

function panImage(evt) {
    const { pan, panning } = userState.editor;
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
    const { editor } = userState;
    editor.panning = false;
}

function handleEditorMouseWheel(evt) {
    evt.preventDefault();

    const oldScale = userState.editor.scale;

    userState.editor.scale += evt.deltaY * -0.005;
    userState.editor.scale = Math.min(
        Math.max(0.15, userState.editor.scale),
        3
    );

    const delta = userState.editor.scale - oldScale;

    EventEditorZoom.fire({
        pageX: evt.pageX,
        pageY: evt.pageY,
        scale: userState.editor.scale,
        delta,
    });
}

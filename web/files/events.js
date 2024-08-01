function handleImageUpload(evt) {
    const files = evt.target.files;
    if (!files.length) return;

    appState.source = null;

    const [file] = files;
    const fr = new FileReader();
    fr.onload = function () {
        image.src = fr.result;
    };
    fr.onprogress = function (evt) {
        EventFileUploadProgress.fire({ progress: evt.loaded / evt.total });
    };

    const image = new Image();
    image.onload = () => {
        EventFileUploadEnd.fire({ progress: 1 });
        appState.source = {
            image,
            width: image.naturalWidth,
            height: image.naturalHeight,
            size: file.size,
            name: file.name,
        };
    };

    EventFileUploadStart.fire({ progress: 0 });
    fr.readAsDataURL(file);
}

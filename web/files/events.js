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
        const imgCanvas = image2Canvas(image);
        appState.source = {
            image: imgCanvas,
            width: image.width,
            height: image.height,
            name: file.name,
            size: file.size,
        };
    };

    EventFileUploadStart.fire({ progress: 0 });
    fr.readAsDataURL(file);
}

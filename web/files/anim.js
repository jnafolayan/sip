function stepFileUploadAnimation({ progress }) {
    if (progress == 0) {
        uploadButton.classList.add("hide");
        uploadProgress.classList.remove("hide");
    }

    progress *= 100;
    uploadProgress.innerHTML = `Loading <span>${(progress).toFixed(0)}%<span>`;
}

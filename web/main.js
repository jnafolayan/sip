let uploadButton, fileInput, uploadProgress;

// VIEWS
let uploadView, editorView;
let editorOriginalImage, editorCompressedImage;

// EVENTS
const [EventFileUploadStart, EventFileUploadProgress, EventFileUploadEnd] = [
    new AppEvent("FILE_UPLOAD_START"),
    new AppEvent("FILE_UPLOAD_PROGRESS"),
    new AppEvent("FILE_UPLOAD_END"),
];
const EventEditorOpened = new AppEvent("EDITOR_OPENED");

// STATE
let userState;

window.onload = setup;

function setup() {
    uploadButton = document.getElementById("uploadButton");
    uploadProgress = document.getElementById("uploadProgress");
    fileInput = document.getElementById("imageUpload");

    uploadView = document.querySelector(".view__upload");
    editorView = document.querySelector(".view__editor");

    editorOriginalImage = document.querySelector(".editor__images__original");
    editorCompressedImage = document.querySelector(
        ".editor__images__compressed"
    );

    userState = createUserState();

    setupEvents();
    subscribeToAppEvents();
}

function setupEvents() {
    fileInput.addEventListener("change", handleImageUpload);
}

function subscribeToAppEvents() {
    // File upload
    EventFileUploadStart.subscribe(stepFileUploadAnimation);
    EventFileUploadProgress.subscribe(stepFileUploadAnimation);
    EventFileUploadEnd.subscribe(stepFileUploadAnimation);
    EventFileUploadEnd.subscribe(() => {
        setTimeout(() => {
            uploadView.classList.add("hide");
            editorView.classList.remove("hide");
            EventEditorOpened.fire();
        }, 500);
    });

    // View
    EventEditorOpened.subscribe(setupEditor);
}

function handleImageUpload(evt) {
    const files = evt.target.files;
    if (!files.length) return;

    const [file] = files;
    const fr = new FileReader();
    fr.onload = function () {
        EventFileUploadEnd.fire({ progress: 1 });
        userState.sourceImage = fr.result;
    };
    fr.onprogress = function (evt) {
        EventFileUploadProgress.fire({ progress: evt.loaded / evt.total });
    };
    EventFileUploadStart.fire({ progress: 0 });
    fr.readAsDataURL(file);
}

// state
function createUserState() {
    return {
        sourceImage: null,
        compressionResult: null,
    };
}

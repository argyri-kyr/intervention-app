let mediaRecorder;
let recordedChunks = [];

function startRecording() {
    navigator.mediaDevices.getUserMedia({ audio: true })
        .then(stream => {
            mediaRecorder = new MediaRecorder(stream);
            mediaRecorder.start();

            mediaRecorder.ondataavailable = function(e) {
                recordedChunks.push(e.data);
            };

            mediaRecorder.onstop = function() {
                let blob = new Blob(recordedChunks, {
                    type: 'audio/ogg; codecs=opus'
                });
                let url = URL.createObjectURL(blob);
                // Do something with the URL here, like download it or play it back
                //Play it back
                let audio = new Audio(url);
                audio.play();
            };
        })
        .catch(err => console.error('getUserMedia() failed: ', err));
}

function stopRecording() {
    if (mediaRecorder) {
        mediaRecorder.stop();
    }
}

let recordButtons = document.querySelectorAll('.button-record');
recordButtons.forEach(button => {
    button.addEventListener('click', function() {
        let img = button.querySelector('img');
        if (mediaRecorder && mediaRecorder.state === 'recording') {
            stopRecording();
            img.src = '/assets/icons/mic.svg'; // Change the image to mic.svg when not recording
        } else {
            startRecording();
            img.src = '/assets/icons/stop.svg'; // Change the image to stop.svg when recording
        }
    });
});
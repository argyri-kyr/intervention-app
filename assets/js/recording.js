let mediaRecorder;
let recordedChunks = [];
let audioData = null;

function toggleRecording(btn) {
    console.log(btn);
    let img = btn.querySelector('img');
    console.log(img);
    if (mediaRecorder && mediaRecorder.state === 'recording') {
        mediaRecorder.stop();
        img.src = '/assets/icons/mic.svg'; 
    } else {
        img.src = '/assets/icons/stop.svg';
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
                let audioURL = URL.createObjectURL(blob);
                // Do something with the URL here, like download it or play it back
                //Play it back
                let audio = new Audio(audioURL);
                audio.play();

                let reader = new FileReader();
                reader.onloadend = function() {
                    audioData = reader.result;
                    console.log(audioData);
                }
                reader.readAsDataURL(blob);
            };
        })
        .catch(err => console.error('getUserMedia() failed: ', err));
    }
}
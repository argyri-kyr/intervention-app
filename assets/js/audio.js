
function speak(text) {
    var msg = new SpeechSynthesisUtterance(" "+text);
    msg.lang = 'el-GR';
    var voices = window.speechSynthesis.getVoices();
    // console.log(voices);
    msg.pitch = 2;
    msg.rate = 1;
    msg.volume = 1;
    window.speechSynthesis.speak(msg);
}

function speakonce(text, element) {
    var msg = new SpeechSynthesisUtterance(" "+text);
    msg.lang = 'el-GR';
    var voices = window.speechSynthesis.getVoices();
    // console.log(voices);
    msg.pitch = 2;
    msg.rate = 1;
    msg.volume = 1;
    window.speechSynthesis.speak(msg);
    element.className = 'speak-word-once';
}
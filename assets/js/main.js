let START_EXERCISE = 0;
let answers = [];
let currentExercise = 0;
let templateDivs = document.querySelectorAll('.template');
let totalExercises = templateDivs.length;
const stopwatches = {};
let currentQuestion = 0;
let totalQuestions = getTotalQuestions();
const displayTime = 4000; // 4 seconds

hideExercises();

let answeredPairs = [];

// Add event listeners to word elements to handle selection
let wordElements = document.querySelectorAll('.speak-word');
wordElements.forEach(word => {
    word.addEventListener('click', function() {
        // Remove 'selected' class from all other words
        wordElements.forEach(w => w.classList.remove('selected'));
        // Add 'selected' class to the clicked word
        word.classList.add('selected');
    });
});

// Add event listeners to the dropzones
var dropzones = document.querySelectorAll('.dropzone');
dropzones.forEach(function(dropzone) {
    dropzone.addEventListener('dragover', allowDrop);
    dropzone.addEventListener('drop', drop);
});

updateProgressBar();
// Start the stopwatches for DecodingTextReading template
startStopwatchesForTemplate('.decoding-text-reading-item');

// Start the stopwatches for MemoryVisualSingleChoice template
startStopwatchesForTemplate('.memory-visual-single-choice-item');

// Check if the desired template is present outside the DOMContentLoaded event listener
const templates = document.querySelectorAll('.RecallImageChoice.template')
let isVisible = false;
templates.forEach(template => {
    if (template.style.display === 'block') {
        isVisible = true;
    }
});

// This is for the cases where i want to start with the memory block as first exercise in order to be initialized
if (templates!=null && templates.length > 0 && isVisible) {
    initializeRecallImageChoice();
}
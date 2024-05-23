let START_EXERCISE = 0;

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
    
// function submitAnswer(question, answer) {
//     fetch('/answer', {
//         method: 'POST',
//         headers: {
//             'Content-Type': 'application/x-www-form-urlencoded',
//         },
//         body: new URLSearchParams({
//             'question': question,
//             'answer': answer
//         })
//     }).then(response => {
//         if (!response.ok) {
//             throw new Error('Network response was not ok');
//         }
//         return response.json();
//     }).then(data => {
//         // Handle the response data here
//     }).catch(error => {
//         console.error('There has been a problem with your fetch operation:', error);
//     });
// }


// Function to extract the index from the class name
function extractIndexFromClassName(className) {
    const match = className.match(/i_(\d+)/);
    return match ? parseInt(match[1], 10) : null;
}

function submitAnswers() {
    // Convert answeredPairs and answers arrays to JSON strings
    var answeredPairsJson = JSON.stringify(answeredPairs);
    var answersJson = JSON.stringify(answers);
    console.log("answeredPairsJson = ", answeredPairsJson)
    console.log("answersJson = ", answersJson)

    // Concatenate the JSON strings with some separator
    var content = "Answered Pairs:\n" + answeredPairsJson + "\n\nAnswers:\n" + answersJson;
    
    // Create a Blob with the content
    var blob = new Blob([content], { type: "text/plain" });
    
    // Create a File from the Blob
    var file = new File([blob], "answers.txt", { type: "text/plain" });
    
    // Create a download link
    var downloadLink = document.createElement("a");
    downloadLink.href = URL.createObjectURL(file);
    downloadLink.download = "answers.txt"; // Set the file name
    
    // Programmatically click the download link
    document.body.appendChild(downloadLink);
    downloadLink.click();
    document.body.removeChild(downloadLink);
}

function hideExercises() {
    let selectedIndex = START_EXERCISE;
    console.log("START_EXERCISE = ",START_EXERCISE)
    let templateDivs = document.querySelectorAll('.template');
    // console.log(templateDivs);
    for(let i = 0; i < templateDivs.length; i++) {
        if(i === selectedIndex) {
            templateDivs[i].style.display = 'block';
        } else {
            templateDivs[i].style.display = 'none';
        }
    }
}


hideExercises();
let currentExercise = 0;
let templateDivs = document.querySelectorAll('.template');
let totalExercises = templateDivs.length;
const stopwatches = {};

function startStopwatch(id, index) {
    let timer;
    console.log(index);
    if (index === 0) {
        timer = 0;
    } else {
        timer = -1;
    }
        
    const updateStopwatch = () => {
        timer++;
        const hours = String(Math.floor(timer / 3600)).padStart(2, '0');
        const minutes = String(Math.floor((timer % 3600) / 60)).padStart(2, '0');
        const seconds = String(timer % 60).padStart(2, '0');
        document.getElementById(id).textContent = `${hours}:${minutes}:${seconds}`;
    };
    stopwatches[id] = setInterval(updateStopwatch, 1000);
}


function stopStopwatch(id) {
    clearInterval(stopwatches[id]);
}

function nextExercise() {
    console.log('Next exercise');
    let templateDivs = document.querySelectorAll('.template');
    totalExercises = templateDivs.length;
    for(let i = 0; i < templateDivs.length; i++) {
        if (templateDivs[i].style.display === 'block') {
            if (i === totalExercises - 1) {
                console.log('Last exercise, go Home');
                //Go to Home
                window.location.href = '/main.html';
            }
            currentExercise = i+1;
            console.log('got Current exercise from block:', currentExercise);
            break;
        }
    }

        templateDivs[currentExercise-1].style.display = 'none';
        templateDivs[currentExercise].style.display = 'block';
    
}

function getTotalQuestions() {
    let templateDivs = document.querySelectorAll('.template');
    let visibleTemplateDivs = Array.from(templateDivs).filter(div => window.getComputedStyle(div).display !== 'none');
    let childDivs = Array.from(visibleTemplateDivs[0].children).filter(child => child.tagName === 'DIV');

    return childDivs.length;
}

let currentQuestion = 0;
let totalQuestions = getTotalQuestions();

function updateQuestion(back = false, templateClass, currentQuestion) {
    let questions = document.querySelectorAll(templateClass);
    let curQ, nexQ = 0;
    if (back) {
        curQ = currentQuestion + 1;
        nexQ = currentQuestion;
    } else {
        curQ = currentQuestion - 1;
        nexQ = currentQuestion;
    }
    console.log(`Questions: ${questions.length} next:${nexQ}`);
    totalQuestions = questions.length;
    let currQuestion = document.querySelector(`${templateClass}.i_${curQ}`);
    let nextQuestion = document.querySelector(`${templateClass}.i_${nexQ}`);

    // Stop the current stopwatch
    const currStopwatch = currQuestion.querySelector('.stopwatch');
    if (currStopwatch) {
        stopStopwatch(currStopwatch.id);
    }

    //Hide the current question
    if (back) {
        currQuestion.style.animation = 'slideOutToRight 0.6s forwards';
        setTimeout(() => {
            nextQuestion.style.display = 'block';
            nextQuestion.style.animation = 'slideInFromLeft 0.6s backwards';
        }, 600);
    } else {
        currQuestion.style.animation = 'slideOutToLeft 0.6s forwards';
        setTimeout(() => {
            nextQuestion.style.display = 'block';
            nextQuestion.style.animation = 'slideInFromRight 0.6s backwards';
        }, 600);
    }
    
    setTimeout(() => {
        currQuestion.style.display = 'none';
    }, 600);
    
    //progressIndicator.textContent = `${currentQuestion + 1} / ${questions.length}`;

    // Start the next stopwatch
    const nextStopwatch = nextQuestion.querySelector('.stopwatch');
    if (nextStopwatch) {
        const index = extractIndexFromClassName(nextQuestion.className);
        startStopwatch(nextStopwatch.id,index);
    }
}

function nextQuestion(itemSelector) {
    let templateClass = itemSelector;
    let questions = document.querySelectorAll(templateClass);
    
    if (currentQuestion === questions.length - 1) {
        currentQuestion = 0;
        nextExercise();
        totalQuestions = getTotalQuestions();
        updateProgressBar();
        return;
    }
    currentQuestion = (currentQuestion + 1) % questions.length;
    updateQuestion(false, templateClass, currentQuestion);
    updateProgressBar();

    // Call storeAnswer function
    // storeAnswer(button, itemSelector, question, index);
}

let answers = [];
function storeAnswer(button, itemSelector, question, index) {
    let item = button.closest(itemSelector); 
    
    // Determine the selected word
    let selectedWord = null;
    let wordElements = item.querySelectorAll('.speak-word');
    wordElements.forEach(word => {
        if (word.classList.contains('selected')) {
            selectedWord = word.textContent;
        }
    });
    //item: item.outerHTML
    answers.push({question:question, index: index, selectedWord: selectedWord });
    console.log(answers);
} 

let answeredPairs = [];
function storeAnswersPairs(selector, question) {
    // Get the container element using the provided selector
    var container = document.querySelector(selector);

    // Loop through each dropzone in the container
    container.querySelectorAll('.dropzone').forEach(function(dropzone) {
        var word1 = dropzone.querySelector('.speak-word').innerText;
        var word2 = dropzone.querySelector('.speak-word[draggable="true"]').innerText;

        // var pair = {};
        // pair[word1] = word2;

        answeredPairs.push({question:question, "column1":word1, "column2":word2});
    });

    // Display the pairs in the console for testing
    console.log(answeredPairs);

    // You can now do whatever you want with the pairs array, like sending it to the server
}

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

function selectAnswer(element) {
    let parent = element.parentNode;
    let sameClassElements = parent.getElementsByClassName(element.className);
    for (let i = 0; i < sameClassElements.length; i++) {
        sameClassElements[i].classList.remove('selected');
    }
    element.classList.add('selected');
}

function drag(ev) {
    ev.dataTransfer.setData("text", ev.target.id);
}

function allowDrop(ev) {
    ev.preventDefault();
}

function drop(ev) {
    ev.preventDefault();
    var data = ev.dataTransfer.getData("text"); 
    ev.target.appendChild(document.getElementById(data));    
}

// Add event listeners to the dropzones
var dropzones = document.querySelectorAll('.dropzone');
dropzones.forEach(function(dropzone) {
    dropzone.addEventListener('dragover', allowDrop);
    dropzone.addEventListener('drop', drop);
});

updateProgressBar();

// Function to calculate progress
function calculateProgress(current, total) {
    return (current / total) * 100;
}

// Function to update progress bar
function updateProgressBar() {
    let questionProgress = calculateProgress(currentQuestion, totalQuestions);
    document.getElementById('questionProgressBar').style.width = questionProgress + '%';
    document.getElementById('exerciseProgressText').innerText = `Άσκηση: ${currentExercise+1} / ${totalExercises}`;
    document.getElementById('questionProgressText').innerText = `Ερώτηση: ${currentQuestion+1} / ${totalQuestions}`;
}



// Function to start the stopwatch for a given template
function startStopwatchesForTemplate(templateClass) {
    // Start the initial stopwatches
    document.querySelectorAll(templateClass).forEach(item => {
        const stopwatch = item.querySelector('.stopwatch');
        if (stopwatch && item.style.display !== 'none') {
            const className = item.className; // Get the class attribute
            const index = extractIndexFromClassName(className); // Extract the index
            if (index !== null) {
                startStopwatch(stopwatch.id, index); // Pass the id and extracted index
            }
        }
    });
}


// Start the stopwatches for DecodingTextReading template
startStopwatchesForTemplate('.decoding-text-reading-item');


// Start the stopwatches for MemoryVisualSingleChoice template
startStopwatchesForTemplate('.memory-visual-single-choice-item');

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// Function to initialize voice recognition
function initializeVoiceRecognition() {
    let recognition;
    const recordedChunks = [];

    const startRecording = () => {
        recordedChunks.length = 0;
        recognition = new webkitSpeechRecognition();
        recognition.lang = 'el-GR'; // Set language to Greek

        recognition.onresult = (event) => {
            recordedChunks.push(event.results[0][0].transcript);
            document.querySelector('.recognized-text').textContent = recordedChunks.join(' ');
        };

        recognition.start();
    };

    const stopRecording = () => {
        if (recognition) {
            recognition.stop();
            const blob = new Blob(recordedChunks, { type: 'text/plain' });
            const url = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'recorded_audio.txt';
            document.body.appendChild(a);
            a.click();
            setTimeout(() => {
                document.body.removeChild(a);
                URL.revokeObjectURL(url);
            }, 0);
        }
    };

    document.querySelector('.start-voice-recognition').addEventListener('click', startRecording);
    document.querySelector('.stop-voice-recognition').addEventListener('click', stopRecording);
}

// Initialize voice recognition when the DOM content is loaded
document.addEventListener("DOMContentLoaded", function() {
    // Check if the desired template is present
    const template = document.querySelector('.phonologic-rhyme-production-sentence.template');
    if (template!=null && template.length > 0) {
        // Run the initialization function
        initializeVoiceRecognition();
    }
});
///////////////////////////////////////////////////////////////////////////////////////////////////////////////


/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
function switchToImages() {
    let templateDivs = document.querySelectorAll('.phonologic-rhyme-match.template');
    if (templateDivs.length === 0) {
        console.error('No phonologic-rhyme-match.template elements found');
        return;
    }

    let question = templateDivs[0].querySelector('h1')?.innerHTML;
    if (!question) {
        console.error('No question found in the first template div');
        return;
    }

    //let visibleTemplateDivs = Array.from(templateDivs).filter(div => window.getComputedStyle(div).display !== 'none');
    //let childDivs = Array.from(visibleTemplateDivs[0].children).filter(child => child.tagName === 'DIV');
    console.log(`Question: ${question}`);
    if (question.includes("εικόνα")) {
        let answerDivs = document.querySelectorAll('.phonologic-rhyme-match-answers');
        answerDivs.forEach(answerDiv => {
            let speakWordDivs = Array.from(answerDiv.children).filter(child => child.className === 'speak-word');

            speakWordDivs.forEach(div => {
                let answer = div.innerText.trim();
                let originalHTML = div.innerHTML;
                let img = document.createElement('img');
                img.src = `/assets/words/${answer}.png`;
                img.alt = answer;
                img.onerror = function() {
                    console.error(`Image /assets/words/${answer}.png does not exist.`);
                    div.innerHTML = originalHTML; // restore original HTML if image doesn't exist
                };
                img.onload = function() {
                    div.innerHTML = ''; // clear the div's content when image loads successfully
                    div.appendChild(img); // append the image to the div
                };
            });
        });
    }
}


document.addEventListener("DOMContentLoaded", function() {
    // Check if the desired template is present
    const template = document.querySelectorAll('.phonologic-rhyme-match.template')
    if (template!=null && template.length > 0) {
        // console.log(template)
        // Run the initialization function
        switchToImages();
    }
});
////////////////////////////////////////////////////////////////////////////////////////////////////////////


/////////////////////////////////////////////////////////////////////////////////////////////////////////////

const displayTime = 4000; // 4 seconds

function updateMemoryChoiceTemplate(choice, automatic, nextButton, templateIndex) {
    const items = document.querySelectorAll(`.memory-visual-single-choice-item.t_${templateIndex}.i_${currentQuestion}`);
    const item = items[0];

    const choice1 = item.querySelectorAll('.memory-visual-single-choice-choice1')[0];
    const choice2 = item.querySelectorAll('.memory-visual-single-choice-choice2')[0];
    
    if (!item) return;
    

    // const choice1 = item.querySelector('.memory-visual-single-choice-choice1');
    // const choice2 = item.querySelector('.memory-visual-single-choice-choice2');
    

    if (choice === 1) {
        choice1.style.display = 'block';
        choice2.style.display = 'none';
        nextButton.forEach((button, index) => {
            button.style.display = 'none';
        });
        
        if (automatic) {
            setTimeout(() => {
                updateMemoryChoiceTemplate(2, automatic, nextButton, templateIndex);
            }, displayTime);
        }
    } else {
        choice1.style.display = 'none';
        choice2.style.display = 'block';
        nextButton.forEach((button, index) => {
            button.style.display = 'block';
        });
    }
}

document.addEventListener("DOMContentLoaded", function() {
    // Define displayTime within a scope accessible to updateMemoryChoiceTemplat
    let showMemorySingleChoiceTimer;

    function nextQuestionMemory(itemSelector, templateIndex, nextButton) {
        clearTimeout(showMemorySingleChoiceTimer);

        let templateClass = itemSelector;
        let questions = document.querySelectorAll(templateClass);
        
        if (currentQuestion === questions.length - 1) {
            currentQuestion = 0;
            templateIndex ++;

            nextExercise();
            
            // style in a block the following exercise
            const items = document.querySelectorAll(`.memory-visual-single-choice-item.t_${templateIndex}.i_${currentQuestion}`);
            const item = items[0];
            item.style.display = 'block'
            updateMemoryChoiceTemplate(1, true, nextButtons, templateIndex);

            totalQuestions = getTotalQuestions();
            updateProgressBar();
            return;
        } 
        
        currentQuestion = (currentQuestion + 1) % questions.length;
        
        updateMemoryChoiceTemplate(1, true, nextButton, templateIndex);
        updateQuestion(false, itemSelector, currentQuestion);
        updateProgressBar();
    }

    // Check if the desired template is present
    const template = document.querySelectorAll('.memory-visual-single-choice.template');
    const nextButtons = document.querySelectorAll('.nextButton-textreading');
    

    if (template != null && template.length > 0) {
        // Ensure items are hidden initially
        const items = document.querySelectorAll('.memory-visual-single-choice-item');
        items.forEach((item, index) => {
            item.style.display = 'none';
        });
        

        // Display the first item
        const firstItem = items[currentQuestion];
        if (firstItem) {
            firstItem.style.display = 'block';
            var templateIndex = 0;
            updateMemoryChoiceTemplate(1, true, nextButtons, templateIndex);
        }
    }

    // Bind nextQuestion function to the next button
    nextButtons.forEach(button => {
        button.addEventListener('click', function() {
            let templateIndex = button.closest('.template').getAttribute('data-template-index');
            const itemSelector = `.memory-visual-single-choice-item.t_${templateIndex}`;
            button.style.display = 'none';
            nextQuestionMemory(itemSelector, templateIndex, nextButtons);
        });
    });

});

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
let START_EXERCISE = 0;
let answers = [];

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
    
function submitAnswers() {
    fetch('/answers', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(answers)
    }).then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    }).then(data => {
        // Handle the response data here
    }).catch(error => {
        console.error('There has been a problem with your fetch operation:', error);
    });
}


// Function to extract the index from the class name
function extractIndexFromClassName(className) {
    const match = className.match(/i_(\d+)/);
    return match ? parseInt(match[1], 10) : null;
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


    // Edw orizv na mpainw stin updateMemoryChoiceTemplate 1h fora meta apo allou typou askhsh
    if (templateDivs[currentExercise].className==='RecallImageChoice template') {
        if(currentQuestion===0){
            initializeRecallImageChoice();
        }
    }
    
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


function chosenOption(templateClass) {
    // Find the selected option
    var selectedAnswer = document.querySelector(`${templateClass}.i_${currentQuestion} .selected`);
    console.log("selectedAnswer = ", selectedAnswer)
    console.log("templateClass = ", templateClass)
    
    // Define an object to store the selected option and image name
    var result = {
        mainWordText: null,
        imageName: null
    };

    // Deselect all options
    // var options = document.querySelectorAll('.option');
    // options.forEach(function(opt) {
    //     opt.classList.remove('selected');
    // });


    if (selectedAnswer) {
        selectedAnswer.classList.add('selected');

        // Retrieve the main word text
        var mainWordTextElement = selectedAnswer.closest('.SelectOption-item').querySelector('.main-sentence, .speak-word');
        if (mainWordTextElement) {
            result.mainWordText = mainWordTextElement.textContent.trim();
        }
        
        console.log("selectedAnswer = ", selectedAnswer)
        var isImage = selectedAnswer.querySelector('img') !== null;

        if (isImage) {
            var selectedImageSrc = selectedAnswer.querySelector('img').getAttribute('src');
            var imageName = selectedImageSrc.substring(selectedImageSrc.lastIndexOf('/') + 1);
            console.log("Selected image name:", imageName);
            result.imageName = decodeURIComponent(imageName);
            
        } else {
            console.log("Selected item is not an image");
        }
    }

    return result;
}



function nextQuestion(itemSelector, mainwordtext, answer) {
    let templateClass = itemSelector;
    let questions = document.querySelectorAll(templateClass);
    var selectedAnswer = document.querySelector(`${templateClass}.i_${currentQuestion} .selected`);

    //Check if next button is pressed without give an answer
    if ((templateClass.startsWith(".AudioRecord")) && (audioData === null)) {
        alert("Please record your answer");
        return;
    }
    if ((templateClass.startsWith(".SelectOption")) && (selectedAnswer===null)) {
        alert("Please give an answer");
        return;
    }
    //if (templateClass.startsWith(".CorrectOption"))
    //     ((!templateClass.startsWith(".CorrectOption")) && (!selectedAnswer))) {
    //     alert("Please answer");
    //     return;
    // }

    storeAnswerNew(templateClass, mainwordtext, answer);

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


    // if (templateClass === '.memory-visual-single-choice-item') {
    //     storeAnswerCorrectOption(questions[currentQuestion].querySelector('h1').textContent, currentQuestion, selectedAnswer.textContent);
    // }
    // storeAnswerCorrectOption
}

function storeAnswerNew(templateClass, mainwordtext, answer) {
    let qdiv = document.querySelector(`${templateClass}.i_${currentQuestion}`);
    var question = qdiv.parentElement.querySelector('h1').textContent;
    
    if (templateClass.startsWith(".CorrectOption")) {
        // storeAnswerOption(question, mainwordtext, answer);
        answers.push({question:question, mainwordtext: mainwordtext, answer: answer});
    } 
    else if (templateClass.startsWith(".AudioRecord")) {
        answers.push({question:question, mainwordtext: mainwordtext, answer: audioData});
        // storeAnswerAudioRecord(question, currentQuestion, audioData);
        //submitAnswers();
        //audioData = null;
    }
    else if (templateClass.startsWith(".SelectOption") || templateClass.startsWith(".RecallImageChoice") ) {
        var result = chosenOption(templateClass);
        var imageName = result.imageName;
        var mainwordoption = result.mainWordText;
        console.log("result = ", result)
        console.log("imageName = ", imageName)
        console.log("mainwordoption = ", mainwordoption)
        answers.push({question:question, mainwordtext: mainwordoption, answer: imageName});
        // storeAnswerOption(question, mainwordoption, imageName)
    }

    console.log("answers = ", answers)
    
}


function storeAnswerAudioRecord(question, mainwordtext, answer) {
    answers.push({question:question, mainwordtext: mainwordtext, answer: answer});
}

function storeAnswerOption(question, mainwordtext, answer) {
    answers.push({question:question, mainwordtext: mainwordtext, answer: answer});
}

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


function selectOption(option) {
    // Deselect all options
    var options = document.querySelectorAll('.option');
    options.forEach(function(opt) {
        opt.classList.remove('selected');
    });
    
    // Select the clicked option
    option.classList.add('selected');
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
    let selectedCount = 0;
    let cName = element.className
    // For this specific template, we need to allow 2 selected options
    if (cName.includes('RecallImageChoice-choice2-options')) {
        // If the clicked element is already selected, remove the 'selected' class
        if (element.classList.contains('selected')) {
            element.classList.remove('selected');
            return;
        }
        for (let i = 0; i < sameClassElements.length; i++) {
            if (sameClassElements[i].classList.contains('selected')) {
                selectedCount++;
            }
        }
        if (selectedCount < 2) {
            element.classList.add('selected');
        }
    } else {
    for (let i = 0; i < sameClassElements.length; i++) {
        sameClassElements[i].classList.remove('selected');
    }
    element.classList.add('selected');
    }
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

// Define your updateMemoryChoiceTemplate function
function updateMemoryChoiceTemplate(choice, automatic, nextButtons) {
    const items = document.querySelectorAll(`.RecallImageChoice-item.i_${currentQuestion}`);
    const item = items[0];

    if (!item) return;

    const choice1 = item.querySelector('.RecallImageChoice-choice1');
    const choice2 = item.querySelector('.RecallImageChoice-choice2');
    
    if (choice === 1) {
        choice1.style.display = 'block';
        choice2.style.display = 'none';
        nextButtons.forEach(button => {
            button.style.display = 'none';
        });
        
        if (automatic) {
            setTimeout(() => {
                updateMemoryChoiceTemplate(2, automatic, nextButtons);
            }, displayTime);
        }
    } else {
        choice1.style.display = 'none';
        choice2.style.display = 'block';
        nextButtons.forEach(button => {
            button.style.display = 'block';
        });
    }
}

// Define your nextQuestionMemory function
function nextQuestionMemory(itemSelector, mainwordtext, answer) {
    console.log("mainwordtext = ", mainwordtext)
    let templateClass = itemSelector;
    let questions = document.querySelectorAll(templateClass);
    const nextButtons = document.querySelectorAll('.nextButton-textreading');
    var selectedAnswer = document.querySelector(`${templateClass}.i_${currentQuestion} .selected`);

    if ((templateClass.startsWith(".RecallImageChoice")) && (selectedAnswer===null)) {
        alert("Please give an answer");
        return;
    }

    storeAnswerNew(templateClass, mainwordtext, answer);

    if (currentQuestion === questions.length - 1) {
        currentQuestion = 0;
        
        nextExercise(); // You need to define this function

        // Display the first item of the next template
        const items = document.querySelectorAll(`.RecallImageChoice-item.i_${currentQuestion}`);
        const item = items[0];
        item.style.display = 'block';
        updateMemoryChoiceTemplate(1, true, nextButtons);

        updateProgressBar(); // You need to define this function
        return;
    }
    
    currentQuestion = (currentQuestion + 1) % questions.length;
    console.log()
    updateMemoryChoiceTemplate(1, true, nextButtons);
    updateQuestion(false, itemSelector, currentQuestion); // You need to define this function
    updateProgressBar(); // You need to define this function
}

// Define your initialization function
function initializeRecallImageChoice() {
    // Check if the desired template is present
    const template = document.querySelectorAll('.RecallImageChoice-item');
    const nextButtons = document.querySelectorAll('.nextButton-textreading');

    
    if (template.length > 0) {
        // Ensure items are hidden initially
        template.forEach(item => {
            item.style.display = 'none';
        });

        // Display the first item
        const firstItem = template[currentQuestion];
        if (firstItem) {
            firstItem.style.display = 'block';
            updateMemoryChoiceTemplate(1, true, nextButtons);
        }
    }

    // // Bind nextQuestion function to the next button
    // nextButtons.forEach(button => {
    //     button.addEventListener('click', function() {
    //         const itemSelector = `.RecallImageChoice-item`;
    //         button.style.display = 'none';
    //         nextQuestionMemory(itemSelector, nextButtons);
    //     });
    // });
}


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

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
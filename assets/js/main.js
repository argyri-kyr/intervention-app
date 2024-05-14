let START_EXERCISE = 5;

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
    console.log(templateDivs);
    for(let i = 0; i < templateDivs.length; i++) {
        if(i === selectedIndex) {
            templateDivs[i].style.display = 'block';
        } else {
            templateDivs[i].style.display = 'none';
        }
    }
}


hideExercises();

function nextExercise() {
    console.log('Next exercise');
    let templateDivs = document.querySelectorAll('.template');
    let currentTemplate = 0;
    for(let i = 0; i < templateDivs.length; i++) {
        if (templateDivs[i].style.display === 'block') {
            currentTemplate = i;
            console.log(currentTemplate);
            break;
        }
    }
    templateDivs[currentTemplate].style.display = 'none';
    templateDivs[(currentTemplate + 1) % templateDivs.length].style.display = 'block';
}

let currentQuestion = 0;
function updateQuestion(back = false, templateClass) {
    let questions = document.querySelectorAll(templateClass);
    console.log(`Template: ${templateClass}`);
    let curQ, nexQ = 0;
    if (back) {
        curQ = currentQuestion + 1;
        nexQ = currentQuestion;
    } else {
        curQ = currentQuestion - 1;
        nexQ = currentQuestion;
    }
    console.log(`Questions: ${questions.length} next:${nexQ}`);
    let currQuestion = document.querySelector(`${templateClass}.i_${curQ}`);
    let nextQuestion = document.querySelector(`${templateClass}.i_${nexQ}`);

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
}

function nextQuestion(itemSelector) {
    let templateClass = itemSelector;
    console.log(templateClass);
    let questions = document.querySelectorAll(templateClass);
    if (currentQuestion === questions.length - 1) {
        currentQuestion = 0;
        nextExercise();
        return;
    }
    currentQuestion = (currentQuestion + 1) % questions.length;
    updateQuestion(false, templateClass);


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

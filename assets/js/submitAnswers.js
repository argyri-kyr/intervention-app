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

function storeAnswer(templateClass, mainwordtext, answer) {
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
    else if (templateClass.startsWith(".SelectOption")) {
        var result = chosenOption(templateClass);
        console.log("result = ", result)
        answers.push({question: question, mainwordtext: mainwordtext, answer: result});
        // storeAnswerOption(question, mainwordoption, imageName)
    }
    else if (templateClass.startsWith(".MatchMultipleOptions")) {
        answers.push({question:question, mainwordtext: mainwordtext, answer: matchedPairs});
        matchedPairs = {};
        // storeAnswerOption(question, mainwordoption, imageName)
    }
    else if (templateClass.startsWith(".RecallImageChoice")) {
        var result = chosenOption(templateClass);
        answers.push({question: question, mainwordtext: mainwordtext, answer: result.join(',')});
    }

    console.log("answers = ", answers)
}

function storeAnswerAudioRecord(question, mainwordtext, answer) {
    answers.push({question:question, mainwordtext: mainwordtext, answer: answer});
}

function storeAnswerOption(question, mainwordtext, answer) {
    answers.push({question:question, mainwordtext: mainwordtext, answer: answer});
}

// function storeAnswer(button, itemSelector, question, index) {
//     let item = button.closest(itemSelector); 
    
//     // Determine the selected word
//     let selectedWord = null;
//     let wordElements = item.querySelectorAll('.speak-word');
//     wordElements.forEach(word => {
//         if (word.classList.contains('selected')) {
//             selectedWord = word.textContent;
//         }
//     });
    
//     answers.push({question:question, index: index, selectedWord: selectedWord });
//     console.log(answers);
// } 

function storeAnswersPairs(selector, question) {
    // Get the container element using the provided selector
    var container = document.querySelector(selector);

    // Loop through each dropzone in the container
    container.querySelectorAll('.dropzone').forEach(function(dropzone) {
        var word1 = dropzone.querySelector('.speak-word').innerText;
        var word2 = dropzone.querySelector('.speak-word[draggable="true"]').innerText;

        answeredPairs.push({question:question, "column1":word1, "column2":word2});
    });

    // Display the pairs in the console for testing
    console.log(answeredPairs);

    // You can now do whatever you want with the pairs array, like sending it to the server
}

function submitAnswers(gameNum) {
    let pathArray = window.location.pathname.split('/');
    let lastPath = pathArray[pathArray.length - 1];
    console.log('Submitting answers for session ', lastPath);
    console.log('Submitting answers for game ', gameNum);
    console.log('Answers: ', JSON.stringify(answers));
    fetch(`/answers/${lastPath}/${gameNum}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(answers)
    }).then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.text();
    }).then(data => {
        if (data) {
            return JSON.parse(data);
        } else {
            console.log('No data returned from server');
        }
    }).catch(error => {
        console.error('There has been a problem with your fetch operation:', error);
    });
    //RESET ANSWERS
    answers = [];
}


// Function to retrieve matched pairs
function getMatchedPairs() {
    return matchedPairs;
}

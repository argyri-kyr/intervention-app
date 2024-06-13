
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
    submitAnswers(currentExercise);

    // Edw orizv na mpainw stin updateMemoryChoiceTemplate 1h fora meta apo allou typou askhsh
    if (templateDivs[currentExercise].className==='RecallImageChoice template') {
        if(currentQuestion===0){
            initializeRecallImageChoice();
        }
    }
}

function updateQuestion(back = false, templateClass, currentQuestion) {
    // let questions = document.querySelectorAll(templateClass);
    let parentTemplateName = templateClass.split('-')[0];
    let parentTemplates = document.querySelectorAll(parentTemplateName);

    // Iterate over each parent template to find the visible one
    let visibleParentTemplate = Array.from(parentTemplates).find(parentTemplate => {
        return window.getComputedStyle(parentTemplate).display !== "none";
    });

    questions = visibleParentTemplate.querySelectorAll(templateClass);
    
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
    let currQuestion = visibleParentTemplate.querySelector(`${templateClass}.i_${curQ}`);
    let nextQuestion = visibleParentTemplate.querySelector(`${templateClass}.i_${nexQ}`);

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


function chosenOption(itemSelector, templateClass) {
    // Find the selected option
    var selectedAnswers = document.querySelectorAll(`${templateClass}.i_${currentQuestion} .selected`);
    
    // Define an object to store the selected option and image name
    var result = []

    if (selectedAnswers) {
        // selectedAnswer.classList.add('selected');
        selectedAnswers.forEach(selectedAnswer => {

            selectedAnswer.classList.add('selected');
            var isImage = selectedAnswer.querySelector('img') !== null;

            if (isImage) {
                var selectedImageSrc = selectedAnswer.querySelector('img').getAttribute('src');
                var imageName = selectedImageSrc.substring(selectedImageSrc.lastIndexOf('/') + 1);
                imageName = decodeURIComponent(imageName);
                result.push(imageName)
                
            } else {
                console.log("Selected item is not an image");
            }
        
        });
    }

    return result;
}



function nextQuestion(itemSelector, mainwordtext, answer) {
    let templateClass = itemSelector;
    var selectedAnswer = document.querySelector(`${templateClass}.i_${currentQuestion} .selected`);

    let parentTemplateName = itemSelector.split('-')[0];
    let parentTemplates = document.querySelectorAll(parentTemplateName);

    // Iterate over each parent template to find the visible one
    let visibleParentTemplate = Array.from(parentTemplates).find(parentTemplate => {
        return window.getComputedStyle(parentTemplate).display !== "none";
    });

    questions = visibleParentTemplate.querySelectorAll(itemSelector);

    console.log("templateClass = ", templateClass)
    console.log("visibleParentTemplate = ", visibleParentTemplate)

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

    storeAnswerNew(itemSelector, templateClass, mainwordtext, answer);

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
    let templateClass = itemSelector;
    let questions = document.querySelectorAll(templateClass);
    const nextButtons = document.querySelectorAll('.nextButton-textreading');
    var selectedAnswer = document.querySelector(`${templateClass}.i_${currentQuestion} .selected`);

    if ((templateClass.startsWith(".RecallImageChoice")) && (selectedAnswer===null)) {
        alert("Please give an answer");
        return;
    }

    storeAnswerNew(itemSelector, templateClass, mainwordtext, answer);

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


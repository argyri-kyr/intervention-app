
function getTotalQuestions() {
    let templateDivs = document.querySelectorAll('.template');
    let visibleTemplateDivs = Array.from(templateDivs).filter(div => window.getComputedStyle(div).display !== 'none');
    let childDivs = Array.from(visibleTemplateDivs[0].children).filter(child => child.tagName === 'DIV');

    return childDivs.length;
}

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

// Function to extract the index from the class name
function extractIndexFromClassName(className) {
    const match = className.match(/i_(\d+)/);
    return match ? parseInt(match[1], 10) : null;
}


// Function to handle drag start
function drag(ev) {
    ev.dataTransfer.setData("text", ev.target.id);
}

// Function to allow drop
function allowDrop(ev) {
    ev.preventDefault();
}

// Function to handle drop
function drop(ev) {
    ev.preventDefault();
    var data = ev.dataTransfer.getData("text");
    var dropTarget = ev.target;

    // Check if the drop area already contains a child element
    if (dropTarget.children.length === 1) {
        // Append the dropped element
        dropTarget.appendChild(document.getElementById(data));

        storePairsFromDrop(dropTarget);

    } else {
        // Prevent dropping another element
        // console.log("Drop area already contains an element.");
    }
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

function selectOption(option) {
    // Deselect all options
    var options = document.querySelectorAll('.option');
    options.forEach(function(opt) {
        opt.classList.remove('selected');
    });
    
    // Select the clicked option
    option.classList.add('selected');
}

function storePairsFromDrop(dropTarget) {
    var texts = Array.from(dropTarget.children).map(child => {
        var onclickValue = child.getAttribute('onclick');
        var textMatch = onclickValue.match(/speak\('([^']+)'\)/);
        return textMatch ? textMatch[1] : '';
    });

    // Iterate over matchedPairs to find and delete the matching pair
    for (let key in matchedPairs) {
        if (matchedPairs[key] === texts[1]) {
            delete matchedPairs[key];
        }
    }

    matchedPairs[texts[0]] = texts[1];

}
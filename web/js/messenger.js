const eventSource = new EventSource('/events');

eventSource.onmessage = function(event) {
    let str = event.data
    let newStr = str.replace(/\\n/g, '\n');
    document.querySelector("textarea[name='message']").value = newStr;
};
let timeout = null;

const textarea = document.getElementById('message');

textarea.addEventListener('keyup', function() {
    clearTimeout(timeout);
    timeout = setTimeout(() => {
        sendMessage(textarea.value);
    }, 500);
});

function sendMessage(message) {
    fetch('/', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({ message: message })
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.text();
        })
}

document.getElementById("copyButton").addEventListener("click", function() {
    const textarea = document.getElementById("message");

    navigator.clipboard.writeText(textarea.value)
});

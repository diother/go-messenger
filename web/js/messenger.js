const eventSource = new EventSource('/events');

eventSource.onmessage = function(event) {
    console.log("New message: ", event.data);
    document.querySelector("textarea[name='message']").value = event.data;
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
        .then(data => {
            console.log('Success:', data);
        })
        .catch((error) => {
            console.error('Error:', error);
        });
}

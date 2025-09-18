document.getElementById('notificationForm').addEventListener('submit', async function(event) {
    event.preventDefault();

    const text = document.getElementById('text').value.trim();
    const telegramId = parseInt(document.getElementById('telegram_id').value.trim(), 10);
    const sendAtInput = document.getElementById('send_at').value;

    if (!text || isNaN(telegramId) || !sendAtInput) {
        displayResponse('Please fill in all fields correctly.', true);
        return;
    }


    const sendAt = sendAtInput + ":00+03:00";

    const payload = {
        text: text,
        telegram_id: telegramId,
        send_at: sendAt
    };

    try {
        const response = await fetch('/notify', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(payload)
        });

        if (response.ok) {
            const data = await response.json();
            displayResponse('Notification created successfully! ID: ' + data.id, false);
            this.reset();
        } else {
            const errorData = await response.json();
            displayResponse('Error: ' + (errorData.error || 'Unknown error'), true);
        }
    } catch (error) {
        displayResponse('Network error: ' + error.message, true);
    }
});

function displayResponse(message, isError) {
    const responseDiv = document.getElementById('response');
    responseDiv.textContent = message;
    responseDiv.style.color = isError ? 'red' : 'green';
}

document.getElementById('loadNotificationsBtn').addEventListener('click', async function() {
    const container = document.getElementById('notificationsContainer');
    container.innerHTML = 'Loading notifications...';
    container.style.color = 'black';

    try {
        const response = await fetch('/notify');
        if (response.ok) {
            const notifications = await response.json();
            if (notifications.length === 0) {
                container.textContent = 'No notifications found.';
                return;
            }

            const list = document.createElement('ul');
            list.style.listStyleType = 'none';
            list.style.padding = '0';

            notifications.forEach(notif => {
                const item = document.createElement('li');
            item.classList.add('notification-card');

            const text = document.createElement('p');
            const sendAtDate = new Date(parseInt(notif.send_at)).toLocaleString();
            text.innerHTML = `ID: ${notif.id}<br>Text: ${notif.text}<br>Telegram ID: ${notif.telegram_id}<br>Send At: ${sendAtDate}<br>Status: ${notif.status}`;
            text.classList.add('notification-text');

            item.appendChild(text);
            list.appendChild(item);
            });

            container.innerHTML = '';
            container.appendChild(list);
        } else {
            const errorData = await response.json();
            container.textContent = 'Error: ' + (errorData.error || 'Unknown error');
            container.style.color = 'red';
        }
    } catch (error) {
        container.textContent = 'Network error: ' + error.message;
        container.style.color = 'red';
    }
});

document.getElementById('cancelNotificationForm').addEventListener('submit', async function(event) {
    event.preventDefault();

    const cancelId = parseInt(document.getElementById('cancel_id').value.trim(), 10);
    const cancelResponseDiv = document.getElementById('cancelResponse');

    if (isNaN(cancelId)) {
        cancelResponseDiv.textContent = 'Please enter a valid notification ID.';
        cancelResponseDiv.style.color = 'red';
        return;
    }

    try {
        const response = await fetch('/notify/' + cancelId, {
            method: 'DELETE'
        });

        if (response.ok) {
            const data = await response.json();
            cancelResponseDiv.textContent = data.status || 'Notification cancelled successfully.';
            cancelResponseDiv.style.color = 'green';
            this.reset();
        } else {
            const errorData = await response.json();
            cancelResponseDiv.textContent = 'Error: ' + (errorData.error || 'Unknown error');
            cancelResponseDiv.style.color = 'red';
        }
    } catch (error) {
        cancelResponseDiv.textContent = 'Network error: ' + error.message;
        cancelResponseDiv.style.color = 'red';
    }
});

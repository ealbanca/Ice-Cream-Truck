// Event form submission
const eventForm = document.getElementById('event-scheduler-form');
if (eventForm) {
    eventForm.onsubmit = async function(e) {
        e.preventDefault();
        const data = {
            date: document.getElementById('event-date').value,
            time: document.getElementById('event-time').value,
            description: document.getElementById('event-description').value
        };
        const res = await fetch('/api/events', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(data)
        });
        if (res.ok) {
            alert('Event scheduled!');
            this.reset();
        } else {
            alert('Error scheduling event');
        }
    };
}

// Contact form submission
const contactForm = document.querySelector('.contact-form');
if (contactForm) {
    contactForm.onsubmit = async function(e) {
        e.preventDefault();
        const data = {
            name: document.getElementById('name').value,
            email: document.getElementById('email').value,
            phone: document.getElementById('phone').value,
            reason: document.getElementById('reason').value,
            message: document.getElementById('yourmessage').value
        };
        const res = await fetch('/api/contact', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(data)
        });
        if (res.ok) {
            alert('Message sent!');
            this.reset();
        } else {
            alert('Error sending message');
        }
    };
}


// Disable Sundays in the date picker
document.addEventListener('DOMContentLoaded', function() {
    var dateInput = document.getElementById('event-date');
    if (dateInput) {
        dateInput.addEventListener('input', function() {
            var d = new Date(this.value);
            // Only show alert if the date is valid and the day is Sunday (0)
            if (this.value && d.getDay() === 6) {
                alert('Scheduling on Sunday is not allowed. Please select another date.');
                this.value = '';
            }
        });
    }
});

// Dynamically populate start/end time options based on day
function setTimeOptions() {
    const dateInput = document.getElementById('event-date');
    const startSelect = document.getElementById('event-start-time');
    const endSelect = document.getElementById('event-end-time');
    if (!dateInput || !startSelect || !endSelect) return;
    let date = dateInput.value;
    let day = date ? new Date(date).getDay() : null;
    let startTimes = [], endTimes = [];
    if (day === 5) { // Saturday (JS: 5)
        for (let h = 10; h <= 21; h++) {
            let label = (h < 12 ? h + ' am' : (h === 12 ? '12 pm' : (h-12) + ' pm'));
            let val = (h < 10 ? '0' : '') + h + ':00';
            startTimes.push({val, label});
        }
        for (let h = 10; h <= 21; h++) {
            let label = (h < 12 ? h + ' am' : (h === 12 ? '12 pm' : (h-12) + ' pm'));
            let val = (h < 10 ? '0' : '') + h + ':00';
            endTimes.push({val, label});
        }
    } else if (day >= 0 && day <= 5) { // Mon-Fri (Mon=1, Fri=5)
        for (let h = 17; h <= 21; h++) {
            let label = (h === 12 ? '12 pm' : (h > 12 ? (h-12) + ' pm' : h + ' am'));
            let val = h + ':00';
            startTimes.push({val, label});
        }
        for (let h = 17; h <= 21; h++) {
            let label = (h === 12 ? '12 pm' : (h > 12 ? (h-12) + ' pm' : h + ' am'));
            let val = h + ':00';
            endTimes.push({val, label});
        }
    } else { // Sunday (0) or no date
        startTimes = [{val: '', label: 'Select a date first'}];
        endTimes = [{val: '', label: 'Select a date first'}];
    }
    startSelect.innerHTML = startTimes.map(t => `<option value="${t.val}">${t.label}</option>`).join('');
    endSelect.innerHTML = endTimes.map(t => `<option value="${t.val}">${t.label}</option>`).join('');
}
document.getElementById('event-date').addEventListener('change', setTimeOptions);
window.addEventListener('DOMContentLoaded', setTimeOptions);

// Event form submission
const eventForm = document.getElementById('event-scheduler-form');
if (eventForm) {
    eventForm.onsubmit = async function(e) {
        e.preventDefault();
        const data = {
            name: document.getElementById('event-name').value,
            email: document.getElementById('event-email').value,
            phone: document.getElementById('event-phone').value,
            date: document.getElementById('event-date').value,
            start_time: document.getElementById('event-start-time').value,
            end_time: document.getElementById('event-end-time').value,
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
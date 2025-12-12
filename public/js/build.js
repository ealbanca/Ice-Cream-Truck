// --- Size selection (single) ---
let maxFlavors = 1;
const sizeChoices = document.querySelectorAll('#size-choices .build-choice');
sizeChoices.forEach(el => {
    el.addEventListener('click', function() {
        sizeChoices.forEach(c => c.classList.remove('selected'));
        this.classList.add('selected');
        document.getElementById('size-input').value = this.dataset.value;
        // Set maxFlavors based on size
        const label = this.dataset.label.toLowerCase();
        if (label === 'small') maxFlavors = 1;
        else if (label === 'medium') maxFlavors = 2;
        else if (label === 'large') maxFlavors = 3;
        updateFlavorLimitMsg();
        clearFlavorCounts();
    });
});

// --- Flavor selection (number boxes, up to maxFlavors total) ---
const flavorCountInputs = document.querySelectorAll('#flavor-choices .flavor-count');
const flavorChoices = document.querySelectorAll('#flavor-choices .build-choice');
function updateFlavorInput() {
    let total = 0;
    let values = [];
    document.querySelectorAll('#flavor-choices .build-choice').forEach(c => {
        const count = parseInt(c.querySelector('.flavor-count').value) || 0;
        for (let i = 0; i < count; i++) values.push(c.dataset.value);
        total += count;
    });
    document.getElementById('flavor-input').value = values.join(',');
    return total;
}
function updateFlavorLimitMsg() {
    document.getElementById('flavor-limit-msg').textContent = `You can select exactly ${maxFlavors} flavor${maxFlavors>1?'s':''} (repeat flavors if you wish).`;
}
function clearFlavorCounts() {
    flavorCountInputs.forEach(input => input.value = 0);
    updateFlavorInput();
}
flavorCountInputs.forEach(input => {
    input.addEventListener('input', function() {
        // Enforce maxFlavors total
        let total = updateFlavorInput();
        if (total > maxFlavors) {
            this.value = Math.max(0, parseInt(this.value) - (total - maxFlavors));
            updateFlavorInput();
        }
    });
});
// Click on flavor always increments (never deselects)
flavorChoices.forEach(choice => {
    choice.addEventListener('click', function(e) {
        // Only increment if not clicking the input itself
        if (e.target.classList.contains('flavor-count')) return;
        const input = this.querySelector('.flavor-count');
        let val = parseInt(input.value) || 0;
        if (val < 3) {
            input.value = val + 1;
            input.dispatchEvent(new Event('input'));
        }
    });
});
updateFlavorLimitMsg();

// --- Toppings selection (number boxes, up to 3 total, none required) ---
const toppingCountInputs = document.querySelectorAll('#toppings-choices .topping-count');
const toppingChoices = document.querySelectorAll('#toppings-choices .build-choice');
function updateToppingInput() {
    let total = 0;
    let values = [];
    document.querySelectorAll('#toppings-choices .build-choice').forEach(c => {
        const count = parseInt(c.querySelector('.topping-count').value) || 0;
        for (let i = 0; i < count; i++) values.push(c.dataset.value);
        total += count;
    });
    document.getElementById('toppings-input').value = values.join(',');
    return total;
}
function updateToppingLimitMsg() {
    document.getElementById('topping-limit-msg').textContent = `You can select up to 3 toppings (repeat toppings if you wish, all optional).`;
}
function clearToppingCounts() {
    toppingCountInputs.forEach(input => input.value = 0);
    updateToppingInput();
}
toppingCountInputs.forEach(input => {
    input.addEventListener('input', function() {
        // Enforce max 3 total
        let total = updateToppingInput();
        if (total > 3) {
            this.value = Math.max(0, parseInt(this.value) - (total - 3));
            updateToppingInput();
        }
    });
});
// Click on topping always increments (never deselects)
toppingChoices.forEach(choice => {
    choice.addEventListener('click', function(e) {
        // Only increment if not clicking the input itself
        if (e.target.classList.contains('topping-count')) return;
        const input = this.querySelector('.topping-count');
        let val = parseInt(input.value) || 0;
        if (val < 3) {
            input.value = val + 1;
            input.dispatchEvent(new Event('input'));
        }
    });
});
updateToppingLimitMsg();

// --- Prevent form submit if required not selected or flavor count wrong or too many toppings ---
const form = document.getElementById('build-form');
form.addEventListener('submit', function(e) {
    if (!document.getElementById('size-input').value) {
        alert('Please select a size.');
        e.preventDefault();
        return;
    }
    const selectedFlavors = document.getElementById('flavor-input').value.split(',').filter(Boolean);
    if (selectedFlavors.length !== maxFlavors) {
        alert(`Please select exactly ${maxFlavors} flavor${maxFlavors>1?'s':''}.`);
        e.preventDefault();
        return;
    }
    const selectedToppings = document.getElementById('toppings-input').value.split(',').filter(Boolean);
    if (selectedToppings.length > 3) {
        alert('You can select up to 3 toppings.');
        e.preventDefault();
        return;
    }
});

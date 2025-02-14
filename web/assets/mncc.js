document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('mnccForm');
    const errors = document.querySelectorAll('.error');
    const shoes = document.getElementById('shoes');
    const shoeSizeContainer = document.getElementById('shoeSizeContainer');
    const shoeSize = document.getElementById('shoeSize');
    const shoeSizeError = document.getElementById('shoeSizeError');

    // Show/hide shoe size based on selection
    shoeSizeContainer.classList.toggle('hidden', !shoes.checked);
    shoes.addEventListener('change', function() {
        shoeSizeContainer.classList.toggle('hidden', !this.checked);
    });

    // Prevent multiple form submissions
    let submissionInProgress = false;

    // Handle form submission
    form.addEventListener('submit', async function(event) {
        event.preventDefault();
        let isValid = true;

        // Exit early if a submission is already in progress
        if (submissionInProgress) {
            return;
        }

        // Reset errors
        errors.forEach(error => error.classList.add('hidden'));

        // Validate required fields
        ['name', 'phoneNumber', 'email'].forEach(field => {
            const input = document.getElementById(field);
            if (!input.value.trim()) {
                document.getElementById(field + 'Error').classList.remove('hidden');
                isValid = false;
                return;
            }

            if (input.dataset.validation && !new RegExp(input.dataset.validation).test(input.value)) {
                document.getElementById(field + 'FormatError').classList.remove('hidden');
                isValid = false;
            }
        });

        // Validate waiver checkbox
        if (document.getElementById('waiver').checked === false) {
            document.getElementById('waiverError').classList.remove('hidden');
            isValid = false;
        }

        // Validate shoe size if climbing shoes are needed
        if (shoes.checked && !shoeSize.value.trim()) {
            shoeSizeError.classList.remove('hidden');
            isValid = false;
        }

        if (isValid && !submissionInProgress) {
            submissionInProgress = true;
            document.getElementById('submitButton').disabled = true;

            const formSlug = document.getElementById('formSlug').value;
            const response = await fetch(`/api/v1/submission/${formSlug}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    name: document.getElementById('name').value,
                    pronouns: document.getElementById('pronouns').value,
                    phone_number: document.getElementById('phoneNumber').value,
                    email: document.getElementById('email').value,
                    dietary_restrictions: document.getElementById('dietaryRestrictions').value,
                    waiver_completed: document.getElementById('waiver').checked ? 'Yes' : 'No',
                    shoes_needed: shoes.checked ? 'Yes' : 'No',
                    shoe_size: shoeSize.value,
                    chalk_needed: document.getElementById('chalk').checked ? 'Yes' : 'No'
                })
            });

            if (response.status === 201) {

            } else {
                submissionInProgress = false;
            }
        }
    });
});
document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('mnccForm');
    const errors = document.querySelectorAll('.error');

    // Prevent multiple form submissions
    let submissionInProgress = false;
    
    // Ensure discord username is lowercase
    document.getElementById('discordUsername').addEventListener('blur', function(event) {
        event.target.value = event.target.value.trim().toLowerCase();
    });

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
        document.getElementById('errorMessage').innerText = '';
        document.getElementById('errorMessage').classList.add('hidden');

        // Validate required fields
        ['name', 'phoneNumber', 'email', 'discordUsername'].forEach(field => {
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

        // Validate locations checkbox
        const locations = document.querySelectorAll('[name=locations]:checked');
        if (locations.length === 0) {
            document.getElementById('locationsError').classList.remove('hidden');
            isValid = false;
        }

        // Validate gear dropdown
        if (document.getElementById('gear').value.length === 0) {
            document.getElementById('gearError').classList.remove('hidden');
            isValid = false;
        }

        // Validate employee benefits dropdown
        if (document.getElementById('employeeBenefits').value.length === 0) {
            document.getElementById('employeeBenefitsError').classList.remove('hidden');
            isValid = false;
        }

        if (isValid && !submissionInProgress) {
            grecaptcha.ready(function() {
                grecaptcha.execute(document.getElementById('recaptchaSiteKey').value, {action: 'submit'}).then(async function(token) {
                    submissionInProgress = true;
                    document.getElementById('submitButton').disabled = true;

                    const formSlug = document.getElementById('formSlug').value;
                    try {
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
                                discord_username: document.getElementById('discordUsername').value,
                                memberships: Array.from(document.querySelectorAll('[name=memberships]:checked')).map(location => location.value).join(', '),
                                locations: Array.from(document.querySelectorAll('[name=locations]:checked')).map(location => location.value).join(', '),
                                gear: document.getElementById('gear').value,
                                skills: Array.from(document.querySelectorAll('[name=skills]:checked')).map(location => location.value).join(', '),
                                benefits: document.getElementById('employeeBenefits').value,
                                recaptcha_token: token
                            })
                        });
    
                        if (response.status === 201) {
                            document.getElementById('successMessage').classList.remove('hidden');
                            form.classList.add('hidden');
                        } else {
                            const errorData = await response.json();
                            if (errorData.error) {
                                document.getElementById('errorMessage').innerText = 'An error occurred while submitting the form. Please try again. (' + response.status + ' - ' + errorData.error + ')';
                            } else {
                                document.getElementById('errorMessage').innerText = 'An error occurred while submitting the form. Please try again. (' + response.status + ')';
                            }
                            document.getElementById('errorMessage').classList.remove('hidden');
                        }
                    }
                    finally {
                        submissionInProgress = false;
                        document.getElementById('submitButton').disabled = false;
                    }
                });
            });
        }
    });
});
document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementsByTagName('form')[0];
    const fields = Array.from(document.querySelectorAll('input, select, textarea')).reduce((acc, field) => {
        if (acc[field.name]) {
            if (Array.isArray(acc[field.name])) {
                acc[field.name].push(field);
            }
            else {
                acc[field.name] = [acc[field.name], field];
            }
        } else {
            acc[field.name] = field;
        }
        
        return acc;
    }, {});
    const errors = document.querySelectorAll('.error');

    // Prevent multiple form submissions
    let submissionInProgress = false;

    // Handle form submission
    form.addEventListener('submit', async function(event) {
        event.preventDefault();
        
        let isValid = true;
        let firstError = null;

        // Exit early if a submission is already in progress
        if (submissionInProgress) {
            return;
        }

        // Reset errors
        errors.forEach(error => error.classList.add('hidden'));

        // Validate fields
        Object.values(fields).forEach(field => {
            if ((field.required || field?.dataset?.required === 'true') && !field.value.trim()) {
                document.getElementById(field.name + 'Error').classList.remove('hidden');
                if (!firstError) firstError = field;
                isValid = false;
            } else if (field.value && field.value.length > 0 && field.dataset.validation && !new RegExp(field.dataset.validation).test(field.value)) {
                document.getElementById(field.name + 'FormatError').classList.remove('hidden');
                if (!firstError) firstError = field;
                isValid = false;
            }
        });

        if (isValid && !submissionInProgress) {
            grecaptcha.ready(function() {
                grecaptcha.execute(document.getElementById('recaptchaSiteKey').value, {action: 'submit'}).then(async function(token) {
                    submissionInProgress = true;
                    document.getElementById('submitButton').disabled = true;

                    const formSlug = document.getElementById('formSlug').value;
                    const response = await fetch(`/api/v1/submission/${formSlug}`, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify({
                            fname: fields['fname'].value,
                            lname: fields['lname'].value,
                            pronouns: fields['pronouns'].value,
                            email: fields['email'].value,
                            need: fields['need'].filter((e) => e.checked).map((e) => e.value).join(', '),
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
                
                    submissionInProgress = false;
                });
            });
        } else if (firstError) {
            // Focus on the first error field and scroll into the view.
            firstError.focus();
            firstError.scrollIntoView();
        }
    });
});